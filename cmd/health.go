package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/health"
	vaultclient "github.com/yourusername/vaultpull/internal/vault"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check Vault connectivity and token validity",
	Long:  `Verifies that Vault is reachable, unsealed, and that the configured token is valid.`,
	RunE:  runHealth,
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().Duration("timeout", 10*time.Second, "request timeout for health check")
}

func runHealth(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	rawClient, err := vaultclient.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	checker := health.New(rawClient.API())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	status := checker.Check(ctx)

	fmt.Fprintf(cmd.OutOrStdout(), "Reachable:     %v\n", status.Reachable)
	fmt.Fprintf(cmd.OutOrStdout(), "Sealed:        %v\n", status.Sealed)
	fmt.Fprintf(cmd.OutOrStdout(), "Authenticated: %v\n", status.Authenticated)
	fmt.Fprintf(cmd.OutOrStdout(), "Latency:       %v\n", status.Latency.Round(time.Millisecond))

	if !status.IsHealthy() {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n", status.Error)
		os.Exit(1)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Status:        OK")
	return nil
}
