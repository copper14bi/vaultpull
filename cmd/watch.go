package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

f13/cobra"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/watch"
)

var watchInterval string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch config file and re-sync secrets on",
	Long: `Polls the vaultpull config file for modifications.
When a change is detected, secrets are automatically re-synced from Vault.`,
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().StringVarP(&watchInterval, "interval", "i", "5s", "poll interval (e.g. 5s, 1m)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, _ []string) error {
	interval, err := time.ParseDuration(watchInterval)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", watchInterval, err)
	}

	cfgPath := cmd.Root().PersistentFlags().Lookup("config").Value.String()

	handler := func(path string) error {
		cfg, err := config.Load(path)
		if err != nil {
			return fmt.Errorf("reload config: %w", err)
		}
		s := sync.New(cfg)
		results, err := s.Run(cmd.Context())
		if err != nil {
			return fmt.Errorf("sync: %w", err)
		}
		fmt.Fprintf(os.Stdout, "synced: %s\n", results.Summary())
		return nil
	}

	w := watch.New(cfgPath, interval, handler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stdout, "watching %s (interval: %s) — press Ctrl+C to stop\n", cfgPath, interval)
	if err := w.Start(ctx); err != nil && err != context.Canceled {
		return err
	}
	return nil
}
