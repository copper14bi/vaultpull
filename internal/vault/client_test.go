package vault

import (
	"testing"
)

func TestToStringMap_BasicTypes(t *testing.T) {
	input := map[string]interface{}{
		"db_password": "s3cr3t",
		"port":        "5432",
		"retries":     3,
	}

	result, err := toStringMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", result["DB_PASSWORD"])
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", result["PORT"])
	}
	if result["RETRIES"] != "3" {
		t.Errorf("expected RETRIES=3, got %q", result["RETRIES"])
	}
}

func TestToStringMap_KeysUppercased(t *testing.T) {
	input := map[string]interface{}{
		"my_secret_key": "value",
	}

	result, err := toStringMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result["MY_SECRET_KEY"]; !ok {
		t.Error("expected key MY_SECRET_KEY to exist")
	}
	if _, ok := result["my_secret_key"]; ok {
		t.Error("expected lowercase key to be absent")
	}
}

func TestExtractData_KVv2Nested(t *testing.T) {
	raw := map[string]interface{}{
		"data": map[string]interface{}{
			"api_key": "abc123",
		},
		"metadata": map[string]interface{}{
			"version": 1,
		},
	}

	result, err := extractData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", result["API_KEY"])
	}
	if _, ok := result["METADATA"]; ok {
		t.Error("metadata should not appear in extracted data")
	}
}

func TestExtractData_KVv1Flat(t *testing.T) {
	raw := map[string]interface{}{
		"username": "admin",
		"password": "hunter2",
	}

	result, err := extractData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["USERNAME"] != "admin" {
		t.Errorf("expected USERNAME=admin, got %q", result["USERNAME"])
	}
	if result["PASSWORD"] != "hunter2" {
		t.Errorf("expected PASSWORD=hunter2, got %q", result["PASSWORD"])
	}
}
