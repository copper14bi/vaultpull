package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/audit"
	"vaultpull/internal/cache"
	"vaultpull/internal/config"
	"vaultpull/internal/diff"
	"vaultpull/internal/filter"
	"vaultpull/internal/output"
	"vaultpull/internal/prompt"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var (
	flagDryRun    bool
	flagForce     bool
	flagFormat    string
	flagCacheTTL  time.Duration
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault into local .env files",
	Long: `Pull secrets from HashiCorp Vault and write them to local .env files.

Secrets are fetched from the configured paths and merged into the target
env file. Existing keys not present in Vault are preserved by default.
Use --force to overwrite all keys.

Example:
  vaultpull sync
  vaultpull sync --dry-run
  vaultpull sync --format json`,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Preview changes without writing to disk")
	syncCmd.Flags().BoolVar(&flagForce, "force", false, "Overwrite existing keys without prompting")
	syncCmd.Flags().StringVar(&flagFormat, "format", "text", "Output format: text or json")
	syncCmd.Flags().DurationVar(&flagCacheTTL, "cache-ttl", 5*time.Minute, "How long to cache Vault responses (0 to disable)")

	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Set up output formatter.
	fmt := output.New(flagFormat)

	// Set up Vault client with optional cache.
	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	var vaultReader sync.SecretReader = vaultClient
	if flagCacheTTL > 0 {
		c := cache.New(flagCacheTTL)
		vaultReader = cache.WrapReader(c, vaultClient)
	}

	// Build filter from config rules.
	f := filter.New(cfg.Filter.Include, cfg.Filter.Exclude)

	// Set up audit logger.
	auditLog, err := audit.NewLogger(cfg.AuditLog)
	if err != nil {
		return fmt.Errorf("creating audit logger: %w", err)
	}
	defer auditLog.Close()

	// Set up confirmer for interactive prompts.
	confirmer := prompt.New(os.Stdin, os.Stderr)

	// Build and run the syncer.
	s := sync.NewWithClient(cfg, vaultReader, f, auditLog)

	result, err := s.Run(cmd.Context())
	if err != nil {
		fmt.Print(output.Result{Error: err})
		return err
	}

	// Show diff and optionally prompt before writing.
	changes := diff.Compare(result.Previous, result.Fetched)
	if flagDryRun {
		diff.Print(changes, os.Stdout)
		fmt.Print(output.Result{Summary: result.Summary(), DryRun: true})
		return nil
	}

	if diff.HasChanges(changes) && !flagForce {
		diff.Print(changes, os.Stdout)
		ok, promptErr := confirmer.Ask("Apply these changes?", true)
		if promptErr != nil {
			return fmt.Errorf("prompt: %w", promptErr)
		}
		if !ok {
			fmt.Print(output.Result{Summary: "sync cancelled by user"})
			return nil
		}
	}

	if err := result.Write(); err != nil {
		fmt.Print(output.Result{Error: err})
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Print(output.Result{Summary: result.Summary()})
	return nil
}
