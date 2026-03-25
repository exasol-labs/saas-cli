# Mission: exasol-saas CLI

> A command-line interface for managing Exasol SaaS resources — databases and clusters — via the Exasol SaaS REST API.

## Problem Statement

Exasol SaaS exposes a REST API for managing cloud databases and compute clusters. However, operations teams, data engineers, and developers need to perform these operations from the terminal — in scripts, CI pipelines, and day-to-day workflows — without writing custom HTTP clients. Existing web console tools are not scriptable. The `exasol-saas` CLI bridges this gap by wrapping the Exasol SaaS API with a familiar, composable command-line interface.

## Target Users

| Persona | Goal | Key Workflow |
|---------|------|--------------|
| Developer / Data Engineer | Query and manage databases from the terminal | List databases, check status, start/stop as needed |
| DevOps / Platform Engineer | Automate provisioning via scripts and CI | Script cluster creation, parse JSON output in pipelines |
| AI Agent | Autonomously spin up and tear down sandbox environments | Create databases, poll status, delete after use — all via CLI commands in tool calls |

## Core Capabilities

1. **Database Management** — Create, list, inspect, update, delete, start, and stop Exasol SaaS databases.
2. **Cluster Management** — Create, list, inspect, update, delete, start, stop, and scale clusters within a database.
3. **Flexible Output** — Display results as human-readable tables by default; switch to JSON output for scripting via `--output=json`.

## Out of Scope

- User authentication flows (login/logout commands); authentication is via Personal Access Token only.
- Extensions management (install, uninstall, instances).
- File operations (upload, download, folder management).
- Usage reporting and billing information.
- IP allowlist and security management.
- User and profile management.

## Domain Glossary

| Term | Definition |
|------|------------|
| Account | The top-level Exasol SaaS tenant, identified by `accountId`. All resources belong to an account. |
| Database | A logical Exasol SaaS database instance. Contains one or more clusters. |
| Cluster | A compute cluster attached to a database. Can be independently started, stopped, and scaled. |
| PAT | Personal Access Token — a long-lived API credential used to authenticate all CLI requests. |
| DLHC | Data Lakehouse Connector — an optional database feature (out of scope for initial version). |
| Schedule | An automated action (start/stop) configured on a database (out of scope for initial version). |

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Language | Go | CLI implementation |
| Module | `github.com/exasol-labs/saas-cli` | Go module path |
| CLI Framework | [cobra](https://github.com/spf13/cobra) | Command/subcommand structure, flags, help |
| HTTP Client | Go standard library `net/http` | API communication |
| Output formatting | [tablewriter](https://github.com/olekukonko/tablewriter) or similar | Human-readable table output |
| Testing | Go `testing` + integration tests | Integration tests against real or mock API |
| Task runner | [Task](https://taskfile.dev) | Build, test, lint, and e2e automation via `Taskfile.yml` |

## Commands

```bash
task build     # build binary to bin/exasol-saas
task test      # run unit and integration tests
task lint      # go vet + gofmt check
task e2e       # run e2e tests via bats
task coverage  # generate HTML coverage report
```

## Project Structure

```
saas-cli/
├── bin/                    # compiled binaries (gitignored)
├── cmd/
│   └── exasol-saas/        # main package — CLI entry point
├── e2e/                    # bats e2e tests
├── internal/
│   ├── api/                # Exasol SaaS HTTP client and request/response types
│   ├── cmd/                # Cobra command definitions (database, cluster subcommands)
│   └── output/             # Table and JSON formatters
├── specs/                  # Feature specs and plans
├── Taskfile.yml            # task runner
└── go.mod
```

## Architecture

The CLI follows a layered architecture:

- **Command layer** (`internal/cmd`): Cobra commands parse flags and arguments, delegate to the API client.
- **API client layer** (`internal/api`): Thin HTTP client wrapping the Exasol SaaS REST API. Handles authentication (PAT via env var `EXASOL_SAAS_TOKEN` or `--token` flag) and JSON serialization.
- **Output layer** (`internal/output`): Formats API responses as tables or JSON depending on `--output` flag.

Data flows: CLI args → Cobra command → API client → Exasol SaaS REST API → response → formatter → stdout.

## Constraints

- **Technical**: Must produce a single self-contained binary. No external runtime required.
- **Authentication**: All requests authenticated via PAT. Token sourced from `EXASOL_SAAS_TOKEN` environment variable or `--token` global flag.
- **Portability**: Must run on Linux, macOS, and Windows.
- **Provider**: Defaults to AWS, not exposed as an option.
- **CLI completeness**: Every API option must be exposed as a CLI flag. Commands must fail with a clear error if a required option is missing. Help text must explicitly label each flag as `(required)` or `(optional)`.

## External Dependencies

| Service | Purpose | Failure Impact |
|---------|---------|----------------|
| Exasol SaaS REST API (`cloud.exasol.com`) | All resource operations | CLI is non-functional without API access |
