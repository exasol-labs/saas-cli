# Spec: resources-api

CLI commands for managing the two core Exasol SaaS resources — databases and clusters.

## How to implement

Use the `/implement-api-resource` skill for each resource in order:

```
/implement-api-resource database
/implement-api-resource cluster
/implement-api-resource security
```

Implement `database` first. `cluster` depends on it because every cluster path requires a `--database-id`. `security` is independent and can be implemented last.

## Resources

### database

Manage Exasol SaaS databases.

**Scenarios**

- **List databases** — `database list` returns a table of all databases in the account
- **Get database status** — `database status <id>` returns details for a single database
- **Create database** — `database create` creates a new database with an initial cluster
- **Update database** — `database update <id> --name <name>` renames a database
- **Delete database** — `database delete <id>` deletes a database
- **Start database** — `database start <id>` starts a stopped database
- **Stop database** — `database stop <id>` stops a running database

### cluster

Manage clusters within a database. All commands require `--database-id`.

**Scenarios**

- **List clusters** — `cluster --database-id <id> list` returns all clusters in a database
- **Get cluster status** — `cluster --database-id <id> status <id>` returns details for a single cluster
- **Create cluster** — `cluster --database-id <id> create` creates a new cluster
- **Update cluster** — `cluster --database-id <id> update <id>` updates name or settings
- **Delete cluster** — `cluster --database-id <id> delete <id>` deletes a cluster
- **Scale cluster** — `cluster --database-id <id> scale <id> --size <size>` resizes a cluster
- **Start cluster** — `cluster --database-id <id> start <id>` starts a stopped cluster
- **Stop cluster** — `cluster --database-id <id> stop <id>` stops a running cluster
- **Get connection info** — `cluster --database-id <id> connect <id>` returns DNS, port, JDBC, and credentials

### security

Manage the account IP allowlist. Each entry permits traffic from a named CIDR range.

API path: `/api/v1/accounts/{accountId}/security/allowlist_ip`

**Schemas**

`AllowedIP` (response): `id`, `name`, `cidrIp`, `createdAt`, `createdBy`, `deletedAt`\*, `deletedBy`\* (\* optional)

`CreateAllowedIP` (request): `name` (required), `cidrIp` (required)

`UpdateAllowedIP` (request): `name` (required), `cidrIp` (required)

**Scenarios**

- **List allowed IPs** — `security list` returns all IP allowlist entries for the account
- **Get allowed IP** — `security status <id>` returns details for a single allowlist entry
- **Create allowed IP** — `security create --name <name> --cidr-ip <cidr>` adds a new IP range
- **Update allowed IP** — `security update <id> --name <name> --cidr-ip <cidr>` replaces an entry
- **Delete allowed IP** — `security delete <id>` removes an IP allowlist entry
