// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package clientapi

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/edgelesssys/marblerun/coordinator/constants"
	"github.com/edgelesssys/marblerun/coordinator/crypto"
	"github.com/edgelesssys/marblerun/coordinator/manifest"
	"github.com/edgelesssys/marblerun/coordinator/seal"
	"github.com/edgelesssys/marblerun/coordinator/state"
	"github.com/edgelesssys/marblerun/coordinator/store"
	"github.com/edgelesssys/marblerun/coordinator/store/request"
	"github.com/edgelesssys/marblerun/coordinator/store/stdstore"
	"github.com/edgelesssys/marblerun/coordinator/store/wrapper"
	"github.com/edgelesssys/marblerun/coordinator/store/wrapper/testutil"
	"github.com/edgelesssys/marblerun/coordinator/user"
	"github.com/edgelesssys/marblerun/test"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestGetCertQuote(t *testing.T) {
	someErr := errors.New("failed")
	// these are not actually root and intermediate certs
	// but we don't care for this test
	rootCert, intermediateCert := test.MustSetupTestCerts(test.RecoveryPrivateKey)

	prepareDefaultStore := func() store.Store {
		s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
		require.NoError(t, wrapper.New(s).PutCertificate(constants.SKCoordinatorRootCert, rootCert))
		require.NoError(t, wrapper.New(s).PutCertificate(constants.SKCoordinatorIntermediateCert, intermediateCert))
		return s
	}

	testCases := map[string]struct {
		store   store.Store
		core    *fakeCore
		wantErr bool
	}{
		"success state accepting Marbles": {
			store: prepareDefaultStore(),
			core: &fakeCore{
				state: state.AcceptingMarbles,
				quote: []byte("quote"),
			},
		},
		"success state accepting manifest": {
			store: prepareDefaultStore(),
			core: &fakeCore{
				state: state.AcceptingManifest,
				quote: []byte("quote"),
			},
		},
		"success state recovery": {
			store: prepareDefaultStore(),
			core: &fakeCore{
				state: state.Recovery,
				quote: []byte("quote"),
			},
		},
		"unsupported state": {
			store: prepareDefaultStore(),
			core: &fakeCore{
				state: state.Uninitialized,
				quote: []byte("quote"),
			},
			wantErr: true,
		},
		"error getting state": {
			store: prepareDefaultStore(),
			core: &fakeCore{
				requireStateErr: someErr,
				quote:           []byte("quote"),
			},
			wantErr: true,
		},
		"root certificate not set": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutCertificate(constants.SKCoordinatorIntermediateCert, intermediateCert))
				return s
			}(),
			core: &fakeCore{
				state: state.AcceptingMarbles,
				quote: []byte("quote"),
			},
			wantErr: true,
		},
		"intermediate certificate not set": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutCertificate(constants.SKCoordinatorRootCert, rootCert))
				return s
			}(),
			core: &fakeCore{
				state: state.AcceptingMarbles,
				quote: []byte("quote"),
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			log, err := zap.NewDevelopment()
			require.NoError(err)
			defer log.Sync()

			api := &ClientAPI{
				core:     tc.core,
				txHandle: tc.store,
				log:      log,
			}

			var intermediateCert, rootCert *x509.Certificate
			if !tc.wantErr {
				intermediateCert = testutil.GetCertificate(t, tc.store, constants.SKCoordinatorIntermediateCert)
				rootCert = testutil.GetCertificate(t, tc.store, constants.SKCoordinatorRootCert)
			}

			cert, quote, err := api.GetCertQuote(context.Background())
			if tc.wantErr {
				assert.Error(err)
				return
			}

			require.NoError(err)
			assert.Equal(tc.core.quote, quote)
			assert.Equal(mustEncodeToPem(t, intermediateCert)+mustEncodeToPem(t, rootCert), cert)
		})
	}
}

func TestGetManifestSignature(t *testing.T) {
	testCases := map[string]struct {
		store   store.Store
		wantErr bool
	}{
		"success": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, s.Put(request.Manifest, []byte("manifest")))
				require.NoError(t, s.Put(request.ManifestSignature, []byte("signature")))
				return s
			}(),
		},
		"GetRawManifest fails": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, s.Put(request.ManifestSignature, []byte("signature")))
				return s
			}(),
			wantErr: true,
		},
		"GetManifestSignature fails": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, s.Put(request.Manifest, []byte("manifest")))
				return s
			}(),
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			log, err := zap.NewDevelopment()
			require.NoError(err)
			defer log.Sync()

			api := &ClientAPI{
				txHandle: tc.store,
				log:      log,
			}

			var rawManifest, manifestSignature, manifestHash []byte
			if !tc.wantErr {
				rawManifest = testutil.GetRawManifest(t, tc.store)
				manifestSignature = testutil.GetManifestSignature(t, tc.store)
				h := sha256.Sum256(rawManifest)
				manifestHash = h[:]
			}

			signature, hash, manifest := api.GetManifestSignature(context.Background())
			if tc.wantErr {
				assert.Nil(signature)
				assert.Nil(hash)
				assert.Nil(manifest)
				return
			}
			assert.Equal(rawManifest, manifest)
			assert.Equal(manifestHash, hash)
			assert.Equal(manifestSignature, signature)
		})
	}
}

