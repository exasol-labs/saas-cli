package cmd

import (
	"fmt"

	"github.com/exasol-labs/saas-cli/internal/api"
	"github.com/exasol-labs/saas-cli/internal/output"
	"github.com/spf13/cobra"
)

// newDatabaseCmd returns the "database" subcommand and its children.
func newDatabaseCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Manage Exasol SaaS databases",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newDatabaseListCmd(cfg),
		newDatabaseStatusCmd(cfg),
		newDatabaseCreateCmd(cfg),
		newDatabaseUpdateCmd(cfg),
		newDatabaseDeleteCmd(cfg),
		newDatabaseStartCmd(cfg),
		newDatabaseStopCmd(cfg),
	)

	return cmd
}

func newDatabaseListCmd(cfg *Config) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all databases",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			dbs, err := client.ListDatabases(cfg.AccountID)
			if err != nil {
				return err
			}
			if len(dbs) == 0 {
				fmt.Fprint(cmd.OutOrStdout(), "No databases found.\n")
				return nil
			}
			return output.PrintDatabases(cmd.OutOrStdout(), dbs, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newDatabaseStatusCmd(cfg *Config) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Get the status of a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			db, err := client.GetDatabase(cfg.AccountID, args[0])
			if err != nil {
				return err
			}
			return output.PrintDatabase(cmd.OutOrStdout(), *db, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newDatabaseCreateCmd(cfg *Config) *cobra.Command {
	var (
		name              string
		region            string
		clusterName       string
		clusterSize       string
		clusterFamily     string
		numNodes          int
		streamType        string
		autoStopEnabled   bool
		autoStopIdleTime  int
		offloadEnabled    bool
		offloadTimeoutMin int
		outputFlag        string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)

			cluster := api.InitialCluster{
				Name:   clusterName,
				Size:   clusterSize,
				Family: clusterFamily,
			}

			if cmd.Flags().Changed("cluster-auto-stop-enabled") || cmd.Flags().Changed("cluster-auto-stop-idle-time") {
				cluster.AutoStop = &api.AutoStop{
					Enabled:  autoStopEnabled,
					IdleTime: autoStopIdleTime,
				}
			}

			if cmd.Flags().Changed("cluster-offload-enabled") || cmd.Flags().Changed("cluster-offload-timeout-min") {
				cluster.Settings = &api.ClusterSettingsUpdate{
					OffloadEnabled:    offloadEnabled,
					OffloadTimeoutMin: offloadTimeoutMin,
				}
			}

			req := api.CreateDatabaseRequest{
				Name:           name,
				Provider:       api.DefaultProvider,
				Region:         region,
				InitialCluster: cluster,
				NumNodes:       numNodes,
				StreamType:     streamType,
			}

			db, err := client.CreateDatabase(cfg.AccountID, req)
			if err != nil {
				return err
			}
			return output.PrintDatabase(cmd.OutOrStdout(), *db, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Database name (required)")
	cmd.Flags().StringVar(&region, "region", "", "Cloud region, e.g. us-east-1 (required)")
	cmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name for the initial cluster (required)")
	cmd.Flags().StringVar(&clusterSize, "cluster-size", "", "Size of the initial cluster, e.g. XS, S, M, L (required)")
	cmd.Flags().StringVar(&clusterFamily, "cluster-family", "", "Cluster family for the initial cluster")
	cmd.Flags().IntVar(&numNodes, "num-nodes", 0, "Number of nodes in the database")
	cmd.Flags().StringVar(&streamType, "stream-type", "", "Stream type for the database")
	cmd.Flags().BoolVar(&autoStopEnabled, "cluster-auto-stop-enabled", false, "Enable auto-stop for the initial cluster when idle")
	cmd.Flags().IntVar(&autoStopIdleTime, "cluster-auto-stop-idle-time", 0, "Idle time in minutes before auto-stop triggers")
	cmd.Flags().BoolVar(&offloadEnabled, "cluster-offload-enabled", false, "Enable query offloading for the initial cluster")
	cmd.Flags().IntVar(&offloadTimeoutMin, "cluster-offload-timeout-min", 0, "Query offload timeout in minutes")
	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("region")
	_ = cmd.MarkFlagRequired("cluster-name")
	_ = cmd.MarkFlagRequired("cluster-size")

	return cmd
}

func newDatabaseUpdateCmd(cfg *Config) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			err := client.UpdateDatabase(cfg.AccountID, args[0], api.UpdateDatabaseRequest{Name: name})
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Database updated.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New database name (required)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newDatabaseDeleteCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.DeleteDatabase(cfg.AccountID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Database deleted.\n")
			return nil
		},
	}
}

func newDatabaseStartCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start <id>",
		Short: "Start a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.StartDatabase(cfg.AccountID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Start request submitted.\n")
			return nil
		},
	}
}

func newDatabaseStopCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <id>",
		Short: "Stop a database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.StopDatabase(cfg.AccountID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Stop request submitted.\n")
			return nil
		},
	}
}
