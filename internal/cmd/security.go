package cmd

import (
	"fmt"

	"github.com/exasol-labs/saas-cli/internal/api"
	"github.com/exasol-labs/saas-cli/internal/output"
	"github.com/spf13/cobra"
)

// newSecurityCmd returns the "security" subcommand and its children.
func newSecurityCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Manage the account IP allowlist",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newSecurityListCmd(cfg),
		newSecurityStatusCmd(cfg),
		newSecurityCreateCmd(cfg),
		newSecurityUpdateCmd(cfg),
		newSecurityDeleteCmd(cfg),
	)

	return cmd
}

func newSecurityListCmd(cfg *Config) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all allowed IP entries for the account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			ips, err := client.ListAllowedIPs(cfg.AccountID)
			if err != nil {
				return err
			}
			if len(ips) == 0 {
				fmt.Fprint(cmd.OutOrStdout(), "No allowed IPs found.\n")
				return nil
			}
			return output.PrintAllowedIPs(cmd.OutOrStdout(), ips, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newSecurityStatusCmd(cfg *Config) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Get details of an allowed IP entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			ip, err := client.GetAllowedIP(cfg.AccountID, args[0])
			if err != nil {
				return err
			}
			return output.PrintAllowedIP(cmd.OutOrStdout(), *ip, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newSecurityCreateCmd(cfg *Config) *cobra.Command {
	var (
		name       string
		cidrIp     string
		outputFlag string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Add an IP range to the account allowlist",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			ip, err := client.AddAllowedIP(cfg.AccountID, api.CreateAllowedIPRequest{
				Name:   name,
				CIDRIp: cidrIp,
			})
			if err != nil {
				return err
			}
			return output.PrintAllowedIP(cmd.OutOrStdout(), *ip, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Descriptive label for this IP range, e.g. office or vpn (required)")
	cmd.Flags().StringVar(&cidrIp, "cidr-ip", "", "IP range in CIDR notation, e.g. 203.0.113.0/24 (required)")
	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("cidr-ip")

	return cmd
}

func newSecurityUpdateCmd(cfg *Config) *cobra.Command {
	var (
		name   string
		cidrIp string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Replace an allowed IP entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.UpdateAllowedIP(cfg.AccountID, args[0], api.UpdateAllowedIPRequest{
				Name:   name,
				CIDRIp: cidrIp,
			}); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Allowed IP updated.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Descriptive label for this IP range, e.g. office or vpn (required)")
	cmd.Flags().StringVar(&cidrIp, "cidr-ip", "", "IP range in CIDR notation, e.g. 203.0.113.0/24 (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("cidr-ip")

	return cmd
}

func newSecurityDeleteCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Remove an IP range from the account allowlist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.DeleteAllowedIP(cfg.AccountID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Allowed IP deleted.\n")
			return nil
		},
	}
}
