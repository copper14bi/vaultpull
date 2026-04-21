package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/template"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	tmplSrc  string
	tmplDst  string
	tmplPerm uint32
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a template file using secrets from Vault",
	Long: `Fetch secrets from Vault and render a Go text/template file,
writing the result to the specified destination path.`,
	RunE: runTemplate,
}

func init() {
	templateCmd.Flags().StringVarP(&tmplSrc, "src", "s", "", "source template file (required)")
	templateCmd.Flags().StringVarP(&tmplDst, "dst", "d", "", "destination output file (required)")
	templateCmd.Flags().Uint32Var(&tmplPerm, "perm", 0o600, "output file permissions (octal)")
	_ = templateCmd.MarkFlagRequired("src")
	_ = templateCmd.MarkFlagRequired("dst")
	RootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets := make(map[string]string)
	for _, path := range cfg.SecretPaths {
		data, err := client.GetSecrets(cmd.Context(), path)
		if err != nil {
			return fmt.Errorf("fetch secrets from %q: %w", path, err)
		}
		for k, v := range data {
			secrets[k] = v
		}
	}

	r := template.New()
	if err := r.RenderFile(tmplSrc, tmplDst, secrets, os.FileMode(tmplPerm)); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "template rendered → %s\n", tmplDst)
	return nil
}
