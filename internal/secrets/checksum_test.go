package secrets

import (
	"testing"
)

func TestChecksum_IndividualKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	res := Checksum(secrets)

	if len(res.Individual) != 2 {
		t.Fatalf("expected 2 individual checksums, got %d", len(res.Individual))
	}
	for k := range secrets {
		if _, ok := res.Individual[k]; !ok {
			t.Errorf("missing checksum for key %q", k)
		}
	}
}

func TestChecksum_CombinedIsDeterministic(t *testing.T) {
	secrets := map[string]string{
		"Z_KEY": "zzz",
		"A_KEY": "aaa",
	}
	res1 := Checksum(secrets)
	res2 := Checksum(secrets)

	if res1.Combined != res2.Combined {
		t.Errorf("combined checksum is not deterministic: %q vs %q", res1.Combined, res2.Combined)
	}
}

func TestChecksum_EmptyMap(t *testing.T) {
	res := Checksum(map[string]string{})
	if res.Combined == "" {
		t.Error("expected non-empty combined checksum for empty map")
	}
	if len(res.Individual) != 0 {
		t.Errorf("expected 0 individual checksums, got %d", len(res.Individual))
	}
}

func TestVerify_MatchingValue(t *testing.T) {
	res := Checksum(map[string]string{"K": "secret"})
	digest := res.Individual["K"]

	if !Verify("secret", digest) {
		t.Error("Verify returned false for matching value")
	}
}

func TestVerify_ModifiedValue(t *testing.T) {
	res := Checksum(map[string]string{"K": "secret"})
	digest := res.Individual["K"]

	if Verify("CHANGED", digest) {
		t.Error("Verify returned true for modified value")
	}
}

func TestVerifyMap_NoMismatches(t *testing.T) {
	secrets := map[string]string{"A": "foo", "B": "bar"}
	res := Checksum(secrets)

	mismatched := VerifyMap(secrets, res.Individual)
	if len(mismatched) != 0 {
		t.Errorf("expected no mismatches, got %v", mismatched)
	}
}

func TestVerifyMap_DetectsTampering(t *testing.T) {
	original := map[string]string{"A": "foo", "B": "bar"}
	res := Checksum(original)

	tampered := map[string]string{"A": "foo", "B": "TAMPERED"}
	mismatched := VerifyMap(tampered, res.Individual)

	if len(mismatched) != 1 || mismatched[0] != "B" {
		t.Errorf("expected [B] as mismatched, got %v", mismatched)
	}
}

func TestVerifyMap_MissingKey(t *testing.T) {
	original := map[string]string{"A": "foo", "B": "bar"}
	res := Checksum(original)

	partial := map[string]string{"A": "foo"} // B is missing
	mismatched := VerifyMap(partial, res.Individual)

	if len(mismatched) != 1 || mismatched[0] != "B" {
		t.Errorf("expected [B] as mismatched due to missing key, got %v", mismatched)
	}
}
