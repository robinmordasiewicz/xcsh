/**
 * Virtual Domain Test Fixtures
 *
 * Mock API response data for testing the virtual domain commands.
 * Based on F5 XC API response format for virtual load balancer resources.
 */

/**
 * Mock HTTP Load Balancer list response
 */
export const mockVirtualList = {
	items: [
		{
			namespace: "default",
			name: "http-lb-1",
			labels: { app: "web", env: "prod" },
			metadata: {
				name: "http-lb-1",
				namespace: "default",
				uid: "uid-001",
				creation_timestamp: "2024-01-15T10:00:00Z",
			},
		},
		{
			namespace: "default",
			name: "http-lb-2",
			labels: { app: "api" },
			metadata: {
				name: "http-lb-2",
				namespace: "default",
				uid: "uid-002",
				creation_timestamp: "2024-01-16T11:00:00Z",
			},
		},
		{
			namespace: "default",
			name: "tcp-lb-1",
			labels: {},
			metadata: {
				name: "tcp-lb-1",
				namespace: "default",
				uid: "uid-003",
				creation_timestamp: "2024-01-17T12:00:00Z",
			},
		},
	],
};

/**
 * Mock HTTP Load Balancer get response (single resource)
 */
export const mockVirtualGet = {
	metadata: {
		name: "http-lb-1",
		namespace: "default",
		uid: "uid-001",
		creation_timestamp: "2024-01-15T10:00:00Z",
		labels: { app: "web", env: "prod" },
	},
	spec: {
		domains: ["example.com", "www.example.com"],
		http: {
			dns_volterra_managed: true,
		},
		default_route_pools: [
			{
				pool: {
					name: "origin-pool-1",
					namespace: "default",
				},
				weight: 1,
			},
		],
		advertise_on_public_default_vip: true,
	},
	system_metadata: {
		tenant: "test-tenant",
		creator_id: "user-123",
	},
};

/**
 * Mock empty list response
 */
export const mockVirtualEmpty = {
	items: [],
};

/**
 * Mock delete response
 */
export const mockVirtualDelete = {
	name: "http-lb-1",
	namespace: "default",
};

/**
 * Mock status response
 */
export const mockVirtualStatus = {
	status: {
		state: "ACTIVE",
		vip_info: {
			vip: "203.0.113.10",
			dns_name: "http-lb-1.example.com",
		},
		health_status: "HEALTHY",
	},
};

/**
 * Mock API error responses
 */
export const mockAPIErrors = {
	unauthorized: {
		code: 401,
		message: "Unauthorized: Invalid or expired token",
		details: [],
	},
	forbidden: {
		code: 403,
		message: "Forbidden: Insufficient permissions to access resource",
		details: [],
	},
	notFound: {
		code: 404,
		message: "Not Found: Resource 'http-lb-missing' does not exist",
		details: [],
	},
	conflict: {
		code: 409,
		message: "Conflict: Resource 'http-lb-1' already exists",
		details: [],
	},
	serverError: {
		code: 500,
		message: "Internal Server Error",
		details: [],
	},
};

/**
 * Mock namespace list response
 */
export const mockNamespaceList = {
	items: [
		{ name: "default", uid: "ns-001" },
		{ name: "production", uid: "ns-002" },
		{ name: "staging", uid: "ns-003" },
	],
};

/**
 * Create mock API response helper
 * @param data - Response data
 * @param status - HTTP status code
 */
export function createMockResponse(data: unknown, status = 200) {
	return {
		ok: status >= 200 && status < 300,
		status,
		statusText: status === 200 ? "OK" : "Error",
		text: () => Promise.resolve(JSON.stringify(data)),
		json: () => Promise.resolve(data),
		headers: new Headers({ "content-type": "application/json" }),
	};
}

/**
 * Create mock error response
 * @param status - HTTP status code
 * @param message - Error message
 */
export function createMockErrorResponse(status: number, message: string) {
	const errorData = {
		code: status,
		message,
		details: [],
	};

	return {
		ok: false,
		status,
		statusText: "Error",
		text: () => Promise.resolve(JSON.stringify(errorData)),
		json: () => Promise.resolve(errorData),
		headers: new Headers({ "content-type": "application/json" }),
	};
}
