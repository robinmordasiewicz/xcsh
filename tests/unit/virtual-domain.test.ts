/**
 * Virtual Domain Tests
 *
 * Comprehensive verification tests for the virtual API-generated domain.
 * Tests output format handling, flag behavior, and error handling.
 */

import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import type { REPLSession } from "../../src/repl/session.js";
import type { APIClient as APIClientType } from "../../src/api/client.js";
import type { ContextPath } from "../../src/repl/context.js";
import {
	mockVirtualList,
	mockVirtualGet,
	mockVirtualEmpty,
	mockVirtualDelete,
	mockAPIErrors,
} from "./virtual-fixtures.js";

// Import executor
import { executeCommand } from "../../src/repl/executor.js";

// Store the current mock response for API client
let mockClientResponse: { data: unknown; ok: boolean; statusCode: number } | null = null;
let mockClientError: Error | null = null;

// Set up mock API response for next client call
function mockAPIResponse(data: unknown, status = 200) {
	mockClientResponse = { data, ok: status >= 200 && status < 300, statusCode: status };
	mockClientError = null;
}

// Mock error response helper
function mockAPIErrorResponse(status: number, message: string) {
	mockClientResponse = null;
	mockClientError = new Error(message);
	(mockClientError as Error & { statusCode: number }).statusCode = status;
}

// Create mock API client that uses the stored response
function createMockAPIClient(authenticated = true): APIClientType {
	return {
		isAuthenticated: () => authenticated,
		isValidated: () => true,
		get: vi.fn().mockImplementation(async () => {
			if (mockClientError) {
				throw mockClientError;
			}
			return mockClientResponse;
		}),
		post: vi.fn().mockImplementation(async () => {
			if (mockClientError) {
				throw mockClientError;
			}
			return mockClientResponse;
		}),
		put: vi.fn().mockImplementation(async () => {
			if (mockClientError) {
				throw mockClientError;
			}
			return mockClientResponse;
		}),
		delete: vi.fn().mockImplementation(async () => {
			if (mockClientError) {
				throw mockClientError;
			}
			return mockClientResponse;
		}),
	} as unknown as APIClientType;
}

// Create mock context - mutable to track domain/action changes
function createMockContext(initialDomain?: string, initialAction?: string): ContextPath {
	let domain: string | null = initialDomain ?? null;
	let action: string | null = initialAction ?? null;

	return {
		get domain() {
			return domain;
		},
		get action() {
			return action;
		},
		resource: null,
		isRoot: () => !domain,
		isAction: () => !!action,
		isDomain: () => !!domain && !action,
		setDomain: (d: string) => {
			domain = d;
		},
		setAction: (a: string) => {
			action = a;
		},
		reset: () => {
			domain = null;
			action = null;
		},
		toString: () =>
			domain ? (action ? `/${domain}/${action}` : `/${domain}`) : "/",
	} as unknown as ContextPath;
}

// Create mock history manager
function createMockHistory() {
	return {
		add: vi.fn(),
		get: vi.fn().mockReturnValue([]),
		getAll: vi.fn().mockReturnValue([]),
		clear: vi.fn(),
		length: 0,
	};
}

// Create mock session - with optional initial domain and action context
function createMockSession(
	outputFormat = "table",
	namespace = "default",
	authenticated = true,
	initialDomain?: string,
	initialAction?: string,
): REPLSession {
	const mockClient = createMockAPIClient(authenticated);
	const mockHistory = createMockHistory();
	const mockContext = createMockContext(initialDomain, initialAction);

	return {
		getOutputFormat: () => outputFormat,
		getNamespace: () => namespace,
		getAPIClient: () => mockClient,
		getContextPath: () => mockContext,
		setOutputFormat: vi.fn(),
		addToHistory: vi.fn(),
		getHistory: () => mockHistory,
	} as unknown as REPLSession;
}

