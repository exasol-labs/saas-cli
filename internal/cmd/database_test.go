package cmd_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/cmd"
)

func executeDatabaseCmd(srv *httptest.Server, args ...string) (string, error) {
	root, cfg := cmd.NewRootCmd()
	cfg.BaseURL = srv.URL
	root.AddCommand(cmd.NewDatabaseCmd(cfg))

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestDatabaseList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"db1","name":"MyDB","status":"RUNNING","provider":"aws","region":"us-east-1"}]`))
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "db1") {
		t.Errorf("expected output to contain %q, got: %q", "db1", out)
	}
	if !strings.Contains(out, "MyDB") {
		t.Errorf("expected output to contain %q, got: %q", "MyDB", out)
	}
}

func TestDatabaseList_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "No databases found") {
		t.Errorf("expected output to contain %q, got: %q", "No databases found", out)
	}
}

func TestDatabaseList_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"db1","name":"MyDB","status":"RUNNING","provider":"aws","region":"us-east-1"}]`))
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "list", "--output=json")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected output to contain %q, got: %q", `"id"`, out)
	}
	if !strings.Contains(out, "db1") {
		t.Errorf("expected output to contain %q, got: %q", "db1", out)
	}
}

func TestDatabaseStatus_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"db1","name":"MyDB","status":"RUNNING","provider":"aws","region":"us-east-1"}`))
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "status", "db1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "db1") {
		t.Errorf("expected output to contain %q, got: %q", "db1", out)
	}
}

func TestDatabaseStatus_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "status", "missing-id")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestDatabaseCreate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		db := map[string]any{"id": "new-db-id", "name": "test", "status": "CREATING", "provider": "aws", "region": "us-east-1"}
		_ = json.NewEncoder(w).Encode(db)
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv,
		"--token", "t", "--account-id", "acc",
		"database", "create",
		"--name", "test",
		"--region", "us-east-1",
		"--cluster-name", "main",
		"--cluster-size", "XS",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "new-db-id") {
		t.Errorf("expected output to contain %q, got: %q", "new-db-id", out)
	}
}

func TestDatabaseCreate_WithOptionalFlags(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "db2", "name": "test"})
	}))
	defer srv.Close()

	_, err := executeDatabaseCmd(srv,
		"--token", "t", "--account-id", "acc",
		"database", "create",
		"--name", "test",
		"--region", "us-east-1",
		"--cluster-name", "main",
		"--cluster-size", "XS",
		"--cluster-family", "memory",
		"--num-nodes", "3",
		"--stream-type", "stream",
		"--cluster-auto-stop-enabled",
		"--cluster-auto-stop-idle-time", "30",
		"--cluster-offload-enabled",
		"--cluster-offload-timeout-min", "60",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if body["numNodes"] != float64(3) {
		t.Errorf("expected numNodes 3, got %v", body["numNodes"])
	}
	cluster, _ := body["initialCluster"].(map[string]any)
	if cluster["family"] != "memory" {
		t.Errorf("expected cluster family 'memory', got %v", cluster["family"])
	}
	autoStop, _ := cluster["autoStop"].(map[string]any)
	if autoStop["enabled"] != true {
		t.Errorf("expected autoStop.enabled true, got %v", autoStop["enabled"])
	}
	if autoStop["idleTime"] != float64(30) {
		t.Errorf("expected autoStop.idleTime 30, got %v", autoStop["idleTime"])
	}
}

func TestDatabaseCreate_MissingFlag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "create")
	if err == nil {
		t.Fatal("expected error when required flags are missing, got nil")
	}
}

func TestDatabaseUpdate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "update", "db1", "--name", "new-name")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "updated") {
		t.Errorf("expected output to contain %q, got: %q", "updated", out)
	}
}

func TestDatabaseDelete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "delete", "db1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected output to contain %q, got: %q", "deleted", out)
	}
}

func TestDatabaseStart_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/start") {
			t.Errorf("expected path to end with /start, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "start", "db1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected output to contain %q, got: %q", "submitted", out)
	}
}

func TestDatabaseStop_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/stop") {
			t.Errorf("expected path to end with /stop, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeDatabaseCmd(srv, "--token", "t", "--account-id", "acc", "database", "stop", "db1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected output to contain %q, got: %q", "submitted", out)
	}
}
