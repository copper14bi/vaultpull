package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/config"
)

var (
	cfgFile string
	appCfg  *config.Config
)

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull pulls secrets from a HashiCorp Vault KV store and writes
them to a local .env file, supporting rotation and selective key mapping.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		appCfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		if err = appCfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands and runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "",
		"config file (default: .vaultpull.yaml in current dir or $HOME)",
	)
}
