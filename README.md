# confy

> 🔐 Easily manage secrets in your files while keeping everything safe for version control.

`confy` is a lightweight Go CLI that checks secrets in and out of configuration files using a KeePass database.

- `checkin`: replaces secret values in files with placeholders which can be safely committed to Git
- `checkout`: restores placeholders back to real secret values
- `add`: adds new secrets to the vault
- `init`: creates a new KeePass DB and optionally a config file

## Table of Contents

1. [Why confy?](#why-confy)
2. [Features](#features)
3. [Quick Start](#quick-start)
4. [Installation](#installation)
5. [Configuration](#configuration)
6. [CLI Commands](#cli-commands)
7. [Workflow Example](#workflow-example)
8. [Security Notes](#security-notes)
9. [Development](#development)
10. [Roadmap](#roadmap)
11. [License](#license)

## Why confy? 🧩

Many projects keep secrets inside files that still need to be versioned, reviewed, and shared. `confy` lets you manage those secrets in a KeePass vault while keeping the files themselves in Git with safe placeholders.

Once a secret is added to the vault, `confy` can replace that secret throughout all matching files so you do not have to handle each occurrence manually.

This is especially useful for Kubernetes manifests, where configuration files are often committed, but secret values should not be stored in plaintext.

Before `checkin`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  template:
    spec:
      containers:
        - name: app
          env:
            - name: API_TOKEN
              value: super-secret-token
```

After `checkin`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  template:
    spec:
      containers:
        - name: app
          env:
            - name: API_TOKEN
              value: confy_secret_api_token
```

and is restored during `checkout`.

## Features ✨

- KeePass-based secret management (`.kdbx`) 🔐
- File extension filters (`.env`, `.yaml`, `.json`, ...) 📄
- Interactive password prompt when not provided via CLI ⌨️
- YAML config file support (`confy.yaml`) for repeatable workflows ⚙️
- Clean and simple CLI for local and CI-like usage 🚀

## Quick Start 🚀

If you do not want to build `confy` yourself, download the latest binary from the GitHub Releases page for your platform and make it executable.

```bash
# 1) Download the latest release binary from GitHub Releases

# 2) Make it executable
chmod +x ./confy

# 3) Show help
./confy help

# 4) Initialize vault + config
./confy init --dbPath confy.kdbx --sourceDir .

# 5) Add a secret
./confy add --entryName db_password --entryValue super-secret

# 6) Remove secrets from files
./confy checkin

# 7) Resolve placeholders back to real secret values
./confy checkout
```

## Installation 📦

### Download a Release Binary

The easiest way to get started is to download the latest binary from the GitHub Releases page.

1. Open the project's [Releases page](https://github.com/x-46/confy/releases).
2. Download the archive that matches your operating system and architecture.
3. Extract the binary.
4. Make it executable on Unix-like systems with `chmod +x ./confy`.
5. Run `./confy help` to verify the installation.

### From Source

#### Prerequisites

- Go `>= 1.26`


```bash
git clone https://github.com/x-46/confy.git
cd confy
go mod download
go build -o confy .
```

Then run either the built binary:

```bash
./confy help
```

or run directly with Go:

```bash
go run . help
```

## Configuration ⚙️

If a `confy.yaml` file exists in your current directory, it is loaded automatically.

Example:

```yaml
sourceDir: .
dbPath: confy.kdbx
fileExtensions:
  - .md
  - .yml
  - .yaml
  - .json
  - .env
```

### Relevant Options

- `sourceDir`: directory processed recursively
- `dbPath`: path to the KeePass database
- `fileExtensions`: list of file extensions to process
- `configFilePath`: explicit path to the config file
- `password`: vault password (interactive prompt if omitted)
- `entryName`: name for a new secret entry (`add`)
- `entryValue`: value for a new secret entry (`add`)

## CLI Commands 🧰

### `help`

Shows command overview or detailed command help.

```bash
confy help
confy checkout --help
```

### `init`

Initializes a new KeePass database and optionally writes a config.

```bash
confy init --dbPath confy.kdbx --configFilePath confy.yaml --sourceDir .
```

### `add`

Creates a secret entry in the vault. Once registered, that secret can be replaced across all matching files during `checkin`.

```bash
confy add --entryName api_key --entryValue supersecret --dbPath confy.kdbx
```

If `--entryValue` is omitted, `confy` asks for it interactively.

### `checkin`

Replaces secret values in files with placeholders like `confy_secret_<entryName>`.

```bash
confy checkin --sourceDir . --dbPath confy.kdbx --fileExtensions .env --fileExtensions .yaml
```

### `checkout`

Replaces placeholders with real secret values from the vault.

```bash
confy checkout --sourceDir . --dbPath confy.kdbx --fileExtensions .env --fileExtensions .yaml
```

## Development 🛠️

### Tests

```bash
go test ./...
```

## Roadmap 🗺️
- Git integration in pre-commit hooks
- Dry-run mode with diff output
- Improved CI integration

## License 📄

This project is licensed under the terms described in [LICENSE](LICENSE).
