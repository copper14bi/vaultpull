package vault

import (
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	api.SetToken(token)

	return &Client{api: api}, nil
}

// GetSecret reads a KV secret from Vault at the given path.
// Supports both KV v1 and KV v2 paths.
func (c *Client) GetSecret(path string) (map[string]string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data, err := extractData(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("extracting data from %q: %w", path, err)
	}

	return data, nil
}

// extractData handles both KV v1 (flat map) and KV v2 (nested under "data") responses.
func extractData(raw map[string]interface{}) (map[string]string, error) {
	if nested, ok := raw["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			return toStringMap(nestedMap)
		}
	}
	return toStringMap(raw)
}

func toStringMap(in map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string, len(in))
	for k, v := range in {
		key := strings.ToUpper(k)
		switch val := v.(type) {
		case string:
			out[key] = val
		case fmt.Stringer:
			out[key] = val.String()
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
	return out, nil
}
