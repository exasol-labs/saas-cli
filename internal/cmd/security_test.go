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

func executeSecurityCmd(srv *httptest.Server, args ...string) (string, error) {
	root, cfg := cmd.NewRootCmd()
	cfg.BaseURL = srv.URL
	root.AddCommand(cmd.NewSecurityCmd(cfg))

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestSecurityList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"ip1","name":"office","cidrIp":"1.2.3.4/32","createdAt":"2024-01-01T00:00:00Z","createdBy":"user@example.com"}]`))
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ip1") {
		t.Errorf("expected output to contain %q, got: %q", "ip1", out)
	}
	if !strings.Contains(out, "office") {
		t.Errorf("expected output to contain %q, got: %q", "office", out)
	}
}

func TestSecurityList_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "No allowed IPs found") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestSecurityStatus_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"ip1","name":"office","cidrIp":"1.2.3.4/32","createdAt":"2024-01-01T00:00:00Z","createdBy":"user@example.com"}`))
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "status", "ip1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ip1") {
		t.Errorf("expected output to contain %q, got: %q", "ip1", out)
	}
}

func TestSecurityStatus_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "status", "missing")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestSecurityCreate_Success(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "ip2", "name": "vpn", "cidrIp": "10.0.0.0/8"})
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv,
		"--token", "t", "--account-id", "acc",
		"security", "create", "--name", "vpn", "--cidr-ip", "10.0.0.0/8",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ip2") {
		t.Errorf("expected output to contain %q, got: %q", "ip2", out)
	}
	if body["name"] != "vpn" {
		t.Errorf("expected name 'vpn', got %v", body["name"])
	}
	if body["cidrIp"] != "10.0.0.0/8" {
		t.Errorf("expected cidrIp '10.0.0.0/8', got %v", body["cidrIp"])
	}
}

func TestSecurityCreate_MissingFlags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "create", "--name", "vpn")
	if err == nil {
		t.Fatal("expected error when required flags are missing, got nil")
	}
}

func TestSecurityUpdate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv,
		"--token", "t", "--account-id", "acc",
		"security", "update", "ip1", "--name", "renamed", "--cidr-ip", "5.6.7.8/32",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "updated") {
		t.Errorf("expected output to contain %q, got: %q", "updated", out)
	}
}

func TestSecurityDelete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	out, err := executeSecurityCmd(srv, "--token", "t", "--account-id", "acc", "security", "delete", "ip1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected output to contain %q, got: %q", "deleted", out)
	}
}
