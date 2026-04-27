package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/lint"
	"github.com/your-org/vaultpull/internal/vault"
)

var lintFlags struct {
	errorOnly bool
	format    string
}

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint secret keys and values for naming and quality issues",
		Long:  "Fetches secrets from Vault and runs lint rules against key names and values.",
		RunE:  runLint,
	}
	lintCmd.Flags().BoolVar(&lintFlags.errorOnly, "errors-only", false, "only report error-severity findings")
	lintCmd.Flags().StringVar(&lintFlags.format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	allSecrets := make(map[string]string)
	for _, path := range cfg.SecretPaths {
		secrets, err := client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("fetch %s: %w", path, err)
		}
		for k, v := range secrets {
			allSecrets[k] = v
		}
	}

	findings := lint.Lint(allSecrets)

	var filtered []lint.Finding
	for _, f := range findings {
		if lintFlags.errorOnly && f.Severity != lint.SeverityError {
			continue
		}
		filtered = append(filtered, f)
	}

	if len(filtered) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ no lint findings")
		return nil
	}

	for _, f := range filtered {
		fmt.Fprintln(cmd.OutOrStdout(), f.String())
	}

	// exit non-zero if any errors found
	for _, f := range filtered {
		if f.Severity == lint.SeverityError {
			os.Exit(1)
		}
	}
	return nil
}
