// Package lint implements naming convention and value quality rules for
// secrets pulled from HashiCorp Vault.
//
// # Rules
//
// The following rules are enforced:
//
//   - key-format (error): keys must match ^[A-Z][A-Z0-9_]*$ to be valid
//     .env identifiers.
//   - empty-value (warning): a secret with an empty value is likely
//     misconfigured or not yet populated in Vault.
//   - placeholder-value (error): common placeholder strings such as
//     "changeme" or "todo" indicate a secret has not been rotated.
//   - key-length (info): keys longer than 64 characters may cause
//     compatibility issues with some tooling.
//
// # Usage
//
//	findings := lint.Lint(secrets)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
package lint
