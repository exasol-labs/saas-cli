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

func executeClusterCmd(srv *httptest.Server, args ...string) (string, error) {
	root, cfg := cmd.NewRootCmd()
	cfg.BaseURL = srv.URL
	root.AddCommand(cmd.NewClusterCmd(cfg))

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestClusterList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"cl1","name":"main","status":"running","size":"XS","familyName":"","mainCluster":true,"createdAt":"2024-01-01T00:00:00Z","createdBy":"user@example.com"}]`))
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "cl1") {
		t.Errorf("expected output to contain %q, got: %q", "cl1", out)
	}
	if !strings.Contains(out, "main") {
		t.Errorf("expected output to contain %q, got: %q", "main", out)
	}
}

func TestClusterList_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "No clusters found") {
		t.Errorf("expected output to contain %q, got: %q", "No clusters found", out)
	}
}

func TestClusterStatus_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"cl1","name":"main","status":"running","size":"XS","familyName":"","mainCluster":true,"createdAt":"2024-01-01T00:00:00Z","createdBy":"user@example.com"}`))
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "status", "cl1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "cl1") {
		t.Errorf("expected output to contain %q, got: %q", "cl1", out)
	}
}

func TestClusterStatus_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "status", "missing")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestClusterCreate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "cl2", "name": "new-cluster", "size": "XS"})
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv,
		"--token", "t", "--account-id", "acc",
		"cluster", "--database-id", "db1",
		"create", "--name", "new-cluster", "--size", "XS",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "cl2") {
		t.Errorf("expected output to contain %q, got: %q", "cl2", out)
	}
}

func TestClusterCreate_WithOptionalFlags(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "cl3", "name": "test"})
	}))
	defer srv.Close()

	_, err := executeClusterCmd(srv,
		"--token", "t", "--account-id", "acc",
		"cluster", "--database-id", "db1",
		"create",
		"--name", "test", "--size", "S",
		"--family", "memory",
		"--auto-stop-enabled", "--auto-stop-idle-time", "30",
		"--offload-enabled", "--offload-timeout-min", "60",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if body["family"] != "memory" {
		t.Errorf("expected family 'memory', got %v", body["family"])
	}
	autoStop, _ := body["autoStop"].(map[string]any)
	if autoStop["enabled"] != true {
		t.Errorf("expected autoStop.enabled true, got %v", autoStop["enabled"])
	}
}

func TestClusterCreate_MissingFlag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "create", "--name", "test")
	if err == nil {
		t.Fatal("expected error when required flags are missing, got nil")
	}
}

func TestClusterUpdate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "update", "cl1", "--name", "renamed")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "updated") {
		t.Errorf("expected output to contain %q, got: %q", "updated", out)
	}
}

func TestClusterUpdate_NoFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "update", "cl1")
	if err != nil {
		t.Fatalf("expected no error with no flags, got: %v", err)
	}
	if !strings.Contains(out, "updated") {
		t.Errorf("expected output to contain %q, got: %q", "updated", out)
	}
}

func TestClusterDelete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "delete", "cl1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected output to contain %q, got: %q", "deleted", out)
	}
}

func TestClusterScale_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/scale") {
			t.Errorf("expected path to end with /scale, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "scale", "cl1", "--size", "M")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected output to contain %q, got: %q", "submitted", out)
	}
}

func TestClusterStart_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/start") {
			t.Errorf("expected path to end with /start, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "start", "cl1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected output to contain %q, got: %q", "submitted", out)
	}
}

func TestClusterStop_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/stop") {
			t.Errorf("expected path to end with /stop, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "stop", "cl1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected output to contain %q, got: %q", "submitted", out)
	}
}

func TestClusterConnect_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/connect") {
			t.Errorf("expected path to end with /connect, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"dns":        "cluster.example.com",
			"port":       8563,
			"jdbc":       "jdbc:exa:cluster.example.com:8563",
			"dbUsername": "admin",
			"ips":        map[string]any{"private": []string{"10.0.0.1"}, "public": []string{"1.2.3.4"}},
		})
	}))
	defer srv.Close()

	out, err := executeClusterCmd(srv, "--token", "t", "--account-id", "acc", "cluster", "--database-id", "db1", "connect", "cl1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "cluster.example.com") {
		t.Errorf("expected output to contain DNS, got: %q", out)
	}
}
