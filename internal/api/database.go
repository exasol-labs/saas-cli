package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DatabaseClusters holds cluster count information for a database.
type DatabaseClusters struct {
	Total   int `json:"total"`
	Running int `json:"running"`
}

// DatabaseIntegration represents an integration linked to a database.
type DatabaseIntegration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Database represents an Exasol SaaS database.
type Database struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Status       string                `json:"status"`
	Provider     string                `json:"provider"`
	Region       string                `json:"region"`
	Clusters     DatabaseClusters      `json:"clusters"`
	Integrations []DatabaseIntegration `json:"integrations,omitempty"`
	CreatedAt    string                `json:"createdAt"`
	CreatedBy    string                `json:"createdBy"`
	DeletedAt    string                `json:"deletedAt,omitempty"`
	DeletedBy    string                `json:"deletedBy,omitempty"`
}

// AutoStop configures automatic stopping of a cluster after an idle period.
type AutoStop struct {
	Enabled  bool `json:"enabled"`
	IdleTime int  `json:"idleTime"`
}

// ClusterSettingsUpdate holds updatable cluster settings.
type ClusterSettingsUpdate struct {
	OffloadEnabled    bool `json:"offloadEnabled,omitempty"`
	OffloadTimeoutMin int  `json:"offloadTimeoutMin,omitempty"`
}

// InitialCluster describes the first cluster to provision with a new database.
type InitialCluster struct {
	Name     string                 `json:"name"`
	Size     string                 `json:"size"`
	Family   string                 `json:"family,omitempty"`
	AutoStop *AutoStop              `json:"autoStop,omitempty"`
	Settings *ClusterSettingsUpdate `json:"settings,omitempty"`
}

// CreateDatabaseRequest holds the parameters for creating a new database.
type CreateDatabaseRequest struct {
	Name           string         `json:"name"`
	Provider       string         `json:"provider"`
	Region         string         `json:"region"`
	InitialCluster InitialCluster `json:"initialCluster"`
	NumNodes       int            `json:"numNodes,omitempty"`
	StreamType     string         `json:"streamType,omitempty"`
}

// UpdateDatabaseRequest holds the parameters for updating an existing database.
type UpdateDatabaseRequest struct {
	Name string `json:"name"`
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}

func databasesPath(accountID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/databases", accountID)
}

func databasePath(accountID, id string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/databases/%s", accountID, id)
}

func (c *Client) ListDatabases(accountID string) ([]Database, error) {
	resp, err := c.do(http.MethodGet, databasesPath(accountID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var databases []Database
	if err := json.NewDecoder(resp.Body).Decode(&databases); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return databases, nil
}

func (c *Client) GetDatabase(accountID, id string) (*Database, error) {
	resp, err := c.do(http.MethodGet, databasePath(accountID, id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var db Database
	if err := json.NewDecoder(resp.Body).Decode(&db); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &db, nil
}

// CreateDatabase creates a new database and returns it, including the server-assigned ID.
func (c *Client) CreateDatabase(accountID string, req CreateDatabaseRequest) (*Database, error) {
	resp, err := c.do(http.MethodPost, databasesPath(accountID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var db Database
	if err := json.NewDecoder(resp.Body).Decode(&db); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &db, nil
}

func (c *Client) UpdateDatabase(accountID, id string, req UpdateDatabaseRequest) error {
	resp, err := c.do(http.MethodPut, databasePath(accountID, id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) DeleteDatabase(accountID, id string) error {
	resp, err := c.do(http.MethodDelete, databasePath(accountID, id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) StartDatabase(accountID, id string) error {
	resp, err := c.do(http.MethodPut, databasePath(accountID, id)+"/start", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) StopDatabase(accountID, id string) error {
	resp, err := c.do(http.MethodPut, databasePath(accountID, id)+"/stop", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}
