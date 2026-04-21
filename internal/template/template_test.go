package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/template"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
}

func TestRender_BasicSubstitution(t *testing.T) {
	r := template.New()
	out, err := r.Render(`host={{ index . "DB_HOST" }} pass={{ index . "DB_PASSWORD" }}`, testSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "host=localhost pass=s3cr3t"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestRender_MissingKeyReturnsError(t *testing.T) {
	r := template.New()
	_, err := r.Render(`{{ index . "MISSING_KEY" }}`, testSecrets)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRender_CustomDelimiters(t *testing.T) {
	r := template.NewWithDelims("[[", "]]")
	out, err := r.Render(`key=[[ index . "API_KEY" ]]`, testSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "key=abc123" {
		t.Errorf("got %q, want %q", out, "key=abc123")
	}
}

func TestRenderFile_WritesOutput(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tmpl.txt")
	dst := filepath.Join(dir, "out.txt")

	content := `DB_HOST={{ index . "DB_HOST" }}`
	if err := os.WriteFile(src, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	r := template.New()
	if err := r.RenderFile(src, dst, testSecrets, 0o600); err != nil {
		t.Fatalf("RenderFile error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "DB_HOST=localhost" {
		t.Errorf("got %q, want %q", string(got), "DB_HOST=localhost")
	}
}

func TestRenderFile_MissingSrcReturnsError(t *testing.T) {
	r := template.New()
	err := r.RenderFile("/nonexistent/tmpl.txt", "/tmp/out.txt", testSecrets, 0o600)
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestRenderFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tmpl.txt")
	dst := filepath.Join(dir, "out.txt")

	if err := os.WriteFile(src, []byte(`static`), 0o600); err != nil {
		t.Fatal(err)
	}

	r := template.New()
	if err := r.RenderFile(src, dst, testSecrets, 0o400); err != nil {
		t.Fatalf("RenderFile error: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o400 {
		t.Errorf("got perm %o, want %o", info.Mode().Perm(), 0o400)
	}
}
