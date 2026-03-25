package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AllowedIP represents an IP allowlist entry for an Exasol SaaS account.
type AllowedIP struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CIDRIp    string `json:"cidrIp"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	DeletedAt string `json:"deletedAt,omitempty"`
	DeletedBy string `json:"deletedBy,omitempty"`
}

// CreateAllowedIPRequest holds the parameters for adding a new allowed IP entry.
type CreateAllowedIPRequest struct {
	Name   string `json:"name"`
	CIDRIp string `json:"cidrIp"`
}

// UpdateAllowedIPRequest holds the parameters for replacing an allowed IP entry.
type UpdateAllowedIPRequest struct {
	Name   string `json:"name"`
	CIDRIp string `json:"cidrIp"`
}

func allowedIPsPath(accountID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/security/allowlist_ip", accountID)
}

func allowedIPPath(accountID, allowlistIpID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/security/allowlist_ip/%s", accountID, allowlistIpID)
}

func (c *Client) ListAllowedIPs(accountID string) ([]AllowedIP, error) {
	resp, err := c.do(http.MethodGet, allowedIPsPath(accountID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var ips []AllowedIP
	if err := json.NewDecoder(resp.Body).Decode(&ips); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return ips, nil
}

func (c *Client) GetAllowedIP(accountID, allowlistIpID string) (*AllowedIP, error) {
	resp, err := c.do(http.MethodGet, allowedIPPath(accountID, allowlistIpID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var ip AllowedIP
	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &ip, nil
}

func (c *Client) AddAllowedIP(accountID string, req CreateAllowedIPRequest) (*AllowedIP, error) {
	resp, err := c.do(http.MethodPost, allowedIPsPath(accountID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var ip AllowedIP
	if err := json.NewDecoder(resp.Body).Decode(&ip); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &ip, nil
}

func (c *Client) UpdateAllowedIP(accountID, allowlistIpID string, req UpdateAllowedIPRequest) error {
	resp, err := c.do(http.MethodPut, allowedIPPath(accountID, allowlistIpID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) DeleteAllowedIP(accountID, allowlistIpID string) error {
	resp, err := c.do(http.MethodDelete, allowedIPPath(accountID, allowlistIpID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}
