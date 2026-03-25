package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/api"
	"github.com/exasol-labs/saas-cli/internal/output"
)

var testDatabases = []api.Database{
	{ID: "db1", Name: "production", Status: "running", Provider: api.DefaultProvider, Region: "us-east-1"},
	{ID: "db2", Name: "staging", Status: "stopped", Provider: "GCP", Region: "eu-west-1"},
}

var testDatabase = api.Database{
	ID: "db1", Name: "production", Status: "running", Provider: api.DefaultProvider, Region: "us-east-1",
}

func TestPrintDatabases_Table_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabases(&buf, testDatabases, output.Table)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"ID", "Name", "Status", "Provider", "Region"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected header column %q in table output, got:\n%s", col, out)
		}
	}
}

func TestPrintDatabases_Table_ContainsRows(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabases(&buf, testDatabases, output.Table)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "production") || !strings.Contains(out, "staging") {
		t.Errorf("expected row data in table output, got:\n%s", out)
	}
}

func TestPrintDatabases_JSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabases(&buf, testDatabases, output.JSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []api.Database
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if len(result) != 2 {
		t.Errorf("expected 2 databases in JSON, got %d", len(result))
	}
}

func TestPrintDatabase_Table_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabase(&buf, testDatabase, output.Table)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"ID", "Name", "Status", "Provider", "Region"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected header column %q in table output, got:\n%s", col, out)
		}
	}
}

func TestPrintDatabase_Table_ContainsRow(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabase(&buf, testDatabase, output.Table)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected db name in table output, got:\n%s", out)
	}
}

func TestPrintDatabase_JSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintDatabase(&buf, testDatabase, output.JSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result api.Database
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if result.ID != "db1" {
		t.Errorf("expected ID db1, got %q", result.ID)
	}
}
