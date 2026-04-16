package sync

import (
	"fmt"
	"path/filepath"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/env"
	"github.com/user/vaultpull/internal/vault"
)

// VaultReader abstracts secret fetching from Vault.
type VaultReader interface {
	GetSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing env files.
type Syncer struct {
	client  VaultReader
	cfg     *config.Config
	backup  bool
}

// New creates a Syncer with the provided config.
func New(cfg *config.Config, backup bool) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	return &Syncer{client: client, cfg: cfg, backup: backup}, nil
}

// NewWithClient creates a Syncer with an injected VaultReader (useful for testing).
func NewWithClient(cfg *config.Config, client VaultReader, backup bool) *Syncer {
	return &Syncer{client: client, cfg: cfg, backup: backup}
}

// Run iterates over all configured secret mappings and writes env files.
func (s *Syncer) Run() ([]Result, error) {
	var results []Result
	for _, mapping := range s.cfg.Mappings {
		secrets, err := s.client.GetSecrets(mapping.SecretPath)
		if err != nil {
			return results, fmt.Errorf("fetch %q: %w", mapping.SecretPath, err)
		}

		outPath := filepath.Clean(mapping.EnvFile)
		writer := env.NewWriter(outPath, s.backup)

		merged := env.Merge(outPath, secrets, mapping.Overwrite)
		if err := writer.Write(merged); err != nil {
			return results, fmt.Errorf("write %q: %w", outPath, err)
		}

		results = append(results, Result{
			SecretPath: mapping.SecretPath,
			EnvFile:    outPath,
			Count:      len(merged),
		})
	}
	return results, nil
}
