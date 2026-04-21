package health_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/vaultpull/internal/health"
)

func newTestClient(t *testing.T, srv *httptest.Server) *vaultapi.Client {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	require.NoError(t, err)
	client.SetToken("test-token")
	return client
}

func TestCheck_UnreachableVault(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "connection refused", http.StatusServiceUnavailable)
	}))
	srv.Close() // immediately close to simulate unreachable

	client := newTestClient(t, srv)
	checker := health.New(client)
	status := checker.Check(context.Background())

	assert.False(t, status.Reachable)
	assert.False(t, status.IsHealthy())
	assert.NotEmpty(t, status.Error)
}

func TestCheck_SealedVault(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"sealed":true,"initialized":true}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	checker := health.New(client)
	status := checker.Check(context.Background())

	assert.True(t, status.Reachable)
	assert.True(t, status.Sealed)
	assert.False(t, status.IsHealthy())
}

func TestIsHealthy_AllGood(t *testing.T) {
	s := health.Status{
		Reachable:     true,
		Authenticated: true,
		Sealed:        false,
	}
	assert.True(t, s.IsHealthy())
}

func TestIsHealthy_NotAuthenticated(t *testing.T) {
	s := health.Status{
		Reachable:     true,
		Authenticated: false,
		Sealed:        false,
	}
	assert.False(t, s.IsHealthy())
}
