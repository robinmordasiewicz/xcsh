package types

// Code generated from OpenAPI specifications. DO NOT EDIT.
// This file contains 268 resource types parsed from F5 XC API specs

func init() {
	registerGeneratedResources()
}

func registerGeneratedResources() {
	Register(&ResourceType{
		Name:              "address_allocator",
		CLIName:           "address-allocator",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/address_allocators",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "advertise_policy",
		CLIName:           "advertise-policy",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/advertise_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "ai_assistant",
		CLIName:           "ai-assistant",
		Description:       "Custom handler for ai assistant related microservice",
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
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "alert",
		CLIName:           "alert",
		Description:       "Alert may be generated based on the metrics or based on severity level in the logs. All alerts are scoped by \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ntenant and namespace and tagged with the following default labels that can be used to fetch\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nthe desired alerts.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"alertname\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" - Name of the alert. This uniquely identifies the alert rule/configuration that generated the alert.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"type\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" - Type of the alert. Type is used to associate alert to a configuration object or any user visible entity.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n         For example, virtual host, virtual network, app_type, etc.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"identifier\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" - Identifier of the alert. For virtual-network, this would be the name of the virtual-network.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"severity\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" - Indicates the severity of the alert. Valid values are minor, major, critical.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAlert may have additional labels associated depending on the labels associated with the metric used to configure the alert rule.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAlerts can be queried by specifying one or more of the above labels in the match filter.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nIf the match filter is not specified, then all the alerts for the tenant and corresponding namespace in the request will be returned in the response.",
		APIPath:           "/api/data/namespaces/{namespace}/alerts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "alert_policy",
		CLIName:           "alert-policy",
		Description:       "F5 Distributed Cloud Statistics API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/alert_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "statistics",
		Domains:           []string{"statistics"},
	})

	Register(&ResourceType{
		Name:              "alert_receiver",
		CLIName:           "alert-receiver",
		Description:       "F5 Distributed Cloud Statistics API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/alert_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "statistics",
		Domains:           []string{"statistics"},
	})

	Register(&ResourceType{
		Name:              "api_credential",
		CLIName:           "api-credential",
		Description:       "F5 Distributed Cloud Api API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/api_credentials",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "api",
		Domains:           []string{"api", "authentication"},
	})

	Register(&ResourceType{
		Name:              "api_definition",
		CLIName:           "api-definition",
		Description:       "F5 Distributed Cloud Api API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/api_definitions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "api",
		Domains:           []string{"api"},
	})

	Register(&ResourceType{
		Name:              "api_group",
		CLIName:           "api-group",
		Description:       "The api_group construct provides a mechanism to classify the universal set of request APIs into a much smaller number of logical groups in order to make it\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\neasier to author and maintain API level access control policies.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nA api_group object consists of an unordered list of api group elements. The method and path from the input request API are matched against all elements in\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nan api_group to determine if the request API belongs to the api group in question. The request API belongs to an api group if it matches at least one element\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nin the api group.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAn api group object can only be created in the 'shared' namespace of a tenant or in the 'shared' namespace of the ves-io tenant. Input request APIs from a\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ntenant are matched against all api groups in that tenant and in the ves-io tenant to determine the set of api groups for that request. The names of the api\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ngroups to which a given request API belongs are subsequently used as input to to check the api group predicate in a service policy or service policy rule.",
		APIPath:           "/api/web/namespaces/{namespace}/api_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "api_group_element",
		CLIName:           "api-group-element",
		Description:       "A api_group_element object consists of an unordered list of HTTP methods and a path regular expression. The method and path from the input request API are\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nmatched against all elements in an api_group to determine if the request API belongs to the api group in question. The match of an input request API against\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nan element is considered to be successful if the input method belongs to the list of HTTP methods in the element and the input path matches the path regex in\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nthe element. The request API belongs to an api group if it matches at least one element in the api group.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAn api group element object can only be created in the 'shared' namespace of a tenant or in the 'shared' namespace of the ves-io tenant. Note that any given\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nelement can belong to one or more api groups. Input request APIs from a tenant are matched against all api groups in that tenant and in the ves-io tenant to\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ndetermine the set of api groups for that request. The names of the api groups to which a given request API belongs are subsequently used as input to to check\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nthe api group predicate in a service policy or service policy rule.",
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
		Name:              "api_sec_rule_suggestion",
		CLIName:           "api-sec-rule-suggestion",
		Description:       "API Security rule suggestions",
		APIPath:           "/api/config/namespaces/{namespace}/rule_suggestions",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "app_api_group",
		CLIName:           "app-api-group",
		Description:       "F5 Distributed Cloud Api API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/app_api_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "api",
		Domains:           []string{"api"},
	})

	Register(&ResourceType{
		Name:              "app_firewall",
		CLIName:           "app-firewall",
		Description:       "F5 Distributed Cloud Waf API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/app_firewalls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "waf",
		Domains:           []string{"waf"},
	})

	Register(&ResourceType{
		Name:              "app_security",
		CLIName:           "app-security",
		Description:       "API to create API endpoint protection rule suggestion from App Security Monitoring pages",
		APIPath:           "/api/config/namespaces/{namespace}/app_securitys",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "app_setting",
		CLIName:           "app-setting",
		Description:       "F5 Distributed Cloud Service Mesh API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/app_settings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "service_mesh",
		Domains:           []string{"service_mesh"},
	})

	Register(&ResourceType{
		Name:              "app_type",
		CLIName:           "app-type",
		Description:       "F5 Distributed Cloud Service Mesh API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/app_types",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "service_mesh",
		Domains:           []string{"service_mesh"},
	})

	Register(&ResourceType{
		Name:              "authentication",
		CLIName:           "authentication",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/authentications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "tenant_and_identity",
		Domains:           []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "aws_tgw_site",
		CLIName:           "aws-tgw-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/aws_tgw_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "aws_vpc_site",
		CLIName:           "aws-vpc-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/aws_vpc_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "azure_vnet_site",
		CLIName:           "azure-vnet-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/azure_vnet_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "bgp",
		CLIName:           "bgp",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/bgps",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "bgp_asn_set",
		CLIName:           "bgp-asn-set",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/bgp_asn_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "bgp_routing_policy",
		CLIName:           "bgp-routing-policy",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/bgp_routing_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
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
		PrimaryDomain:     "bigip",
		Domains:           []string{"bigip"},
	})

	Register(&ResourceType{
		Name:              "bigip_irule",
		CLIName:           "bigip-irule",
		Description:       "F5 Distributed Cloud Bigip API specifications",
		APIPath:           "/api/bigipconnector/namespaces/{namespace}/bigip_irules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "bigip",
		Domains:           []string{"bigip"},
	})

	Register(&ResourceType{
		Name:              "bigip_virtual_server",
		CLIName:           "bigip-virtual-server",
		Description:       "BIG-IP virtual server view repesents the internal virtual host corresponding to the virtual-servers discovered from BIG-IPs\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nIt exposes parameters to enable API discovery and other WAAP security features on the virtual server.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nView will create following child objects.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Virtual-host\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* API-inventory\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* App-type\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* App-setting",
		APIPath:           "/api/config/namespaces/{namespace}/bigip_virtual_servers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "billing_payment_method",
		CLIName:           "billing-payment-method",
		Description:       "Billing payment method management",
		APIPath:           "/api/web/namespaces/{namespace}/payment_methods",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "billing_plan_transition",
		CLIName:           "billing-plan-transition",
		Description:       "Billing plan transition management",
		APIPath:           "/api/web/namespaces/{namespace}/plan_transitions",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "bot_defense_app_infrastructure",
		CLIName:           "bot-defense-app-infrastructure",
		Description:       "F5 Distributed Cloud Bot And Threat Defense API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/bot_defense_app_infrastructures",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "bot_and_threat_defense",
		Domains:           []string{"bot_and_threat_defense"},
	})

	Register(&ResourceType{
		Name:              "cdn_cache_rule",
		CLIName:           "cdn-cache-rule",
		Description:       "F5 Distributed Cloud Cdn API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cdn_cache_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cdn",
		Domains:           []string{"cdn"},
	})

	Register(&ResourceType{
		Name:              "cdn_loadbalancer",
		CLIName:           "cdn-loadbalancer",
		Description:       "F5 Distributed Cloud Cdn API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cdn_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cdn",
		Domains:           []string{"cdn"},
	})

	Register(&ResourceType{
		Name:              "certificate",
		CLIName:           "certificate",
		Description:       "F5 Distributed Cloud Certificates API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/certificates",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "certificates",
		Domains:           []string{"certificates"},
	})

	Register(&ResourceType{
		Name:              "certificate_chain",
		CLIName:           "certificate-chain",
		Description:       "F5 Distributed Cloud Certificates API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/certificate_chains",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "certificates",
		Domains:           []string{"certificates"},
	})

	Register(&ResourceType{
		Name:              "certified_hardware",
		CLIName:           "certified-hardware",
		Description:       "Certified Hardware object represents physical hardware or cloud instance type that will be used to instantiate\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\na volterra software appliance instance for the F5XC sites (Customer edge site). It has following information\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Type\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Vendor\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Model\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  List of devices supported\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Image name\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Latest image release as status\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nCertified Hardware objects are only available in volterra shared namespace (ves-io/shared).\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nThese are created by volterra. It serves following purpose.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Let user know supported hardware and devices on given hardware.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  How they are used and configured at boot strap\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  *  Image in which boot strap config  and any custom scripts are bundled\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nThis is required so that zero touch provisioning would work. If a generic image is used then, user will have to login into the hardware and\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nprovide bootstrap config.",
		APIPath:           "/api/config/namespaces/{namespace}/certified_hardwares",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "cloud_connect",
		CLIName:           "cloud-connect",
		Description:       "F5 Distributed Cloud Cloud Infrastructure API specifications",
		APIPath:           "/api/config/namespaces/system/edge_credentials",
		SupportsNamespace: false,
		Operations:        AllOperations(),
		PrimaryDomain:     "cloud_infrastructure",
		Domains:           []string{"cloud_infrastructure"},
	})

	Register(&ResourceType{
		Name:              "cloud_credentials",
		CLIName:           "cloud-credentials",
		Description:       "F5 Distributed Cloud Cloud Infrastructure API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_credentialss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cloud_infrastructure",
		Domains:           []string{"cloud_infrastructure"},
	})

	Register(&ResourceType{
		Name:              "cloud_elastic_ip",
		CLIName:           "cloud-elastic-ip",
		Description:       "F5 Distributed Cloud Cloud Infrastructure API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_elastic_ips",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cloud_infrastructure",
		Domains:           []string{"cloud_infrastructure"},
	})

	Register(&ResourceType{
		Name:              "cloud_link",
		CLIName:           "cloud-link",
		Description:       "F5 Distributed Cloud Cloud Infrastructure API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_links",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cloud_infrastructure",
		Domains:           []string{"cloud_infrastructure"},
	})

	Register(&ResourceType{
		Name:              "cloud_region",
		CLIName:           "cloud-region",
		Description:       "Cloud Region contains tenant specific configuration\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nobject. Users cannot create/delete these objects. They will be internally created\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nor deleted whenever the corresponding cloud_region_region object is created/deleted \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nand possibly based on tenant configuration (e.g. Cloud Region feature may be disabled\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nfor some tenants)",
		APIPath:           "/api/config/namespaces/{namespace}/cloud_regions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "cluster",
		CLIName:           "cluster",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/clusters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration"},
	})

	Register(&ResourceType{
		Name:              "cminstance",
		CLIName:           "cminstance",
		Description:       "F5 Distributed Cloud Marketplace API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/cminstances",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "marketplace",
		Domains:           []string{"marketplace"},
	})

	Register(&ResourceType{
		Name:              "contact",
		CLIName:           "contact",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/contacts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "tenant_and_identity",
		Domains:           []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "container_registry",
		CLIName:           "container-registry",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/container_registrys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration"},
	})

	Register(&ResourceType{
		Name:              "crl",
		CLIName:           "crl",
		Description:       "F5 Distributed Cloud Certificates API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/crls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "certificates",
		Domains:           []string{"certificates"},
	})

	Register(&ResourceType{
		Name:              "customer_support",
		CLIName:           "customer-support",
		Description:       "F5 Distributed Cloud Support API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/customer_supports",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "support",
		Domains:           []string{"support"},
	})

	Register(&ResourceType{
		Name:              "data_privacy_geo_config",
		CLIName:           "data-privacy-geo-config",
		Description:       "Data privacy geo configuration",
		APIPath:           "/api/config/namespaces/{namespace}/geo_configs",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
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
		Description:       "F5 Distributed Cloud Data And Privacy Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/data_types",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "data_and_privacy_security",
		Domains:           []string{"data_and_privacy_security"},
	})

	Register(&ResourceType{
		Name:              "dc_cluster_group",
		CLIName:           "dc-cluster-group",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/dc_cluster_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "discovered_service",
		CLIName:           "discovered-service",
		Description:       "Discovered Services represents the services (virtual-servers, k8s services, etc) which are discovered via the different discovery workflows.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\npackage for Discovered Services",
		APIPath:           "/api/config/namespaces/{namespace}/discovered_services",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "discovery",
		CLIName:           "discovery",
		Description:       "F5 Distributed Cloud Api API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/discoverys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "api",
		Domains:           []string{"api"},
	})

	Register(&ResourceType{
		Name:              "dns_compliance_checks",
		CLIName:           "dns-compliance-checks",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/dns_compliance_checkss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
	})

	Register(&ResourceType{
		Name:              "dns_domain",
		CLIName:           "dns-domain",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/dns_domains",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
	})

	Register(&ResourceType{
		Name:              "dns_lb_health_check",
		CLIName:           "dns-lb-health-check",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_lb_health_checks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
	})

	Register(&ResourceType{
		Name:              "dns_lb_pool",
		CLIName:           "dns-lb-pool",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_lb_pools",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
	})

	Register(&ResourceType{
		Name:              "dns_load_balancer",
		CLIName:           "dns-load-balancer",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_load_balancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
	})

	Register(&ResourceType{
		Name:              "dns_zone",
		CLIName:           "dns-zone",
		Description:       "F5 Distributed Cloud Dns API specifications",
		APIPath:           "/api/config/dns/namespaces/{namespace}/dns_zones",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "dns",
		Domains:           []string{"dns"},
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
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "endpoint",
		CLIName:           "endpoint",
		Description:       "F5 Distributed Cloud Service Mesh API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/endpoints",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "service_mesh",
		Domains:           []string{"service_mesh"},
	})

	Register(&ResourceType{
		Name:              "enhanced_firewall_policy",
		CLIName:           "enhanced-firewall-policy",
		Description:       "F5 Distributed Cloud Waf API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/enhanced_firewall_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "waf",
		Domains:           []string{"waf"},
	})

	Register(&ResourceType{
		Name:              "external_connector",
		CLIName:           "external-connector",
		Description:       "F5 Distributed Cloud Marketplace API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/external_connectors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "marketplace",
		Domains:           []string{"marketplace"},
	})

	Register(&ResourceType{
		Name:              "fast_acl",
		CLIName:           "fast-acl",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/fast_acls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "fast_acl_rule",
		CLIName:           "fast-acl-rule",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/fast_acl_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "filter_set",
		CLIName:           "filter-set",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/filter_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "fleet",
		CLIName:           "fleet",
		Description:       "F5 Distributed Cloud Ce Management API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/fleets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ce_management",
		Domains:           []string{"ce_management"},
	})

	Register(&ResourceType{
		Name:              "flow",
		CLIName:           "flow",
		Description:       "APIs to get Flow records and data",
		APIPath:           "/api/data/namespaces/{namespace}/flows",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
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
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/forward_proxy_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "forwarding_class",
		CLIName:           "forwarding-class",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/forwarding_classs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "gcp_vpc_site",
		CLIName:           "gcp-vpc-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/gcp_vpc_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "geo_location_set",
		CLIName:           "geo-location-set",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/geo_location_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "gia",
		CLIName:           "gia",
		Description:       "GIA operations",
		APIPath:           "/api/config/namespaces/system/gias",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "global_log_receiver",
		CLIName:           "global-log-receiver",
		Description:       "F5 Distributed Cloud Statistics API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/global_log_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "statistics",
		Domains:           []string{"statistics"},
	})

	Register(&ResourceType{
		Name:              "graph_connectivity",
		CLIName:           "graph-connectivity",
		Description:       "Connectivity graph analysis",
		APIPath:           "/api/graph/namespaces/{namespace}/connectivity",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "graph_l3l4",
		CLIName:           "graph-l3l4",
		Description:       "L3/L4 topology graph",
		APIPath:           "/api/graph/namespaces/{namespace}/l3l4",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "graph_service",
		CLIName:           "graph-service",
		Description:       "Service graph queries",
		APIPath:           "/api/graph/namespaces/{namespace}/services",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "graph_site",
		CLIName:           "graph-site",
		Description:       "Site topology graph",
		APIPath:           "/api/graph/namespaces/{namespace}/sites",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "healthcheck",
		CLIName:           "healthcheck",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/healthchecks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "http_loadbalancer",
		CLIName:           "http-loadbalancer",
		Description:       "F5 Distributed Cloud Cdn API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/http_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "cdn",
		Domains:           []string{"cdn", "virtual"},
	})

	Register(&ResourceType{
		Name:              "ike1",
		CLIName:           "ike1",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/ike1s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "ike2",
		CLIName:           "ike2",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/ike2s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "ike_phase1_profile",
		CLIName:           "ike-phase1-profile",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/ike_phase1_profiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "ike_phase2_profile",
		CLIName:           "ike-phase2-profile",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/ike_phase2_profiles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "implicit_label",
		CLIName:           "implicit-label",
		Description:       "Implicit labels are attached to objects implicitly by the system. Users are not allowed to create/update/delete these labels\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nThey are also not allowed to attach/detach these labels to objects. This API is provided to get the implicit labels available\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nto be used in service-policies",
		APIPath:           "/api/config/namespaces/system/implicit_labels",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect",
		CLIName:           "infraprotect",
		Description:       "APIs to get monitoring data for infraprotect.",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotects",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_asn",
		CLIName:           "infraprotect-asn",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_asns",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_asn_prefix",
		CLIName:           "infraprotect-asn-prefix",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_asn_prefixs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_deny_list_rule",
		CLIName:           "infraprotect-deny-list-rule",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_deny_list_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_rule",
		CLIName:           "infraprotect-firewall-rule",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_rule_group",
		CLIName:           "infraprotect-firewall-rule-group",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rule_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_firewall_ruleset",
		CLIName:           "infraprotect-firewall-ruleset",
		Description:       "DDoS transit Firewall Ruleset",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_firewall_rulesets",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "infraprotect_information",
		CLIName:           "infraprotect-information",
		Description:       "Infraprotect information about the current organisation",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_informations",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "infraprotect_internet_prefix_advertisement",
		CLIName:           "infraprotect-internet-prefix-advertisement",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_internet_prefix_advertisements",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "infraprotect_tunnel",
		CLIName:           "infraprotect-tunnel",
		Description:       "F5 Distributed Cloud Ddos API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/infraprotect_tunnels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ddos",
		Domains:           []string{"ddos"},
	})

	Register(&ResourceType{
		Name:              "ip_prefix_set",
		CLIName:           "ip-prefix-set",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/ip_prefix_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "k8s_cluster",
		CLIName:           "k8s-cluster",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_clusters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "k8s_cluster_role",
		CLIName:           "k8s-cluster-role",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_cluster_roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "k8s_cluster_role_binding",
		CLIName:           "k8s-cluster-role-binding",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_cluster_role_bindings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "k8s_pod_security_admission",
		CLIName:           "k8s-pod-security-admission",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_pod_security_admissions",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration"},
	})

	Register(&ResourceType{
		Name:              "k8s_pod_security_policy",
		CLIName:           "k8s-pod-security-policy",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/k8s_pod_security_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration"},
	})

	Register(&ResourceType{
		Name:              "known_label",
		CLIName:           "known-label",
		Description:       "F5 Distributed Cloud Users API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/known_labels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "users",
		Domains:           []string{"users"},
	})

	Register(&ResourceType{
		Name:              "known_label_key",
		CLIName:           "known-label-key",
		Description:       "F5 Distributed Cloud Users API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/known_label_keys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "users",
		Domains:           []string{"users"},
	})

	Register(&ResourceType{
		Name:              "log",
		CLIName:           "log",
		Description:       "Two types of logs are supported, viz, access logs and audit logs.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  * Access logs are sampled records of API calls made to a virtual host. It contains\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    both the request and the response data with more context like application type,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    user, request path, method, request body, response code, source,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    destination, etc.,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  * Audit logs provides audit of all configuration changes made in the system using\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    public APIs provided by Volterra. It contains both the request and response body\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    with additional context necessary for post-mortem analysis such as user, request path,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    method, request body, response code, source, destination service, etc.,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nBoth the access logs and audit logs are used to find \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"who did what and when and what was the result?\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwho - answered by user/user-agent in the log.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwhat - answered by request url/method/body in the log.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwhen - answered by timestamp in the log.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nresult - answered by response code in the log.",
		APIPath:           "/api/data/namespaces/{namespace}/logs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "log_receiver",
		CLIName:           "log-receiver",
		Description:       "F5 Distributed Cloud Statistics API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/log_receivers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "statistics",
		Domains:           []string{"statistics"},
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
		Name:              "malicious_user_mitigation",
		CLIName:           "malicious-user-mitigation",
		Description:       "F5 Distributed Cloud Secops And Incident Response API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/malicious_user_mitigations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "secops_and_incident_response",
		Domains:           []string{"secops_and_incident_response"},
	})

	Register(&ResourceType{
		Name:              "malware_protection_subscription",
		CLIName:           "malware-protection-subscription",
		Description:       "Malware protection subscription",
		APIPath:           "/api/config/namespaces/system/malware_protection_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "marketplace_aws_account",
		CLIName:           "marketplace-aws-account",
		Description:       "AWS Marketplace account management",
		APIPath:           "/api/marketplace/aws/accounts",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "marketplace_xc_saas",
		CLIName:           "marketplace-xc-saas",
		Description:       "XC SaaS marketplace management",
		APIPath:           "/api/marketplace/xc_saas",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "module_management",
		CLIName:           "module-management",
		Description:       "Package for managing a module.",
		APIPath:           "/api/config/namespaces/{namespace}/module_managements",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "namespace",
		CLIName:           "namespace",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces",
		SupportsNamespace: false,
		Operations:        AllOperations(),
		DeleteConfig: &DeleteConfig{
			PathSuffix:  "/cascade_delete",
			Method:      "POST",
			IncludeBody: true,
		},
		PrimaryDomain: "tenant_and_identity",
		Domains:       []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "namespace_role",
		CLIName:           "namespace-role",
		Description:       "Namespace role defines a user's role in a namespace.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nNamespace role links a user with a role namespace. Using this object one can assign/remove a role to a user in a namespace.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nNamespace roles are assigned either explicitly by calling this API, or implicitly by creating users or by signing up.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n see role object for information on roles\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n see namespace object for information on namespaces\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n see user object for information on users",
		APIPath:           "/api/web/namespaces/{namespace}/namespace_roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "nat_policy",
		CLIName:           "nat-policy",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/nat_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "network_connector",
		CLIName:           "network-connector",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_connectors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "network_firewall",
		CLIName:           "network-firewall",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_firewalls",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "network_interface",
		CLIName:           "network-interface",
		Description:       "F5 Distributed Cloud Ce Management API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_interfaces",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ce_management",
		Domains:           []string{"ce_management"},
	})

	Register(&ResourceType{
		Name:              "network_policy",
		CLIName:           "network-policy",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "network_policy_rule",
		CLIName:           "network-policy-rule",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "network_policy_set",
		CLIName:           "network-policy-set",
		Description:       "Network policy set implements L3/L4 stateful firewall.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nIt is a list of one or more Network policy references and are applied sequentially in order specified in the list.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nNetwork policy set can be configured via network firewall object or vK8s\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Network Firewall is a fleet object which can take a reference to network policy set.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n  Firewall is applied to VIRTUAL_NETWORK_SITE_LOCAL and VIRTUAL_NETWORK_SITE in corresponding fleet\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* vK8s will inherit network policy set configured in its own namespace and tenant\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nNetwork policy references to be attached to network policy set can be picked from\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Same namespace and tenant as of network policy set\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Shared namespace of the Tenant\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Shared namespace of 'ves.io' tenant",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "network_policy_view",
		CLIName:           "network-policy-view",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/network_policy_views",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "nfv_service",
		CLIName:           "nfv-service",
		Description:       "F5 Distributed Cloud Service Mesh API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/nfv_services",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "service_mesh",
		Domains:           []string{"service_mesh"},
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
		Name:              "nginx_one_subscription",
		CLIName:           "nginx-one-subscription",
		Description:       "NGINX One subscription",
		APIPath:           "/api/config/namespaces/system/nginx_one_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "observability_subscription",
		CLIName:           "observability-subscription",
		Description:       "Observability subscription management",
		APIPath:           "/api/config/namespaces/system/observability_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "oidc_provider",
		CLIName:           "oidc-provider",
		Description:       "F5XC Identity supports identity brokering and third-party identity can be added to\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nprovide Single Sign-On (SSO) login for user to access tenant via VoltConsole. \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nUsing Volterra's OIDC Provider config object API(s), tenant admin can configure and manage\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nSSO providers such as Google, Microsoft, Okta or any provider that supports\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nOpenIDConnect(OIDC) V1.0 protocol. It is required that the OIDC provider application that will \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nbe integrated must have the support enabled for Authorization Code flow as defined by the specification.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nBefore proceeding, admin needs to have access to organization's authentication application\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nor has permission to create a new one. Create API require entering details of well-known\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nOpenID configuration of authentication application and once successful creation, admin should\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nenable the redirect URL provided by volterra identity in application's allowed list of URLs. \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nMore details of this can be found under create request/response.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nWith OIDC provider configured, admin of a tenant can make use of Single Sign On (SSO)\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nfunctionality for users to access F5XC service using same email address and admin has the\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nflexibility to re-use the authentication/identity provider that may be already using \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwithin their organization. Once SSO is enabled, except tenant admin (owner) all users in the tenant \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwill be converted to SSO user and will lose existing email/password login created with Volterra\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nand can only use SSO to login.",
		APIPath:           "/api/config/namespaces/{namespace}/oidc_providers",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

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
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_debug",
		CLIName:           "operate-debug",
		Description:       "Debug operations",
		APIPath:           "/api/operate/namespaces/{namespace}/debug",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
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
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_lte",
		CLIName:           "operate-lte",
		Description:       "LTE operations",
		APIPath:           "/api/operate/namespaces/{namespace}/lte",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_ping",
		CLIName:           "operate-ping",
		Description:       "Ping operation",
		APIPath:           "/api/operate/namespaces/{namespace}/ping",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_route",
		CLIName:           "operate-route",
		Description:       "Route operations",
		APIPath:           "/api/operate/namespaces/{namespace}/route",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_tcpdump",
		CLIName:           "operate-tcpdump",
		Description:       "TCP dump operation",
		APIPath:           "/api/operate/namespaces/{namespace}/tcpdump",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_traceroute",
		CLIName:           "operate-traceroute",
		Description:       "Traceroute operation",
		APIPath:           "/api/operate/namespaces/{namespace}/traceroute",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_usb",
		CLIName:           "operate-usb",
		Description:       "USB operations",
		APIPath:           "/api/operate/namespaces/{namespace}/usb",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "operate_wifi",
		CLIName:           "operate-wifi",
		Description:       "WiFi operations",
		APIPath:           "/api/operate/namespaces/{namespace}/wifi",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "origin_pool",
		CLIName:           "origin-pool",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/origin_pools",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
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
		Name:              "pbac_catalog",
		CLIName:           "pbac-catalog",
		Description:       "PBAC catalog",
		APIPath:           "/api/web/namespaces/system/catalogs",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: false, Get: false, List: false, Update: true, Delete: false, Status: false},
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
		Description:       "F5 Distributed Cloud Rate Limiting API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/policers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "rate_limiting",
		Domains:           []string{"rate_limiting"},
	})

	Register(&ResourceType{
		Name:              "policy_based_routing",
		CLIName:           "policy-based-routing",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/policy_based_routings",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
	})

	Register(&ResourceType{
		Name:              "protocol_inspection",
		CLIName:           "protocol-inspection",
		Description:       "F5 Distributed Cloud Waf API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/protocol_inspections",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "waf",
		Domains:           []string{"waf"},
	})

	Register(&ResourceType{
		Name:              "protocol_policer",
		CLIName:           "protocol-policer",
		Description:       "F5 Distributed Cloud Rate Limiting API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/protocol_policers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "rate_limiting",
		Domains:           []string{"rate_limiting"},
	})

	Register(&ResourceType{
		Name:              "proxy",
		CLIName:           "proxy",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/proxys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "public_ip",
		CLIName:           "public-ip",
		Description:       "public_ip object represents a public IP address that is available on a set of virtual sites",
		APIPath:           "/api/config/namespaces/{namespace}/public_ips",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "quota",
		CLIName:           "quota",
		Description:       "F5 Distributed Cloud Billing And Usage API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/quotas",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
		PrimaryDomain:     "billing_and_usage",
		Domains:           []string{"billing_and_usage"},
	})

	Register(&ResourceType{
		Name:              "rate_limiter",
		CLIName:           "rate-limiter",
		Description:       "F5 Distributed Cloud Rate Limiting API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/rate_limiters",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "rate_limiting",
		Domains:           []string{"rate_limiting"},
	})

	Register(&ResourceType{
		Name:              "rate_limiter_policy",
		CLIName:           "rate-limiter-policy",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/rate_limiter_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "rbac_policy",
		CLIName:           "rbac-policy",
		Description:       "A rbac_policy object consists of list of rbac policy rules that when assigned to a user via Role object,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nit controls access of an user to list of APIs defined under the API Group name.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nEach rule under rbac_policy consist of a name of the API Group.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nBy default, the access is set to allow for the API Group.",
		APIPath:           "/api/web/namespaces/{namespace}/rbac_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "registration",
		CLIName:           "registration",
		Description:       "F5 Distributed Cloud Ce Management API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/registrations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ce_management",
		Domains:           []string{"ce_management"},
	})

	Register(&ResourceType{
		Name:              "report",
		CLIName:           "report",
		Description:       "Report configuration contains the information like\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    Time at which the report was last sent to object store.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    Report ID.",
		APIPath:           "/api/config/namespaces/{namespace}/reports",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "report_config",
		CLIName:           "report-config",
		Description:       "F5 Distributed Cloud Statistics API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/report_configs",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "statistics",
		Domains:           []string{"statistics"},
	})

	Register(&ResourceType{
		Name:              "role",
		CLIName:           "role",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/roles",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "tenant_and_identity",
		Domains:           []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "route",
		CLIName:           "route",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/routes",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "scim",
		CLIName:           "scim",
		Description:       "This schema specification details Volterra's support for SCIM protocol.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAdmin can use SCIM feature on top of SSO to enable automated provisioning of\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nuser and user groups from external identity provider into the F5 saas platform.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nWith this feature, complete life cycle management of user and groups can be\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nachieved from single source of truth which is managed by tenant's admin.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ncurrent protocol support is using schema version v2.0 https://datatracker.ietf.org/doc/html/rfc7643 \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nSCIM feature can be enabled part of SSO configuration (using RPC `UpdateScimIntegration` under oidc_provider resource)\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nBy default, F5XC will not sync groups and users. Admin is required to set object identifier of group \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nin external identity provider to corresponding user_group resource in volterra. Users with corresponding\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ngroup membership if exist in external identity provider will be synced.",
		APIPath:           "/api/scim/v2",
		SupportsNamespace: false,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_management",
		CLIName:           "secret-management",
		Description:       "F5XC Secret Management service serves APIs for information required for offline secret encryption such as getting the public key and getting the secret policy document.",
		APIPath:           "/api/config/namespaces/{namespace}/secret_managements",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "secret_management_access",
		CLIName:           "secret-management-access",
		Description:       "F5 Distributed Cloud Blindfold API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/secret_management_accesss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "blindfold",
		Domains:           []string{"blindfold"},
	})

	Register(&ResourceType{
		Name:              "secret_policy",
		CLIName:           "secret-policy",
		Description:       "F5 Distributed Cloud Blindfold API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/secret_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "blindfold",
		Domains:           []string{"blindfold"},
	})

	Register(&ResourceType{
		Name:              "secret_policy_rule",
		CLIName:           "secret-policy-rule",
		Description:       "F5 Distributed Cloud Blindfold API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/secret_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "blindfold",
		Domains:           []string{"blindfold"},
	})

	Register(&ResourceType{
		Name:              "securemesh_site",
		CLIName:           "securemesh-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/securemesh_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "securemesh_site_v2",
		CLIName:           "securemesh-site-v2",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/securemesh_site_v2s",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "segment",
		CLIName:           "segment",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/segments",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security"},
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
		Description:       "F5 Distributed Cloud Data And Privacy Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/sensitive_data_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "data_and_privacy_security",
		Domains:           []string{"data_and_privacy_security"},
	})

	Register(&ResourceType{
		Name:              "service_policy",
		CLIName:           "service-policy",
		Description:       "F5 Distributed Cloud Network Security API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/service_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network_security",
		Domains:           []string{"network_security", "virtual"},
	})

	Register(&ResourceType{
		Name:              "service_policy_rule",
		CLIName:           "service-policy-rule",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/service_policy_rules",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "service_policy_set",
		CLIName:           "service-policy-set",
		Description:       "A service_policy_set object consists of an ordered list of references to service_policy objects. The policies are evaluated in the specified order against\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\na set of input fields that are extracted from or derived from an L7 request API. The evaluation of the policy set terminates when the request API matches a\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\npolicy which results in an \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"allow\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" or \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"deny\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" action. If the request API does not match a policy, the next policy in the list is evaluated. If the request API\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ndoes not match any of the policies in the set, the result is a \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"default_deny\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" action.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nThere can be no more than one service policy set configured under a namespace. That policy set is applied to requests destined to all virtual_hosts in that\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nnamespace.",
		APIPath:           "/api/config/namespaces/{namespace}/service_policy_sets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_allowlist_policy",
		CLIName:           "shape-bot-defense-bot-allowlist-policy",
		Description:       "Shape bot defense allowlist policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_allowlist_policys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_endpoint_policy",
		CLIName:           "shape-bot-defense-bot-endpoint-policy",
		Description:       "Shape bot defense endpoint policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_endpoint_policys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_infrastructure",
		CLIName:           "shape-bot-defense-bot-infrastructure",
		Description:       "Shape bot defense infrastructure",
		APIPath:           "/api/config/namespaces/{namespace}/bot_infrastructures",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_bot_network_policy",
		CLIName:           "shape-bot-defense-bot-network-policy",
		Description:       "Shape bot defense network policy",
		APIPath:           "/api/config/namespaces/{namespace}/bot_network_policys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_instance",
		CLIName:           "shape-bot-defense-instance",
		Description:       "Shape Bot Defense Instance is the main configuration for a Shape Integration.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nIt defines various configuration parameters needed to use Shape SSEs.",
		APIPath:           "/api/config/namespaces/{namespace}/shape_bot_defense_instances",
		SupportsNamespace: true,
		Operations:        AllOperations(),
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
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_defense_subscription",
		CLIName:           "shape-bot-defense-subscription",
		Description:       "Shape bot defense subscription",
		APIPath:           "/api/config/namespaces/system/bot_defense_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_detection_rule",
		CLIName:           "shape-bot-detection-rule",
		Description:       "Shape bot detection rule",
		APIPath:           "/api/config/namespaces/{namespace}/bot_detection_rules",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_bot_detection_update",
		CLIName:           "shape-bot-detection-update",
		Description:       "Shape bot detection update",
		APIPath:           "/api/config/namespaces/{namespace}/bot_detection_updates",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

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
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense",
		CLIName:           "shape-client-side-defense",
		Description:       "Shape client side defense",
		APIPath:           "/api/config/namespaces/system/client_side_defenses",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_allowed_domain",
		CLIName:           "shape-client-side-defense-allowed-domain",
		Description:       "Shape client side defense allowed domain",
		APIPath:           "/api/config/namespaces/{namespace}/allowed_domains",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_mitigated_domain",
		CLIName:           "shape-client-side-defense-mitigated-domain",
		Description:       "Shape client side defense mitigated domain",
		APIPath:           "/api/config/namespaces/{namespace}/mitigated_domains",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_protected_domain",
		CLIName:           "shape-client-side-defense-protected-domain",
		Description:       "Shape client side defense protected domain",
		APIPath:           "/api/config/namespaces/{namespace}/protected_domains",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_client_side_defense_subscription",
		CLIName:           "shape-client-side-defense-subscription",
		Description:       "Shape client side defense subscription",
		APIPath:           "/api/config/namespaces/system/client_side_defense_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_data_delivery",
		CLIName:           "shape-data-delivery",
		Description:       "Shape data delivery",
		APIPath:           "/api/config/namespaces/{namespace}/data_deliverys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
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
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_device_id",
		CLIName:           "shape-device-id",
		Description:       "Shape device ID",
		APIPath:           "/api/config/namespaces/system/device_ids",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_mobile_app_shield_subscription",
		CLIName:           "shape-mobile-app-shield-subscription",
		Description:       "Shape mobile app shield subscription",
		APIPath:           "/api/config/namespaces/system/mobile_app_shield_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_mobile_integrator_subscription",
		CLIName:           "shape-mobile-integrator-subscription",
		Description:       "Shape mobile integrator subscription",
		APIPath:           "/api/config/namespaces/system/mobile_integrator_subscriptions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_recognize",
		CLIName:           "shape-recognize",
		Description:       "Shape recognize",
		APIPath:           "/api/shape/recognize",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_safe",
		CLIName:           "shape-safe",
		Description:       "Shape safe",
		APIPath:           "/api/shape/safe",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "shape_safeap",
		CLIName:           "shape-safeap",
		Description:       "Shape safe AP",
		APIPath:           "/api/shape/safeap",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "signup",
		CLIName:           "signup",
		Description:       "Use this API to signup for F5XC service.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\none can signup to use volterra service as an individual/free account or\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nas a team account more suited for enterprise customers. \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nfor more details on what each type of account features, visit - https://console.ves.volterra.io/signup/usage_plan\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nsince signup flow includes more complex selections and passing in secure payment processing,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nwe recommend using web UI for this process https://console.ves.volterra.io/signup/start",
		APIPath:           "/api/signup",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "site",
		CLIName:           "site",
		Description:       "Site represent physical/cloud cluster of volterra processing elements. There are two types of sites currently.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n   Regional Edge (RE)\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    Regional Edge sites are network edge sites owned and operated by volterra edge cloud. RE can be used to\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    host VIPs, run API gateway or any application at network edge.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n   Customer Edge (CE)\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    Customer Edge sites are edge sites owned by customer and operated by volterra cloud. CE can be as application gateway\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    to connect applications in multiple clusters or clouds. CE can also run applications at customer premise.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n   \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n   Nginx One\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n     Nginx One sites are sites owned and operated by Nginx One SaaS Console. Nginx One site can be used to configure service\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n     discovery which allows customer to bring their NGINX inventory visibility into the core XC workspaces.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nSite configuration contains the information like\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    Site locations\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    parameters to override fleet config\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    IP Addresses to be used by automatic vip assignments for default networks\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n    etc\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n Sites are automatically created by registration mechanism. They can be modified to add location or description and they can be deleted.",
		APIPath:           "/api/config/namespaces/{namespace}/sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "site_mesh_group",
		CLIName:           "site-mesh-group",
		Description:       "F5 Distributed Cloud Service Mesh API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/site_mesh_groups",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "service_mesh",
		Domains:           []string{"service_mesh"},
	})

	Register(&ResourceType{
		Name:              "srv6_network_slice",
		CLIName:           "srv6-network-slice",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/srv6_network_slices",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "status_at_site",
		CLIName:           "status-at-site",
		Description:       "Any user configured object in F5XC Edge Cloud has a status object associated with that it.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nAn object may be created in multiple sites and therefore it is desirable to have the ability\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nto get the current status of the configured object in a given site.",
		APIPath:           "/api/config/namespaces/{namespace}/status_at_sites",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "stored_object",
		CLIName:           "stored-object",
		Description:       "Stored object management",
		APIPath:           "/api/config/namespaces/{namespace}/stored_objects",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "subnet",
		CLIName:           "subnet",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/subnets",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "subscription",
		CLIName:           "subscription",
		Description:       "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nRepresents addon subscription \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nEywa will create the schema.pbac.addon_subscription object (SUBSCRIPTION_PENDING)\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nSRE/Support team member via f5xc-support tenant changes the status of the addon_subscription object (SUBSCRIPTION_ENABLE)",
		APIPath:           "/api/usage/namespaces/{namespace}/subscriptions",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "synthetic_monitor",
		CLIName:           "synthetic-monitor",
		Description:       "Custom handler for DNS Monitor and HTTP Monitor",
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
		Name:              "tcp_loadbalancer",
		CLIName:           "tcp-loadbalancer",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tcp_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "tenant",
		CLIName:           "tenant",
		Description:       "Package for working with Tenant representation.",
		APIPath:           "/api/config/namespaces/{namespace}/tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_configuration",
		CLIName:           "tenant-configuration",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tenant_configurations",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "tenant_and_identity",
		Domains:           []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "tenant_management",
		CLIName:           "tenant-management",
		Description:       "\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nTenant profile objects are required for creating child tenant using Child Tenant Manager as part of MSP.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nTenant Profile is the template which defines the child tenant configuration properties e.g., Name, plan,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nchild tenant groups, allowed groups, log receiver, alert receiver etc.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nWhile defining tenant profile, admin can choose PBAC plan to be subscribed for child tenant, user groups\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nneeds to be created on child tenant and allowed groups which can be mapped to user groups from root MSP\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ntenant to allow access to child tenant. It also stores log and alert receiver configuration for streaming\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nlogs and sending alert notification.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nThis feature may not be enabled by default and will require subscribing to additional addon service\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n`Tenant Management` depending upon the tenant's plan",
		APIPath:           "/api/config/namespaces/{namespace}/tenant_managements",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_allowed_tenant",
		CLIName:           "tenant-management-allowed-tenant",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/allowed_tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_child_tenant",
		CLIName:           "tenant-management-child-tenant",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/child_tenants",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_child_tenant_manager",
		CLIName:           "tenant-management-child-tenant-manager",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/web/namespaces/{namespace}/child_tenant_managers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
	})

	Register(&ResourceType{
		Name:              "tenant_management_managed_tenant",
		CLIName:           "tenant-management-managed-tenant",
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
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
		Description:       "View will create following child objects.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* Virtual-host\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* API-inventory\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* App-type\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n* App-setting",
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
		Name:              "token",
		CLIName:           "token",
		Description:       "F5 Distributed Cloud Users API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tokens",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "users",
		Domains:           []string{"users"},
	})

	Register(&ResourceType{
		Name:              "topology",
		CLIName:           "topology",
		Description:       "APIs to get topology of all the resources associated/connected to a site such as other CEs (Customer Edge),\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nREs (Regional Edge), VPCs (Virtual Private Cloud) networks, etc., and the associated metrics. Relationship between\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nthe resources associated with a site is represented as a graph, where each resource/entity is represented as a node\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n(example: CE, RE, VPC, Subnet, etc.,) and their association is represented as edge (example: CE - RE connection,\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nNetwork - Subnets association, etc.,)",
		APIPath:           "/api/topology",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: true, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "tpm_api_key",
		CLIName:           "tpm-api-key",
		Description:       "F5 Distributed Cloud Bot And Threat Defense API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_api_keys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: true, Delete: false, Status: false},
		PrimaryDomain:     "bot_and_threat_defense",
		Domains:           []string{"bot_and_threat_defense"},
	})

	Register(&ResourceType{
		Name:              "tpm_category",
		CLIName:           "tpm-category",
		Description:       "F5 Distributed Cloud Bot And Threat Defense API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_categorys",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: true, Delete: false, Status: false},
		PrimaryDomain:     "bot_and_threat_defense",
		Domains:           []string{"bot_and_threat_defense"},
	})

	Register(&ResourceType{
		Name:              "tpm_manager",
		CLIName:           "tpm-manager",
		Description:       "F5 Distributed Cloud Bot And Threat Defense API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tpm_managers",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
		PrimaryDomain:     "bot_and_threat_defense",
		Domains:           []string{"bot_and_threat_defense"},
	})

	Register(&ResourceType{
		Name:              "tpm_provision",
		CLIName:           "tpm-provision",
		Description:       "TPM Provisioning APIs used to generate F5XC certificates\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nto program device TPM.",
		APIPath:           "/api/config/namespaces/system/tpm_provisions",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "trusted_ca_list",
		CLIName:           "trusted-ca-list",
		Description:       "F5 Distributed Cloud Certificates API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/trusted_ca_lists",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "certificates",
		Domains:           []string{"certificates"},
	})

	Register(&ResourceType{
		Name:              "tunnel",
		CLIName:           "tunnel",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/tunnels",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network"},
	})

	Register(&ResourceType{
		Name:              "udp_loadbalancer",
		CLIName:           "udp-loadbalancer",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/udp_loadbalancers",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
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
		Name:              "upgrade_status",
		CLIName:           "upgrade-status",
		Description:       "Upgrade status custom APIs",
		APIPath:           "/api/config/namespaces/{namespace}/upgrade_statuss",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "usage",
		CLIName:           "usage",
		Description:       "Usage plan related RPCs. Used for billing and onboarding.",
		APIPath:           "/api/web/namespaces/{namespace}/hourly_usage_details",
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

	Register(&ResourceType{
		Name:              "usage_plan",
		CLIName:           "usage-plan",
		Description:       "Usage plan information",
		APIPath:           "/api/usage/namespaces/{namespace}/plans",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "usb_policy",
		CLIName:           "usb-policy",
		Description:       "F5 Distributed Cloud Ce Management API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/usb_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "ce_management",
		Domains:           []string{"ce_management"},
	})

	Register(&ResourceType{
		Name:              "user",
		CLIName:           "user",
		Description:       "This API can be used to manage various attributes of the user like\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nrole, contact information etc.",
		APIPath:           "/api/web/namespaces/system/users",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: true, Status: false},
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
		Description:       "F5 Distributed Cloud Tenant And Identity API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/user_identifications",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "tenant_and_identity",
		Domains:           []string{"tenant_and_identity"},
	})

	Register(&ResourceType{
		Name:              "user_setting",
		CLIName:           "user-setting",
		Description:       "User settings management",
		APIPath:           "/api/web/namespaces/system/user_settings",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: false, Get: true, List: true, Update: true, Delete: true, Status: false},
	})

	Register(&ResourceType{
		Name:              "views_terraform_parameters",
		CLIName:           "views-terraform-parameters",
		Description:       "Terraform parameters view",
		APIPath:           "/api/config/namespaces/{namespace}/terraform_parameters",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "views_view_internal",
		CLIName:           "views-view-internal",
		Description:       "Internal view",
		APIPath:           "/api/config/namespaces/{namespace}/view_internals",
		SupportsNamespace: true,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "virtual_appliance",
		CLIName:           "virtual-appliance",
		Description:       "Upgrade status custom APIs",
		APIPath:           "/api/config/namespaces/system/virtual_appliances",
		SupportsNamespace: false,
		Operations:        ResourceOperations{Create: true, Get: false, List: false, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "virtual_host",
		CLIName:           "virtual-host",
		Description:       "F5 Distributed Cloud Virtual API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_hosts",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "virtual",
		Domains:           []string{"virtual"},
	})

	Register(&ResourceType{
		Name:              "virtual_k8s",
		CLIName:           "virtual-k8s",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_k8ss",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration", "sites"},
	})

	Register(&ResourceType{
		Name:              "virtual_network",
		CLIName:           "virtual-network",
		Description:       "F5 Distributed Cloud Network API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_networks",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "network",
		Domains:           []string{"network", "service_mesh"},
	})

	Register(&ResourceType{
		Name:              "virtual_site",
		CLIName:           "virtual-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/virtual_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "voltshare",
		CLIName:           "voltshare",
		Description:       "F5XC VoltShare service serves APIs for users to securing the secrets to share it with each other.",
		APIPath:           "/api/config/namespaces/{namespace}/voltshares",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "voltshare_admin_policy",
		CLIName:           "voltshare-admin-policy",
		Description:       "F5 Distributed Cloud Blindfold API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/voltshare_admin_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "blindfold",
		Domains:           []string{"blindfold"},
	})

	Register(&ResourceType{
		Name:              "voltstack_site",
		CLIName:           "voltstack-site",
		Description:       "F5 Distributed Cloud Sites API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/voltstack_sites",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "sites",
		Domains:           []string{"sites"},
	})

	Register(&ResourceType{
		Name:              "waf",
		CLIName:           "waf",
		Description:       "APIs to get monitoring information about WAF instances on virtual-host basis. \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\nIt gets data for a given virtual host based on any WAF instance attached to virtual host or route used by virtual host.",
		APIPath:           "/api/config/namespaces/{namespace}/wafs",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "waf_exclusion_policy",
		CLIName:           "waf-exclusion-policy",
		Description:       "F5 Distributed Cloud Waf API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/waf_exclusion_policys",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "waf",
		Domains:           []string{"waf"},
	})

	Register(&ResourceType{
		Name:              "waf_signatures_changelog",
		CLIName:           "waf-signatures-changelog",
		Description:       "WAF Signatures Changelog custom APIs",
		APIPath:           "/api/config/namespaces/{namespace}/waf_signatures_changelogs",
		SupportsNamespace: true,
		Operations:        ResourceOperations{Create: true, Get: true, List: true, Update: false, Delete: false, Status: false},
	})

	Register(&ResourceType{
		Name:              "was_user_token",
		CLIName:           "was-user-token",
		Description:       "WAS user token",
		APIPath:           "/api/web/namespaces/system/was_user_tokens",
		SupportsNamespace: false,
		Operations:        ReadOnlyOperations(),
	})

	Register(&ResourceType{
		Name:              "workload",
		CLIName:           "workload",
		Description:       "F5 Distributed Cloud Kubernetes API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/workloads",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes",
		Domains:           []string{"kubernetes", "kubernetes_and_orchestration"},
	})

	Register(&ResourceType{
		Name:              "workload_flavor",
		CLIName:           "workload-flavor",
		Description:       "F5 Distributed Cloud Kubernetes And Orchestration API specifications",
		APIPath:           "/api/config/namespaces/{namespace}/workload_flavors",
		SupportsNamespace: true,
		Operations:        AllOperations(),
		PrimaryDomain:     "kubernetes_and_orchestration",
		Domains:           []string{"kubernetes_and_orchestration"},
	})

}
