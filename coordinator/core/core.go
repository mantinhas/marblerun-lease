// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package core provides the core functionality for the Coordinator object including state transition, APIs for marbles and clients, handling of manifests and the sealing functionalities.
package core

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	providerRPC "github.com/edgelesssys/marblerun/cli/rpc"
	"github.com/edgelesssys/marblerun/coordinator/constants"
	corecrypto "github.com/edgelesssys/marblerun/coordinator/crypto"
	"github.com/edgelesssys/marblerun/coordinator/events"
	"github.com/edgelesssys/marblerun/coordinator/manifest"
	"github.com/edgelesssys/marblerun/coordinator/quote"
	"github.com/edgelesssys/marblerun/coordinator/recovery"
	"github.com/edgelesssys/marblerun/coordinator/rpc"
	"github.com/edgelesssys/marblerun/coordinator/seal"
	"github.com/edgelesssys/marblerun/coordinator/state"
	"github.com/edgelesssys/marblerun/coordinator/store"
	"github.com/edgelesssys/marblerun/coordinator/store/request"
	"github.com/edgelesssys/marblerun/coordinator/store/stdstore"
	"github.com/edgelesssys/marblerun/coordinator/store/wrapper"
	"github.com/edgelesssys/marblerun/util"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

// Core implements the core logic of the Coordinator.

type LeaseState struct {
	start                   time.Time
	afterStartAliveDuration time.Duration
	currentLeaseTimer       *time.Timer
	rwtex                   sync.RWMutex
}

func (c *Core) GetAllowedLeaseTime() (time.Duration, error) {
	c.leaseState.rwtex.RLock()
	defer c.leaseState.rwtex.RUnlock()

	var result time.Duration = c.leaseState.afterStartAliveDuration - time.Since(c.leaseState.start)

	if result <= 0 {
		return 0, errors.New("allowed lease time is non positive")
	}
	return result, nil

}

type Core struct {
	mux sync.Mutex

	quote []byte
	qv    quote.Validator
	qi    quote.Issuer

	recovery recovery.Recovery
	metrics  *coreMetrics

	txHandle transactionHandle

	log      *zap.Logger
	eventlog *events.Log

	leaseState LeaseState

	rpc.UnimplementedMarbleServer
}

// RequireState checks if the Coordinator is in one of the given states.
// This function locks the Core's mutex and therefore should be paired with `defer c.mux.Unlock()`.
func (c *Core) RequireState(ctx context.Context, states ...state.State) error {
	c.mux.Lock()

	getter, rollback, _, err := wrapper.WrapTransaction(ctx, c.txHandle)
	if err != nil {
		return err
	}
	defer rollback()

	curState, err := getter.GetState()
	if err != nil {
		return err
	}
	for _, s := range states {
		if s == curState {
			return nil
		}
	}
	return errors.New("server is not in expected state")
}

// AdvanceState advances the state of the Coordinator.
func (c *Core) AdvanceState(newState state.State, tx interface {
	PutState(state.State) error
	GetState() (state.State, error)
},
) error {
	curState, err := tx.GetState()
	if err != nil {
		return err
	}
	if !(curState < newState && newState < state.Max) {
		panic(fmt.Errorf("cannot advance from %d to %d", curState, newState))
	}
	return tx.PutState(newState)
}

// Unlock the Core's mutex.
func (c *Core) Unlock() {
	c.mux.Unlock()
}

