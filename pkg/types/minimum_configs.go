// Package types provides minimum configuration definitions for common resources.
package types

// MinimumConfigurations provides copy-paste ready minimal configurations
// for common F5 XC resources. These are designed for AI assistants to
// generate deterministic, working configurations.
var MinimumConfigurations = map[string]*MinimumConfigSpec{
	"http_loadbalancer": {
		Description: "Minimum viable HTTP load balancer with public internet access",
		RequiredFields: []string{
			"metadata.name",
			"metadata.namespace",
			"spec.domains",
			"spec.http OR spec.https OR spec.https_auto_cert",
			"spec.advertise_on_public_default_vip",
		},
		ExampleYAML: `metadata:
  name: example-http-lb
  namespace: shared
spec:
  domains:
    - example.example.com
  http:
    port: 80
  advertise_on_public_default_vip: {}`,
		ExampleCommand: "xcsh cdn create http_loadbalancer -n shared -i http-lb.yaml",
		Domain:         "cdn",
	},
	"origin_pool": {
		Description: "Minimum viable origin pool with a single backend server",
		RequiredFields: []string{
			"metadata.name",
			"metadata.namespace",
			"spec.origin_servers",
			"spec.port",
		},
		ExampleYAML: `metadata:
  name: example-origin-pool
  namespace: shared
spec:
  origin_servers:
    - public_ip:
        ip: 192.168.1.100
  port: 8080
  loadbalancer_algorithm: ROUND_ROBIN`,
		ExampleCommand: "xcsh cdn create origin_pool -n shared -i origin-pool.yaml",
		Domain:         "cdn",
	},
	"healthcheck": {
		Description: "Minimum viable HTTP health check",
		RequiredFields: []string{
			"metadata.name",
			"metadata.namespace",
			"spec.http_health_check OR spec.tcp_health_check",
		},
		ExampleYAML: `metadata:
  name: example-healthcheck
  namespace: shared
spec:
  http_health_check:
    path: /health
  interval: 30
  timeout: 10`,
		ExampleCommand: "xcsh cdn create healthcheck -n shared -i healthcheck.yaml",
		Domain:         "cdn",
	},
	"tcp_loadbalancer": {
		Description: "Minimum viable TCP load balancer",
		RequiredFields: []string{
			"metadata.name",
			"metadata.namespace",
			"spec.listen_port",
			"spec.origin_pools",
		},
		ExampleYAML: `metadata:
  name: example-tcp-lb
  namespace: shared
spec:
  listen_port: 3306
  origin_pools:
    - pool:
        name: example-origin-pool
        namespace: shared`,
		ExampleCommand: "xcsh virtual create tcp_loadbalancer -n shared -i tcp-lb.yaml",
		Domain:         "virtual",
	},
	"app_firewall": {
		Description: "Minimum viable Web Application Firewall policy",
		RequiredFields: []string{
			"metadata.name",
			"metadata.namespace",
		},
		ExampleYAML: `metadata:
  name: example-waf
  namespace: shared
spec:
  allow_all_response_codes: true
  default_anonymization: true
  use_default_blocking_page: true`,
		ExampleCommand: "xcsh waf create app_firewall -n shared -i waf.yaml",
		Domain:         "waf",
	},
}

// GetMinimumConfiguration returns the minimum configuration for a resource type.
// Returns nil if no minimum configuration is defined for the resource.
func GetMinimumConfiguration(resourceName string) *MinimumConfigSpec {
	if config, ok := MinimumConfigurations[resourceName]; ok {
		return config
	}
	return nil
}

// init adds minimum configurations to the resource schemas
func init() {
	for resourceName, minConfig := range MinimumConfigurations {
		if schema, ok := ResourceSchemas[resourceName]; ok {
			schema.MinimumConfiguration = minConfig
			ResourceSchemas[resourceName] = schema
		}
	}
}