describe("virtual domain output formatting", () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		vi.clearAllMocks();
	});

	describe("list command", () => {
		it("returns JSON output with --output json", async () => {
			mockAPIResponse(mockVirtualList);
			// Start in "virtual > list" action context so command executes immediately
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--output json", session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(parsed.items).toBeDefined();
			expect(parsed.items.length).toBe(3);
		});

		it("returns YAML output with --output yaml", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--output yaml", session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("items:");
			expect(yamlOutput).toContain("name:");
		});

		it("returns TSV output with --output tsv", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--output tsv", session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns table output by default", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			// Execute with a namespace flag to trigger list without specifying format
			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			const tableOutput = result.output.join("\n");
			// Table format has headers
			expect(tableOutput).toMatch(/NAMESPACE|NAME|LABELS/i);
		});

		it("returns empty array with --output none", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--output none", session);

			expect(result.error).toBeUndefined();
			// With none format, output should be minimal
			expect(
				result.output.length === 0 ||
					result.output[0] === "(no output)",
			).toBe(true);
		});

		it("respects session default format", async () => {
			mockAPIResponse(mockVirtualList);
			// Session default is json, start in virtual > list context
			const session = createMockSession(
				"json",
				"default",
				true,
				"virtual",
				"list",
			);

			// Execute with namespace flag to trigger without specifying output format
			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			// Should be valid JSON since session default is json
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
		});

		it("handles empty results", async () => {
			mockAPIResponse(mockVirtualEmpty);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("", session);

			expect(result.error).toBeUndefined();
		});
	});

	describe("get command", () => {
		it("returns JSON output with --output json", async () => {
			mockAPIResponse(mockVirtualGet);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			const result = await executeCommand("http-lb-1 --output json", session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(parsed.metadata).toBeDefined();
			expect(parsed.metadata.name).toBe("http-lb-1");
		});

		it("returns YAML output with --output yaml", async () => {
			mockAPIResponse(mockVirtualGet);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			const result = await executeCommand("http-lb-1 --output yaml", session);

			expect(result.error).toBeUndefined();
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("metadata:");
			expect(yamlOutput).toContain("name:");
		});

		it("returns table output by default", async () => {
			mockAPIResponse(mockVirtualGet);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			const result = await executeCommand("http-lb-1", session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("handles missing resource (404)", async () => {
			mockAPIErrorResponse(404, "Resource 'http-lb-missing' not found");
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			const result = await executeCommand("http-lb-missing", session);

			expect(result.error).toBeDefined();
			const output = result.output.join("\n");
			expect(output).toContain("not found");
		});

		it("requires name parameter", async () => {
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			// Use a flag command to trigger execution without a name
			const result = await executeCommand("--output json", session);

			expect(result.error).toBeDefined();
			expect(result.output.join("\n")).toContain("name");
		});
	});

	describe("delete command", () => {
		it("returns success message", async () => {
			mockAPIResponse(mockVirtualDelete);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"delete",
			);

			const result = await executeCommand("http-lb-1", session);

			expect(result.error).toBeUndefined();
		});

		it("handles missing resource (404)", async () => {
			mockAPIErrorResponse(404, "Resource 'http-lb-missing' not found");
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"delete",
			);

			const result = await executeCommand("http-lb-missing", session);

			expect(result.error).toBeDefined();
		});

		it("requires name parameter", async () => {
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"delete",
			);

			// Use a flag command to trigger execution without a name
			const result = await executeCommand("--output json", session);

			expect(result.error).toBeDefined();
			expect(result.output.join("\n")).toContain("name");
		});
	});

	describe("namespace handling", () => {
		it("uses --namespace flag when provided", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"json",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand(
				"--namespace production --output json",
				session,
			);

			expect(result.error).toBeUndefined();
			// Verify output contains expected data
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
		});

		it("falls back to session namespace", async () => {
			mockAPIResponse(mockVirtualList);
			const session = createMockSession(
				"json",
				"staging",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--output json", session);

			expect(result.error).toBeUndefined();
			// Verify output contains expected data
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
		});

		it("handles invalid namespace", async () => {
			mockAPIErrorResponse(404, "Namespace 'nonexistent' not found");
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--namespace nonexistent", session);

			expect(result.error).toBeDefined();
		});
	});

	describe("error handling", () => {
		it("formats 401 unauthorized error", async () => {
			mockAPIErrorResponse(401, mockAPIErrors.unauthorized.message);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeDefined();
			const output = result.output.join("\n");
			expect(output).toContain("Unauthorized");
		});

		it("formats 403 forbidden error", async () => {
			mockAPIErrorResponse(403, mockAPIErrors.forbidden.message);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeDefined();
			const output = result.output.join("\n");
			expect(output).toContain("Forbidden");
		});

		it("formats 404 not found error", async () => {
			mockAPIErrorResponse(404, mockAPIErrors.notFound.message);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"get",
			);

			const result = await executeCommand("http-lb-missing", session);

			expect(result.error).toBeDefined();
			const output = result.output.join("\n");
			expect(output).toContain("Not Found");
		});

		it("formats 500 server error", async () => {
			mockAPIErrorResponse(500, mockAPIErrors.serverError.message);
			const session = createMockSession(
				"table",
				"default",
				true,
				"virtual",
				"list",
			);

			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeDefined();
			const output = result.output.join("\n");
			expect(output).toContain("Internal Server Error");
		});
	});

	describe("authentication handling", () => {
		it("returns error when not connected", async () => {
			// Create session with null client but in virtual > list action context
			const mockContext = createMockContext("virtual", "list");
			const mockHistory = createMockHistory();
			const session = {
				getOutputFormat: () => "table",
				getNamespace: () => "default",
				getAPIClient: () => null, // No client = not connected
				getContextPath: () => mockContext,
				addToHistory: vi.fn(),
				getHistory: () => mockHistory,
			} as unknown as REPLSession;

			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeDefined();
			expect(result.output.join("\n")).toContain("Not connected");
		});

		it("returns error when not authenticated", async () => {
			// Not authenticated, but in virtual > list action context
			const session = createMockSession(
				"table",
				"default",
				false,
				"virtual",
				"list",
			);

			const result = await executeCommand("--namespace default", session);

			expect(result.error).toBeDefined();
			expect(result.output.join("\n")).toContain("Not authenticated");
		});
	});
});
