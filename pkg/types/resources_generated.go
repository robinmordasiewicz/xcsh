package types

// Code generated from OpenAPI specifications. DO NOT EDIT.
// This file contains 269 resource types parsed from F5 XC API specs

func init() {
	registerGeneratedResources()
}

func registerGeneratedResources() {
	Register(&ResourceType{
		Name:              "address_allocator",
		CLIName:           "address-allocator",
		Description:       "Address Allocator object is used to allocate an address or a subnet from a given",
		APIPath:           "/api/config/namespaces/{namespace}/address_allocators",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "advertise_policy",
		CLIName:           "advertise-policy",
		Description:       "advertise_policy object controls how and where a service represented by a given ",
		APIPath:           "/api/config/namespaces/{namespace}/advertise_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "alert_policy",
		CLIName:           "alert-policy",
		Description:       "Alert Policy is used to specify a set of routes to match the incoming alert and ",
		APIPath:           "/api/config/namespaces/{namespace}/alert_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "alert_receiver",
		CLIName:           "alert-receiver",
		Description:       "Alert Receiver is used to specify a receiver (slack, pagerDuty, etc.,) to send t",
		APIPath:           "/api/config/namespaces/{namespace}/alert_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_credential",
		CLIName:           "api-credential",
		Description:       "F5XC supports 2 variation of credentials - ",
		APIPath:           "/api/web/namespaces/{namespace}/api_credentials",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_definition",
		CLIName:           "api-definition",
		Description:       "The api_definition construct provides a mechanism to create api_groups based on ",
		APIPath:           "/api/config/namespaces/{namespace}/api_definitions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_group",
		CLIName:           "api-group",
		Description:       "The api_group construct provides a mechanism to classify the universal set of re",
		APIPath:           "/api/web/namespaces/{namespace}/api_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_group_element",
		CLIName:           "api-group-element",
		Description:       "A api_group_element object consists of an unordered list of HTTP methods and a p",
		APIPath:           "/api/web/namespaces/{namespace}/api_group_elements",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_sec_api_crawler",
		CLIName:           "api-sec-api-crawler",
		Description:       "This is the api crawler type",
		APIPath:           "/api/config/namespaces/{namespace}/api_crawlers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_sec_api_discovery",
		CLIName:           "api-sec-api-discovery",
		Description:       "The api_discovery contains settings for API discovery",
		APIPath:           "/api/config/namespaces/{namespace}/api_discoverys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_sec_api_testing",
		CLIName:           "api-sec-api-testing",
		Description:       "This is the api testing type",
		APIPath:           "/api/config/namespaces/{namespace}/api_testings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_sec_code_base_integration",
		CLIName:           "api-sec-code-base-integration",
		Description:       "Code base integration",
		APIPath:           "/api/config/namespaces/{namespace}/code_base_integrations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "app_api_group",
		CLIName:           "app-api-group",
		Description:       "The app_api_group construct provides a mechanism to classify the universal set o",
		APIPath:           "/api/config/namespaces/{namespace}/app_api_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "app_firewall",
		CLIName:           "app-firewall",
		Description:       "WAF Configuration",
		APIPath:           "/api/config/namespaces/{namespace}/app_firewalls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "app_setting",
		CLIName:           "app-setting",
		Description:       "\"App Setting\" controls advanced monitoring of applications defined by \"App type\"",
		APIPath:           "/api/config/namespaces/{namespace}/app_settings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "app_type",
		CLIName:           "app-type",
		Description:       "App Type object defines a application profile type from an advanced monitoring/s",
		APIPath:           "/api/config/namespaces/{namespace}/app_types",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "authentication",
		CLIName:           "authentication",
		Description:       "Authentication Object contains authentication specific config . This includes",
		APIPath:           "/api/config/namespaces/{namespace}/authentications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "aws_tgw_site",
		CLIName:           "aws-tgw-site",
		Description:       "AWS TGW site view defines a required parameters that can be used in CRUD, to cre",
		APIPath:           "/api/config/namespaces/{namespace}/aws_tgw_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "aws_vpc_site",
		CLIName:           "aws-vpc-site",
		Description:       "AWS VPC site view defines a required parameters that can be used in CRUD, to cre",
		APIPath:           "/api/config/namespaces/{namespace}/aws_vpc_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "azure_vnet_site",
		CLIName:           "azure-vnet-site",
		Description:       "Azure VNet site view defines a required parameters that can be used in CRUD, to ",
		APIPath:           "/api/config/namespaces/{namespace}/azure_vnet_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bgp",
		CLIName:           "bgp",
		Description:       "BGP object represents configuration of bgp protocol on given network interface o",
		APIPath:           "/api/config/namespaces/{namespace}/bgps",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bgp_asn_set",
		CLIName:           "bgp-asn-set",
		Description:       "An unordered set of RFC 6793 defined 4-byte AS numbers that can be used to creat",
		APIPath:           "/api/config/namespaces/{namespace}/bgp_asn_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bgp_routing_policy",
		CLIName:           "bgp-routing-policy",
		Description:       "BGP Routing Policy is a list of rules, which contains match criteria and",
		APIPath:           "/api/config/namespaces/{namespace}/bgp_routing_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bigcne_data_group",
		CLIName:           "bigcne-data-group",
		Description:       "A data group is a group of related items - IP addresses/subnets, strings, or int",
		APIPath:           "/api/config/namespaces/{namespace}/data_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bigcne_irule",
		CLIName:           "bigcne-irule",
		Description:       "iRule object defines the iRule that can be used in CRUD, to create and manage iR",
		APIPath:           "/api/config/namespaces/{namespace}/irules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bigip_apm",
		CLIName:           "bigip-apm",
		Description:       "BIG-IP APM Service handles the life-cycle management of BIG-IP appliances.",
		APIPath:           "/api/config/namespaces/{namespace}/apms",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bigip_virtual_server",
		CLIName:           "bigip-virtual-server",
		Description:       "BIG-IP virtual server view repesents the internal virtual host corresponding to ",
		APIPath:           "/api/config/namespaces/{namespace}/bigip_virtual_servers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "bot_defense_app_infrastructure",
		CLIName:           "bot-defense-app-infrastructure",
		Description:       "Bot Defense App Infrastructure is the main configuration for a Bot Defense Advan",
		APIPath:           "/api/config/namespaces/{namespace}/bot_defense_app_infrastructures",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cdn_cache_rule",
		CLIName:           "cdn-cache-rule",
		Description:       "CDN cache rule view defines a required parameters that can be used in CRUD, to c",
		APIPath:           "/api/config/namespaces/{namespace}/cdn_cache_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cdn_loadbalancer",
		CLIName:           "cdn-loadbalancer",
		Description:       "CDN Loadbalancer view defines a required parameters that can be used in CRUD, to",
		APIPath:           "/api/config/namespaces/{namespace}/cdn_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "certificate",
		CLIName:           "certificate",
		Description:       "Certificate represents a client or server certificate.",
		APIPath:           "/api/config/namespaces/{namespace}/certificates",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "certificate_chain",
		CLIName:           "certificate-chain",
		Description:       "Certificate chain is list of certificates used to establish chain of trust from ",
		APIPath:           "/api/config/namespaces/{namespace}/certificate_chains",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "certified_hardware",
		CLIName:           "certified-hardware",
		Description:       "Certified Hardware object represents physical hardware or cloud instance type th",
		APIPath:           "/api/config/namespaces/{namespace}/certified_hardwares",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_connect",
		CLIName:           "cloud-connect",
		Description:       "Cloud Connect Represents connection endpoint for cloud.",
		APIPath:           "/api/config/namespaces/system/edge_credentials",
		SupportsNamespace: false,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_credentials",
		CLIName:           "cloud-credentials",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_credentialss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_elastic_ip",
		CLIName:           "cloud-elastic-ip",
		Description:       "Cloud Elastic IP object represents a cloud elastic IP address that are created f",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_elastic_ips",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_link",
		CLIName:           "cloud-link",
		Description:       "CloudLink is used to establish private connectivity from customer network to Clo",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_links",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_region",
		CLIName:           "cloud-region",
		Description:       "Cloud Region contains tenant specific configuration",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_regions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cluster",
		CLIName:           "cluster",
		Description:       "cluster object represent common set of endpoints (providers of service) that can",
		APIPath:           "/api/config/namespaces/{namespace}/clusters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cminstance",
		CLIName:           "cminstance",
		Description:       "cminsatnce object can be used to enable connectivity between ce site and bigip c",
		APIPath:           "/api/config/namespaces/{namespace}/cminstances",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "contact",
		CLIName:           "contact",
		Description:       "Customer or tenant contact details.",
		APIPath:           "/api/web/namespaces/{namespace}/contacts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "container_registry",
		CLIName:           "container-registry",
		Description:       "Container registry is the container or docker registry information",
		APIPath:           "/api/config/namespaces/{namespace}/container_registrys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "crl",
		CLIName:           "crl",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/crls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "customer_support",
		CLIName:           "customer-support",
		Description:       "Handles creation and listing of support issues (by tenant and user).",
		APIPath:           "/api/web/namespaces/{namespace}/customer_supports",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "data_privacy_lma_region",
		CLIName:           "data-privacy-lma-region",
		Description:       "LMA Region.",
		APIPath:           "/api/config/namespaces/{namespace}/lma_regions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "data_type",
		CLIName:           "data-type",
		Description:       "A data_type is defined by a set of rules. these rules include the patterns for w",
		APIPath:           "/api/config/namespaces/{namespace}/data_types",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dc_cluster_group",
		CLIName:           "dc-cluster-group",
		Description:       "A DC Cluster Group represents a collection of sites that",
		APIPath:           "/api/config/namespaces/{namespace}/dc_cluster_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "discovery",
		CLIName:           "discovery",
		Description:       "Service discovery in F5XC performs following",
		APIPath:           "/api/config/namespaces/{namespace}/discoverys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_compliance_checks",
		CLIName:           "dns-compliance-checks",
		Description:       "DNS Compliance Checks view defines the required parameters that can be used in C",
		APIPath:           "/api/config/namespaces/{namespace}/dns_compliance_checkss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_domain",
		CLIName:           "dns-domain",
		Description:       "DNS Domain object is used for delegating DNS sub domain to volterra. It can also",
		APIPath:           "/api/config/namespaces/{namespace}/dns_domains",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "endpoint",
		CLIName:           "endpoint",
		Description:       "Endpoint object represent the actual endpoint that provides the service (Origin ",
		APIPath:           "/api/config/namespaces/{namespace}/endpoints",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "enhanced_firewall_policy",
		CLIName:           "enhanced-firewall-policy",
		Description:       "Enhanced Firewall Policy defined firewall rules applied in the site",
		APIPath:           "/api/config/namespaces/{namespace}/enhanced_firewall_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "external_connector",
		CLIName:           "external-connector",
		Description:       "External Connector configuration mainly includes the following:",
		APIPath:           "/api/config/namespaces/{namespace}/external_connectors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "fast_acl",
		CLIName:           "fast-acl",
		Description:       "Fast ACL provides destination and specifies rules to protect the site from denia",
		APIPath:           "/api/config/namespaces/{namespace}/fast_acls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "fast_acl_rule",
		CLIName:           "fast-acl-rule",
		Description:       " Fast ACL rule",
		APIPath:           "/api/config/namespaces/{namespace}/fast_acl_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "filter_set",
		CLIName:           "filter-set",
		Description:       "Filter Set is a set of saved filtering criteria used in the Console. This allows",
		APIPath:           "/api/config/namespaces/{namespace}/filter_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "fleet",
		CLIName:           "fleet",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/fleets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "flow_anomaly",
		CLIName:           "flow-anomaly",
		Description:       "Flow Anomaly.",
		APIPath:           "/api/config/namespaces/{namespace}/flow_anomalys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "forward_proxy_policy",
		CLIName:           "forward-proxy-policy",
		Description:       "Forward Proxy policy defines access control rules for connections going via forw",
		APIPath:           "/api/config/namespaces/{namespace}/forward_proxy_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "forwarding_class",
		CLIName:           "forwarding-class",
		Description:       "In Policy Based Routing(forwarding) (PBR) PBR policy can select Forwarding Class",
		APIPath:           "/api/config/namespaces/{namespace}/forwarding_classs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "gcp_vpc_site",
		CLIName:           "gcp-vpc-site",
		Description:       "GCP VPC site view defines a required parameters that can be used in CRUD, to cre",
		APIPath:           "/api/config/namespaces/{namespace}/gcp_vpc_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "global_log_receiver",
		CLIName:           "global-log-receiver",
		Description:       "Global Log Receiver is used to specify a receiver (s3 bucket, etc.) for periodic",
		APIPath:           "/api/config/namespaces/{namespace}/global_log_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "healthcheck",
		CLIName:           "healthcheck",
		Description:       "Health check configuration for a given cluster.",
		APIPath:           "/api/config/namespaces/{namespace}/healthchecks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "http_loadbalancer",
		CLIName:           "http-loadbalancer",
		Description:       "HTTP Load Balancer view defines a required parameters that can be used in CRUD, ",
		APIPath:           "/api/config/namespaces/{namespace}/http_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ike_phase1_profile",
		CLIName:           "ike-phase1-profile",
		Description:       "IKE Phase1 profile mainly includes the following",
		APIPath:           "/api/config/namespaces/{namespace}/ike_phase1_profiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ike_phase2_profile",
		CLIName:           "ike-phase2-profile",
		Description:       "IKE Phase2 profile mainly includes the following",
		APIPath:           "/api/config/namespaces/{namespace}/ike_phase2_profiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ike1",
		CLIName:           "ike1",
		Description:       "IKE Phase1 profile mainly includes the following",
		APIPath:           "/api/config/namespaces/{namespace}/ike1s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ike2",
		CLIName:           "ike2",
		Description:       "IKE Phase2 profile mainly includes the following",
		APIPath:           "/api/config/namespaces/{namespace}/ike2s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "implicit_label",
		CLIName:           "implicit-label",
		Description:       "Implicit labels are attached to objects implicitly by the system. Users are not ",
		APIPath:           "/api/config/namespaces/system/implicit_labels",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "ip_prefix_set",
		CLIName:           "ip-prefix-set",
		Description:       "An ip prefix set contains an unordered list of IP prefixes. It can can be used t",
		APIPath:           "/api/config/namespaces/{namespace}/ip_prefix_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "k8s_cluster",
		CLIName:           "k8s-cluster",
		Description:       "K8s cluster represents the real physical K8s cluster on the site. It can be used",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_clusters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "k8s_cluster_role",
		CLIName:           "k8s-cluster-role",
		Description:       "K8s Cluster Role is K8s ClusterRole object, which represents set of permissions ",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_cluster_roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "k8s_cluster_role_binding",
		CLIName:           "k8s-cluster-role-binding",
		Description:       "Cluster role binding allows administrator to assign cluster wide cluster role to",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_cluster_role_bindings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "k8s_pod_security_admission",
		CLIName:           "k8s-pod-security-admission",
		Description:       "Pod security admission allows users to enforce Pod Security Standards",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_pod_security_admissions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "k8s_pod_security_policy",
		CLIName:           "k8s-pod-security-policy",
		Description:       "Pod Security Policies enable fine-grained authorization of pod creation and upda",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_pod_security_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "known_label",
		CLIName:           "known-label",
		Description:       "Known labels serves two purposes",
		APIPath:           "/api/config/namespaces/{namespace}/known_labels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "known_label_key",
		CLIName:           "known-label-key",
		Description:       "Known label key serves two purposes",
		APIPath:           "/api/config/namespaces/{namespace}/known_label_keys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "log_receiver",
		CLIName:           "log-receiver",
		Description:       "Log Receiver is used to specify a receiver (syslog, splunk, datadog etc.,) to se",
		APIPath:           "/api/config/namespaces/{namespace}/log_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "malicious_user_mitigation",
		CLIName:           "malicious-user-mitigation",
		Description:       "A malicious_user_mitigation object consists of settings that specify the actions",
		APIPath:           "/api/config/namespaces/{namespace}/malicious_user_mitigations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "namespace",
		CLIName:           "namespace",
		Description:       "namespace creates logical independent workspace within a tenant. Within a namesp",
		APIPath:           "/api/web/namespaces",
		SupportsNamespace: false,
		Operations:        AllOperations(),
		DeleteConfig: &DeleteConfig{
			PathSuffix:  "/cascade_delete",
			Method:      "POST",
			IncludeBody: true,
		},
	})

	Register(&ResourceType{
		Name:              "namespace_role",
		CLIName:           "namespace-role",
		Description:       "Namespace role defines a user's role in a namespace.",
		APIPath:           "/api/web/namespaces/{namespace}/namespace_roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nat_policy",
		CLIName:           "nat-policy",
		Description:       "NAT Policy object represents the configuration of Network Address Translation pa",
		APIPath:           "/api/config/namespaces/{namespace}/nat_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_connector",
		CLIName:           "network-connector",
		Description:       "Network Connector is used to create connection between two virtual networks on a",
		APIPath:           "/api/config/namespaces/{namespace}/network_connectors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_firewall",
		CLIName:           "network-firewall",
		Description:       "Network Firewall is applicable when referred to by a Fleet. The Network Firewall",
		APIPath:           "/api/config/namespaces/{namespace}/network_firewalls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_interface",
		CLIName:           "network-interface",
		Description:       "Network Interface object represents the configuration of a network device in a f",
		APIPath:           "/api/config/namespaces/{namespace}/network_interfaces",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_policy",
		CLIName:           "network-policy",
		Description:       "Network Policy is applied to all IP packets to and from a given endpoint (called",
		APIPath:           "/api/config/namespaces/{namespace}/network_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_policy_rule",
		CLIName:           "network-policy-rule",
		Description:       "Network Policy Rule is applied to given remote endpoints to and from a given loc",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_policy_set",
		CLIName:           "network-policy-set",
		Description:       "Network policy set implements L3/L4 stateful firewall.",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_policy_view",
		CLIName:           "network-policy-view",
		Description:       "Network policy site view defines a required parameters that can be used in CRUD,",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_views",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nfv_service",
		CLIName:           "nfv-service",
		Description:       "NFV Service manages the lifecycle  of the NFV appliance, which includes the func",
		APIPath:           "/api/config/namespaces/{namespace}/nfv_services",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nginx_one_nginx_csg",
		CLIName:           "nginx-one-nginx-csg",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/nginx_csgs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nginx_one_nginx_instance",
		CLIName:           "nginx-one-nginx-instance",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/nginx_instances",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nginx_one_nginx_server",
		CLIName:           "nginx-one-nginx-server",
		Description:       "",
		APIPath:           "/api/config/namespaces/{namespace}/nginx_dataplane_servers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nginx_one_nginx_service_discovery",
		CLIName:           "nginx-one-nginx-service-discovery",
		Description:       "NGINX Service discovery in F5XC",
		APIPath:           "/api/config/namespaces/{namespace}/nginx_service_discoverys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "origin_pool",
		CLIName:           "origin-pool",
		Description:       "Origin pool is a view to create cluster and endpoints that can be used in HTTP l",
		APIPath:           "/api/config/namespaces/{namespace}/origin_pools",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "pbac_addon_service",
		CLIName:           "pbac-addon-service",
		Description:       "Basic unit of logical representation of a F5XC service.",
		APIPath:           "/api/web/namespaces/{namespace}/addon_services",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "pbac_addon_subscription",
		CLIName:           "pbac-addon-subscription",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/addon_subscriptions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "pbac_navigation_tile",
		CLIName:           "pbac-navigation-tile",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/navigation_tiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "pbac_plan",
		CLIName:           "pbac-plan",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/plans",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "policer",
		CLIName:           "policer",
		Description:       "* Policer objects enforces traffic rate limits",
		APIPath:           "/api/config/namespaces/{namespace}/policers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "policy_based_routing",
		CLIName:           "policy-based-routing",
		Description:       "Policy based routing is used to control how different classes of traffic is forw",
		APIPath:           "/api/config/namespaces/{namespace}/policy_based_routings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "protocol_inspection",
		CLIName:           "protocol-inspection",
		Description:       "Protocol Inspection view defines the required parameters that can be used in CRU",
		APIPath:           "/api/config/namespaces/{namespace}/protocol_inspections",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "protocol_policer",
		CLIName:           "protocol-policer",
		Description:       "Protocol policer has set or network protocol fields and flags to be match on",
		APIPath:           "/api/config/namespaces/{namespace}/protocol_policers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "proxy",
		CLIName:           "proxy",
		Description:       "Proxy view defines a required parameters that can be used in CRUD, to create and",
		APIPath:           "/api/config/namespaces/{namespace}/proxys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "public_ip",
		CLIName:           "public-ip",
		Description:       "public_ip object represents a public IP address that is available on a set of vi",
		APIPath:           "/api/config/namespaces/{namespace}/public_ips",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "quota",
		CLIName:           "quota",
		Description:       "Quota object is used to configure the limits on how many of a resource type can ",
		APIPath:           "/api/web/namespaces/{namespace}/quotas",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "rate_limiter",
		CLIName:           "rate-limiter",
		Description:       "A rate_limiter specifies a list of rate limit unit periods and the corresponding",
		APIPath:           "/api/config/namespaces/{namespace}/rate_limiters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "rate_limiter_policy",
		CLIName:           "rate-limiter-policy",
		Description:       "Rate limiter policy defines parameters that can be used for fine-grained control",
		APIPath:           "/api/config/namespaces/{namespace}/rate_limiter_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "rbac_policy",
		CLIName:           "rbac-policy",
		Description:       "A rbac_policy object consists of list of rbac policy rules that when assigned to",
		APIPath:           "/api/web/namespaces/{namespace}/rbac_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "role",
		CLIName:           "role",
		Description:       "Defines the role the user has in a namespace. There are two kinds of roles:",
		APIPath:           "/api/web/namespaces/{namespace}/roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "route",
		CLIName:           "route",
		Description:       "route object is used to configuring L7 routing decision. route is made of three ",
		APIPath:           "/api/config/namespaces/{namespace}/routes",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_management_access",
		CLIName:           "secret-management-access",
		Description:       "secret_management_access object is used to define configuration on how to connec",
		APIPath:           "/api/config/namespaces/{namespace}/secret_management_accesss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "securemesh_site",
		CLIName:           "securemesh-site",
		Description:       "Secure Mesh site defines a required parameters that can be used in CRUD, to crea",
		APIPath:           "/api/config/namespaces/{namespace}/securemesh_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "securemesh_site_v2",
		CLIName:           "securemesh-site-v2",
		Description:       "Secure Mesh site defines a required parameters that can be used in CRUD, to crea",
		APIPath:           "/api/config/namespaces/{namespace}/securemesh_site_v2s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "segment",
		CLIName:           "segment",
		Description:       "Network Segment.",
		APIPath:           "/api/config/namespaces/{namespace}/segments",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "segment_connection",
		CLIName:           "segment-connection",
		Description:       "Configure a Segment Connector to allow network traffic between Segments",
		APIPath:           "/api/config/namespaces/{namespace}/segment_connections",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "sensitive_data_policy",
		CLIName:           "sensitive-data-policy",
		Description:       "The sensitive_data_policy is a policy defined by the user to discover the releva",
		APIPath:           "/api/config/namespaces/{namespace}/sensitive_data_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "service_policy",
		CLIName:           "service-policy",
		Description:       "A service_policy object consists of an unordered list of predicates and a list o",
		APIPath:           "/api/config/namespaces/{namespace}/service_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "service_policy_rule",
		CLIName:           "service-policy-rule",
		Description:       "A service_policy_rule object consists of an unordered list of predicates and an ",
		APIPath:           "/api/config/namespaces/{namespace}/service_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "service_policy_set",
		CLIName:           "service-policy-set",
		Description:       "A service_policy_set object consists of an ordered list of references to service",
		APIPath:           "/api/config/namespaces/{namespace}/service_policy_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_instance",
		CLIName:           "shape-bot-defense-instance",
		Description:       "Shape Bot Defense Instance is the main configuration for a Shape Integration.",
		APIPath:           "/api/config/namespaces/{namespace}/shape_bot_defense_instances",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "site",
		CLIName:           "site",
		Description:       "Site represent physical/cloud cluster of volterra processing elements. There are",
		APIPath:           "/api/config/namespaces/{namespace}/sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "site_mesh_group",
		CLIName:           "site-mesh-group",
		Description:       "Site mesh group is a configuration tool to provide Site to Site",
		APIPath:           "/api/config/namespaces/{namespace}/site_mesh_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "srv6_network_slice",
		CLIName:           "srv6-network-slice",
		Description:       "An srv6_network_slice represents a network slice in an operator network that use",
		APIPath:           "/api/config/namespaces/{namespace}/srv6_network_slices",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "subnet",
		CLIName:           "subnet",
		Description:       "Subnet object is used to support VMs/pods with multiple interfaces,",
		APIPath:           "/api/config/namespaces/{namespace}/subnets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tcp_loadbalancer",
		CLIName:           "tcp-loadbalancer",
		Description:       "TCP load balancer view defines a required parameters that can be used in CRUD, t",
		APIPath:           "/api/config/namespaces/{namespace}/tcp_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_configuration",
		CLIName:           "tenant-configuration",
		Description:       "Tenant configuration consists of three main parts:",
		APIPath:           "/api/config/namespaces/{namespace}/tenant_configurations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_allowed_tenant",
		CLIName:           "tenant-management-allowed-tenant",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/allowed_tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_child_tenant",
		CLIName:           "tenant-management-child-tenant",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/child_tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_child_tenant_manager",
		CLIName:           "tenant-management-child-tenant-manager",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/child_tenant_managers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_managed_tenant",
		CLIName:           "tenant-management-managed-tenant",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/managed_tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_tenant_profile",
		CLIName:           "tenant-management-tenant-profile",
		Description:       "",
		APIPath:           "/api/web/namespaces/{namespace}/tenant_profiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "third_party_application",
		CLIName:           "third-party-application",
		Description:       "View will create following child objects.",
		APIPath:           "/api/config/namespaces/{namespace}/third_party_applications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ticket_management_ticket_tracking_system",
		CLIName:           "ticket-management-ticket-tracking-system",
		Description:       "Public Custom APIs for Ticket Tracking System related operations",
		APIPath:           "/api/web/namespaces/{namespace}/ticket_tracking_systems",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "trusted_ca_list",
		CLIName:           "trusted-ca-list",
		Description:       "A Root CA Certificate represents list of trusted root CAs",
		APIPath:           "/api/config/namespaces/{namespace}/trusted_ca_lists",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tunnel",
		CLIName:           "tunnel",
		Description:       "Tunnel configuration allows user to specify parameters for configuring static tu",
		APIPath:           "/api/config/namespaces/{namespace}/tunnels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "udp_loadbalancer",
		CLIName:           "udp-loadbalancer",
		Description:       "UDP load balancer view defines a required parameters that can be used in CRUD, t",
		APIPath:           "/api/config/namespaces/{namespace}/udp_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "ui_static_component",
		CLIName:           "ui-static-component",
		Description:       "stores information about the UI Components in key-value pair",
		APIPath:           "/api/web/namespaces/{namespace}/static_components",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "usage",
		CLIName:           "usage",
		Description:       "Resource usage and pricing custom APIs",
		APIPath:           "/api/web/namespaces/{namespace}/hourly_usage_details",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "usb_policy",
		CLIName:           "usb-policy",
		Description:       "USB policy is used to specify list of USB devices allowed to be attached to node",
		APIPath:           "/api/config/namespaces/{namespace}/usb_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "user_group",
		CLIName:           "user-group",
		Description:       "Represents group for a given tenant",
		APIPath:           "/api/web/namespaces/{namespace}/user_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "user_identification",
		CLIName:           "user-identification",
		Description:       "A user_identification object consists of an ordered list of rules. The rules are",
		APIPath:           "/api/config/namespaces/{namespace}/user_identifications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "virtual_host",
		CLIName:           "virtual-host",
		Description:       "Virtual host is main anchor configuration for a proxy. Primary application for v",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_hosts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "virtual_k8s",
		CLIName:           "virtual-k8s",
		Description:       "Virtual K8s object exposes a Kubernetes API endpoint in the namespace that opera",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_k8ss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "virtual_network",
		CLIName:           "virtual-network",
		Description:       "Virtual network is an isolated L3 network. A virtual network can contain",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_networks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "virtual_site",
		CLIName:           "virtual-site",
		Description:       "Virtual site object is mechanism to create arbitrary set of sites",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "voltstack_site",
		CLIName:           "voltstack-site",
		Description:       "App Stack site defines a required parameters that can be used in CRUD, to create",
		APIPath:           "/api/config/namespaces/{namespace}/voltstack_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "waf_exclusion_policy",
		CLIName:           "waf-exclusion-policy",
		Description:       "WAF Exclusion Policy record",
		APIPath:           "/api/config/namespaces/{namespace}/waf_exclusion_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "workload",
		CLIName:           "workload",
		Description:       "Workload is used to configure and deploy a workload in Virtual Kubernetes. A wor",
		APIPath:           "/api/config/namespaces/{namespace}/workloads",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "workload_flavor",
		CLIName:           "workload-flavor",
		Description:       "Workload flavor is used to assign CPU, memory, and storage resources to workload",
		APIPath:           "/api/config/namespaces/{namespace}/workload_flavors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// ============================================================
	// Additional API Resources (119 missing specs)
	// ============================================================

	// AI & Machine Learning
	Register(&ResourceType{
		Name:              "ai_assistant",
		CLIName:           "ai-assistant",
		Description:       "AI Assistant for query and feedback operations",
		APIPath:           "/api/gen-ai/namespaces/{namespace}/query",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "ai_data_bfdp",
		CLIName:           "ai-data-bfdp",
		Description:       "AI Data BFDP operations",
		APIPath:           "/api/ai_data/bfdp",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "ai_data_bfdp_subscription",
		CLIName:           "ai-data-bfdp-subscription",
		Description:       "AI Data BFDP subscription management",
		APIPath:           "/api/ai_data/bfdp/subscription",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	// Alerts & Monitoring
	Register(&ResourceType{
		Name:              "alert",
		CLIName:           "alert",
		Description:       "Alert management and history",
		APIPath:           "/api/data/namespaces/{namespace}/alerts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_sec_rule_suggestion",
		CLIName:           "api-sec-rule-suggestion",
		Description:       "API Security rule suggestions",
		APIPath:           "/api/config/namespaces/{namespace}/rule_suggestions",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "maintenance_status",
		CLIName:           "maintenance-status",
		Description:       "Maintenance status information",
		APIPath:           "/api/config/namespaces/system/maintenance_status",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "upgrade_status",
		CLIName:           "upgrade-status",
		Description:       "Upgrade status information",
		APIPath:           "/api/config/namespaces/{namespace}/upgrade_statuss",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	// Application Security
	Register(&ResourceType{
		Name:              "app_security",
		CLIName:           "app-security",
		Description:       "Application security suggestions and configuration",
		APIPath:           "/api/config/namespaces/{namespace}/app_securitys",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "bigip_irule",
		CLIName:           "bigip-irule",
		Description:       "BIG-IP iRule management",
		APIPath:           "/api/bigipconnector/namespaces/{namespace}/bigip_irules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// Billing & Commerce
	Register(&ResourceType{
		Name:              "billing_payment_method",
		CLIName:           "billing-payment-method",
		Description:       "Billing payment method management",
		APIPath:           "/api/web/namespaces/{namespace}/payment_methods",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "billing_plan_transition",
		CLIName:           "billing-plan-transition",
		Description:       "Billing plan transition management",
		APIPath:           "/api/web/namespaces/{namespace}/plan_transitions",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "marketplace_aws_account",
		CLIName:           "marketplace-aws-account",
		Description:       "AWS Marketplace account management",
		APIPath:           "/api/marketplace/aws/accounts",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "marketplace_xc_saas",
		CLIName:           "marketplace-xc-saas",
		Description:       "XC SaaS marketplace management",
		APIPath:           "/api/marketplace/xc_saas",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	// Data Management & Privacy
	Register(&ResourceType{
		Name:              "data_privacy_geo_config",
		CLIName:           "data-privacy-geo-config",
		Description:       "Data privacy geo configuration",
		APIPath:           "/api/config/namespaces/{namespace}/geo_configs",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_management",
		CLIName:           "secret-management",
		Description:       "Secret management operations",
		APIPath:           "/api/config/namespaces/{namespace}/secret_managements",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_policy",
		CLIName:           "secret-policy",
		Description:       "Secret policy configuration",
		APIPath:           "/api/config/namespaces/{namespace}/secret_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_policy_rule",
		CLIName:           "secret-policy-rule",
		Description:       "Secret policy rule configuration",
		APIPath:           "/api/config/namespaces/{namespace}/secret_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// DNS Services
	Register(&ResourceType{
		Name:              "dns_load_balancer",
		CLIName:           "dns-load-balancer",
		Description:       "DNS load balancer configuration",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_load_balancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_lb_health_check",
		CLIName:           "dns-lb-health-check",
		Description:       "DNS load balancer health check configuration",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_lb_health_checks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_lb_pool",
		CLIName:           "dns-lb-pool",
		Description:       "DNS load balancer pool configuration",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_lb_pools",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_zone",
		CLIName:           "dns-zone",
		Description:       "DNS zone configuration",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_zones",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_zone_rrset",
		CLIName:           "dns-zone-rrset",
		Description:       "DNS zone resource record set management",
		APIPath:           "/api/config/dns/dns_zones/{dns_zone_name}/rrsets",
		SupportsNamespace: false,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "dns_zone_subscription",
		CLIName:           "dns-zone-subscription",
		Description:       "DNS zone subscription management",
		APIPath:           "/api/config/dns/subscription",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	// Flow & Network Analytics
	Register(&ResourceType{
		Name:              "flow",
		CLIName:           "flow",
		Description:       "Network flow data queries",
		APIPath:           "/api/data/namespaces/{namespace}/flows",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "discovered_service",
		CLIName:           "discovered-service",
		Description:       "Discovered service information",
		APIPath:           "/api/config/namespaces/{namespace}/discovered_services",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	// Graph & Visualization
	Register(&ResourceType{
		Name:              "graph_connectivity",
		CLIName:           "graph-connectivity",
		Description:       "Connectivity graph analysis",
		APIPath:           "/api/graph/namespaces/{namespace}/connectivity",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "graph_l3l4",
		CLIName:           "graph-l3l4",
		Description:       "L3/L4 topology graph",
		APIPath:           "/api/graph/namespaces/{namespace}/l3l4",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "graph_service",
		CLIName:           "graph-service",
		Description:       "Service graph queries",
		APIPath:           "/api/graph/namespaces/{namespace}/services",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "graph_site",
		CLIName:           "graph-site",
		Description:       "Site topology graph",
		APIPath:           "/api/graph/namespaces/{namespace}/sites",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	// Geo-Location
	Register(&ResourceType{
		Name:              "geo_location_set",
		CLIName:           "geo-location-set",
		Description:       "Geo location set configuration",
		APIPath:           "/api/config/namespaces/{namespace}/geo_location_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "gia",
		CLIName:           "gia",
		Description:       "GIA operations",
		APIPath:           "/api/config/namespaces/system/gias",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
			Delete: true,
		},
	})

	// Infrastructure Protection
	Register(&ResourceType{
		Name:              "infraprotect",
		CLIName:           "infraprotect",
		Description:       "Infrastructure protection configuration",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotects",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_asn",
		CLIName:           "infraprotect-asn",
		Description:       "Infrastructure protection ASN configuration",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_asns",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_asn_prefix",
		CLIName:           "infraprotect-asn-prefix",
		Description:       "Infrastructure protection ASN prefix configuration",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_asn_prefixs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_deny_list_rule",
		CLIName:           "infraprotect-deny-list-rule",
		Description:       "Infrastructure protection deny list rule",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_deny_list_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_rule",
		CLIName:           "infraprotect-firewall-rule",
		Description:       "Infrastructure protection firewall rule",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_rule_group",
		CLIName:           "infraprotect-firewall-rule-group",
		Description:       "Infrastructure protection firewall rule group",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rule_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_ruleset",
		CLIName:           "infraprotect-firewall-ruleset",
		Description:       "Infrastructure protection firewall ruleset",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rulesets",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "infraprotect_information",
		CLIName:           "infraprotect-information",
		Description:       "Infrastructure protection information",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_informations",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_internet_prefix_advertisement",
		CLIName:           "infraprotect-internet-prefix-advertisement",
		Description:       "Infrastructure protection internet prefix advertisement",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_internet_prefix_advertisements",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_tunnel",
		CLIName:           "infraprotect-tunnel",
		Description:       "Infrastructure protection tunnel configuration",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_tunnels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// Logging & Events
	Register(&ResourceType{
		Name:              "log",
		CLIName:           "log",
		Description:       "Log data queries",
		APIPath:           "/api/data/namespaces/{namespace}/logs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// Management & Administration
	Register(&ResourceType{
		Name:              "module_management",
		CLIName:           "module-management",
		Description:       "Module management information",
		APIPath:           "/api/config/namespaces/{namespace}/module_managements",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "oidc_provider",
		CLIName:           "oidc-provider",
		Description:       "OIDC provider configuration",
		APIPath:           "/api/config/namespaces/{namespace}/oidc_providers",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "registration",
		CLIName:           "registration",
		Description:       "Registration management",
		APIPath:           "/api/config/namespaces/{namespace}/registrations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "scim",
		CLIName:           "scim",
		Description:       "SCIM user and group management",
		APIPath:           "/api/scim/v2",
		SupportsNamespace: false,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant",
		CLIName:           "tenant",
		Description:       "Tenant management",
		APIPath:           "/api/config/namespaces/{namespace}/tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management",
		CLIName:           "tenant-management",
		Description:       "Tenant management operations",
		APIPath:           "/api/config/namespaces/{namespace}/tenant_managements",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "token",
		CLIName:           "token",
		Description:       "Token management",
		APIPath:           "/api/config/namespaces/{namespace}/tokens",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tpm_api_key",
		CLIName:           "tpm-api-key",
		Description:       "TPM API key management",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_api_keys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "tpm_category",
		CLIName:           "tpm-category",
		Description:       "TPM category management",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_categorys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "tpm_manager",
		CLIName:           "tpm-manager",
		Description:       "TPM manager operations",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_managers",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "tpm_provision",
		CLIName:           "tpm-provision",
		Description:       "TPM provisioning operations",
		APIPath:           "/api/config/namespaces/system/tpm_provisions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "user",
		CLIName:           "user",
		Description:       "User management",
		APIPath:           "/api/web/namespaces/system/users",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "user_setting",
		CLIName:           "user-setting",
		Description:       "User settings management",
		APIPath:           "/api/web/namespaces/system/user_settings",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
			Delete: true,
		},
	})

	// Monitoring & Observability
	Register(&ResourceType{
		Name:              "observability_subscription",
		CLIName:           "observability-subscription",
		Description:       "Observability subscription management",
		APIPath:           "/api/config/namespaces/system/observability_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "synthetic_monitor",
		CLIName:           "synthetic-monitor",
		Description:       "Synthetic monitoring configuration",
		APIPath:           "/api/config/namespaces/{namespace}/synthetic_monitors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "synthetic_monitor_dns",
		CLIName:           "synthetic-monitor-dns",
		Description:       "DNS synthetic monitor configuration",
		APIPath:           "/api/config/namespaces/{namespace}/v1_dns_monitors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "synthetic_monitor_http",
		CLIName:           "synthetic-monitor-http",
		Description:       "HTTP synthetic monitor configuration",
		APIPath:           "/api/config/namespaces/{namespace}/v1_http_monitors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "report",
		CLIName:           "report",
		Description:       "Report management",
		APIPath:           "/api/config/namespaces/{namespace}/reports",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "report_config",
		CLIName:           "report-config",
		Description:       "Report configuration",
		APIPath:           "/api/config/namespaces/{namespace}/report_configs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	// Network Operations
	Register(&ResourceType{
		Name:              "operate_bgp",
		CLIName:           "operate-bgp",
		Description:       "BGP operational data",
		APIPath:           "/api/operate/namespaces/{namespace}/bgp",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "operate_crl",
		CLIName:           "operate-crl",
		Description:       "CRL operations",
		APIPath:           "/api/operate/namespaces/{namespace}/crl",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_debug",
		CLIName:           "operate-debug",
		Description:       "Debug operations",
		APIPath:           "/api/operate/namespaces/{namespace}/debug",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_dhcp",
		CLIName:           "operate-dhcp",
		Description:       "DHCP operational status",
		APIPath:           "/api/operate/namespaces/{namespace}/dhcp",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "operate_flow",
		CLIName:           "operate-flow",
		Description:       "Flow operations",
		APIPath:           "/api/operate/namespaces/{namespace}/flow",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_lte",
		CLIName:           "operate-lte",
		Description:       "LTE operations",
		APIPath:           "/api/operate/namespaces/{namespace}/lte",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_ping",
		CLIName:           "operate-ping",
		Description:       "Ping operation",
		APIPath:           "/api/operate/namespaces/{namespace}/ping",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_route",
		CLIName:           "operate-route",
		Description:       "Route operations",
		APIPath:           "/api/operate/namespaces/{namespace}/route",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_tcpdump",
		CLIName:           "operate-tcpdump",
		Description:       "TCP dump operation",
		APIPath:           "/api/operate/namespaces/{namespace}/tcpdump",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_traceroute",
		CLIName:           "operate-traceroute",
		Description:       "Traceroute operation",
		APIPath:           "/api/operate/namespaces/{namespace}/traceroute",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_usb",
		CLIName:           "operate-usb",
		Description:       "USB operations",
		APIPath:           "/api/operate/namespaces/{namespace}/usb",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "operate_wifi",
		CLIName:           "operate-wifi",
		Description:       "WiFi operations",
		APIPath:           "/api/operate/namespaces/{namespace}/wifi",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	// Platform Services
	Register(&ResourceType{
		Name:              "signup",
		CLIName:           "signup",
		Description:       "User signup operations",
		APIPath:           "/api/signup",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "status_at_site",
		CLIName:           "status-at-site",
		Description:       "Status at site information",
		APIPath:           "/api/config/namespaces/{namespace}/status_at_sites",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "topology",
		CLIName:           "topology",
		Description:       "Topology queries and generation",
		APIPath:           "/api/topology",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "stored_object",
		CLIName:           "stored-object",
		Description:       "Stored object management",
		APIPath:           "/api/config/namespaces/{namespace}/stored_objects",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "subscription",
		CLIName:           "subscription",
		Description:       "Subscription management",
		APIPath:           "/api/usage/namespaces/{namespace}/subscriptions",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "usage_plan",
		CLIName:           "usage-plan",
		Description:       "Usage plan information",
		APIPath:           "/api/usage/namespaces/{namespace}/plans",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "usage_invoice",
		CLIName:           "usage-invoice",
		Description:       "Usage invoice information",
		APIPath:           "/api/usage/namespaces/{namespace}/invoices",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	// Shape Bot Defense
	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_allowlist_policy",
		CLIName:           "shape-bot-defense-bot-allowlist-policy",
		Description:       "Shape bot defense allowlist policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_allowlist_policys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_endpoint_policy",
		CLIName:           "shape-bot-defense-bot-endpoint-policy",
		Description:       "Shape bot defense endpoint policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_endpoint_policys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_infrastructure",
		CLIName:           "shape-bot-defense-bot-infrastructure",
		Description:       "Shape bot defense infrastructure",
		APIPath:           "/api/config/namespaces/{namespace}/bot_infrastructures",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_network_policy",
		CLIName:           "shape-bot-defense-bot-network-policy",
		Description:       "Shape bot defense network policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_network_policys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_mobile_base_config",
		CLIName:           "shape-bot-defense-mobile-base-config",
		Description:       "Shape bot defense mobile base configuration",
		APIPath:           "/api/config/namespaces/{namespace}/mobile_base_configs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_mobile_sdk",
		CLIName:           "shape-bot-defense-mobile-sdk",
		Description:       "Shape bot defense mobile SDK",
		APIPath:           "/api/config/namespaces/{namespace}/mobile_sdks",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_protected_application",
		CLIName:           "shape-bot-defense-protected-application",
		Description:       "Shape bot defense protected application",
		APIPath:           "/api/config/namespaces/{namespace}/protected_applications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_reporting",
		CLIName:           "shape-bot-defense-reporting",
		Description:       "Shape bot defense reporting",
		APIPath:           "/api/config/namespaces/{namespace}/bot_defense_reporting",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_subscription",
		CLIName:           "shape-bot-defense-subscription",
		Description:       "Shape bot defense subscription",
		APIPath:           "/api/config/namespaces/system/bot_defense_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_detection_rule",
		CLIName:           "shape-bot-detection-rule",
		Description:       "Shape bot detection rule",
		APIPath:           "/api/config/namespaces/{namespace}/bot_detection_rules",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_bot_detection_update",
		CLIName:           "shape-bot-detection-update",
		Description:       "Shape bot detection update",
		APIPath:           "/api/config/namespaces/{namespace}/bot_detection_updates",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	// Shape BRM Alerts
	Register(&ResourceType{
		Name:              "shape_brmalerts_alert_gen_policy",
		CLIName:           "shape-brmalerts-alert-gen-policy",
		Description:       "Shape BRM alerts generation policy",
		APIPath:           "/api/config/namespaces/{namespace}/alert_gen_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_brmalerts_alert_template",
		CLIName:           "shape-brmalerts-alert-template",
		Description:       "Shape BRM alerts template",
		APIPath:           "/api/config/namespaces/{namespace}/alert_templates",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	// Shape Client Side Defense
	Register(&ResourceType{
		Name:              "shape_client_side_defense",
		CLIName:           "shape-client-side-defense",
		Description:       "Shape client side defense",
		APIPath:           "/api/config/namespaces/system/client_side_defenses",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_allowed_domain",
		CLIName:           "shape-client-side-defense-allowed-domain",
		Description:       "Shape client side defense allowed domain",
		APIPath:           "/api/config/namespaces/{namespace}/allowed_domains",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_mitigated_domain",
		CLIName:           "shape-client-side-defense-mitigated-domain",
		Description:       "Shape client side defense mitigated domain",
		APIPath:           "/api/config/namespaces/{namespace}/mitigated_domains",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_protected_domain",
		CLIName:           "shape-client-side-defense-protected-domain",
		Description:       "Shape client side defense protected domain",
		APIPath:           "/api/config/namespaces/{namespace}/protected_domains",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_subscription",
		CLIName:           "shape-client-side-defense-subscription",
		Description:       "Shape client side defense subscription",
		APIPath:           "/api/config/namespaces/system/client_side_defense_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_device_id",
		CLIName:           "shape-device-id",
		Description:       "Shape device ID",
		APIPath:           "/api/config/namespaces/system/device_ids",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
			Delete: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_mobile_app_shield_subscription",
		CLIName:           "shape-mobile-app-shield-subscription",
		Description:       "Shape mobile app shield subscription",
		APIPath:           "/api/config/namespaces/system/mobile_app_shield_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_mobile_integrator_subscription",
		CLIName:           "shape-mobile-integrator-subscription",
		Description:       "Shape mobile integrator subscription",
		APIPath:           "/api/config/namespaces/system/mobile_integrator_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	// Shape Data Delivery
	Register(&ResourceType{
		Name:              "shape_data_delivery",
		CLIName:           "shape-data-delivery",
		Description:       "Shape data delivery",
		APIPath:           "/api/config/namespaces/{namespace}/data_deliverys",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_data_delivery_receiver",
		CLIName:           "shape-data-delivery-receiver",
		Description:       "Shape data delivery receiver",
		APIPath:           "/api/config/namespaces/{namespace}/data_delivery_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_data_delivery_subscription",
		CLIName:           "shape-data-delivery-subscription",
		Description:       "Shape data delivery subscription",
		APIPath:           "/api/config/namespaces/system/data_delivery_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_recognize",
		CLIName:           "shape-recognize",
		Description:       "Shape recognize",
		APIPath:           "/api/shape/recognize",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_safe",
		CLIName:           "shape-safe",
		Description:       "Shape safe",
		APIPath:           "/api/shape/safe",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "shape_safeap",
		CLIName:           "shape-safeap",
		Description:       "Shape safe AP",
		APIPath:           "/api/shape/safeap",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Get:    true,
			Create: true,
		},
	})

	// WAF & Application Protection
	Register(&ResourceType{
		Name:              "waf",
		CLIName:           "waf",
		Description:       "WAF configuration",
		APIPath:           "/api/config/namespaces/{namespace}/wafs",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "waf_signatures_changelog",
		CLIName:           "waf-signatures-changelog",
		Description:       "WAF signatures changelog",
		APIPath:           "/api/config/namespaces/{namespace}/waf_signatures_changelogs",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	// Miscellaneous
	Register(&ResourceType{
		Name:              "malware_protection_subscription",
		CLIName:           "malware-protection-subscription",
		Description:       "Malware protection subscription",
		APIPath:           "/api/config/namespaces/system/malware_protection_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "nginx_one_subscription",
		CLIName:           "nginx-one-subscription",
		Description:       "NGINX One subscription",
		APIPath:           "/api/config/namespaces/system/nginx_one_subscriptions",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "pbac_catalog",
		CLIName:           "pbac-catalog",
		Description:       "PBAC catalog",
		APIPath:           "/api/web/namespaces/system/catalogs",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Update: true,
		},
	})

	Register(&ResourceType{
		Name:              "virtual_appliance",
		CLIName:           "virtual-appliance",
		Description:       "Virtual appliance",
		APIPath:           "/api/config/namespaces/system/virtual_appliances",
		SupportsNamespace: false,
		Operations: ResourceOperations{
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "voltshare",
		CLIName:           "voltshare",
		Description:       "Voltshare",
		APIPath:           "/api/config/namespaces/{namespace}/voltshares",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "voltshare_admin_policy",
		CLIName:           "voltshare-admin-policy",
		Description:       "Voltshare admin policy",
		APIPath:           "/api/config/namespaces/{namespace}/voltshare_admin_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "was_user_token",
		CLIName:           "was-user-token",
		Description:       "WAS user token",
		APIPath:           "/api/web/namespaces/system/was_user_tokens",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	// Views/Read-Only Resources
	Register(&ResourceType{
		Name:              "views_terraform_parameters",
		CLIName:           "views-terraform-parameters",
		Description:       "Terraform parameters view",
		APIPath:           "/api/config/namespaces/{namespace}/terraform_parameters",
		SupportsNamespace: true,
		Operations: ResourceOperations{
			Get:    true,
			List:   true,
			Create: true,
		},
	})

	Register(&ResourceType{
		Name:              "views_view_internal",
		CLIName:           "views-view-internal",
		Description:       "Internal view",
		APIPath:           "/api/config/namespaces/{namespace}/view_internals",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

}

