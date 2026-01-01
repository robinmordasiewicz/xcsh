/**
 * Output Test Fixtures
 * Shared test data for output format testing
 */

/**
 * Standard F5 XC resource for testing
 */
export const standardResource = {
	name: "test-resource",
	namespace: "default",
	status: "ACTIVE",
	created: "2024-01-15T10:30:00Z",
};

/**
 * Resource list for table testing
 */
export const resourceList = [
	{
		name: "resource-1",
		namespace: "default",
		status: "ACTIVE",
		labels: { env: "prod", team: "platform" },
	},
	{
		name: "resource-2",
		namespace: "staging",
		status: "PENDING",
		labels: { env: "staging" },
	},
	{
		name: "resource-3",
		namespace: "default",
		status: "ERROR",
		labels: {},
	},
];

/**
 * Single resource for single-item tests
 */
export const singleResource = {
	name: "single-item",
	namespace: "system",
	status: "ACTIVE",
	labels: { component: "core" },
};

/**
 * Data with items wrapper (API response format)
 */
export const wrappedItems = {
	items: resourceList,
	metadata: {
		total: 3,
		page: 1,
	},
};

/**
 * Unicode test data with emoji and CJK characters
 */
export const unicodeData = {
	name: "Test 🚀",
	status: "✅ Active",
	japanese: "日本語テスト",
	chinese: "中文测试",
	labels: { emoji: "🎉", type: "unicode" },
};

/**
 * Unicode resource list for table testing
 */
export const unicodeResourceList = [
	{ name: "Resource 🚀", namespace: "default", status: "✅ Active", labels: {} },
	{
		name: "日本語リソース",
		namespace: "system",
		status: "⏳ Pending",
		labels: {},
	},
	{ name: "中文资源", namespace: "staging", status: "❌ Error", labels: {} },
];

/**
 * Special characters in values
 */
export const specialCharsData = {
	name: "path/with/slashes",
	path: '/path/with spaces/and"quotes',
	regex: "pattern.*[a-z]+",
	json: '{"nested": "value"}',
};

/**
 * Empty and null data variations
 */
export const emptyData = {
	emptyArray: [],
	emptyObject: {},
	emptyItems: { items: [] },
	nullValue: null,
	undefinedField: { field: undefined },
};

/**
 * Nested data for JSON/YAML testing
 */
export const nestedData = {
	level1: {
		level2: {
			level3: {
				level4: {
					value: "deeply nested",
				},
			},
		},
	},
	array: [{ nested: true }, { also: "nested" }],
	mixed: {
		string: "value",
		number: 42,
		boolean: true,
		null: null,
		array: [1, 2, 3],
	},
};

/**
 * Large dataset generator for performance testing
 */
export function generateLargeDataset(
	count: number,
): Array<Record<string, unknown>> {
	return Array.from({ length: count }, (_, i) => ({
		name: `resource-${i + 1}`,
		namespace: i % 3 === 0 ? "default" : i % 3 === 1 ? "staging" : "prod",
		status: i % 4 === 0 ? "ACTIVE" : i % 4 === 1 ? "PENDING" : "ERROR",
		created: new Date(Date.now() - i * 86400000).toISOString(),
		labels: { index: String(i), batch: String(Math.floor(i / 10)) },
	}));
}

/**
 * Wide dataset with many columns
 */
export function generateWideDataset(
	columns: number,
): Array<Record<string, unknown>> {
	const row: Record<string, unknown> = {};
	for (let i = 0; i < columns; i++) {
		row[`column_${i + 1}`] = `value_${i + 1}`;
	}
	return [row, { ...row }, { ...row }];
}

/**
 * Long text for wrapping tests
 */
export const longTextData = {
	name: "resource-with-very-long-name-that-should-wrap",
	description:
		"This is a very long description that should definitely wrap across multiple lines when displayed in a table with reasonable column widths.",
	labels: {
		"very-long-label-key-name": "very-long-label-value-content",
		another: "label",
	},
};

/**
 * API error response for error formatting tests
 */
export const apiErrorResponses = {
	unauthorized: {
		statusCode: 401,
		body: { message: "Invalid API token", code: "UNAUTHORIZED" },
		operation: "list resources",
	},
	forbidden: {
		statusCode: 403,
		body: { message: "Access denied to namespace", code: "FORBIDDEN" },
		operation: "get resource",
	},
	notFound: {
		statusCode: 404,
		body: { message: "Resource not found", code: "NOT_FOUND" },
		operation: "get resource",
	},
	conflict: {
		statusCode: 409,
		body: { message: "Resource already exists", code: "CONFLICT" },
		operation: "create resource",
	},
	rateLimit: {
		statusCode: 429,
		body: { message: "Rate limit exceeded", code: "RATE_LIMITED" },
		operation: "list resources",
	},
	serverError: {
		statusCode: 500,
		body: { message: "Internal server error", details: "Database timeout" },
		operation: "update resource",
	},
};

/**
 * Key-value data for formatKeyValueBox tests
 */
export const keyValueData = [
	{ label: "User", value: "admin@example.com" },
	{ label: "Tenant", value: "acme-corp" },
	{ label: "Namespace", value: "default" },
	{ label: "API URL", value: "https://api.example.com" },
];

/**
 * Column definitions for custom table tests
 */
export const customColumns = [
	{ header: "ID", accessor: "id", width: 10 },
	{ header: "NAME", accessor: "name", minWidth: 15, maxWidth: 30 },
	{ header: "STATUS", accessor: "status", align: "center" as const },
	{ header: "COUNT", accessor: "count", align: "right" as const },
];

/**
 * Data for custom column tests
 */
export const customColumnData = [
	{ id: "001", name: "First Resource", status: "Active", count: 42 },
	{ id: "002", name: "Second Resource", status: "Pending", count: 17 },
	{ id: "003", name: "Third Resource with Long Name", status: "Error", count: 0 },
];