func TestGetSecrets(t *testing.T) {
	newUserWithPermissions := func(name string, secretNames ...string) *user.User {
		u := user.NewUser(name, nil)
		u.Assign(user.NewPermission(user.PermissionReadSecret, secretNames))
		return u
	}

	testCases := map[string]struct {
		store   store.Store
		core    *fakeCore
		request []string
		user    *user.User
		wantErr bool
	}{
		"success": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutSecret("secret1", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				require.NoError(t, wrapper.New(s).PutSecret("secret2", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				return s
			}(),
			core: &fakeCore{state: state.AcceptingMarbles},
			request: []string{
				"secret1",
				"secret2",
			},
			user: newUserWithPermissions("test", "secret1", "secret2"),
		},
		"wrong state": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutSecret("secret1", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				require.NoError(t, wrapper.New(s).PutSecret("secret2", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				return s
			}(),
			core: &fakeCore{state: state.AcceptingManifest},
			request: []string{
				"secret1",
				"secret2",
			},
			user:    newUserWithPermissions("test", "secret1", "secret2"),
			wantErr: true,
		},
		"user is missing permissions": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutSecret("secret1", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				require.NoError(t, wrapper.New(s).PutSecret("secret2", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				return s
			}(),
			core: &fakeCore{state: state.AcceptingMarbles},
			request: []string{
				"secret1",
				"secret2",
			},
			user:    newUserWithPermissions("test", "secret2"), // only permission for secret2
			wantErr: true,
		},
		"secret does not exist": {
			store: func() store.Store {
				s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
				require.NoError(t, wrapper.New(s).PutSecret("secret1", manifest.Secret{
					Type:    manifest.SecretTypePlain,
					Private: []byte("secret"),
				}))
				return s
			}(),
			core: &fakeCore{state: state.AcceptingMarbles},
			request: []string{
				"secret1",
				"secret2",
			},
			user:    newUserWithPermissions("test", "secret1", "secret2"),
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			log, err := zap.NewDevelopment()
			require.NoError(err)
			defer log.Sync()

			api := &ClientAPI{
				txHandle: tc.store,
				core:     tc.core,
				log:      log,
			}

			storedSecrets := testutil.GetSecretMap(t, tc.store)

			secrets, err := api.GetSecrets(context.Background(), tc.request, tc.user)
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			for name, secret := range secrets {
				assert.Equal(storedSecrets[name], secret)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	testCases := map[string]struct {
		core    *fakeCore
		wantErr bool
	}{
		"success": {
			core: &fakeCore{state: state.AcceptingManifest},
		},
		"error": {
			core: &fakeCore{
				state:       state.AcceptingManifest,
				getStateErr: errors.New("failed"),
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			log, err := zap.NewDevelopment()
			require.NoError(err)
			defer log.Sync()

			api := &ClientAPI{
				core: tc.core,
				log:  log,
			}

			status, _, err := api.GetStatus(context.Background())
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			assert.Equal(tc.core.state, status)
		})
	}
}

func TestGetUpdateLog(t *testing.T) {
	t.Log("WARNING: Missing unit Test for GetUpdateLog")
}

func TestRecover(t *testing.T) {
	someErr := errors.New("failed")
	_, rootCert := test.MustSetupTestCerts(test.RecoveryPrivateKey)
	defaultStore := func() store.Store {
		s := stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), "")
		require.NoError(t, wrapper.New(s).PutCertificate(constants.SKCoordinatorRootCert, rootCert))
		return s
	}

	testCases := map[string]struct {
		store    *fakeStore
		recovery *stubRecovery
		core     *fakeCore
		wantErr  bool
	}{
		"success": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state: state.Recovery,
			},
		},
		"more than one key required": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{
				recoverKeysLeft: 1,
			},
			core: &fakeCore{
				state: state.Recovery,
			},
		},
		"SetRecoveryData fails does not result in error": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{
				setRecoveryDataErr: someErr,
			},
			core: &fakeCore{
				state: state.Recovery,
			},
		},
		"Coordinator not in recovery state": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state: state.AcceptingManifest,
			},
			wantErr: true,
		},
		"RecoverKey fails": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{
				recoverKeyErr: someErr,
			},
			core: &fakeCore{
				state: state.Recovery,
			},
			wantErr: true,
		},
		"LoadState fails": {
			store: &fakeStore{
				store:        defaultStore(),
				loadStateErr: someErr,
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state: state.Recovery,
			},
			wantErr: true,
		},
		"SetEncryptionKey fails": {
			store: &fakeStore{
				store:               defaultStore(),
				setEncryptionKeyErr: someErr,
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state: state.Recovery,
			},
			wantErr: true,
		},
		"GetCertificate fails": {
			store: &fakeStore{
				store: stdstore.New(&seal.MockSealer{}, afero.NewMemMapFs(), ""),
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state: state.Recovery,
			},
			wantErr: true,
		},
		"GenerateQuote fails": {
			store: &fakeStore{
				store: defaultStore(),
			},
			recovery: &stubRecovery{},
			core: &fakeCore{
				state:            state.Recovery,
				generateQuoteErr: someErr,
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			log, err := zap.NewDevelopment()
			require.NoError(err)
			defer log.Sync()

			api := &ClientAPI{
				txHandle: tc.store,
				recovery: tc.recovery,
				core:     tc.core,
				log:      log,
			}

			keysLeft, err := api.Recover(context.Background(), []byte("recoveryKey"))
			if tc.wantErr {
				assert.Error(err)
				return
			}

			require.NoError(err)
			assert.True(tc.core.unlockCalled)
			assert.Equal(tc.recovery.recoverKeysLeft, keysLeft)
			if keysLeft > 0 {
				assert.False(tc.store.loadCalled)
				return
			}
			assert.True(tc.store.setRecoveryDataCalled)
			assert.True(tc.store.loadCalled)
			assert.NotNil(tc.core.quote)
		})
	}
}

