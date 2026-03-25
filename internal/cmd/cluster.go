package cmd

import (
	"fmt"

	"github.com/exasol-labs/saas-cli/internal/api"
	"github.com/exasol-labs/saas-cli/internal/output"
	"github.com/spf13/cobra"
)

// newClusterCmd returns the "cluster" subcommand and its children.
func newClusterCmd(cfg *Config) *cobra.Command {
	var databaseID string

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage Exasol SaaS clusters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&databaseID, "database-id", "", "ID of the database containing the clusters (required)")
	_ = cmd.MarkPersistentFlagRequired("database-id")

	cmd.AddCommand(
		newClusterListCmd(cfg, &databaseID),
		newClusterStatusCmd(cfg, &databaseID),
		newClusterCreateCmd(cfg, &databaseID),
		newClusterUpdateCmd(cfg, &databaseID),
		newClusterDeleteCmd(cfg, &databaseID),
		newClusterScaleCmd(cfg, &databaseID),
		newClusterStartCmd(cfg, &databaseID),
		newClusterStopCmd(cfg, &databaseID),
		newClusterConnectCmd(cfg, &databaseID),
	)

	return cmd
}

func newClusterListCmd(cfg *Config, databaseID *string) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all clusters in a database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			clusters, err := client.ListClusters(cfg.AccountID, *databaseID)
			if err != nil {
				return err
			}
			if len(clusters) == 0 {
				fmt.Fprint(cmd.OutOrStdout(), "No clusters found.\n")
				return nil
			}
			return output.PrintClusters(cmd.OutOrStdout(), clusters, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newClusterStatusCmd(cfg *Config, databaseID *string) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Get the status of a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			cluster, err := client.GetCluster(cfg.AccountID, *databaseID, args[0])
			if err != nil {
				return err
			}
			return output.PrintCluster(cmd.OutOrStdout(), *cluster, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}

func newClusterCreateCmd(cfg *Config, databaseID *string) *cobra.Command {
	var (
		name              string
		size              string
		family            string
		autoStopEnabled   bool
		autoStopIdleTime  int
		offloadEnabled    bool
		offloadTimeoutMin int
		outputFlag        string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new cluster in a database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)

			req := api.CreateClusterRequest{
				Name:   name,
				Size:   size,
				Family: family,
			}

			if cmd.Flags().Changed("auto-stop-enabled") || cmd.Flags().Changed("auto-stop-idle-time") {
				req.AutoStop = &api.AutoStop{
					Enabled:  autoStopEnabled,
					IdleTime: autoStopIdleTime,
				}
			}

			if cmd.Flags().Changed("offload-enabled") || cmd.Flags().Changed("offload-timeout-min") {
				req.Settings = &api.ClusterSettingsUpdate{
					OffloadEnabled:    offloadEnabled,
					OffloadTimeoutMin: offloadTimeoutMin,
				}
			}

			cluster, err := client.CreateCluster(cfg.AccountID, *databaseID, req)
			if err != nil {
				return err
			}
			return output.PrintCluster(cmd.OutOrStdout(), *cluster, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Cluster name (required)")
	cmd.Flags().StringVar(&size, "size", "", "Cluster size, e.g. XS, S, M, L (required)")
	cmd.Flags().StringVar(&family, "family", "", "Cluster family")
	cmd.Flags().BoolVar(&autoStopEnabled, "auto-stop-enabled", false, "Enable auto-stop when the cluster is idle")
	cmd.Flags().IntVar(&autoStopIdleTime, "auto-stop-idle-time", 0, "Idle time in minutes before auto-stop triggers")
	cmd.Flags().BoolVar(&offloadEnabled, "offload-enabled", false, "Enable query offloading")
	cmd.Flags().IntVar(&offloadTimeoutMin, "offload-timeout-min", 0, "Query offload timeout in minutes")
	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("size")

	return cmd
}

func newClusterUpdateCmd(cfg *Config, databaseID *string) *cobra.Command {
	var (
		name              string
		autoStopEnabled   bool
		autoStopIdleTime  int
		offloadEnabled    bool
		offloadTimeoutMin int
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)

			req := api.UpdateClusterRequest{}
			if cmd.Flags().Changed("name") {
				req.Name = name
			}
			if cmd.Flags().Changed("auto-stop-enabled") || cmd.Flags().Changed("auto-stop-idle-time") {
				req.AutoStop = &api.AutoStop{
					Enabled:  autoStopEnabled,
					IdleTime: autoStopIdleTime,
				}
			}
			if cmd.Flags().Changed("offload-enabled") || cmd.Flags().Changed("offload-timeout-min") {
				req.Settings = &api.ClusterSettingsUpdate{
					OffloadEnabled:    offloadEnabled,
					OffloadTimeoutMin: offloadTimeoutMin,
				}
			}

			if err := client.UpdateCluster(cfg.AccountID, *databaseID, args[0], req); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Cluster updated.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New cluster name")
	cmd.Flags().BoolVar(&autoStopEnabled, "auto-stop-enabled", false, "Enable auto-stop when the cluster is idle")
	cmd.Flags().IntVar(&autoStopIdleTime, "auto-stop-idle-time", 0, "Idle time in minutes before auto-stop triggers")
	cmd.Flags().BoolVar(&offloadEnabled, "offload-enabled", false, "Enable query offloading")
	cmd.Flags().IntVar(&offloadTimeoutMin, "offload-timeout-min", 0, "Query offload timeout in minutes")

	return cmd
}

func newClusterDeleteCmd(cfg *Config, databaseID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.DeleteCluster(cfg.AccountID, *databaseID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Cluster deleted.\n")
			return nil
		},
	}
}

func newClusterScaleCmd(cfg *Config, databaseID *string) *cobra.Command {
	var (
		size   string
		family string
	)

	cmd := &cobra.Command{
		Use:   "scale <id>",
		Short: "Scale a cluster to a different size",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			req := api.ScaleClusterRequest{
				Size:   size,
				Family: family,
			}
			if err := client.ScaleCluster(cfg.AccountID, *databaseID, args[0], req); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Scale request submitted.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&size, "size", "", "Target cluster size, e.g. XS, S, M, L (required)")
	cmd.Flags().StringVar(&family, "family", "", "Target cluster family")
	_ = cmd.MarkFlagRequired("size")

	return cmd
}

func newClusterStartCmd(cfg *Config, databaseID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "start <id>",
		Short: "Start a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.StartCluster(cfg.AccountID, *databaseID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Start request submitted.\n")
			return nil
		},
	}
}

func newClusterStopCmd(cfg *Config, databaseID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <id>",
		Short: "Stop a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			if err := client.StopCluster(cfg.AccountID, *databaseID, args[0]); err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), "Stop request submitted.\n")
			return nil
		},
	}
}

func newClusterConnectCmd(cfg *Config, databaseID *string) *cobra.Command {
	var outputFlag string

	cmd := &cobra.Command{
		Use:   "connect <id>",
		Short: "Get connection details for a cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(cfg.Token, cfg.BaseURL)
			conn, err := client.GetClusterConnection(cfg.AccountID, *databaseID, args[0])
			if err != nil {
				return err
			}
			return output.PrintClusterConnection(cmd.OutOrStdout(), *conn, output.Format(outputFlag))
		},
	}

	cmd.Flags().StringVar(&outputFlag, "output", "table", "Output format: table or json")

	return cmd
}