// NewCore creates and initializes a new Core object.
func NewCore(
	dnsNames []string, qv quote.Validator, qi quote.Issuer, txHandle transactionHandle,
	recovery recovery.Recovery, zapLogger *zap.Logger, promFactory *promauto.Factory, eventlog *events.Log,
) (*Core, error) {
	c := &Core{
		qv:       qv,
		qi:       qi,
		recovery: recovery,
		txHandle: txHandle,
		log:      zapLogger,
		eventlog: eventlog,
	}
	c.metrics = newCoreMetrics(promFactory, c, "coordinator")

	zapLogger.Info("Loading state")
	recoveryData, loadErr := txHandle.LoadState()
	if err := c.recovery.SetRecoveryData(recoveryData); err != nil {
		c.log.Error("Could not retrieve recovery data from state. Recovery will be unavailable", zap.Error(err))
	}

	transaction, rollback, commit, err := wrapper.WrapTransaction(context.Background(), c.txHandle)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// set core to uninitialized if no state is set
	var cur_state state.State
	if cur_state, err = transaction.GetState(); err != nil {
		if errors.Is(err, store.ErrValueUnset) {
			if err := transaction.PutState(state.Uninitialized); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if loadErr != nil {
		var keyErr *seal.EncryptionKeyError
		if !errors.As(loadErr, &keyErr) {
			return nil, loadErr
		}
		// sealed state was found, but couldn't be decrypted, go to recovery mode or reset manifest
		c.log.Error("Failed to decrypt sealed state. Proceeding with a new state. Use the /recover API endpoint to load an old state, or submit a new manifest to overwrite the old state. Look up the documentation for more information on how to proceed.")
		if err := c.setCAData(dnsNames, transaction); err != nil {
			return nil, err
		}
		if err := c.AdvanceState(state.Recovery, transaction); err != nil {
			return nil, err
		}
	} else if _, err := transaction.GetRawManifest(); errors.Is(err, store.ErrValueUnset) {
		// no state was found, wait for manifest
		c.log.Info("No sealed state found. Proceeding with new state.")
		if err := c.setCAData(dnsNames, transaction); err != nil {
			return nil, err
		}
		if err := transaction.PutState(state.AcceptingManifest); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// recovered from a sealed state, reload components and finish the store transaction
		txHandle.SetRecoveryData(recoveryData)
	}

	rootCert, err := transaction.GetCertificate(constants.SKCoordinatorRootCert)
	if err != nil {
		return nil, err
	}

	if err := commit(context.Background()); err != nil {
		return nil, fmt.Errorf("committing state: %w", err)
	}

	err = c.GenerateQuote(rootCert.Raw)

	if cur_state == state.AcceptingMarbles {
		if manifest, err := transaction.GetManifest(); err != nil {
			return nil, err
		} else {

			privk, err := transaction.GetPrivateKey(constants.SKCoordinatorRootKey)
			if err != nil {
				return nil, fmt.Errorf("loading root private key from store: %w", err)
			}

			rootCert, err := transaction.GetCertificate(constants.SKCoordinatorRootCert)
			if err != nil {
				return nil, fmt.Errorf("loading root certificate from store: %w", err)
			}

			rootCertString, quote, err := c.GetCertQuote(transaction)
			if err != nil {
				return nil, err
			}

			trust_protocol := manifest.DeactivationSettings["Coordinator"].TrustProtocol

			c.log.Info("Trust protocol", zap.String("trust_protocol", trust_protocol))

			switch trust_protocol {
			case "ping":
				url, manifestCertificate, retries, pingInterval, retryInterval := c.ExtractKeepAliveSettings(manifest.DeactivationSettings["Coordinator"])
				if err := c.SetupKeepAlive(url, manifestCertificate, retries, pingInterval, retryInterval, rootCertString, quote, privk, rootCert); err != nil {
					return nil, err
				}
			case "lease":
				url, manifestCertificate, retries, leaseInterval, retryInterval := c.ExtractLeaseKeepAliveSettings(manifest.DeactivationSettings["Coordinator"])
				if err := c.SetupLeaseKeepAlive(url, manifestCertificate, retries, leaseInterval, retryInterval, rootCertString, quote, privk, rootCert); err != nil {
					return nil, err
				}
			}
		}
	}

	return c, err
}

// NewCoreWithMocks creates a new core object with quote and seal mocks for testing.
func NewCoreWithMocks() *Core {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	validator := quote.NewMockValidator()
	issuer := quote.NewMockIssuer()
	sealer := &seal.MockSealer{}
	recovery := recovery.NewSinglePartyRecovery()
	core, err := NewCore([]string{"localhost"}, validator, issuer, stdstore.New(sealer, afero.Afero{Fs: afero.NewMemMapFs()}, ""), recovery, zapLogger, nil, nil)
	if err != nil {
		panic(err)
	}
	return core
}

// inSimulationMode returns true if we operate in OE_SIMULATION mode.
func (c *Core) inSimulationMode() bool {
	return len(c.quote) == 0
}

// GetTLSConfig gets the core's TLS configuration.
func (c *Core) GetTLSConfig() (*tls.Config, error) {
	return &tls.Config{
		GetCertificate: c.GetTLSRootCertificate,
		ClientAuth:     tls.RequestClientCert,
	}, nil
}

// GetTLSRootCertificate creates a TLS certificate for the Coordinators self-signed x509 certificate.
//
// This function initializes a read transaction and should not be called from other functions with ongoing transactions.
func (c *Core) GetTLSRootCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	data, rollback, _, err := wrapper.WrapTransaction(clientHello.Context(), c.txHandle)
	if err != nil {
		return nil, err
	}
	defer rollback()

	curState, err := data.GetState()
	if err != nil {
		return nil, err
	}
	if curState == state.Uninitialized {
		return nil, errors.New("don't have a cert yet")
	}

	rootCert, err := data.GetCertificate(constants.SKCoordinatorRootCert)
	if err != nil {
		return nil, err
	}
	rootPrivK, err := data.GetPrivateKey(constants.SKCoordinatorRootKey)
	if err != nil {
		return nil, err
	}

	return util.TLSCertFromDER(rootCert.Raw, rootPrivK), nil
}

// GetTLSMarbleRootCertificate creates a TLS certificate for the Coordinator's x509 marbleRoot certificate.
//
// This function initializes a read transaction and should not be called from other functions with ongoing transactions.
func (c *Core) GetTLSMarbleRootCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	data, rollback, _, err := wrapper.WrapTransaction(clientHello.Context(), c.txHandle)
	if err != nil {
		return nil, err
	}
	defer rollback()

	curState, err := data.GetState()
	if err != nil {
		return nil, err
	}
	if curState == state.Uninitialized {
		return nil, errors.New("don't have a cert yet")
	}

	marbleRootCert, err := data.GetCertificate(constants.SKMarbleRootCert)
	if err != nil {
		return nil, err
	}
	intermediatePrivK, err := data.GetPrivateKey(constants.SKCoordinatorIntermediateKey)
	if err != nil {
		return nil, err
	}

	return util.TLSCertFromDER(marbleRootCert.Raw, intermediatePrivK), nil
}