func TestSetManifest(t *testing.T) {
	t.Log("WARNING: Missing unit Test for SetManifest")
}

func TestUpdateManifest(t *testing.T) {
	t.Log("WARNING: Missing unit Test for UpdateManifest")
}

func TestVerifyUser(t *testing.T) {
	t.Log("WARNING: Missing unit Test for VerifyUser")
}

func TestWriteSecrets(t *testing.T) {
	t.Log("WARNING: Missing unit Test for WriteSecrets")
}

type fakeCore struct {
	unlockCalled    bool
	state           state.State
	requireStateErr error
	advanceStateErr error
	getStateMsg     string
	// recoveryErr        error
	quote                         []byte
	generateQuoteErr              error
	getStateErr                   error
	generatedSecrets              map[string]manifest.Secret
	generateSecretsErr            error
	setupKeepAlive                error
	setupLeaseKeepAlive           error
	extractKeepAliveSettings      func(manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration)
	extractLeaseKeepAliveSettings func(manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration)
	deactivateMarbles             error
}

func (c *fakeCore) Unlock() {
	c.unlockCalled = true
}

func (c *fakeCore) RequireState(_ context.Context, states ...state.State) error {
	if c.requireStateErr != nil {
		return c.requireStateErr
	}

	for _, state := range states {
		if state == c.state {
			return nil
		}
	}
	return errors.New("core is not in expected state")
}

func (c *fakeCore) AdvanceState(newState state.State, _ interface {
	PutState(state.State) error
	GetState() (state.State, error)
},
) error {
	if c.advanceStateErr != nil {
		return c.advanceStateErr
	}

	if !(c.state < newState && newState < state.Max) {
		panic("invalid state transition")
	}
	c.state = newState
	return nil
}

func (c *fakeCore) GetState(_ context.Context) (state.State, string, error) {
	return c.state, c.getStateMsg, c.getStateErr
}

func (c *fakeCore) GenerateSecrets(
	newSecrets map[string]manifest.Secret, _ uuid.UUID, rootCert *x509.Certificate, privK *ecdsa.PrivateKey, _ *ecdsa.PrivateKey,
) (map[string]manifest.Secret, error) {
	if c.generateSecretsErr != nil || c.generatedSecrets != nil {
		return c.generatedSecrets, c.generateSecretsErr
	}

	secrets := make(map[string]manifest.Secret, len(newSecrets))
	for name, secret := range newSecrets {
		switch secret.Type {
		case manifest.SecretTypeSymmetricKey:
			secret.Public = bytes.Repeat([]byte{0x00}, 32)
			secret.Private = bytes.Repeat([]byte{0x01}, 32)
		case manifest.SecretTypeCertECDSA, manifest.SecretTypeCertED25519, manifest.SecretTypeCertRSA:
			cert, key, err := crypto.GenerateCert([]string{"localhost"}, name, nil, rootCert, privK)
			if err != nil {
				return nil, err
			}
			secret.Public, err = x509.MarshalPKIXPublicKey(key.Public())
			if err != nil {
				return nil, err
			}
			secret.Private, err = x509.MarshalPKCS8PrivateKey(key)
			if err != nil {
				return nil, err
			}
			secret.Cert = manifest.Certificate(*cert)
		}
		secrets[name] = secret
	}
	return secrets, nil
}

