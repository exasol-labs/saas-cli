package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/api"
)

func TestListAllowedIPs(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.AllowedIP{
			{ID: "ip1", Name: "office", CIDRIp: "1.2.3.4/32", CreatedAt: "2024-01-01T00:00:00Z", CreatedBy: "user@example.com"},
		})
	})

	ips, err := client.ListAllowedIPs("acct1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 1 || ips[0].ID != "ip1" {
		t.Errorf("unexpected result: %+v", ips)
	}
}

func TestGetAllowedIP(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.AllowedIP{
			ID:        "ip1",
			Name:      "office",
			CIDRIp:    "1.2.3.4/32",
			CreatedAt: "2024-01-01T00:00:00Z",
			CreatedBy: "user@example.com",
		})
	})

	ip, err := client.GetAllowedIP("acct1", "ip1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip.ID != "ip1" {
		t.Errorf("expected ID ip1, got %q", ip.ID)
	}
}

func TestGetAllowedIP_NotFound(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.GetAllowedIP("acct1", "missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestAddAllowedIP(t *testing.T) {
	var received api.CreateAllowedIPRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.AllowedIP{ID: "ip2", Name: "vpn", CIDRIp: "10.0.0.0/8"})
	})

	ip, err := client.AddAllowedIP("acct1", api.CreateAllowedIPRequest{
		Name:   "vpn",
		CIDRIp: "10.0.0.0/8",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip.ID != "ip2" {
		t.Errorf("expected ID ip2, got %q", ip.ID)
	}
	if received.Name != "vpn" || received.CIDRIp != "10.0.0.0/8" {
		t.Errorf("unexpected request body: %+v", received)
	}
}

func TestUpdateAllowedIP(t *testing.T) {
	var received api.UpdateAllowedIPRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.UpdateAllowedIP("acct1", "ip1", api.UpdateAllowedIPRequest{
		Name:   "renamed",
		CIDRIp: "5.6.7.8/32",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Name != "renamed" || received.CIDRIp != "5.6.7.8/32" {
		t.Errorf("unexpected request body: %+v", received)
	}
}

func TestDeleteAllowedIP(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteAllowedIP("acct1", "ip1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
