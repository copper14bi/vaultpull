package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/secrets"
)

var ageCmd = &cobra.Command{
	Use:   "age",
	Short: "Check secret age against rotation thresholds",
	Long:  `Reads a JSON map of key:RFC3339-timestamp pairs and reports which secrets are fresh, approaching expiry, or overdue for rotation.`,
	RunE:  runAge,
}

func init() {
	ageCmd.Flags().String("timestamps", "", "path to JSON file containing key:timestamp map (RFC3339)")
	ageCmd.Flags().Int("warn-days", 60, "days after which a secret is considered stale")
	ageCmd.Flags().Int("expire-days", 90, "days after which a secret is considered expired")
	ageCmd.Flags().String("output", "text", "output format: text or json")
	_ = ageCmd.MarkFlagRequired("timestamps")
	rootCmd.AddCommand(ageCmd)
}

func runAge(cmd *cobra.Command, _ []string) error {
	tsFile, _ := cmd.Flags().GetString("timestamps")
	warnDays, _ := cmd.Flags().GetInt("warn-days")
	expireDays, _ := cmd.Flags().GetInt("expire-days")
	outFmt, _ := cmd.Flags().GetString("output")

	data, err := os.ReadFile(tsFile)
	if err != nil {
		return fmt.Errorf("reading timestamps file: %w", err)
	}

	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parsing timestamps JSON: %w", err)
	}

	timestamps := make(map[string]time.Time, len(raw))
	for k, v := range raw {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return fmt.Errorf("invalid timestamp for %q: %w", k, err)
		}
		timestamps[k] = t
	}

	opts := secrets.AgeOptions{
		WarnAfter:   time.Duration(warnDays) * 24 * time.Hour,
		ExpireAfter: time.Duration(expireDays) * 24 * time.Hour,
	}

	results := secrets.CheckAgeMap(timestamps, opts)

	if outFmt == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	for _, r := range results {
		var icon string
		switch r.Status {
		case secrets.AgeFresh:
			icon = "✓"
		case secrets.AgeWarning:
			icon = "!"
		case secrets.AgeExpired:
			icon = "✗"
		}
		fmt.Fprintf(os.Stdout, "%s  %-30s %s\n", icon, r.Key, r.Message)
	}
	return nil
}
