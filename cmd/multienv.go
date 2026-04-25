package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/multienv"
	"github.com/yourusername/vaultpull/internal/sync"
)

var multienvCmd = &cobra.Command{
	Use:   "multienv",
	Short: "Sync secrets into multiple .env files by prefix",
	Long: `Reads secrets from Vault and distributes them into separate .env files
based on key prefixes defined in the configuration's multi_targets section.`,
	RunE: runMultienv,
}

func init() {
	multienvCmd.Flags().StringP("config", "c", ".vaultpull.yaml", "config file path")
	multienvCmd.Flags().Bool("dry-run", false, "print targets and key counts without writing files")
	rootCmd.AddCommand(multienvCmd)
}

func runMultienv(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.MultiTargets) == 0 {
		return fmt.Errorf("no multi_targets defined in config; add at least one target")
	}

	syncer := sync.New(cfg)
	secrets, err := syncer.FetchAll(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	targets := make([]multienv.Target, 0, len(cfg.MultiTargets))
	for _, mt := range cfg.MultiTargets {
		targets = append(targets, multienv.Target{
			OutputFile: mt.OutputFile,
			Prefixes:   mt.Prefixes,
			BackupDir:  cfg.BackupDir,
		})
	}

	if dryRun {
		fmt.Fprintln(os.Stdout, "Dry-run mode — no files will be written:")
		for _, t := range targets {
			fmt.Fprintf(os.Stdout, "  %s  prefixes=%v\n", t.OutputFile, t.Prefixes)
		}
		return nil
	}

	w := multienv.New(targets)
	counts, err := w.WriteAll(secrets)
	if err != nil {
		return fmt.Errorf("write env files: %w", err)
	}

	for file, n := range counts {
		fmt.Fprintf(os.Stdout, "wrote %d key(s) → %s\n", n, file)
	}
	return nil
}