// GetQuote returns the quote of the Coordinator.
func (c *Core) GetQuote() []byte {
	return c.quote
}

// GenerateQuote generates a quote for the Coordinator using the given certificate.
// If no quote can be generated due to the system not supporting SGX, no error is returned,
// and the Coordinator proceeds to run in simulation mode.
func (c *Core) GenerateQuote(cert []byte) error {
	c.log.Info("Generating quote")
	quote, err := c.qi.Issue(cert)
	if err != nil {
		if err.Error() == "OE_UNSUPPORTED" {
			c.log.Warn("Failed to get quote. Proceeding in simulation mode.", zap.Error(err))
			// If we run in SimulationMode we get OE_UNSUPPORTED error here
			// For testing purpose we do not want to just fail here
			// Instead we store an empty quote that will make it transparent to the client that the integrity of the mesh can not be guaranteed.
			return nil
		}
		return QuoteError{err}
	}

	c.quote = quote

	return nil
}

func getClientTLSCert(ctx context.Context) *x509.Certificate {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil
	}
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	// the following check is just for safety (not for security)
	if !ok || len(tlsInfo.State.PeerCertificates) == 0 {
		return nil
	}
	return tlsInfo.State.PeerCertificates[0]
}

// GetState returns the current state of the Coordinator.
func (c *Core) GetState(ctx context.Context) (state.State, string, error) {
	data, rollback, _, err := wrapper.WrapTransaction(ctx, c.txHandle)
	if err != nil {
		return -1, "Cannot determine coordinator status.", fmt.Errorf("initializing read transaction: %w", err)
	}
	defer rollback()

	curState, err := data.GetState()
	if err != nil {
		return -1, "Cannot determine coordinator status.", err
	}

	var status string

	switch curState {
	case state.Recovery:
		status = "Coordinator is in recovery mode. Either upload a key to unseal the saved state, or set a new manifest. For more information on how to proceed, consult the documentation."
	case state.AcceptingManifest:
		status = "Coordinator is ready to accept a manifest."
	case state.AcceptingMarbles:
		status = "Coordinator is running correctly and ready to accept marbles."
	default:
		return -1, "Cannot determine coordinator status.", errors.New("cannot determine coordinator status")
	}

	return curState, status, nil
}

