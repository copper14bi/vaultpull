package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vaultpull/internal/secrets"
)

func TestNewPolicy_Defaults(t *testing.T) {
	p := secrets.NewPolicy(nil)
	require.NotNil(t, p)
}

func TestPolicy_Enforce_AllowsPublic(t *testing.T) {
	p := secrets.NewPolicy(nil)

	secrets := map[string]string{
		"APP_NAME": "myapp",
		"LOG_LEVEL": "info",
	}

	err := p.Enforce(secrets)
	assert.NoError(t, err)
}

func TestPolicy_Enforce_RejectsWeakSecret(t *testing.T) {
	opts := &secrets.PolicyOptions{
		MinEntropyBits: 40.0,
		BlockWeakSecrets: true,
	}
	p := secrets.NewPolicy(opts)

	secrets := map[string]string{
		"DB_PASSWORD": "password",
	}

	err := p.Enforce(secrets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestPolicy_Enforce_AllowsStrongSecret(t *testing.T) {
	opts := &secrets.PolicyOptions{
		MinEntropyBits: 40.0,
		BlockWeakSecrets: true,
	}
	p := secrets.NewPolicy(opts)

	secrets := map[string]string{
		"DB_PASSWORD": "xK9#mP2$qL8nRv5@wZ3",
	}

	err := p.Enforce(secrets)
	assert.NoError(t, err)
}

func TestPolicy_Enforce_RejectsPlaceholder(t *testing.T) {
	opts := &secrets.PolicyOptions{
		BlockPlaceholders: true,
	}
	p := secrets.NewPolicy(opts)

	secrets := map[string]string{
		"API_KEY": "CHANGEME",
	}

	err := p.Enforce(secrets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API_KEY")
}

func TestPolicy_Enforce_MultipleViolations(t *testing.T) {
	opts := &secrets.PolicyOptions{
		MinEntropyBits: 40.0,
		BlockWeakSecrets: true,
		BlockPlaceholders: true,
	}
	p := secrets.NewPolicy(opts)

	secrets := map[string]string{
		"DB_PASSWORD": "weak",
		"API_SECRET": "TODO",
		"SAFE_KEY":   "xK9#mP2$qL8nRv5@wZ3",
	}

	err := p.Enforce(secrets)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
	assert.Contains(t, err.Error(), "API_SECRET")
	assert.NotContains(t, err.Error(), "SAFE_KEY")
}

func TestPolicy_Enforce_EmptySecretsAlwaysPasses(t *testing.T) {
	opts := &secrets.PolicyOptions{
		MinEntropyBits: 40.0,
		BlockWeakSecrets: true,
		BlockPlaceholders: true,
	}
	p := secrets.NewPolicy(opts)

	err := p.Enforce(map[string]string{})
	assert.NoError(t, err)
}

func TestPolicy_Enforce_SkipsNonSensitiveKeys(t *testing.T) {
	opts := &secrets.PolicyOptions{
		MinEntropyBits: 40.0,
		BlockWeakSecrets: true,
	}
	p := secrets.NewPolicy(opts)

	// Non-sensitive keys like APP_NAME should not be entropy-checked
	secrets := map[string]string{
		"APP_NAME": "hi",
		"PORT":     "8080",
	}

	err := p.Enforce(secrets)
	assert.NoError(t, err)
}
