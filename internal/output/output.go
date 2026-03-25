// Package output handles formatting and printing of CLI results.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/exasol-labs/saas-cli/internal/api"
)

// Format specifies how results are rendered.
type Format string

const (
	// Table renders results as a human-readable table.
	Table Format = "table"
	// JSON renders results as JSON.
	JSON Format = "json"
)

// PrintDatabases writes a list of databases to w in the given format.
func PrintDatabases(w io.Writer, databases []api.Database, format Format) error {
	if format == JSON {
		return printJSON(w, databases)
	}
	return printDatabaseTable(w, databases)
}

// PrintDatabase writes a single database to w in the given format.
func PrintDatabase(w io.Writer, db api.Database, format Format) error {
	if format == JSON {
		return printJSON(w, db)
	}
	return printDatabaseTable(w, []api.Database{db})
}

func printDatabaseTable(w io.Writer, databases []api.Database) error {
	tw := newTabWriter(w)
	fmt.Fprintln(tw, "ID\tName\tStatus\tProvider\tRegion\tCreatedAt")
	for _, db := range databases {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", db.ID, db.Name, db.Status, db.Provider, db.Region, db.CreatedAt)
	}
	return tw.Flush()
}

// PrintClusters writes a list of clusters to w in the given format.
func PrintClusters(w io.Writer, clusters []api.Cluster, format Format) error {
	if format == JSON {
		return printJSON(w, clusters)
	}
	return printClusterTable(w, clusters)
}

// PrintCluster writes a single cluster to w in the given format.
func PrintCluster(w io.Writer, cluster api.Cluster, format Format) error {
	if format == JSON {
		return printJSON(w, cluster)
	}
	return printClusterTable(w, []api.Cluster{cluster})
}

// PrintClusterConnection writes cluster connection details to w in the given format.
func PrintClusterConnection(w io.Writer, conn api.ClusterConnection, format Format) error {
	if format == JSON {
		return printJSON(w, conn)
	}
	tw := newTabWriter(w)
	fmt.Fprintln(tw, "DNS\tPort\tJDBC\tUsername")
	fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", conn.DNS, conn.Port, conn.JDBC, conn.DBUsername)
	return tw.Flush()
}

func printClusterTable(w io.Writer, clusters []api.Cluster) error {
	tw := newTabWriter(w)
	fmt.Fprintln(tw, "ID\tName\tStatus\tSize\tFamily\tMain\tCreatedAt")
	for _, cl := range clusters {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%v\t%s\n",
			cl.ID, cl.Name, cl.Status, cl.Size, cl.FamilyName, cl.MainCluster, cl.CreatedAt)
	}
	return tw.Flush()
}

// PrintAllowedIPs writes a list of allowed IP entries to w in the given format.
func PrintAllowedIPs(w io.Writer, ips []api.AllowedIP, format Format) error {
	if format == JSON {
		return printJSON(w, ips)
	}
	return printAllowedIPTable(w, ips)
}

// PrintAllowedIP writes a single allowed IP entry to w in the given format.
func PrintAllowedIP(w io.Writer, ip api.AllowedIP, format Format) error {
	if format == JSON {
		return printJSON(w, ip)
	}
	return printAllowedIPTable(w, []api.AllowedIP{ip})
}

func printAllowedIPTable(w io.Writer, ips []api.AllowedIP) error {
	tw := newTabWriter(w)
	fmt.Fprintln(tw, "ID\tName\tCIDR\tCreatedAt\tCreatedBy")
	for _, ip := range ips {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", ip.ID, ip.Name, ip.CIDRIp, ip.CreatedAt, ip.CreatedBy)
	}
	return tw.Flush()
}

func newTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
}

func printJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
