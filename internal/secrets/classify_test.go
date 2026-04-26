package secrets

import (
	"testing"
)

func TestClassify_SecretTier(t *testing.T) {
	cases := []string{"DB_PASSWORD", "API_TOKEN", "PRIVATE_KEY", "APP_SECRET"}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			got := Classify(key)
			if got != ClassSecret {
				t.Errorf("Classify(%q) = %s, want secret", key, got)
			}
		})
	}
}

func TestClassify_ConfidentialTier(t *testing.T) {
	cases := []string{"TLS_CERT", "DATABASE_DSN", "AUTH_HEADER"}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			got := Classify(key)
			if got != ClassConfidential {
				t.Errorf("Classify(%q) = %s, want confidential", key, got)
			}
		})
	}
}

func TestClassify_InternalTier(t *testing.T) {
	cases := []string{"SERVICE_URL", "REDIS_HOST", "API_ENDPOINT"}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			got := Classify(key)
			if got != ClassInternal {
				t.Errorf("Classify(%q) = %s, want internal", key, got)
			}
		})
	}
}

func TestClassify_PublicTier(t *testing.T) {
	cases := []string{"HTTP_PORT", "DEBUG", "APP_ENV"}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			got := Classify(key)
			if got != ClassPublic {
				t.Errorf("Classify(%q) = %s, want public", key, got)
			}
		})
	}
}

func TestClassify_DefaultsToInternal(t *testing.T) {
	got := Classify("SOME_UNKNOWN_VAR")
	if got != ClassInternal {
		t.Errorf("Classify(unknown) = %s, want internal", got)
	}
}

func TestClassification_String(t *testing.T) {
	tests := []struct {
		c    Classification
		want string
	}{
		{ClassPublic, "public"},
		{ClassInternal, "internal"},
		{ClassConfidential, "confidential"},
		{ClassSecret, "secret"},
	}
	for _, tt := range tests {
		if got := tt.c.String(); got != tt.want {
			t.Errorf"Classification(%d).String() = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestClassifyMap_ReturnsAllKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"HTTP_PORT":   "8080",
		"SERVICE_URL": "http://example.com",
	}
	result := ClassifyMap(secrets)
	if len(result) != len(secrets) {
		t.Fatalf("ClassifyMap returned %d entries, want %d", len(result), len(secrets))
	}
	if result["DB_PASSWORD"] != ClassSecret {
		t.Errorf("DB_PASSWORD should be ClassSecret, got %s", result["DB_PASSWORD"])
	}
	if result["HTTP_PORT"] != ClassPublic {
		t.Errorf("HTTP_PORT should be ClassPublic, got %s", result["HTTP_PORT"])
	}
}
