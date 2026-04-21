// Package health provides Vault connectivity and authentication checks.
package health

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Status holds the result of a health check.
type Status struct {
	Reachable    bool
	Authenticated bool
	Sealed       bool
	Latency      time.Duration
	Error        string
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	client *vaultapi.Client
}

// New creates a Checker using the provided Vault API client.
func New(client *vaultapi.Client) *Checker {
	return &Checker{client: client}
}

// Check performs a connectivity and auth validation against Vault.
func (c *Checker) Check(ctx context.Context) Status {
	start := time.Now()

	sys := c.client.Sys()
	health, err := sys.HealthWithContext(ctx)
	latency := time.Since(start)

	if err != nil {
		return Status{
			Reachable: false,
			Latency:   latency,
			Error:     fmt.Sprintf("vault unreachable: %v", err),
		}
	}

	if health.Sealed {
		return Status{
			Reachable: true,
			Sealed:    true,
			Latency:   latency,
			Error:     "vault is sealed",
		}
	}

	// Validate token by calling /auth/token/lookup-self
	_, err = c.client.Auth().Token().LookupSelfWithContext(ctx)
	if err != nil {
		return Status{
			Reachable:     true,
			Authenticated: false,
			Latency:       latency,
			Error:         fmt.Sprintf("token invalid or expired: %v", err),
		}
	}

	return Status{
		Reachable:     true,
		Authenticated: true,
		Latency:       latency,
	}
}

// IsHealthy returns true only when Vault is reachable, unsealed, and authenticated.
func (s Status) IsHealthy() bool {
	return s.Reachable && s.Authenticated && !s.Sealed
}
