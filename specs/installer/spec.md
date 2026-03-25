# Spec: installer

A one-liner shell script that downloads and installs the correct pre-built `exasol-saas` binary for the current platform directly from the latest GitHub release.

## Output file

`install.sh` (repo root)

## Purpose

Allow users to install the CLI without cloning the repo, building from source, or manually selecting the right binary. A single `curl | sh` command handles everything.

## Usage

```bash
curl -fsSL https://raw.githubusercontent.com/exasol-labs/saas-cli/main/install.sh | sh
```

Users may also override the install directory:

```bash
curl -fsSL .../install.sh | INSTALL_DIR=/usr/local/bin sh
```

## Platform detection

Detect OS and architecture at runtime using `uname`.

### Supported platforms

| OS | Arch | Binary name |
|---|---|---|
| macOS | x86_64 / amd64 | `exasol-saas-darwin-amd64` |
| macOS | arm64 (Apple Silicon) | `exasol-saas-darwin-arm64` |
| Linux | x86_64 / amd64 | `exasol-saas-linux-amd64` |
| Linux | aarch64 / arm64 | `exasol-saas-linux-arm64` |

Windows is not supported by this script (Windows users should download manually from the releases page).

If the detected OS/arch combination is unsupported, print a clear error with a link to the releases page and exit non-zero:
```
Error: unsupported platform: <OS>/<arch>
Download manually from: https://github.com/exasol-labs/saas-cli/releases
```

## Version resolution

By default, resolve the latest release version from the GitHub API:

```bash
VERSION=$(curl -fsSL https://api.github.com/repos/exasol-labs/saas-cli/releases/latest \
  | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
```

Users may pin a specific version via the `VERSION` environment variable:

```bash
curl -fsSL .../install.sh | VERSION=v0.1.0 sh
```

If version resolution fails (e.g. no network access to GitHub API, no releases published yet), print a clear error and exit non-zero.

## Download

Construct the download URL from the resolved version and binary name:

```
https://github.com/exasol-labs/saas-cli/releases/download/<VERSION>/<BINARY>
```

Download to a temporary file using `curl -fsSL`. On failure, print a clear error and exit non-zero.

## Install directory

Default install directory: `/usr/local/bin`.

Override via the `INSTALL_DIR` environment variable.

After downloading, move the binary to `$INSTALL_DIR/exasol-saas` and make it executable (`chmod +x`).

If the install directory is not writable, the `mv` will fail and bash's `-e` flag exits with a clear OS error. Do not attempt to use `sudo` automatically — the user should either run with appropriate permissions or override `INSTALL_DIR` to a writable path.

## Verification

After installing, run `exasol-saas --version` to confirm the binary works. Print the output.

On success, print:
```
exasol-saas <version> installed to <install-dir>/exasol-saas
```

## Error handling

- Use `set -euo pipefail`.
- Use `curl -fsSL` for all downloads so HTTP errors cause immediate failure.
- Clean up the temporary file on exit (success or failure) using a `trap`.

## Prerequisites

The script requires only `curl` and basic POSIX tools (`uname`, `mv`, `chmod`). Both are available by default on macOS and all major Linux distributions. No additional dependencies.

## Example output

```
Detecting platform... darwin/arm64
Resolving latest version... v0.1.0
Downloading exasol-saas-darwin-arm64...
Installing to /usr/local/bin/exasol-saas...
exasol-saas v0.1.0 installed to /usr/local/bin/exasol-saas
```