// GenerateSecrets generates secrets for the given manifest and parent certificate.
func (c *Core) GenerateSecrets(
	secrets map[string]manifest.Secret, id uuid.UUID,
	parentCertificate *x509.Certificate, parentPrivKey *ecdsa.PrivateKey, rootPrivK *ecdsa.PrivateKey,
) (map[string]manifest.Secret, error) {
	// Create a new map so we do not overwrite the entries in the manifest
	newSecrets := make(map[string]manifest.Secret)

	// Generate secrets
	for name, secret := range secrets {
		// Skip user defined secrets, these will be uploaded by a user
		if secret.UserDefined {
			continue
		}

		// Skip secrets from wrong context
		if secret.Shared != (id == uuid.Nil) {
			continue
		}

		c.log.Info("Generating secret", zap.String("name", name), zap.String("type", secret.Type), zap.Uint("size", secret.Size))
		switch secret.Type {
		// Raw = Symmetric Key
		case manifest.SecretTypeSymmetricKey:
			// Check secret size
			if secret.Size == 0 || secret.Size%8 != 0 {
				return nil, fmt.Errorf("invalid secret size: %v", name)
			}

			var generatedValue []byte
			// If a secret is shared, we generate a completely random key. If a secret is constrained to a marble, we derive a key from the core's private key.
			if secret.Shared {
				generatedValue = make([]byte, secret.Size/8)
				_, err := rand.Read(generatedValue)
				if err != nil {
					return nil, err
				}
			} else {
				salt := id.String() + name
				secretKeyDerive := rootPrivK.D.Bytes()
				var err error
				generatedValue, err = util.DeriveKey(secretKeyDerive, []byte(salt), secret.Size/8)
				if err != nil {
					return nil, err
				}
			}

			// Get secret object from manifest, create a copy, modify it and put in in the new map so we do not overwrite the manifest entries
			secret.Private = generatedValue
			secret.Public = generatedValue

			newSecrets[name] = secret

		case manifest.SecretTypeCertRSA:
			// Generate keys
			privKey, err := rsa.GenerateKey(rand.Reader, int(secret.Size))
			if err != nil {
				c.log.Error("Failed to generate RSA key", zap.Error(err))
				return nil, err
			}

			// Generate certificate
			newSecrets[name], err = c.generateCertificateForSecret(secret, parentCertificate, parentPrivKey, privKey, &privKey.PublicKey)
			if err != nil {
				return nil, err
			}

		case manifest.SecretTypeCertED25519:
			if secret.Size != 0 {
				return nil, fmt.Errorf("invalid secret size for cert-ed25519, none is expected. given: %v", name)
			}

			// Generate keys
			pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				c.log.Error("Failed to generate ed25519 key", zap.Error(err))
				return nil, err
			}

			// Generate certificate
			newSecrets[name], err = c.generateCertificateForSecret(secret, parentCertificate, parentPrivKey, privKey, pubKey)
			if err != nil {
				return nil, err
			}

		case manifest.SecretTypeCertECDSA:
			var curve elliptic.Curve

			switch secret.Size {
			case 224:
				curve = elliptic.P224()
			case 256:
				curve = elliptic.P256()
			case 384:
				curve = elliptic.P384()
			case 521:
				curve = elliptic.P521()
			default:
				c.log.Error("ECDSA secrets only support P224, P256, P384 and P521 as curve. Check the supplied size.", zap.String("name", name), zap.String("type", secret.Type), zap.Uint("size", secret.Size))
				return nil, fmt.Errorf("unsupported size %d: does not map to a supported curve", secret.Size)
			}

			// Generate keys
			privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				c.log.Error("Failed to generate ECSDA key", zap.Error(err))
				return nil, err
			}

			// Generate certificate
			newSecrets[name], err = c.generateCertificateForSecret(secret, parentCertificate, parentPrivKey, privKey, &privKey.PublicKey)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unsupported secret of type %s", secret.Type)
		}
	}

	return newSecrets, nil
}

