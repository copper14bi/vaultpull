package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/snapshot"
	"github.com/your-org/vaultpull/internal/vault"
)

var snapshotFile string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture or diff a snapshot of Vault secrets",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current Vault secrets to a snapshot file",
	RunE:  runSnapshotSave,
}

var snapshotDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare current Vault secrets against a saved snapshot",
	RunE:  runSnapshotDiff,
}

func init() {
	snapshotCmd.PersistentFlags().StringVarP(&snapshotFile, "file", "f", ".vaultpull.snapshot.json", "snapshot file path")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotDiffCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}
	allSecrets := make(map[string]string)
	for _, p := range cfg.SecretPaths {
		secrets, err := client.GetSecrets(cmd.Context(), p)
		if err != nil {
			return fmt.Errorf("vault read %s: %w", p, err)
		}
		for k, v := range secrets {
			allSecrets[k] = v
		}
	}
	snap := snapshot.New(cfg.SecretPaths[0], allSecrets)
	if err := snap.Save(snapshotFile); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "snapshot saved to %s (%d keys)\n", snapshotFile, len(allSecrets))
	return nil
}

func runSnapshotDiff(cmd *cobra.Command, _ []string) error {
	snap, err := snapshot.Load(snapshotFile)
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}
	current := make(map[string]string)
	for _, p := range cfg.SecretPaths {
		secrets, err := client.GetSecrets(cmd.Context(), p)
		if err != nil {
			return fmt.Errorf("vault read %s: %w", p, err)
		}
		for k, v := range secrets {
			current[k] = v
		}
	}
	d := snap.Compare(current)
	if !d.HasDrift() {
		fmt.Fprintln(os.Stdout, "no drift detected")
		return nil
	}
	printDriftSection("added", d.Added)
	printDriftSection("removed", d.Removed)
	printDriftSection("changed", d.Changed)
	return nil
}

func printDriftSection(label string, keys []string) {
	if len(keys) == 0 {
		return
	}
	sort.Strings(keys)
	fmt.Fprintf(os.Stdout, "%s (%d):\n", label, len(keys))
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "  %s\n", k)
	}
}