func (c *fakeCore) GetQuote() []byte {
	return c.quote
}

func (c *fakeCore) GenerateQuote(quoteData []byte) error {
	if c.generateQuoteErr != nil {
		return c.generateQuoteErr
	}

	quote := sha256.Sum256(quoteData)
	c.quote = quote[:]
	return nil
}
func (c *fakeCore) SetupKeepAlive(string, *x509.Certificate, int, time.Duration, time.Duration, string, []byte, *ecdsa.PrivateKey, *x509.Certificate) error {
	if c.setupKeepAlive != nil {
		return c.setupKeepAlive
	}

	return nil
}

func (c *fakeCore) SetupLeaseKeepAlive(string, *x509.Certificate, int, time.Duration, time.Duration, string, []byte, *ecdsa.PrivateKey, *x509.Certificate, manifest.Deactivation) error {
	if c.setupLeaseKeepAlive != nil {
		return c.setupLeaseKeepAlive
	}

	return nil
}

func (c *fakeCore) ExtractKeepAliveSettings(d manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration) {
	if c.extractKeepAliveSettings != nil {
		s, c, i, d1, d2 := c.extractKeepAliveSettings(d)
		return s, c, i, d1, d2
	}

	// Return some default values if the function is not set
	return "", nil, 0, 0, 0
}
func (c *fakeCore) ExtractLeaseKeepAliveSettings(d manifest.Deactivation) (string, *x509.Certificate, int, time.Duration, time.Duration) {
	if c.extractLeaseKeepAliveSettings != nil {
		s, c, i, d1, d2 := c.extractLeaseKeepAliveSettings(d)
		return s, c, i, d1, d2
	}

	// Return some default values if the function is not set
	return "", nil, 0, 0, 0
}

func (c *fakeCore) DeactivateMarbles(_ context.Context) error {
	if c.deactivateMarbles != nil {
		return c.deactivateMarbles
	}

	return nil
}

func (c *fakeCore) GetCertQuote(wrapper wrapper.Wrapper) (cert string, certQuote []byte, err error) {
	if c.getStateErr != nil {
		return "", nil, c.getStateErr
	}

	return "", c.quote, nil
}

type fakeStore struct {
	store                 store.Store
	beginTransactionErr   error
	setRecoveryDataCalled bool
	recoveryData          []byte
	encryptionKey         []byte
	setEncryptionKeyErr   error
	loadStateRes          []byte
	loadStateErr          error
	loadCalled            bool
}

func (s *fakeStore) BeginTransaction(ctx context.Context) (store.Transaction, error) {
	if s.beginTransactionErr != nil {
		return nil, s.beginTransactionErr
	}
	return s.store.BeginTransaction(ctx)
}

func (s *fakeStore) SetEncryptionKey(key []byte) error {
	if s.setEncryptionKeyErr != nil {
		return s.setEncryptionKeyErr
	}
	s.encryptionKey = key
	return nil
}

func (s *fakeStore) SetRecoveryData(recoveryData []byte) {
	s.setRecoveryDataCalled = true
	s.recoveryData = recoveryData
}

func (s *fakeStore) LoadState() ([]byte, error) {
	s.loadCalled = true
	return s.loadStateRes, s.loadStateErr
}

type stubRecovery struct {
	generateEncryptionKeyRes []byte
	generateEncryptionKeyErr error
	generateRecoveryDataRes  map[string][]byte
	generateRecoveryDataErr  error
	recoverKeyRes            []byte
	recoverKeyErr            error
	recoverKeysLeft          int
	getRecoveryDataRes       []byte
	getRecoveryDataErr       error
	setRecoveryDataErr       error
}

func (s *stubRecovery) GenerateEncryptionKey(_ map[string]string) ([]byte, error) {
	return s.generateEncryptionKeyRes, s.generateEncryptionKeyErr
}

func (s *stubRecovery) GenerateRecoveryData(_ map[string]string) (map[string][]byte, []byte, error) {
	return s.generateRecoveryDataRes, nil, s.generateRecoveryDataErr
}

func (s *stubRecovery) RecoverKey(_ []byte) (int, []byte, error) {
	return s.recoverKeysLeft, s.recoverKeyRes, s.recoverKeyErr
}

func (s *stubRecovery) GetRecoveryData() ([]byte, error) {
	return s.getRecoveryDataRes, s.getRecoveryDataErr
}

func (s *stubRecovery) SetRecoveryData(_ []byte) error {
	return s.setRecoveryDataErr
}

func mustEncodeToPem(t *testing.T, cert *x509.Certificate) string {
	t.Helper()

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if pemCert == nil {
		t.Fatal("failed to encode certificate to PEM")
	}
	return string(pemCert)
}
