package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// NamespaceListResult holds the response from listing namespaces
type NamespaceListResult struct {
	Items []struct {
		Name string `json:"name"`
	} `json:"items"`
}

// ResourceListResult holds the response from listing resources
type ResourceListResult struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"items"`
}

// ListNamespaces returns all available namespaces for completion
func (c *Client) ListNamespaces(ctx context.Context) (*NamespaceListResult, error) {
	// Use existing namespace list endpoint
	resp, err := c.Get(ctx, "/api/web/namespaces", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list namespaces: status %d", resp.StatusCode)
	}

	var result NamespaceListResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListResources returns resource names for completion (lightweight query)
func (c *Client) ListResources(ctx context.Context, rt *types.ResourceType, namespace string) (*ResourceListResult, error) {
	// Build API path
	path := rt.APIPath
	if rt.SupportsNamespace {
		path = fmt.Sprintf(path, namespace)
	}

	// Use lightweight query - only fetch names, not full configs
	query := url.Values{}
	query.Set("select", "metadata.name")

	resp, err := c.Get(ctx, path, query)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list resources: status %d", resp.StatusCode)
	}

	var result ResourceListResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetName returns the name of a resource at the given index
func (r *ResourceListResult) GetName(idx int) string {
	if idx >= len(r.Items) {
		return ""
	}
	return r.Items[idx].Metadata.Name
}
