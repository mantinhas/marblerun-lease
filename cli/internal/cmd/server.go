// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
		Short: "Initiates a provider side connection to the coordinator. Can run a ping server, or a lease server.",
		Long: `
Gets the root CA of the coordinator and opens a GRPC server to accept either ping or lease requests.`,
		Example: "marblerun server ping example.com:4433",
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(NewPingCmd())
	cmd.AddCommand(NewLeaseCmd())

	return cmd
}

func NewPingCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ping",
		Short: "Starts a ping server",
		Long:  `Starts a ping server that listens for incoming ping requests.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runPingServer,
	}
	cmd.Flags().StringP("cert", "c", "", "PEM encoded admin certificate file (required)")
	cmd.MarkFlagRequired("cert")
	cmd.Flags().StringP("key", "k", "", "PEM encoded admin key file (required)")
	cmd.MarkFlagRequired("key")
	cmd.Flags().IntP("port", "p", 50051, "Port for the coordinator's gRPC server. Different from the port used to connect to the coordinator.")

	return cmd
}

func NewLeaseCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "lease",
		Short: "Starts a lease server",
		Long:  `Starts a lease server that listens for incoming lease requests.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runLeaseServer,
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

func runPingServer(cmd *cobra.Command, args []string) error {
	return runServer(cmd, args, "ping")
}

func runLeaseServer(cmd *cobra.Command, args []string) error {
	return runServer(cmd, args, "lease")
}

func runServer(cmd *cobra.Command, args []string, servertype string) error {
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
	switch servertype {
	case "ping":
		rpc.RegisterProviderServer(s, &pingServer{})
	case "lease":
		rpc.RegisterProviderServer(s, &leaseServer{})
	default:
		log.Fatalf("Unknown server type %s", servertype)
	}

	log.Printf("%s grpc server listening at %v", servertype, lis.Addr())
	go func() {
		log.Fatalf("failed to serve: %v", s.Serve(lis))
	}()

	if servertype == "lease" {
		displayHelpMessage()
		userInputCommands()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	log.Println("Shutting down the server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped.")

	return nil
}

func displayHelpMessage() {
	fmt.Println("--- Lease Server CLI ---")
	fmt.Println("List of commands:")
	fmt.Println("\tchange [duration] - Change the lease duration. Duration is duration notation, e.g. 10s, 5m, 2h, 1d")
	fmt.Println("\thelp - Display this help message")
}

func userInputCommands() {
	for {
		in := bufio.NewReader(os.Stdin)
		command, _ := in.ReadString('\n')
		command = strings.TrimSuffix(command, "\n")
		commandParts := strings.Split(command, " ")
		// print command parts
		switch commandParts[0] {
		case "change":
			if len(commandParts) != 2 {
				fmt.Println("Incorect number of arguments. Usage: change [duration]")
				continue
			}

			duration, err := time.ParseDuration(commandParts[1])
			if err != nil {
				fmt.Printf("Invalid duration \"%s\". Example: 12h30m15s\n", commandParts[1])
				continue
			}
			lease.defaultLeaseTime = duration
			fmt.Printf("Lease duration changed to %s\n", lease.defaultLeaseTime)
		case "help":
			displayHelpMessage()
		default:
			fmt.Printf("Unknown %s command. Type 'help' for a list of commands.\n", commandParts[0])
		}

	}
}

type pingServer struct {
	rpc.UnimplementedProviderServer
}
type leaseServer struct {
	rpc.UnimplementedProviderServer
}

func (s *pingServer) Ping(ctx context.Context, in *rpc.PingReq) (*rpc.PingResp, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		fmt.Printf("Received ping from address %s\n", p.Addr.String())
	} else {
		fmt.Println("Received ping from an unknown address")
	}

	return &rpc.PingResp{Ok: true}, nil
}

type leaseState struct {
	lastLeaseStartTime time.Time
	lastLeaseEndTime   time.Time
	defaultLeaseTime   time.Duration
}

var lease = leaseState{defaultLeaseTime: time.Duration(10) * time.Second}

func (lease *leaseState) addLeaseToState() error {
	if time.Now().Before(lease.lastLeaseStartTime) {
		return fmt.Errorf("A renewal has already been granted. You cannot renew a lease while another offered lease is in queue.")
	}
	leasesHaveBegun := !lease.lastLeaseStartTime.IsZero() || !lease.lastLeaseEndTime.IsZero()

	if !leasesHaveBegun {
		lease.lastLeaseStartTime = time.Now()
	} else {
		lease.lastLeaseStartTime = lease.lastLeaseEndTime
	}
	lease.lastLeaseEndTime = lease.lastLeaseStartTime.Add(lease.defaultLeaseTime)
	fmt.Printf("Lease granted. Start time: %s, End time: %s\n", lease.lastLeaseStartTime.Format(time.RFC3339), lease.lastLeaseEndTime.Format(time.RFC3339))
	return nil
}

func (s *leaseServer) Lease(ctx context.Context, in *rpc.LeaseReq) (*rpc.LeaseOffer, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		fmt.Printf("Received lease request from address %s\n", p.Addr.String())
	} else {
		fmt.Println("Received lease request from an unknown address\n")
	}

	err := lease.addLeaseToState()

	if err != nil {
		return &rpc.LeaseOffer{Ok: false}, nil
	}
	return &rpc.LeaseOffer{Ok: true, LeaseDuration: lease.defaultLeaseTime.String()}, nil
}
