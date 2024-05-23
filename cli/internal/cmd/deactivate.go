// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/edgelesssys/marblerun/cli/internal/rest"
)

func NewDeactivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivates the MarbleRun Coordinator",
		Long: `
Deactivates the MarbleRun Coordinator. This will stop the Coordinator from
serving any requests and will also stop the MarbleRun from running any
applications.`,
		Example: "deactivate example.com:4433",
		Args:    cobra.ExactArgs(1),
		RunE:    runDeactivate,
	}

	cmd.Flags().StringP("cert", "c", "", "PEM encoded admin certificate file (required)")
	cmd.MarkFlagRequired("cert")
	cmd.Flags().StringP("key", "k", "", "PEM encoded admin key file (required)")
	cmd.MarkFlagRequired("key")

	return cmd
}

func runDeactivate(cmd *cobra.Command, args []string) error {
	hostname := args[0]

	client, err := rest.NewAuthenticatedClient(cmd, hostname)
	if err != nil {
		return err
	}

	cmd.Println("Successfully verified Coordinator, now deactivating")

	_, err = client.Post(cmd.Context(), rest.DeactivateEndpoint, rest.ContentJSON, nil)
	if err != nil {
		return fmt.Errorf("deactivating : %w", err)
	}
	cmd.Println("Deactivated coordinator")

	return nil
}
