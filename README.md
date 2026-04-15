# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with rotation support.

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/vaultpull/releases) page.

---

## Usage

Authenticate with Vault and pull secrets into a local `.env` file:

```bash
# Set your Vault address and token
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a Vault path into a .env file
vaultpull pull --path secret/data/myapp --out .env

# Pull and rotate secrets (re-generates dynamic credentials)
vaultpull pull --path secret/data/myapp --out .env --rotate

# Watch for changes and auto-sync
vaultpull watch --path secret/data/myapp --out .env --interval 60s
```

Your `.env` file will be created or updated with the latest secrets from Vault:

```
DB_PASSWORD=s3cr3t
API_KEY=abc123
TOKEN=xyz789
```

---

## Configuration

`vaultpull` can also be configured via a `vaultpull.yaml` file in your project root:

```yaml
vault_addr: https://vault.example.com
path: secret/data/myapp
out: .env
rotate: false
```

---

## License

[MIT](LICENSE) © 2024 yourusername