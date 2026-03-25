package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/api"
)

func TestListClusters(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.Cluster{
			{ID: "cl1", Name: "main", Status: "running", Size: "XS", CreatedAt: "2024-01-01T00:00:00Z"},
		})
	})

	clusters, err := client.ListClusters("acct1", "db1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clusters) != 1 || clusters[0].ID != "cl1" {
		t.Errorf("unexpected clusters: %+v", clusters)
	}
}

func TestGetCluster(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.Cluster{
			ID:        "cl1",
			Name:      "main",
			Size:      "XS",
			CreatedAt: "2024-01-01T00:00:00Z",
			CreatedBy: "user@example.com",
		})
	})

	cluster, err := client.GetCluster("acct1", "db1", "cl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cluster.ID != "cl1" {
		t.Errorf("expected cluster ID cl1, got %q", cluster.ID)
	}
}

func TestGetCluster_NotFound(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.GetCluster("acct1", "db1", "missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestCreateCluster(t *testing.T) {
	var received api.CreateClusterRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(api.Cluster{ID: "cl2", Name: "new-cluster"})
	})

	cluster, err := client.CreateCluster("acct1", "db1", api.CreateClusterRequest{
		Name: "new-cluster",
		Size: "XS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cluster.ID != "cl2" {
		t.Errorf("expected cluster ID cl2, got %q", cluster.ID)
	}
	if received.Name != "new-cluster" || received.Size != "XS" {
		t.Errorf("unexpected request body: %+v", received)
	}
}

func TestCreateCluster_WithOptionalFields(t *testing.T) {
	var received api.CreateClusterRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(api.Cluster{ID: "cl3"})
	})

	_, err := client.CreateCluster("acct1", "db1", api.CreateClusterRequest{
		Name:   "cl3",
		Size:   "S",
		Family: "memory",
		AutoStop: &api.AutoStop{
			Enabled:  true,
			IdleTime: 30,
		},
		Settings: &api.ClusterSettingsUpdate{
			OffloadEnabled:    true,
			OffloadTimeoutMin: 60,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Family != "memory" {
		t.Errorf("expected family 'memory', got %q", received.Family)
	}
	if received.AutoStop == nil || !received.AutoStop.Enabled {
		t.Errorf("expected AutoStop to be set and enabled")
	}
}

func TestUpdateCluster(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.UpdateCluster("acct1", "db1", "cl1", api.UpdateClusterRequest{Name: "renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteCluster(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteCluster("acct1", "db1", "cl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScaleCluster(t *testing.T) {
	var received api.ScaleClusterRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/clusters/cl1/scale" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	})

	err := client.ScaleCluster("acct1", "db1", "cl1", api.ScaleClusterRequest{Size: "M"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Size != "M" {
		t.Errorf("expected size M, got %q", received.Size)
	}
}

func TestStartCluster(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/clusters/cl1/start" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.StartCluster("acct1", "db1", "cl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStopCluster(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/clusters/cl1/stop" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.StopCluster("acct1", "db1", "cl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetClusterConnection(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/clusters/cl1/connect" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.ClusterConnection{
			DNS:        "cluster.example.com",
			Port:       8563,
			JDBC:       "jdbc:exa:cluster.example.com:8563",
			DBUsername: "admin",
			IPs: api.ConnectionIPs{
				Private: []string{"10.0.0.1"},
				Public:  []string{"1.2.3.4"},
			},
		})
	})

	conn, err := client.GetClusterConnection("acct1", "db1", "cl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.DNS != "cluster.example.com" {
		t.Errorf("expected DNS cluster.example.com, got %q", conn.DNS)
	}
	if conn.Port != 8563 {
		t.Errorf("expected port 8563, got %d", conn.Port)
	}
}
