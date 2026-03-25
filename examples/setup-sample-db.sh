#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CLI="$REPO_ROOT/bin/exasol-saas"

# ---------------------------------------------------------------------------
# Mode parsing
# ---------------------------------------------------------------------------

MODE="${1:-create}"
DB_ID="${2:-}"

case "$MODE" in
  create) ;;
  existing)
    if [ -z "$DB_ID" ]; then
      echo "Error: MODE=existing requires a database ID as the second argument"
      echo "Usage: setup-sample-db.sh existing <db-id>"
      exit 1
    fi
    ;;
  *)
    echo "Error: unknown mode '$MODE'. Use 'create' or 'existing'."
    exit 1
    ;;
esac

# ---------------------------------------------------------------------------
# Prerequisites
# ---------------------------------------------------------------------------

if [ ! -f "$CLI" ]; then
  echo "Error: bin/exasol-saas not found. Build it first with:"
  echo "  task build"
  exit 1
fi

if ! command -v exapump &>/dev/null; then
  echo "Error: exapump is required. Install it with:"
  echo "  curl -fsSL https://raw.githubusercontent.com/exasol-labs/exapump/main/install.sh | sh"
  echo "See: https://github.com/exasol-labs/exapump"
  exit 1
fi

if ! command -v jq &>/dev/null; then
  echo "Error: jq is required. Install it with:"
  echo "  macOS:          brew install jq"
  echo "  Debian/Ubuntu:  sudo apt-get install -y jq"
  echo "  RHEL/Fedora:    sudo dnf install -y jq"
  exit 1
fi

# ---------------------------------------------------------------------------
# Required environment variables
# ---------------------------------------------------------------------------

if [ -z "${EXASOL_SAAS_TOKEN:-}" ]; then
  echo "Error: EXASOL_SAAS_TOKEN is required (personal access token for the SaaS API)"
  exit 1
fi

if [ -z "${EXASOL_SAAS_ACCOUNT_ID:-}" ]; then
  echo "Error: EXASOL_SAAS_ACCOUNT_ID is required"
  exit 1
fi

# ---------------------------------------------------------------------------
# create mode: Step 1 — Create the database
# ---------------------------------------------------------------------------

if [ "$MODE" = "create" ]; then
  echo "[1/6] Creating database..."
  DB_ID=$("$CLI" database create \
    --name         "SampleDatabase" \
    --region       eu-west-1 \
    --cluster-name "MainCluster" \
    --cluster-size S \
    --num-nodes    1 \
    --output json | jq -r '.id')
  echo "Created database: $DB_ID"

  # -------------------------------------------------------------------------
  # create mode: Step 2 — Wait for the database to become running
  # -------------------------------------------------------------------------

  echo "[2/6] Waiting for database to start..."
  while true; do
    STATUS=$("$CLI" database status "$DB_ID" --output json | jq -r '.status')
    echo "Database status: $STATUS"
    [ "$STATUS" = "running" ] && break
    sleep 10
  done
  echo "Database is running."
fi

# ---------------------------------------------------------------------------
# Step 3/1 — Get the cluster ID
# ---------------------------------------------------------------------------

if [ "$MODE" = "create" ]; then
  echo "[3/6] Getting cluster ID..."
else
  echo "[1/4] Getting cluster ID..."
fi
CLUSTER_ID=$("$CLI" cluster --database-id "$DB_ID" list --output json | jq -r '.[0].id')
echo "Cluster ID: $CLUSTER_ID"

# ---------------------------------------------------------------------------
# Step 4/2 — Allow public access
# ---------------------------------------------------------------------------

if [ "$MODE" = "create" ]; then
  echo "[4/6] Adding IP allowlist rule..."
else
  echo "[2/4] Adding IP allowlist rule..."
fi
"$CLI" security create --name "public" --cidr-ip "0.0.0.0/0" > /dev/null
echo "IP allowlist rule added."

# ---------------------------------------------------------------------------
# Step 5/3 — Get connection details
# ---------------------------------------------------------------------------

if [ "$MODE" = "create" ]; then
  echo "[5/6] Getting connection details..."
else
  echo "[3/4] Getting connection details..."
fi
CONN=$("$CLI" cluster --database-id "$DB_ID" connect "$CLUSTER_ID" --output json)
sleep 10

DB_HOST=$(echo "$CONN" | jq -r '.dns')
DB_PORT=$(echo "$CONN" | jq -r '.port')
DB_USER=$(echo "$CONN" | jq -r '.dbUsername')

# ---------------------------------------------------------------------------
# Step 6/4 — Verify connectivity with exapump
# ---------------------------------------------------------------------------

if [ "$MODE" = "create" ]; then
  echo "[6/6] Verifying connectivity..."
else
  echo "[4/4] Verifying connectivity..."
fi
DSN="exasol://${DB_USER}:${EXASOL_SAAS_TOKEN}@${DB_HOST}:${DB_PORT}/?tls=true&validateservercertificate=0"
exapump sql --dsn "$DSN" 'SELECT 1'
echo "Connection verified. Database is ready."
