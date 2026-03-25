# Spec: sample-db-setup

A self-contained bash script that demonstrates the full lifecycle of spinning up a new Exasol SaaS database and verifying connectivity via a `SELECT 1`.

## Output file

`examples/setup-sample-db.sh`

## Usage

```
setup-sample-db.sh [MODE] [DB_ID]
```

| Argument | Description |
|---|---|
| `MODE` | `create` (default) or `existing` |
| `DB_ID` | Required when `MODE=existing` — the ID of the already-running database to use |

### Mode: `create` (default)

Runs the full lifecycle: creates a new database, waits for it to start, then connects and verifies.

```bash
./examples/setup-sample-db.sh
# or
./examples/setup-sample-db.sh create
```

### Mode: `existing`

Skips database creation and the wait-for-running poll. Jumps directly to getting the cluster ID, adding the IP allowlist rule, getting connection details, and verifying connectivity.

```bash
./examples/setup-sample-db.sh existing <db-id>
```

If `existing` is given without a `<db-id>`, print an error and exit:
```
Error: MODE=existing requires a database ID as the second argument
Usage: setup-sample-db.sh existing <db-id>
```

Any other value for MODE is an error:
```
Error: unknown mode '<value>'. Use 'create' or 'existing'.
```

## Prerequisites

The script must check for both tools at startup. If either is missing, print a clear error with install instructions and exit non-zero.

### exasol-saas CLI binary

Check for `bin/exasol-saas` relative to the script's directory (i.e. the repo root `bin/`).

If not present, print a clear error with build instructions and exit non-zero:
```
bin/exasol-saas not found. Build it first with:
  task build
```

Use the binary via its path (`bin/exasol-saas`) throughout the script rather than relying on it being in `$PATH`.

### exapump

Check with `command -v exapump`.

Install instructions to show on failure:
```
exapump is required. Install it with:
  curl -fsSL https://raw.githubusercontent.com/exasol-labs/exapump/main/install.sh | sh
See: https://github.com/exasol-labs/exapump
```

### jq

Check with `command -v jq`.

Install instructions to show on failure:
```
jq is required. Install it with:
  macOS:          brew install jq
  Debian/Ubuntu:  sudo apt-get install -y jq
  RHEL/Fedora:    sudo dnf install -y jq
```

## Required environment variables

The script must verify all three are set before proceeding. If any is missing, print a descriptive error and exit.

| Variable | Purpose |
|---|---|
| `EXASOL_SAAS_TOKEN` | Personal access token for the SaaS API (used by the CLI and as the database password) |
| `EXASOL_SAAS_ACCOUNT_ID` | Account ID (used by the CLI) |

The CLI picks up `EXASOL_SAAS_TOKEN` and `EXASOL_SAAS_ACCOUNT_ID` automatically from the environment — no `--token` / `--account-id` flags needed in the script.

`EXASOL_SAAS_TOKEN` doubles as the database password when building the exapump DSN.

## Script flow

Steps 1–2 are skipped in `existing` mode. Steps 3–6 run in both modes.

### Step 1 — Create the database (create mode only)

```
exasol-saas database create \
  --name         "SampleDatabase" \
  --region       eu-west-1 \
  --cluster-name "MainCluster" \
  --cluster-size S \
  --num-nodes    1 \
  --output json
```

Capture the JSON output. Extract the database ID with jq:
```bash
DB_ID=$(... | jq -r '.id')
```

Print: `Created database: $DB_ID`

### Step 2 — Wait for the database to become running (create mode only)

Poll `database status <id>` in a sleep loop (10 s interval). Print progress on each iteration.

```bash
while true; do
  STATUS=$(exasol-saas database status "$DB_ID" --output json | jq -r '.status')
  echo "Database status: $STATUS"
  [ "$STATUS" = "running" ] && break
  sleep 10
done
```

Print: `Database is running.`

### Step 3 — Get the cluster ID

```bash
CLUSTER_ID=$(exasol-saas cluster --database-id "$DB_ID" list --output json | jq -r '.[0].id')
```

Print: `Cluster ID: $CLUSTER_ID`

### Step 4 — Allow public access

Add an IP allowlist rule that permits connections from any IP before attempting to connect:

```bash
"$CLI" security create --name "public" --cidr-ip "0.0.0.0/0"
```

Print: `IP allowlist rule added.`

### Step 5 — Get connection details

```bash
CONN=$(exasol-saas cluster --database-id "$DB_ID" connect "$CLUSTER_ID" --output json)
DB_HOST=$(echo "$CONN" | jq -r '.dns')
DB_PORT=$(echo "$CONN" | jq -r '.port')
DB_USER=$(echo "$CONN" | jq -r '.dbUsername')
```

### Step 6 — Verify connectivity with exapump

Build the DSN and run `SELECT 1`:

```bash
DSN="exasol://${DB_USER}:${EXASOL_SAAS_TOKEN}@${DB_HOST}:${DB_PORT}/?tls=true&validateservercertificate=0"
exapump sql --dsn "$DSN" 'SELECT 1'
```

Print on success: `Connection verified. Database is ready.`

## Error handling

- Use `set -euo pipefail` at the top of the script.
- Any failed CLI call or jq extraction exits the script immediately via `-e`.
- The polling loop is the only intentional wait; all other commands fail fast.

## Step numbering

Step numbers in the output reflect what actually runs:

- `create` mode: steps `[1/6]` through `[6/6]`
- `existing` mode: steps `[1/4]` through `[4/4]` (skips the create and wait steps; renumbers from 1)

## Example output

### create mode

```
[1/6] Creating database...
Created database: abc123
[2/6] Waiting for database to start...
Database status: starting
Database status: starting
Database status: running
Database is running.
[3/6] Getting cluster ID...
Cluster ID: cl456
[4/6] Adding IP allowlist rule...
IP allowlist rule added.
[5/6] Getting connection details...
[6/6] Verifying connectivity...
Connection verified. Database is ready.
```

### existing mode

```
[1/4] Getting cluster ID...
Cluster ID: cl456
[2/4] Adding IP allowlist rule...
IP allowlist rule added.
[3/4] Getting connection details...
[4/4] Verifying connectivity...
Connection verified. Database is ready.
```
