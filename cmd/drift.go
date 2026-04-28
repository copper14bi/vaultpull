package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/secrets"
)

var driftCmd = &cobra.Command{
	Use:   "drift",
	Short: "Detect secrets that have drifted from policy (age, entropy, rotation)",
	RunE:  runDrift,
}

func init() {
	driftCmd.Flags().Int("max-age", 90, "Maximum allowed age in days before a secret is flagged")
	driftCmd.Flags().Float64("min-entropy", 3.5, "Minimum Shannon entropy required for sensitive keys")
	driftCmd.Flags().Bool("require-rotation", true, "Flag secrets that have never been rotated")
	driftCmd.Flags().String("created-at", "", "RFC3339 creation timestamp (default: now minus max-age)")
	driftCmd.Flags().String("rotated-at", "", "RFC3339 last-rotation timestamp (empty = never rotated)")
	rootCmd.AddCommand(driftCmd)
}

func runDrift(cmd *cobra.Command, _ []string) error {
	maxAge, _ := cmd.Flags().GetInt("max-age")
	minEntropy, _ := cmd.Flags().GetFloat64("min-entropy")
	requireRotation, _ := cmd.Flags().GetBool("require-rotation")
	createdAtStr, _ := cmd.Flags().GetString("created-at")
	rotatedAtStr, _ := cmd.Flags().GetString("rotated-at")

	opts := secrets.DriftOptions{
		MaxAgeDays:       maxAge,
		EntropyThreshold: minEntropy,
		RequireRotation:  requireRotation,
	}

	createdAt := time.Now().Add(-time.Duration(maxAge) * 24 * time.Hour)
	if createdAtStr != "" {
		parsed, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return fmt.Errorf("invalid --created-at: %w", err)
		}
		createdAt = parsed
	}

	var rotatedAt time.Time
	if rotatedAtStr != "" {
		parsed, err := time.Parse(time.RFC3339, rotatedAtStr)
		if err != nil {
			return fmt.Errorf("invalid --rotated-at: %w", err)
		}
		rotatedAt = parsed
	}

	// Read secrets from stdin env-style (KEY=VALUE) for pipeline use.
	envSecrets := parseSimpleEnv(os.Stdin)
	if len(envSecrets) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no secrets provided via stdin")
		return nil
	}

	results := secrets.CheckDriftMap(envSecrets, createdAt, rotatedAt, opts)
	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ no drift detected")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "drift detected in %d secret(s):\n", len(results))
	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s — %s\n", r.Severity, r.Key, r.Reason)
	}
	return nil
}