func (c *Core) generateCertificateForSecret(secret manifest.Secret, parentCertificate *x509.Certificate, parentPrivKey *ecdsa.PrivateKey, privKey crypto.PrivateKey, pubKey crypto.PublicKey) (manifest.Secret, error) {
	// Load given information from manifest as template
	template := x509.Certificate(secret.Cert)

	// Define or overwrite some values for sane standards
	if template.DNSNames == nil {
		template.DNSNames = []string{"localhost"}
	}
	if template.IPAddresses == nil {
		template.IPAddresses = util.DefaultCertificateIPAddresses
	}
	if template.KeyUsage == 0 {
		template.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	}
	if template.ExtKeyUsage == nil {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	}
	if template.Subject.CommonName == "" {
		if len(template.DNSNames) == 1 {
			template.Subject.CommonName = template.DNSNames[0]
		} else {
			template.Subject.CommonName = "MarbleRun Generated Certificate"
		}
	}
	var err error
	template.SerialNumber, err = util.GenerateCertificateSerialNumber()
	if err != nil {
		c.log.Error("No serial number supplied; random number generation failed.", zap.Error(err))
		return manifest.Secret{}, err
	}

	template.BasicConstraintsValid = true
	template.NotBefore = time.Now()

	// If NotAfter is not set, we will use ValidFor for the end of the certificate lifetime. This can only happen once on initial manifest set
	if template.NotAfter.IsZero() {
		// User can specify a duration in days, otherwise it's one year by default
		if secret.ValidFor == 0 {
			secret.ValidFor = 365
		}

		template.NotAfter = time.Now().AddDate(0, 0, int(secret.ValidFor))
	} else if secret.ValidFor != 0 {
		// reset expiration date for private secrets
		if !secret.Shared {
			template.NotAfter = time.Now().AddDate(0, 0, int(secret.ValidFor))
		}
	}

	// Generate certificate with given public key
	secretCertRaw, err := x509.CreateCertificate(rand.Reader, &template, parentCertificate, pubKey, parentPrivKey)
	if err != nil {
		c.log.Error("Failed to generate X.509 certificate", zap.Error(err))
		return manifest.Secret{}, err
	}

	cert, err := x509.ParseCertificate(secretCertRaw)
	if err != nil {
		c.log.Error("Failed to parse newly generated X.509 certificate", zap.Error(err))
		return manifest.Secret{}, err
	}

	// Assemble secret object
	secret.Cert = manifest.Certificate(*cert)
	secret.Private, err = x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		c.log.Error("Failed to marshal private key to secret object", zap.Error(err))
		return manifest.Secret{}, err
	}
	secret.Public, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		c.log.Error("Failed to marshal public key to secret object", zap.Error(err))
		return manifest.Secret{}, err
	}

	return secret, nil
}

func (c *Core) setCAData(dnsNames []string, putter interface {
	PutCertificate(name string, cert *x509.Certificate) error
	PutPrivateKey(name string, key *ecdsa.PrivateKey) error
},
) error {
	rootCert, rootPrivK, err := corecrypto.GenerateCert(dnsNames, constants.CoordinatorName, nil, nil, nil)
	if err != nil {
		return err
	}
	// Creating a cross-signed intermediate cert. See https://github.com/edgelesssys/marblerun/issues/175
	intermediateCert, intermediatePrivK, err := corecrypto.GenerateCert(dnsNames, constants.CoordinatorIntermediateName, nil, rootCert, rootPrivK)
	if err != nil {
		return err
	}
	marbleRootCert, _, err := corecrypto.GenerateCert(dnsNames, constants.CoordinatorIntermediateName, intermediatePrivK, nil, nil)
	if err != nil {
		return err
	}

	if err := putter.PutCertificate(constants.SKCoordinatorRootCert, rootCert); err != nil {
		return err
	}
	if err := putter.PutCertificate(constants.SKCoordinatorIntermediateCert, intermediateCert); err != nil {
		return err
	}
	if err := putter.PutCertificate(constants.SKMarbleRootCert, marbleRootCert); err != nil {
		return err
	}
	if err := putter.PutPrivateKey(constants.SKCoordinatorRootKey, rootPrivK); err != nil {
		return err
	}
	if err := putter.PutPrivateKey(constants.SKCoordinatorIntermediateKey, intermediatePrivK); err != nil {
		return err
	}

	return nil
}

