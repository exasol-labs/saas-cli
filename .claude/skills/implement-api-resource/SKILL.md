---
name: implement-api-resource
description: Scaffold a new CLI resource end-to-end — API client, Cobra commands, output formatting, and tests — from the OpenAPI spec.
argument-hint: <resource-name>
disable-model-invocation: true
---

# Implement API Resource

Scaffold a complete new CLI resource for `$ARGUMENTS`, following codebase patterns and the Exasol SaaS OpenAPI spec.

## Phase 1: Read the Codebase

Read **all** of the following before writing any code:

- `specs/mission.md` — constraints (especially: provider never exposed, all flags mandatory/optional labelled)
- Every file in `internal/api/` — client patterns, types, helpers
- Every file in `internal/cmd/` — command patterns, flag conventions, wiring
- Every file in `internal/output/` — output formatting patterns
- `cmd/exasol-saas/` — entry point and binary setup

## Phase 2: Read the OpenAPI Spec

Fetch `https://cloud.exasol.com/openapi.json` and extract for the target resource:

- All REST paths and HTTP methods
- Request body schemas (create, update, scale, etc.) — every property, its type, and whether it is in the `required` array
- Response schema — every property and its requirement status
- Nested schemas referenced via `$ref` (resolve them fully)

The spec is the single source of truth. Do not invent or omit fields.

## Phase 3: Interview

Use `AskUserQuestion` to ask the following in a **single message** before writing any code:

1. **Operations** — default is every operation found in the spec. Should any be excluded?
2. **Parent resource** — does the path require a parent ID (e.g. `--database-id`)? Confirm the flag name.

Wait for the answer before continuing.

## Phase 4: Generate Files

Follow the exact patterns from `database.go`. Apply these flag rules without exception:

| Rule | Detail |
|------|--------|
| Expose every writable field | No API field may be left out |
| `provider` is never a flag | Hardcode to `"AWS"` where the API requires it |
| Required field | `MarkFlagRequired` + description must explain the field's purpose and end with `(required)` |
| Optional field | Description must explain the field's purpose — no label, no `(optional)` suffix |
| Nested objects | Flatten to individual flags (e.g. `--auto-stop-idle-time`) |
| Help text quality | Every flag must have a meaningful description — never leave it empty or use just the field name. Include valid values or examples where helpful (e.g. `"Cloud region, e.g. us-east-1"`, `"Cluster size, e.g. XS, S, M, L"`) |

### Files to generate

| File | Content |
|------|---------|
| `internal/api/<resource>.go` | Types and client methods for all operations |
| `internal/api/<resource>_test.go` | Tests for all API client methods |
| `internal/cmd/<resource>.go` | Cobra subcommands for all operations |
| `internal/cmd/<resource>_test.go` | Tests for all Cobra commands |

### Wire up

1. Add `new<Resource>Cmd(cfg)` to `internal/cmd/root.go` — same pattern as `newDatabaseCmd`
2. If the resource introduces a new Core Capability, add it to `specs/mission.md`

## Phase 5: Verify

```bash
task build
task test
task lint
```

Show the output of each command. Fix any errors before reporting done.

## Phase 6: Report

```
✓ internal/api/<resource>.go       — <N> methods
✓ internal/api/<resource>_test.go  — <N> tests
✓ internal/cmd/<resource>.go       — <N> subcommands (<list>)
✓ internal/cmd/<resource>_test.go  — <N> tests
✓ internal/cmd/root.go             — wired up
✓ specs/mission.md                 — updated / unchanged
```
