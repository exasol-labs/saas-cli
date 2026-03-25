package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ClusterSettings holds the read-only resolved settings for a cluster.
type ClusterSettings struct {
	OffloadEnabled    bool `json:"offloadEnabled"`
	OffloadTimeoutMin int  `json:"offloadTimeoutMin"`
}

// Cluster represents an Exasol SaaS cluster.
type Cluster struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Status      string           `json:"status"`
	Size        string           `json:"size"`
	FamilyName  string           `json:"familyName"`
	MainCluster bool             `json:"mainCluster"`
	CreatedAt   string           `json:"createdAt"`
	CreatedBy   string           `json:"createdBy"`
	DeletedAt   string           `json:"deletedAt,omitempty"`
	DeletedBy   string           `json:"deletedBy,omitempty"`
	AutoStop    *AutoStop        `json:"autoStop,omitempty"`
	Settings    *ClusterSettings `json:"settings,omitempty"`
}

// CreateClusterRequest holds the parameters for creating a new cluster.
type CreateClusterRequest struct {
	Name     string                 `json:"name"`
	Size     string                 `json:"size"`
	Family   string                 `json:"family,omitempty"`
	AutoStop *AutoStop              `json:"autoStop,omitempty"`
	Settings *ClusterSettingsUpdate `json:"settings,omitempty"`
}

// UpdateClusterRequest holds the parameters for updating an existing cluster.
type UpdateClusterRequest struct {
	Name     string                 `json:"name,omitempty"`
	AutoStop *AutoStop              `json:"autoStop,omitempty"`
	Settings *ClusterSettingsUpdate `json:"settings,omitempty"`
}

// ScaleClusterRequest holds the parameters for scaling a cluster.
type ScaleClusterRequest struct {
	Size   string `json:"size"`
	Family string `json:"family,omitempty"`
}

// ConnectionIPs holds the private and public IP addresses for a cluster.
type ConnectionIPs struct {
	Private []string `json:"private"`
	Public  []string `json:"public"`
}

// ClusterConnection holds the connection details for a cluster.
type ClusterConnection struct {
	DNS        string        `json:"dns"`
	Port       int           `json:"port"`
	JDBC       string        `json:"jdbc"`
	IPs        ConnectionIPs `json:"ips"`
	DBUsername string        `json:"dbUsername"`
}

func clustersPath(accountID, databaseID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/databases/%s/clusters", accountID, databaseID)
}

func clusterPath(accountID, databaseID, clusterID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/databases/%s/clusters/%s", accountID, databaseID, clusterID)
}

func (c *Client) ListClusters(accountID, databaseID string) ([]Cluster, error) {
	resp, err := c.do(http.MethodGet, clustersPath(accountID, databaseID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var clusters []Cluster
	if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return clusters, nil
}

func (c *Client) GetCluster(accountID, databaseID, clusterID string) (*Cluster, error) {
	resp, err := c.do(http.MethodGet, clusterPath(accountID, databaseID, clusterID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var cluster Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &cluster, nil
}

func (c *Client) CreateCluster(accountID, databaseID string, req CreateClusterRequest) (*Cluster, error) {
	resp, err := c.do(http.MethodPost, clustersPath(accountID, databaseID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var cluster Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &cluster, nil
}

func (c *Client) UpdateCluster(accountID, databaseID, clusterID string, req UpdateClusterRequest) error {
	resp, err := c.do(http.MethodPut, clusterPath(accountID, databaseID, clusterID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) DeleteCluster(accountID, databaseID, clusterID string) error {
	resp, err := c.do(http.MethodDelete, clusterPath(accountID, databaseID, clusterID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) ScaleCluster(accountID, databaseID, clusterID string, req ScaleClusterRequest) error {
	resp, err := c.do(http.MethodPut, clusterPath(accountID, databaseID, clusterID)+"/scale", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) StartCluster(accountID, databaseID, clusterID string) error {
	resp, err := c.do(http.MethodPut, clusterPath(accountID, databaseID, clusterID)+"/start", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) StopCluster(accountID, databaseID, clusterID string) error {
	resp, err := c.do(http.MethodPut, clusterPath(accountID, databaseID, clusterID)+"/stop", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) GetClusterConnection(accountID, databaseID, clusterID string) (*ClusterConnection, error) {
	resp, err := c.do(http.MethodGet, clusterPath(accountID, databaseID, clusterID)+"/connect", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var conn ClusterConnection
	if err := json.NewDecoder(resp.Body).Decode(&conn); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &conn, nil
}
