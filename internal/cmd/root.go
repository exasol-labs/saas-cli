// Package cmd implements the CLI command tree.
package cmd

import (
	"fmt"
	"os"

	"github.com/exasol-labs/saas-cli/internal/api"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags. Defaults to "dev" for local builds.
var Version = "dev"

// Config holds resolved runtime configuration shared across subcommands.
type Config struct {
	Token     string
	AccountID string
	BaseURL   string
}

// NewRootCmd creates a new root command with all flags configured.
func NewRootCmd() (*cobra.Command, *Config) {
	cfg := &Config{}
	var tokenFlag string
	var accountIDFlag string

	root := &cobra.Command{
		Use:     "exasol-saas",
		Version: Version,
		Short:   "CLI for managing Exasol SaaS resources",
		Long: `exasol-saas is a command-line tool for managing your Exasol SaaS
databases, clusters, and other cloud resources.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.HasParent() {
				return nil
			}

			token := tokenFlag
			if token == "" {
				token = os.Getenv("EXASOL_SAAS_TOKEN")
			}
			if token == "" {
				return fmt.Errorf("authentication token is required: set --token flag or EXASOL_SAAS_TOKEN environment variable")
			}
			cfg.Token = token

			accountID := accountIDFlag
			if accountID == "" {
				accountID = os.Getenv("EXASOL_SAAS_ACCOUNT_ID")
			}
			if accountID == "" {
				return fmt.Errorf("account ID is required: set --account-id flag or EXASOL_SAAS_ACCOUNT_ID environment variable")
			}
			cfg.AccountID = accountID

			if cfg.BaseURL == "" {
				cfg.BaseURL = api.DefaultBaseURL
			}

			return nil
		},
	}

	root.PersistentFlags().StringVar(&tokenFlag, "token", "", "Exasol SaaS personal access token (overrides EXASOL_SAAS_TOKEN env var)")
	root.PersistentFlags().StringVar(&accountIDFlag, "account-id", "", "Exasol SaaS account ID (overrides EXASOL_SAAS_ACCOUNT_ID env var)")

	return root, cfg
}

// Execute is the main entry point for the CLI.
func Execute() {
	root, cfg := NewRootCmd()
	root.AddCommand(newDatabaseCmd(cfg))
	root.AddCommand(newClusterCmd(cfg))
	root.AddCommand(newSecurityCmd(cfg))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