// QuoteError is returned when the quote could not be retrieved.
type QuoteError struct {
	err error
}

// Error returns the error message.
func (e QuoteError) Error() string {
	return fmt.Sprintf("failed to get quote: %v", e.err)
}

type transactionHandle interface {
	BeginTransaction(context.Context) (store.Transaction, error)
	SetEncryptionKey([]byte) error
	SetRecoveryData([]byte)
	LoadState() ([]byte, error)
}

func (c *Core) SetupLeaseKeepAlive(connectionURL string, certificate *x509.Certificate, retries int, leaseTime time.Duration, retryInterval time.Duration, rootCertString string, quote []byte, privk *ecdsa.PrivateKey, rootCert *x509.Certificate) error {
	clientCert := util.TLSCertFromDER(rootCert.Raw, privk)

	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{*clientCert},
		InsecureSkipVerify: true, // We'll handle verification
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("missing server's certificate")
			}
			incomingCert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("failed to parse server's certificate: %v", err)
			}

			if !certificate.Equal(incomingCert) {
				return fmt.Errorf("server's certificate does not match the stored certificate")
			}
			return nil
		},
	}

	creds := credentials.NewTLS(&tlsConfig)

	connection, err := grpc.Dial(connectionURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}

	client := providerRPC.NewProviderClient(connection)

	go func() {
		defer connection.Close()

		var successfulLease bool

		for {
			// Timeout for the lease
			c.log.Info("Starting Lease:", zap.Duration("leaseTime", leaseTime))
			c.leaseState.rwtex.Lock()
			c.leaseState.start = time.Now()
			c.leaseState.afterStartAliveDuration = leaseTime
			c.leaseState.currentLeaseTimer = time.NewTimer(leaseTime)
			c.leaseState.rwtex.Unlock()

			// For this lease, with duration X, after X/2 seconds, the client will send a LeaseReq
			resultChan := make(chan string)
			errorChan := make(chan error)
			done := make(chan bool)
			successfulLease = false
			time.AfterFunc(leaseTime/2, func() {
				for i := 0; i < retries; i++ {
					func() error {
						c.log.Info("Requesting Lease...", zap.Int("retry", i+1), zap.Int("retries", retries))
						ctx, cancel := context.WithTimeout(context.Background(), leaseTime/2)
						defer cancel()

						// Attempt to re-dial if necessary
						if connection.GetState() != connectivity.Ready {
							connection, err = grpc.Dial(connectionURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
							if err != nil {
								c.log.Error("re-dial attempt failed: %v", zap.Error(err))
								return err
							}
							client = providerRPC.NewProviderClient(connection)
						}

						resp, err := client.Lease(ctx, &providerRPC.LeaseReq{})

						if err != nil {
							c.log.Error("LeaseReq failed: %v", zap.Error(err))
							return err
						}
						if !resp.Ok {
							c.log.Error("LeaseReq failed: response not ok")
							return errors.New("LeaseReq failed: response not ok")
						}

						done <- true
						resultChan <- resp.LeaseDuration
						errorChan <- nil
						successfulLease = true
						return nil
					}()
					if successfulLease {
						break
					}
					time.Sleep(retryInterval)
				}
			})

			waitOutLease := func() {
				<-c.leaseState.currentLeaseTimer.C
				c.log.Info("Ended Lease.")
			}

			// defer waiting for the timer to finish if it didnt finish yet
			select {
			case <-done:
				result, err := <-resultChan, <-errorChan
				if err != nil {
					c.log.Error("LeaseReq failed: %v", zap.Error(err))
					waitOutLease()
					os.Exit(1)
					return
				} else {
					c.leaseState.rwtex.Lock()
					leaseTime, err = time.ParseDuration(result)
					c.leaseState.afterStartAliveDuration += leaseTime
					c.leaseState.rwtex.Unlock()
					c.log.Info("Lease offered", zap.Duration("leaseTime", leaseTime))
					waitOutLease()
				}
			case <-c.leaseState.currentLeaseTimer.C:
				// Lease expired
				c.log.Info("Lease expired without renewal after all retries failed. Exiting.")
				os.Exit(1)
				return
			}
		}
	}()

	return nil

}

