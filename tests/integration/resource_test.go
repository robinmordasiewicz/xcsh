package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/robinmordasiewicz/vesctl/pkg/client"
)

// TestResource_ListHTTPLoadBalancers tests listing HTTP load balancers
func TestResource_ListHTTPLoadBalancers(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// List HTTP load balancers in shared namespace
	resp, err := c.Get(ctx, "/api/config/namespaces/shared/http_loadbalancers", nil)
	if err != nil {
		t.Fatalf("List HTTP LBs failed: %v", err)
	}

	t.Logf("HTTP Load Balancers - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d HTTP load balancers", len(items))
				for i, item := range items {
					if i >= 5 {
						t.Logf("... and %d more", len(items)-5)
						break
					}
					if lb, ok := item.(map[string]interface{}); ok {
						if meta, ok := lb["metadata"].(map[string]interface{}); ok {
							t.Logf("  - %s", meta["name"])
						}
					}
				}
			}
		}
	}
}

// TestResource_ListOriginPools tests listing origin pools
func TestResource_ListOriginPools(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/shared/origin_pools", nil)
	if err != nil {
		t.Fatalf("List origin pools failed: %v", err)
	}

	t.Logf("Origin Pools - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d origin pools", len(items))
			}
		}
	}
}

// TestResource_ListHealthchecks tests listing healthchecks
func TestResource_ListHealthchecks(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/shared/healthchecks", nil)
	if err != nil {
		t.Fatalf("List healthchecks failed: %v", err)
	}

	t.Logf("Healthchecks - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d healthchecks", len(items))
			}
		}
	}
}

// TestResource_ListAppFirewalls tests listing app firewalls
func TestResource_ListAppFirewalls(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/shared/app_firewalls", nil)
	if err != nil {
		t.Fatalf("List app firewalls failed: %v", err)
	}

	t.Logf("App Firewalls - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d app firewalls", len(items))
			}
		}
	}
}

// TestResource_ListServicePolicies tests listing service policies
func TestResource_ListServicePolicies(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/shared/service_policys", nil)
	if err != nil {
		t.Fatalf("List service policies failed: %v", err)
	}

	t.Logf("Service Policies - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d service policies", len(items))
			}
		}
	}
}

