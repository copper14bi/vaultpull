// Package template renders Go text/templates populated with secrets fetched
// from HashiCorp Vault.
//
// # Overview
//
// Instead of writing a plain .env file, users can supply a custom template
// (e.g. a docker-compose override or a shell export script) and have
// vaultpull inject the resolved secret values at sync time.
//
// # Usage
//
//	r := template.New()
//	out, err := r.Render(tmplText, secrets)
//
// Or to read from / write to disk:
//
//	err := r.RenderFile("docker-compose.tmpl", "docker-compose.override.yml", secrets, 0o600)
//
// # Delimiters
//
// The default delimiters are {{ and }}.  If your template format already uses
// those characters (e.g. Helm charts), construct the renderer with custom
// delimiters via NewWithDelims.
package template
