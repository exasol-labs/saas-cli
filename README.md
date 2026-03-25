# Exasol SaaS CLI

A command-line interface for managing Exasol SaaS resources — databases, clusters, and account security.

## 🚀 Get Started

### Install

```bash
curl -fsSL https://raw.githubusercontent.com/exasol-labs/saas-cli/main/install.sh | sh
```

### Explore commands

```bash
exasol-saas --help
```

## 🗄️ How to Create a Database

The sample script [`examples/setup-sample-db.sh`](examples/setup-sample-db.sh) walks through the full flow below end-to-end.

### 📦 Prerequisites

- **exapump** — to run queries against the database:
  ```bash
  curl -fsSL https://raw.githubusercontent.com/exasol-labs/exapump/main/install.sh | sh
  ```
- **jq** — to parse JSON output from the CLI:
  ```bash
  brew install jq          # macOS
  sudo apt-get install jq  # Debian/Ubuntu
  sudo dnf install jq      # RHEL/Fedora
  ```

### 🔑 Set up credentials

```bash
export EXASOL_SAAS_TOKEN=your-token
export EXASOL_SAAS_ACCOUNT_ID=your-account-id
```

### ⚡ From zero to first query

**1. 🗄️ Create a database**

```bash
DB_ID=$(exasol-saas database create \
  --name SampleDatabase --region eu-west-1 \
  --cluster-name MainCluster --cluster-size S --num-nodes 1 \
  --output json | jq -r '.id')
```

**2. 🔍 Get the cluster ID**

Once the database status reaches `running`:

```bash
CLUSTER_ID=$(exasol-saas cluster --database-id "$DB_ID" list --output json | jq -r '.[0].id')
```

**3. 🌐 Allow network access**

```bash
exasol-saas security create --name public --cidr-ip 0.0.0.0/0
```

**4. 🔌 Get connection details**

```bash
CONN=$(exasol-saas cluster --database-id "$DB_ID" connect "$CLUSTER_ID" --output json)
DB_HOST=$(echo "$CONN" | jq -r '.dns')
DB_PORT=$(echo "$CONN" | jq -r '.port')
DB_USER=$(echo "$CONN" | jq -r '.dbUsername')
```

**5. 🎉 Run a query**

```bash
exapump sql \
  --dsn "exasol://${DB_USER}:${EXASOL_SAAS_TOKEN}@${DB_HOST}:${DB_PORT}/?tls=true&validateservercertificate=0" \
  'SELECT 1'
```