// TestResource_CRUD_HealthCheck tests create, read, update, delete for healthcheck
func TestResource_CRUD_HealthCheck(t *testing.T) {
	c := getTestClient(t)

	testName := fmt.Sprintf("test-healthcheck-%d", time.Now().Unix())
	namespace := "shared"

	// Create
	t.Run("Create", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		payload := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      testName,
				"namespace": namespace,
				"labels": map[string]string{
					"test": "true",
				},
			},
			"spec": map[string]interface{}{
				"http_health_check": map[string]interface{}{
					"path": "/health",
				},
				"timeout":       5,
				"interval":      30,
				"unhealthy_threshold": 2,
				"healthy_threshold":   3,
			},
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/api/config/namespaces/%s/healthchecks", namespace), payload)
		if err != nil {
			t.Fatalf("Create healthcheck failed: %v", err)
		}

		if resp.StatusCode >= 400 {
			t.Logf("Create failed with status %d: %s", resp.StatusCode, string(resp.Body))
			// Don't fail immediately - some environments may restrict this
			if resp.StatusCode == 403 {
				t.Skip("Permission denied - skipping CRUD test")
			}
			t.Fatalf("Create failed with status %d", resp.StatusCode)
		}

		t.Logf("Created healthcheck %s", testName)
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := c.Get(ctx, fmt.Sprintf("/api/config/namespaces/%s/healthchecks/%s", namespace, testName), nil)
		if err != nil {
			t.Fatalf("Get healthcheck failed: %v", err)
		}

		if resp.StatusCode >= 400 {
			t.Fatalf("Get healthcheck failed with status %d: %s", resp.StatusCode, string(resp.Body))
		}

		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if meta, ok := result["metadata"].(map[string]interface{}); ok {
			if name, ok := meta["name"].(string); ok && name != testName {
				t.Errorf("Expected name %s, got %s", testName, name)
			}
		}

		t.Logf("Read healthcheck: %+v", result)
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		payload := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      testName,
				"namespace": namespace,
				"labels": map[string]string{
					"test":    "true",
					"updated": "true",
				},
			},
			"spec": map[string]interface{}{
				"http_health_check": map[string]interface{}{
					"path": "/healthz",
				},
				"timeout":       10,
				"interval":      60,
				"unhealthy_threshold": 3,
				"healthy_threshold":   2,
			},
		}

		resp, err := c.Put(ctx, fmt.Sprintf("/api/config/namespaces/%s/healthchecks/%s", namespace, testName), payload)
		if err != nil {
			t.Fatalf("Update healthcheck failed: %v", err)
		}

		if resp.StatusCode >= 400 {
			t.Fatalf("Update healthcheck failed with status %d: %s", resp.StatusCode, string(resp.Body))
		}

		t.Logf("Updated healthcheck %s", testName)
	})

	// Delete (cleanup)
	t.Run("Delete", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		resp, err := c.Delete(ctx, fmt.Sprintf("/api/config/namespaces/%s/healthchecks/%s", namespace, testName))
		if err != nil {
			t.Fatalf("Delete healthcheck failed: %v", err)
		}

		if resp.StatusCode >= 400 && resp.StatusCode != 404 {
			t.Fatalf("Delete healthcheck failed with status %d: %s", resp.StatusCode, string(resp.Body))
		}

		t.Logf("Deleted healthcheck %s", testName)
	})

	// Verify deletion
	t.Run("VerifyDeleted", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := c.Get(ctx, fmt.Sprintf("/api/config/namespaces/%s/healthchecks/%s", namespace, testName), nil)
		if err != nil {
			t.Fatalf("Verify deletion failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404 after deletion, got %d", resp.StatusCode)
		} else {
			t.Log("Verified healthcheck was deleted")
		}
	})
}

// TestResource_ListCloudCredentials tests listing cloud credentials
func TestResource_ListCloudCredentials(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/system/cloud_credentialss", nil)
	if err != nil {
		t.Fatalf("List cloud credentials failed: %v", err)
	}

	t.Logf("Cloud Credentials - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d cloud credentials", len(items))
			}
		}
	}
}

// TestResource_ListSites tests listing sites
func TestResource_ListSites(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/system/sites", nil)
	if err != nil {
		t.Fatalf("List sites failed: %v", err)
	}

	t.Logf("Sites - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d sites", len(items))
				for i, item := range items {
					if i >= 5 {
						t.Logf("... and %d more", len(items)-5)
						break
					}
					if site, ok := item.(map[string]interface{}); ok {
						if meta, ok := site["metadata"].(map[string]interface{}); ok {
							t.Logf("  - %s", meta["name"])
						}
					}
				}
			}
		}
	}
}

// TestResource_VirtualHosts tests listing virtual hosts
func TestResource_VirtualHosts(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/config/namespaces/shared/virtual_hosts", nil)
	if err != nil {
		t.Fatalf("List virtual hosts failed: %v", err)
	}

	t.Logf("Virtual Hosts - status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d virtual hosts", len(items))
			}
		}
	}
}

// TestResource_OutputFormats tests different output formats work correctly
func TestResource_OutputFormats(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/web/namespaces/system", nil)
	if err != nil {
		t.Fatalf("Get namespace failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	// Test JSON output
	t.Run("JSON", func(t *testing.T) {
		output, err := client.FormatOutput(resp.Body, "json")
		if err != nil {
			t.Errorf("JSON format failed: %v", err)
		}
		if output == "" {
			t.Error("JSON output is empty")
		}
	})

	// Test YAML output
	t.Run("YAML", func(t *testing.T) {
		output, err := client.FormatOutput(resp.Body, "yaml")
		if err != nil {
			t.Errorf("YAML format failed: %v", err)
		}
		if output == "" {
			t.Error("YAML output is empty")
		}
	})
}
