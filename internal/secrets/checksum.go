package secrets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// ChecksumResult holds the checksum details for a secret map.
type ChecksumResult struct {
	// Individual maps each key to its SHA-256 hex digest.
	Individual map[string]string
	// Combined is a single digest over all key=value pairs (sorted by key).
	Combined string
}

// Checksum computes SHA-256 digests for each secret value and a combined
// digest for the entire map. Keys are always sorted before hashing so the
// combined digest is deterministic regardless of map iteration order.
func Checksum(secrets map[string]string) ChecksumResult {
	individual := make(map[string]string, len(secrets))
	for k, v := range secrets {
		individual[k] = hashValue(v)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, secrets[k])
	}
	combined := hashValue(sb.String())

	return ChecksumResult{
		Individual: individual,
		Combined:   combined,
	}
}

// Verify returns true when the checksum of value matches the expected hex
// digest produced by a previous call to Checksum.
func Verify(value, expectedHex string) bool {
	return hashValue(value) == expectedHex
}

// VerifyMap checks every key present in expected against the corresponding
// value in secrets. It returns a slice of keys whose values no longer match.
func VerifyMap(secrets map[string]string, expected map[string]string) []string {
	var mismatched []string
	for k, digest := range expected {
		v, ok := secrets[k]
		if !ok || !Verify(v, digest) {
			mismatched = append(mismatched, k)
		}
	}
	sort.Strings(mismatched)
	return mismatched
}

func hashValue(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
