package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/secrets"
)

var (
	scanFile      string
	scanEntropyOn bool
)

func init() {
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a .env file for weak or suspicious secrets",
		RunE:  runScan,
	}
	scanCmd.Flags().StringVarP(&scanFile, "file", "f", ".env", "Path to the .env file to scan")
	scanCmd.Flags().BoolVar(&scanEntropyOn, "entropy", true, "Check Shannon entropy of sensitive values")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, _ []string) error {
	w := NewWithWriter(cmd.OutOrStdout())

	data, err := os.ReadFile(scanFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", scanFile, err)
	}

	parsed, err := env.ParseFile(data)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	findings := secrets.Scan(parsed)
	if len(findings) == 0 && !scanEntropyOn {
		w.Success("No suspicious secrets found in " + scanFile)
		return nil
	}

	for _, f := range findings {
		fmt.Fprintf(cmd.OutOrStdout(), "[SUSPICIOUS] key=%s reason=%s\n", f.Key, f.Reason)
	}

	if scanEntropyOn {
		weakEntropy := secrets.CheckEntropyMap(parsed)
		for _, r := range weakEntropy {
			fmt.Fprintf(cmd.OutOrStdout(), "[WEAK ENTROPY] value=%q reason=%s entropy=%.2f\n",
				r.Value, r.Reason, r.Entropy)
		}
		if len(weakEntropy) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d weak secret(s) detected.\n", len(weakEntropy))
		}
	}

	total := len(findings)
	if total == 0 {
		w.Success("No suspicious patterns found in " + scanFile)
	} else {
		w.Error(fmt.Sprintf("%d suspicious finding(s) in %s", total, scanFile))
	}
	return nil
}
