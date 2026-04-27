package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/secrets"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Check TTL status of secrets loaded from a .env file",
	Long: `Reads expiry timestamps from environment variables of the form
<KEY>_EXPIRES_AT (RFC3339) and reports their TTL status.`,
	RunE: runTTL,
}

var ttlWarnHours int

func init() {
	ttlCmd.Flags().StringP("env-file", "f", ".env", "Path to .env file containing _EXPIRES_AT keys")
	ttlCmd.Flags().IntVar(&ttlWarnHours, "warn-hours", 24, "Hours remaining before status becomes 'warning'")
	rootCmd.AddCommand(ttlCmd)
}

func runTTL(cmd *cobra.Command, _ []string) error {
	envFile, _ := cmd.Flags().GetString("env-file")

	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	expiries := parseExpiryKeys(string(data))
	if len(expiries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No expiry keys found (expected format: KEY_EXPIRES_AT=<RFC3339>)")
		return nil
	}

	opts := secrets.TTLOptions{
		WarnThreshold: time.Duration(ttlWarnHours) * time.Hour,
	}

	results := secrets.CheckTTLMap(expiries, opts)
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })

	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "%-30s %-10s %s\n", "KEY", "STATUS", "MESSAGE")
	fmt.Fprintf(w, "%-30s %-10s %s\n", "---", "------", "-------")

	hasExpired := false
	for _, r := range results {
		fmt.Fprintf(w, "%-30s %-10s %s\n", r.Key, r.Status, r.Message)
		if r.Status == secrets.TTLStatusExpired {
			hasExpired = true
		}
	}

	if hasExpired {
		return fmt.Errorf("one or more secrets have expired TTLs")
	}
	return nil
}

// parseExpiryKeys scans env file lines for KEY_EXPIRES_AT=<RFC3339> entries.
func parseExpiryKeys(content string) map[string]time.Time {
	result := map[string]time.Time{}
	for _, line := range splitLines(content) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		key, val, ok := cutEnvLine(line)
		if !ok {
			continue
		}
		if len(key) > 11 && key[len(key)-11:] == "_EXPIRES_AT" {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				baseKey := key[:len(key)-11]
				result[baseKey] = t
			}
		}
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func cutEnvLine(line string) (key, val string, ok bool) {
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}
