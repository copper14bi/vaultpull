package vault

import "fmt"

// MockClient implements a fake Vault client for use in tests.
type MockClient struct {
	// Secrets maps path -> key/value pairs to return.
	Secrets map[string]map[string]string
	// Errors maps path -> error to simulate read failures.
	Errors map[string]error
	// Calls records which paths were requested.
	Calls []string
}

// NewMockClient returns an initialised MockClient.
func NewMockClient() *MockClient {
	return &MockClient{
		Secrets: make(map[string]map[string]string),
		Errors:  make(map[string]error),
	}
}

// GetSecret returns the pre-configured secret for path, or an error if set.
func (m *MockClient) GetSecret(path string) (map[string]string, error) {
	m.Calls = append(m.Calls, path)

	if err, ok := m.Errors[path]; ok {
		return nil, err
	}

	if data, ok := m.Secrets[path]; ok {
		copy := make(map[string]string, len(data))
		for k, v := range data {
			copy[k] = v
		}
		return copy, nil
	}

	return nil, fmt.Errorf("no secret found at path %q", path)
}

// SecretFetcher is the interface satisfied by both Client and MockClient.
type SecretFetcher interface {
	GetSecret(path string) (map[string]string, error)
}