func (c *Core) SetupKeepAlive(connectionURL string, certificate *x509.Certificate, retries int, pingInterval time.Duration, retryInterval time.Duration, rootCertString string, quote []byte, privk *ecdsa.PrivateKey, rootCert *x509.Certificate) error {
	clientCert := util.TLSCertFromDER(rootCert.Raw, privk)

	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{*clientCert},
		InsecureSkipVerify: true, // We'll handle verification
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("missing server's certificate")
			}
			incomingCert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("failed to parse server's certificate: %v", err)
			}

			if !certificate.Equal(incomingCert) {
				return fmt.Errorf("server's certificate does not match the stored certificate")
			}
			return nil
		},
	}

	creds := credentials.NewTLS(&tlsConfig)

	connection, err := grpc.Dial(connectionURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}

	client := providerRPC.NewProviderClient(connection)

	go func() {
		defer connection.Close()

		ticker := time.NewTicker(pingInterval)
		for range ticker.C {
			successfulPing := false
			for i := 0; i < retries; i++ {
				err := func() error {
					ctx, cancel := context.WithTimeout(context.Background(), pingInterval)
					defer cancel()

					// Attempt to re-dial if necessary
					if connection.GetState() != connectivity.Ready {
						connection, err = grpc.Dial(connectionURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err != nil {
							c.log.Error("re-dial attempt failed: %v", zap.Error(err))
							return err
						}
						client = providerRPC.NewProviderClient(connection)
					}

					resp, err := client.Ping(ctx, &providerRPC.PingReq{Cert: rootCertString, Quote: quote})
					if err != nil {
						// Handle connection error
						c.log.Error("ping attempt failed: %v", zap.Error(err))
						return err
					}

					if !resp.Ok {
						// Handle response error
						c.log.Error("ping attempt failed: response not ok")
						return errors.New("ping failed: response not ok")
					}

					c.log.Info("ping ok")
					successfulPing = true
					return nil
				}()

				if err == nil {
					break
				}

				if i < retries-1 { // If it's not the last retry, wait for the retryInterval
					time.Sleep(retryInterval)
				}
			}

			if !successfulPing {
				c.log.Error("exiting after all retries failed")
				os.Exit(1)
			}
		}
	}()

	return nil
}

func (c *Core) ExtractKeepAliveSettings(manifestDeactivation manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration) {
	// Extract and convert values from the manifest
	connectionURL := manifestDeactivation.ConnectionUrl
	// if ConnectionUrl starts with "http://" or "https://", remove it
	connectionURL = strings.TrimPrefix(connectionURL, "http://")
	connectionURL = strings.TrimPrefix(connectionURL, "https://")
	retries := manifestDeactivation.PingSettings.Retries
	pingInterval, _ := time.ParseDuration(manifestDeactivation.PingSettings.RequestInterval) // _ because we know it's valid from previous validation
	retryInterval, _ := time.ParseDuration(manifestDeactivation.PingSettings.RetryInterval)  // _ because we know it's valid from previous validation

	connectionCertificateString := manifestDeactivation.ConnectionCertificate
	block, _ := pem.Decode([]byte(connectionCertificateString))

	var connectionCertificate *x509.Certificate
	if block != nil {
		connectionCertificate, _ = x509.ParseCertificate(block.Bytes)
	}

	return connectionURL, connectionCertificate, retries, pingInterval, retryInterval

}

