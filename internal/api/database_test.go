package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/api"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *api.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return api.NewClient("test-token", srv.URL)
}

func TestListDatabases(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer auth header, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.Database{
			{ID: "db1", Name: "my-db", Status: "running", CreatedAt: "2024-01-01T00:00:00Z", CreatedBy: "user@example.com"},
		})
	})

	dbs, err := client.ListDatabases("acct1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dbs) != 1 || dbs[0].ID != "db1" {
		t.Errorf("unexpected databases: %+v", dbs)
	}
	if dbs[0].CreatedAt != "2024-01-01T00:00:00Z" {
		t.Errorf("expected CreatedAt to be set, got %q", dbs[0].CreatedAt)
	}
}

func TestGetDatabase(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.Database{
			ID:        "db1",
			Name:      "my-db",
			CreatedAt: "2024-01-01T00:00:00Z",
			CreatedBy: "user@example.com",
			Clusters:  api.DatabaseClusters{Total: 2, Running: 1},
		})
	})

	db, err := client.GetDatabase("acct1", "db1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.ID != "db1" {
		t.Errorf("expected db ID db1, got %q", db.ID)
	}
	if db.Clusters.Total != 2 {
		t.Errorf("expected Clusters.Total 2, got %d", db.Clusters.Total)
	}
}

func TestGetDatabase_NotFound(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.GetDatabase("acct1", "db1")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestCreateDatabase(t *testing.T) {
	var received api.CreateDatabaseRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.Database{ID: "db2", Name: "new-db", CreatedAt: "2024-01-01T00:00:00Z"})
	})

	db, err := client.CreateDatabase("acct1", api.CreateDatabaseRequest{
		Name:     "new-db",
		Provider: api.DefaultProvider,
		Region:   "us-east-1",
		InitialCluster: api.InitialCluster{
			Name: "main",
			Size: "XS",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.ID != "db2" {
		t.Errorf("expected db ID db2, got %q", db.ID)
	}
	if received.Provider != api.DefaultProvider {
		t.Errorf("expected provider AWS, got %q", received.Provider)
	}
}

func TestCreateDatabase_WithOptionalFields(t *testing.T) {
	var received api.CreateDatabaseRequest
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(api.Database{ID: "db3", Name: "new-db"})
	})

	_, err := client.CreateDatabase("acct1", api.CreateDatabaseRequest{
		Name:       "new-db",
		Provider:   api.DefaultProvider,
		Region:     "us-east-1",
		NumNodes:   3,
		StreamType: "stream",
		InitialCluster: api.InitialCluster{
			Name:   "main",
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
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.NumNodes != 3 {
		t.Errorf("expected NumNodes 3, got %d", received.NumNodes)
	}
	if received.InitialCluster.AutoStop == nil || !received.InitialCluster.AutoStop.Enabled {
		t.Errorf("expected AutoStop to be set and enabled")
	}
	if received.InitialCluster.Settings == nil || received.InitialCluster.Settings.OffloadTimeoutMin != 60 {
		t.Errorf("expected Settings.OffloadTimeoutMin 60, got %+v", received.InitialCluster.Settings)
	}
}

func TestUpdateDatabase(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.UpdateDatabase("acct1", "db1", api.UpdateDatabaseRequest{Name: "renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteDatabase(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteDatabase("acct1", "db1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStartDatabase(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/start" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.StartDatabase("acct1", "db1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStopDatabase(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/acct1/databases/db1/stop" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := client.StopDatabase("acct1", "db1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
