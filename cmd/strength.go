package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/secrets"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var strengthCmd = &cobra.Command{
	Use:   "strength",
	Short: "Assess the strength of secrets fetched from Vault",
	Long:  `Fetches secrets from Vault and reports a strength score for each value, highlighting weak or placeholder secrets.`,
	RunE:  runStrength,
}

var strengthMinLevel string

func init() {
	strengthCmd.Flags().StringVar(&strengthMinLevel, "min-level", "fair", "Minimum acceptable strength level (weak|fair|strong|excellent)")
	rootCmd.AddCommand(strengthCmd)
}

func runStrength(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	syncer := sync.NewWithClient(cfg, client)
	allSecrets, err := syncer.FetchAll(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	minLevel := parseMinLevel(strengthMinLevel)

	failed := false
	keys := make([]string, 0, len(allSecrets))
	for k := range allSecrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(os.Stdout, "%-40s %-10s %s\n", "KEY", "LEVEL", "SCORE")
	fmt.Fprintf(os.Stdout, "%-40s %-10s %s\n", "---", "-----", "-----")

	for _, key := range keys {
		result := secrets.CheckStrength(allSecrets[key])
		marker := ""
		if result.Level < minLevel {
			marker = " !"
			failed = true
		}
		fmt.Fprintf(os.Stdout, "%-40s %-10s %3d%s\n", key, result.Level, result.Score, marker)
		for _, reason := range result.Reasons {
			fmt.Fprintf(os.Stdout, "  ↳ %s\n", reason)
		}
	}

	if failed {
		fmt.Fprintln(os.Stderr, "\n[!] One or more secrets did not meet the minimum strength level:", strengthMinLevel)
		os.Exit(1)
	}
	return nil
}

func parseMinLevel(s string) secrets.StrengthLevel {
	switch s {
	case "weak":
		return secrets.StrengthWeak
	case "fair":
		return secrets.StrengthFair
	case "strong":
		return secrets.StrengthStrong
	case "excellent":
		return secrets.StrengthExcellent
	default:
		return secrets.StrengthFair
	}
}