func (c *Core) ExtractLeaseKeepAliveSettings(manifestDeactivation manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration) {
	// Extract and convert values from the manifest
	connectionURL := manifestDeactivation.ConnectionUrl
	// if ConnectionUrl starts with "http://" or "https://", remove it
	connectionURL = strings.TrimPrefix(connectionURL, "http://")
	connectionURL = strings.TrimPrefix(connectionURL, "https://")
	retries := manifestDeactivation.LeaseSettings.Retries
	leaseInterval, _ := time.ParseDuration(manifestDeactivation.LeaseSettings.RequestInterval) // _ because we know it's valid from previous validation
	retryInterval, _ := time.ParseDuration(manifestDeactivation.LeaseSettings.RetryInterval)   // _ because we know it's valid from previous validation

	connectionCertificateString := manifestDeactivation.ConnectionCertificate
	block, _ := pem.Decode([]byte(connectionCertificateString))

	var connectionCertificate *x509.Certificate
	if block != nil {
		connectionCertificate, _ = x509.ParseCertificate(block.Bytes)
	}

	return connectionURL, connectionCertificate, retries, leaseInterval, retryInterval

}

// call deactivate on all marbles
func (c *Core) DeactivateMarbles(ctx context.Context) error {
	data, rollback, _, err := wrapper.WrapTransaction(ctx, c.txHandle)
	if err != nil {
		return err
	}
	defer rollback()

	marbleIpIter, err := data.GetIterator(request.MarbleIP)
	if err != nil {
		return err
	}

	marbleRootCertificate, err := data.GetCertificate(constants.SKMarbleRootCert)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(marbleRootCertificate)

	privk, _ := data.GetPrivateKey(constants.SKCoordinatorIntermediateKey)

	creds, _ := util.LoadGRPCTLSCredentials(marbleRootCertificate, privk, certPool, false)

	for marbleIpIter.HasNext() {
		name, err := marbleIpIter.GetNext()
		if err != nil {
			return err
		}
		ip, err := data.GetMarbleIP(name)
		if err != nil {
			return err
		}
		// url is ip with port 50060
		url := fmt.Sprintf("%s%s", ip, constants.MarbleDeactivationPort)
		c.log.Info("Deactivating marble", zap.String("url", url))

		conn, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
		if err != nil {
			return err
		}
		defer conn.Close()

		client := rpc.NewMarbleClient(conn)
		_, err = client.Deactivate(ctx, &rpc.DeactivateReq{})
		if err != nil {
			return err
		}

		c.log.Info("Deactivated marble", zap.String("url", url))

	}

	return nil

}

// GetCertQuote
func (c *Core) GetCertQuote(wrapper wrapper.Wrapper) (cert string, certQuote []byte, err error) {
	rootCert, err := wrapper.GetCertificate(constants.SKCoordinatorRootCert)
	if err != nil {
		return "", nil, fmt.Errorf("loading root certificate from store: %w", err)
	}
	if rootCert == nil {
		return "", nil, errors.New("loaded nil root certificate from store")
	}

	intermediateCert, err := wrapper.GetCertificate(constants.SKCoordinatorIntermediateCert)
	if err != nil {
		return "", nil, fmt.Errorf("loading intermediate certificate from store: %w", err)
	}
	if intermediateCert == nil {
		return "", nil, errors.New("loaded nil intermediate certificate from store")
	}

	pemCertRoot := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCert.Raw})
	if len(pemCertRoot) <= 0 {
		return "", nil, errors.New("pem.EncodeToMemory failed for root certificate")
	}

	// Include intermediate certificate if a manifest has been set
	pemCertIntermediate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: intermediateCert.Raw})
	if len(pemCertIntermediate) <= 0 {
		return "", nil, errors.New("pem.EncodeToMemory failed for intermediate certificate")
	}

	return string(pemCertIntermediate) + string(pemCertRoot), c.GetQuote(), nil
}
