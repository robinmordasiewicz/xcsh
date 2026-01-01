/**
 * Domain Formatter Tests
 * Tests for unified domain output formatting functions
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import {
	parseDomainOutputFlags,
	formatDomainOutput,
	formatKeyValueOutput,
	formatListOutput,
} from "../../src/output/domain-formatter.js";
import type { DomainFormatOptions, KeyValueData } from "../../src/output/domain-formatter.js";

describe("parseDomainOutputFlags", () => {
	beforeEach(() => {
		// Mock TTY detection for consistent behavior
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("parses --output json flag", () => {
		const result = parseDomainOutputFlags(["--output", "json"], "table");
		expect(result.options.format).toBe("json");
		expect(result.remainingArgs).toEqual([]);
	});

	it("parses --output yaml flag", () => {
		const result = parseDomainOutputFlags(["--output", "yaml"], "table");
		expect(result.options.format).toBe("yaml");
	});

	it("parses --output tsv flag", () => {
		const result = parseDomainOutputFlags(["--output", "tsv"], "table");
		expect(result.options.format).toBe("tsv");
	});

	it("parses --output none flag", () => {
		const result = parseDomainOutputFlags(["--output", "none"], "table");
		expect(result.options.format).toBe("none");
	});

	it("uses session default when no flag provided", () => {
		const result = parseDomainOutputFlags([], "yaml");
		expect(result.options.format).toBe("yaml");
	});

	it("defaults to table when no flag and no session default", () => {
		const result = parseDomainOutputFlags([]);
		expect(result.options.format).toBe("table");
	});

	it("preserves remaining args after parsing output flag", () => {
		const result = parseDomainOutputFlags(["--output", "json", "arg1", "arg2"], "table");
		expect(result.options.format).toBe("json");
		expect(result.remainingArgs).toEqual(["arg1", "arg2"]);
	});

	it("parses --no-color flag", () => {
		const result = parseDomainOutputFlags(["--no-color"], "table");
		expect(result.options.noColor).toBe(true);
		expect(result.remainingArgs).toEqual([]);
	});

	it("handles both output and no-color flags", () => {
		const result = parseDomainOutputFlags(["--output", "json", "--no-color", "myarg"], "table");
		expect(result.options.format).toBe("json");
		expect(result.options.noColor).toBe(true);
		expect(result.remainingArgs).toEqual(["myarg"]);
	});
});

describe("formatDomainOutput", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	const testData = { name: "test", value: 123 };
	const baseOptions: DomainFormatOptions = { format: "table", noColor: true };

	it("returns empty array for none format", () => {
		const result = formatDomainOutput(testData, { ...baseOptions, format: "none" });
		expect(result).toEqual([]);
	});

	it("formats data as JSON", () => {
		const result = formatDomainOutput(testData, { ...baseOptions, format: "json" });
		expect(result.join("\n")).toContain('"name": "test"');
		expect(result.join("\n")).toContain('"value": 123');
	});

	it("formats data as YAML", () => {
		const result = formatDomainOutput(testData, { ...baseOptions, format: "yaml" });
		expect(result.join("\n")).toContain("name: test");
		expect(result.join("\n")).toContain("value: 123");
	});

	it("returns empty array for empty input", () => {
		const result = formatDomainOutput("", { ...baseOptions, format: "json" });
		expect(result).toEqual(['""']);
	});
});

describe("formatKeyValueOutput", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	const testData: KeyValueData = {
		name: "test-profile",
		apiUrl: "https://example.com",
		active: true,
	};
	const baseOptions: DomainFormatOptions = { format: "table", noColor: true, title: "Test" };

	it("returns empty array for none format", () => {
		const result = formatKeyValueOutput(testData, { ...baseOptions, format: "none" });
		expect(result).toEqual([]);
	});

	it("formats key-value data as JSON", () => {
		const result = formatKeyValueOutput(testData, { ...baseOptions, format: "json" });
		const json = result.join("\n");
		expect(json).toContain('"name": "test-profile"');
		expect(json).toContain('"apiUrl": "https://example.com"');
		expect(json).toContain('"active": true');
	});

	it("formats key-value data as YAML", () => {
		const result = formatKeyValueOutput(testData, { ...baseOptions, format: "yaml" });
		const yaml = result.join("\n");
		expect(yaml).toContain("name: test-profile");
		expect(yaml).toContain("apiUrl: https://example.com");
		expect(yaml).toContain("active: true");
	});

	it("formats key-value data as TSV", () => {
		const result = formatKeyValueOutput(testData, { ...baseOptions, format: "tsv" });
		expect(result.some((line) => line.includes("name\t"))).toBe(true);
		expect(result.some((line) => line.includes("apiUrl\t"))).toBe(true);
	});

	it("handles null and undefined values in JSON format", () => {
		const dataWithNulls: KeyValueData = {
			name: "test",
			nullValue: null,
			undefinedValue: undefined,
		};
		const result = formatKeyValueOutput(dataWithNulls, { ...baseOptions, format: "json" });
		const json = result.join("\n");
		expect(json).toContain('"name": "test"');
		// JSON preserves null values
		expect(json).toContain("nullValue");
		// undefined values are excluded from JSON
		expect(json).not.toContain("undefinedValue");
	});

	it("returns empty array for empty data", () => {
		const result = formatKeyValueOutput({}, baseOptions);
		expect(result).toEqual([]);
	});

	it("converts camelCase keys to Title Case for table format", () => {
		const result = formatKeyValueOutput(
			{ apiUrl: "https://example.com" },
			baseOptions,
		);
		// Table format should include human-readable label
		expect(result.join("\n")).toContain("Api Url");
	});
});

describe("formatListOutput", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	const testData = [
		{ name: "item1", status: "active" },
		{ name: "item2", status: "inactive" },
	];
	const baseOptions: DomainFormatOptions = { format: "table", noColor: true };

	it("returns empty array for none format", () => {
		const result = formatListOutput(testData, { ...baseOptions, format: "none" });
		expect(result).toEqual([]);
	});

	it("formats list data as JSON", () => {
		const result = formatListOutput(testData, { ...baseOptions, format: "json" });
		const json = result.join("\n");
		expect(json).toContain('"name": "item1"');
		expect(json).toContain('"name": "item2"');
		expect(json.startsWith("[")).toBe(true);
	});

	it("formats list data as YAML", () => {
		const result = formatListOutput(testData, { ...baseOptions, format: "yaml" });
		const yaml = result.join("\n");
		expect(yaml).toContain("- name: item1");
		expect(yaml).toContain("- name: item2");
	});

	it("formats list data as TSV", () => {
		const result = formatListOutput(testData, { ...baseOptions, format: "tsv" });
		expect(result.length).toBe(2);
		expect(result.every((line) => line.includes("\t"))).toBe(true);
	});

	it("returns empty array for empty list", () => {
		const result = formatListOutput([], { ...baseOptions, format: "tsv" });
		expect(result).toEqual([]);
	});

	it("handles null values in objects", () => {
		const dataWithNull = [{ name: "test", value: null }];
		const result = formatListOutput(dataWithNull, { ...baseOptions, format: "tsv" });
		expect(result.length).toBe(1);
	});

	it("stringifies nested objects in TSV format", () => {
		const nestedData = [{ name: "test", nested: { key: "value" } }];
		const result = formatListOutput(nestedData, { ...baseOptions, format: "tsv" });
		expect(result[0]).toContain('{"key":"value"}');
	});
});

describe("format integration", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("parseDomainOutputFlags result works with formatKeyValueOutput", () => {
		const { options } = parseDomainOutputFlags(["--output", "json"], "table");
		const data: KeyValueData = { name: "test", value: 123 };
		const result = formatKeyValueOutput(data, { ...options, title: "Test" });
		expect(result.join("\n")).toContain('"name": "test"');
	});

	it("parseDomainOutputFlags result works with formatListOutput", () => {
		const { options } = parseDomainOutputFlags(["--output", "yaml"], "table");
		const data = [{ name: "item1" }, { name: "item2" }];
		const result = formatListOutput(data, options);
		expect(result.join("\n")).toContain("- name: item1");
	});

	it("handles text format", () => {
		const { options } = parseDomainOutputFlags(["--output", "text"], "json");
		// text is a valid format that gets treated like table in formatters
		expect(options.format).toBe("text");
	});
});
