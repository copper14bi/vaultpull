package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/secrets"
)

var (
	lifecycleMinStrength string
	lifecycleMaxAgeDays  int
	lifecycleWarnAgeDays int
)

func init() {
	lifecycleCmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Evaluate age, TTL, and strength lifecycle of secrets in a .env file",
		RunE:  runLifecycle,
	}
	lifecycleCmd.Flags().StringVar(&lifecycleMinStrength, "min-strength", "weak", "Minimum acceptable strength level (weak|fair|strong|excellent)")
	lifecycleCmd.Flags().IntVar(&lifecycleMaxAgeDays, "max-age-days", 90, "Maximum secret age in days before considered expired")
	lifecycleCmd.Flags().IntVar(&lifecycleWarnAgeDays, "warn-age-days", 60, "Age in days to start warning")
	rootCmd.AddCommand(lifecycleCmd)
}

func runLifecycle(cmd *cobra.Command, args []string) error {
	envFile := ".env"
	if len(args) > 0 {
		envFile = args[0]
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	secretMap := parseSimpleEnv(string(data))
	if len(secretMap) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	minLevel := parseMinLevel(lifecycleMinStrength)
	opts := secrets.DefaultLifecycleOptions()
	opts.MinStrength = minLevel
	opts.Age.MaxAge = time.Duration(lifecycleMaxAgeDays) * 24 * time.Hour
	opts.Age.WarnAge = time.Duration(lifecycleWarnAgeDays) * 24 * time.Hour

	results := secrets.CheckLifecycleMap(secretMap, time.Now(), opts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tSTATUS\tSTRENGTH\tMESSAGES")
	hasIssues := false
	for _, r := range results {
		msg := "-"
		if len(r.Messages) > 0 {
			msg = r.Messages[0]
			hasIssues = true
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Key, r.Status, r.Strength, msg)
	}
	w.Flush()

	if hasIssues {
		return fmt.Errorf("one or more secrets failed lifecycle checks")
	}
	return nil
}

func parseSimpleEnv(content string) map[string]string {
	result := map[string]string{}
	for _, line := range splitLines(content) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if k, v, ok := cutEnvLine(line); ok {
			result[k] = v
		}
	}
	return result
}
