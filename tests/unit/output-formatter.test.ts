/**
 * Output Formatter Tests
 * Tests for main formatter routing and format-specific functions
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import {
	formatOutput,
	formatJSON,
	formatYAML,
	formatTSV,
	formatAPIError,
	parseOutputFormat,
} from "../../src/output/formatter.js";
import {
	standardResource,
	resourceList,
	wrappedItems,
	nestedData,
	emptyData,
	specialCharsData,
	apiErrorResponses,
	unicodeData,
} from "./output-test-fixtures.js";

describe("formatOutput", () => {
	beforeEach(() => {
		// Mock TTY detection for consistent color behavior
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	describe("format routing", () => {
		it("routes to JSON format", () => {
			const result = formatOutput(standardResource, "json");
			expect(result).toBe(JSON.stringify(standardResource, null, 2));
		});

		it("routes to YAML format", () => {
			const result = formatOutput(standardResource, "yaml");
			expect(result).toContain("name: test-resource");
			expect(result).toContain("namespace: default");
		});

		it("routes to table format", () => {
			const result = formatOutput(resourceList, "table", true);
			expect(result).toContain("NAME");
			expect(result).toContain("resource-1");
		});

		it("routes text to table format", () => {
			const result = formatOutput(resourceList, "text", true);
			expect(result).toContain("NAME");
		});

		it("routes to TSV format", () => {
			const result = formatOutput(resourceList, "tsv");
			expect(result).toContain("\t");
			expect(result).toContain("resource-1");
		});

		it("returns empty string for none format", () => {
			const result = formatOutput(standardResource, "none");
			expect(result).toBe("");
		});

		it("routes spec format to JSON", () => {
			const result = formatOutput(standardResource, "spec");
			expect(result).toBe(JSON.stringify(standardResource, null, 2));
		});

		it("defaults to table format", () => {
			const result = formatOutput(resourceList, undefined, true);
			expect(result).toContain("NAME");
		});
	});

	describe("noColor parameter", () => {
		it("respects noColor=true", () => {
			const result = formatOutput(resourceList, "table", true);
			// Should use ASCII box characters when colors disabled
			expect(result).toContain("+");
			expect(result).not.toContain("\x1b[");
		});
	});
});

describe("formatJSON", () => {
	it("formats standard object with 2-space indent", () => {
		const result = formatJSON(standardResource);
		expect(result).toBe(JSON.stringify(standardResource, null, 2));
	});

	it("formats nested objects", () => {
		const result = formatJSON(nestedData);
		expect(result).toContain("level1");
		expect(result).toContain("deeply nested");
	});

	it("formats arrays", () => {
		const result = formatJSON(resourceList);
		expect(result.startsWith("[")).toBe(true);
		expect(result.endsWith("]")).toBe(true);
	});

	it("handles null value", () => {
		const result = formatJSON(null);
		expect(result).toBe("null");
	});

	it("handles empty array", () => {
		const result = formatJSON([]);
		expect(result).toBe("[]");
	});

	it("handles empty object", () => {
		const result = formatJSON({});
		expect(result).toBe("{}");
	});

	it("preserves special characters in strings", () => {
		const result = formatJSON(specialCharsData);
		expect(result).toContain("path/with/slashes");
		expect(result).toContain('"quotes"');
	});

	it("handles unicode characters", () => {
		const result = formatJSON(unicodeData);
		expect(result).toContain("🚀");
		expect(result).toContain("日本語テスト");
	});
});

describe("formatYAML", () => {
	it("formats standard object", () => {
		const result = formatYAML(standardResource);
		expect(result).toContain("name: test-resource");
		expect(result).toContain("namespace: default");
		expect(result).toContain("status: ACTIVE");
	});

	it("formats nested structures", () => {
		const result = formatYAML(nestedData);
		expect(result).toContain("level1:");
		expect(result).toContain("level2:");
		expect(result).toContain("deeply nested");
	});

	it("formats arrays", () => {
		const result = formatYAML(resourceList);
		expect(result).toContain("- name:");
	});

	it("handles null value", () => {
		const result = formatYAML(null);
		expect(result.trim()).toBe("null");
	});

	it("handles empty array", () => {
		const result = formatYAML([]);
		expect(result.trim()).toBe("[]");
	});

	it("handles empty object", () => {
		const result = formatYAML({});
		expect(result.trim()).toBe("{}");
	});

	it("handles unicode characters", () => {
		const result = formatYAML(unicodeData);
		expect(result).toContain("🚀");
		expect(result).toContain("日本語テスト");
	});

	it("escapes special YAML characters", () => {
		const data = { key: "value: with colon" };
		const result = formatYAML(data);
		// YAML should properly handle the colon in the value
		expect(result).toContain("value: with colon");
	});
});

describe("formatTSV", () => {
	it("formats array of objects as TSV", () => {
		const result = formatTSV(resourceList);
		const lines = result.split("\n");
		expect(lines.length).toBe(3);
		expect(lines[0]).toContain("\t");
	});

	it("prioritizes common columns (name, namespace, status)", () => {
		const result = formatTSV(resourceList);
		const firstLine = result.split("\n")[0];
		const columns = firstLine?.split("\t") ?? [];
		// Name should appear early
		expect(columns.some((col) => col?.includes("resource-1"))).toBe(true);
	});

	it("handles items wrapper", () => {
		const result = formatTSV(wrappedItems);
		expect(result).toContain("resource-1");
		expect(result.split("\n").length).toBe(3);
	});

	it("returns empty string for empty data", () => {
		expect(formatTSV(emptyData.emptyArray)).toBe("");
		expect(formatTSV(emptyData.emptyItems)).toBe("");
	});

	it("handles null values in objects", () => {
		const data = [{ name: "test", value: null }];
		const result = formatTSV(data);
		expect(result).toContain("test");
	});

	it("stringifies nested objects", () => {
		const data = [{ name: "test", nested: { key: "value" } }];
		const result = formatTSV(data);
		expect(result).toContain('{"key":"value"}');
	});

	it("handles single item", () => {
		const result = formatTSV(standardResource);
		expect(result).toContain("test-resource");
	});
});

describe("formatAPIError", () => {
	it("formats 401 unauthorized error", () => {
		const { statusCode, body, operation } = apiErrorResponses.unauthorized;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: list resources failed (HTTP 401)");
		expect(result).toContain("Invalid API token");
		expect(result).toContain("UNAUTHORIZED");
		expect(result).toContain("Authentication failed");
	});

	it("formats 403 forbidden error", () => {
		const { statusCode, body, operation } = apiErrorResponses.forbidden;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: get resource failed (HTTP 403)");
		expect(result).toContain("Permission denied");
	});

	it("formats 404 not found error", () => {
		const { statusCode, body, operation } = apiErrorResponses.notFound;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: get resource failed (HTTP 404)");
		expect(result).toContain("Resource not found");
		expect(result).toContain("Verify the name and namespace");
	});

	it("formats 409 conflict error", () => {
		const { statusCode, body, operation } = apiErrorResponses.conflict;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: create resource failed (HTTP 409)");
		expect(result).toContain("Conflict");
	});

	it("formats 429 rate limit error", () => {
		const { statusCode, body, operation } = apiErrorResponses.rateLimit;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: list resources failed (HTTP 429)");
		expect(result).toContain("Rate limited");
	});

	it("formats 500 server error", () => {
		const { statusCode, body, operation } = apiErrorResponses.serverError;
		const result = formatAPIError(statusCode, body, operation);

		expect(result).toContain("ERROR: update resource failed (HTTP 500)");
		expect(result).toContain("Server error");
		expect(result).toContain("Database timeout");
	});

	it("handles error without body", () => {
		const result = formatAPIError(500, null, "test operation");
		expect(result).toContain("ERROR: test operation failed (HTTP 500)");
	});

	it("handles error with only message", () => {
		const result = formatAPIError(400, { message: "Bad request" }, "test");
		expect(result).toContain("Bad request");
	});
});

describe("parseOutputFormat", () => {
	it("parses json format", () => {
		expect(parseOutputFormat("json")).toBe("json");
		expect(parseOutputFormat("JSON")).toBe("json");
	});

	it("parses yaml format", () => {
		expect(parseOutputFormat("yaml")).toBe("yaml");
		expect(parseOutputFormat("YAML")).toBe("yaml");
	});

	it("parses table format", () => {
		expect(parseOutputFormat("table")).toBe("table");
		expect(parseOutputFormat("text")).toBe("table");
		expect(parseOutputFormat("")).toBe("table");
	});

	it("parses tsv format", () => {
		expect(parseOutputFormat("tsv")).toBe("tsv");
		expect(parseOutputFormat("TSV")).toBe("tsv");
	});

	it("parses none format", () => {
		expect(parseOutputFormat("none")).toBe("none");
	});

	it("returns table for invalid format", () => {
		expect(parseOutputFormat("invalid")).toBe("table");
		expect(parseOutputFormat("xml")).toBe("table");
	});
});
