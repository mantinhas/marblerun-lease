// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"github.com/edgelesssys/marblerun/cli/internal/rest"
	"github.com/edgelesssys/marblerun/cli/rpc"
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Initiates a ping server",
		Long: `
Gets the root CA of the coordinator and opens a GRPC server to accept ping requests.`,
		Example: "server example.com:4433",
		Args:    cobra.ExactArgs(1),
		RunE:    runServer,
	}

	cmd.Flags().StringP("cert", "c", "", "PEM encoded admin certificate file (required)")
	cmd.MarkFlagRequired("cert")
	cmd.Flags().StringP("key", "k", "", "PEM encoded admin key file (required)")
	cmd.MarkFlagRequired("key")
	cmd.Flags().IntP("port", "p", 50051, "Port for the coordinator's gRPC server. Different from the port used to connect to the coordinator.")

	return cmd
}

func getRootCA(cmd *cobra.Command, host string) ([]*pem.Block, error) {
	flags, err := rest.ParseFlags(cmd)
	if err != nil {
		return nil, err
	}

	certs, err := rest.VerifyCoordinator(
		cmd.Context(), cmd.OutOrStdout(), host,
		flags.EraConfig, flags.Insecure, flags.AcceptedTCBStatuses,
	)
	if err != nil {
		return nil, fmt.Errorf("retrieving root certificate from Coordinator: %w", err)
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates received from Coordinator")
	}

	cmd.Println("Successfully retrieved Coordinator root CA certificate")

	return certs, nil
}

func runServer(cmd *cobra.Command, args []string) error {
	hostname := args[0]

	rootCAblock, err := getRootCA(cmd, hostname)
	if err != nil {
		return err
	}

	flags, err := rest.ParseAuthenticatedFlags(cmd)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	rootCAcert, err := x509.ParseCertificate(rootCAblock[len(rootCAblock)-1].Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{flags.ClientCert},
		ClientAuth:   tls.RequireAnyClientCert,
		// We have to manually check if the client certificate matches the root CA
		// Note: a server needs to do this manually, compared to a client where this is done automatically
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("missing client's certificate")
			}
			incomingCert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("failed to parse server's certificate: %v", err)
			}

			if !rootCAcert.Equal(incomingCert) {
				return fmt.Errorf("server's certificate does not match the stored certificate")
			}
			return nil
		},
	}

	creds := credentials.NewTLS(tlsConfig)
	port, _ := cmd.Flags().GetInt("port")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.Creds(creds),
	)
	rpc.RegisterProviderServer(s, &server{})
	log.Printf("grpc server listening at %v", lis.Addr())
	go func() {
		log.Fatalf("failed to serve: %v", s.Serve(lis))
	}()

	// Wait for Ctrl+C to exit.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	log.Println("Shutting down the server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped.")

	return nil
}

type server struct {
	rpc.UnimplementedProviderServer
}

func (s *server) Ping(ctx context.Context, in *rpc.PingReq) (*rpc.PingResp, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		fmt.Printf("Received ping from address %s\n", p.Addr.String())
	} else {
		fmt.Println("Received ping from an unknown address")
	}

	return &rpc.PingResp{Ok: true}, nil
}
