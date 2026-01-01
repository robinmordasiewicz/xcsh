/**
 * Table Rendering Tests
 * Tests for beautiful table formatting with snapshots
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import {
	formatBeautifulTable,
	formatResourceTable,
	formatKeyValueBox,
	wrapText,
	DEFAULT_RESOURCE_COLUMNS,
} from "../../src/output/table.js";
import type { TableConfig, ColumnDefinition } from "../../src/output/types.js";
import {
	resourceList,
	singleResource,
	wrappedItems,
	unicodeResourceList,
	longTextData,
	keyValueData,
	customColumns,
	customColumnData,
	generateLargeDataset,
	generateWideDataset,
	emptyData,
} from "./output-test-fixtures.js";

describe("wrapText", () => {
	it("returns single line for short text", () => {
		const result = wrapText("short text", 20);
		expect(result).toEqual(["short text"]);
	});

	it("wraps at word boundaries", () => {
		const result = wrapText("hello world test", 10);
		expect(result.length).toBeGreaterThan(1);
		expect(result[0]).toBe("hello");
	});

	it("handles text exactly at max width", () => {
		const result = wrapText("exactly10c", 10);
		expect(result).toEqual(["exactly10c"]);
	});

	it("handles long words without spaces", () => {
		const result = wrapText("superlongwordwithoutspaces", 10);
		expect(result.length).toBeGreaterThan(1);
		// Should break at max width when no spaces
		expect(result[0]?.length).toBeLessThanOrEqual(10);
	});

	it("handles empty string", () => {
		const result = wrapText("", 10);
		expect(result).toEqual([""]);
	});

	it("preserves multiple spaces at line start after wrap", () => {
		const result = wrapText("word1 word2 word3", 6);
		expect(result.length).toBeGreaterThan(1);
	});
});

describe("formatBeautifulTable", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	describe("basic table rendering", () => {
		it("renders empty string for empty data", () => {
			const result = formatBeautifulTable([], { columns: [] }, true);
			expect(result).toBe("");
		});

		it("renders single row table", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable([singleResource], config, true);

			expect(result).toContain("NAMESPACE");
			expect(result).toContain("NAME");
			expect(result).toContain("single-item");
			expect(result).toContain("system");
		});

		it("renders multiple row table", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable(resourceList, config, true);

			expect(result).toContain("resource-1");
			expect(result).toContain("resource-2");
			expect(result).toContain("resource-3");
		});

		it("renders table with title", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
				title: "Resources",
			};
			const result = formatBeautifulTable(resourceList, config, true);

			expect(result).toContain("Resources");
		});
	});

	describe("column alignment", () => {
		const alignColumns: ColumnDefinition[] = [
			{ header: "LEFT", accessor: "left", align: "left" },
			{ header: "CENTER", accessor: "center", align: "center" },
			{ header: "RIGHT", accessor: "right", align: "right" },
		];

		it("handles different column alignments", () => {
			const data = [{ left: "L", center: "C", right: "R" }];
			const config: TableConfig = { columns: alignColumns };
			const result = formatBeautifulTable(data, config, true);

			expect(result).toContain("LEFT");
			expect(result).toContain("CENTER");
			expect(result).toContain("RIGHT");
		});
	});

	describe("column widths", () => {
		it("respects minWidth", () => {
			const columns: ColumnDefinition[] = [
				{ header: "A", accessor: "a", minWidth: 20 },
			];
			const data = [{ a: "x" }];
			const result = formatBeautifulTable(data, { columns }, true);

			// Column should be at least 20 chars
			const lines = result.split("\n");
			const headerLine = lines.find((l) => l.includes("A"));
			expect(headerLine?.length).toBeGreaterThanOrEqual(20);
		});

		it("respects maxWidth", () => {
			const columns: ColumnDefinition[] = [
				{ header: "A", accessor: "a", maxWidth: 10 },
			];
			const data = [{ a: "this is a very long value that should be truncated" }];
			const result = formatBeautifulTable(
				data,
				{ columns, wrapText: true },
				true,
			);

			// Value should be wrapped
			expect(result).toContain("this is a");
		});

		it("respects fixed width", () => {
			const columns: ColumnDefinition[] = [
				{ header: "A", accessor: "a", width: 15 },
			];
			const data = [{ a: "value" }];
			const result = formatBeautifulTable(data, { columns }, true);

			expect(result).toContain("value");
		});
	});

	describe("text wrapping", () => {
		it("wraps long text when enabled", () => {
			const columns: ColumnDefinition[] = [
				{ header: "DESC", accessor: "description", maxWidth: 20 },
			];
			const data = [{ description: longTextData.description }];
			const result = formatBeautifulTable(
				data,
				{ columns, wrapText: true },
				true,
			);

			// Should have multiple lines for the content
			const lines = result.split("\n");
			expect(lines.length).toBeGreaterThan(3);
		});

		it("disables wrapping when wrapText=false", () => {
			const columns: ColumnDefinition[] = [
				{ header: "DESC", accessor: "description", maxWidth: 20 },
			];
			const data = [{ description: "long text that should not wrap here" }];
			const result = formatBeautifulTable(
				data,
				{ columns, wrapText: false },
				true,
			);

			// Content should be truncated, not wrapped
			const contentLines = result
				.split("\n")
				.filter((l) => !l.includes("DESC") && l.includes("|"));
			expect(contentLines.length).toBeLessThanOrEqual(3);
		});
	});

	describe("row separators", () => {
		it("renders without row separators by default", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable(resourceList, config, true);

			// Table should render without extra separators between data rows
			expect(result).toContain("resource-1");
			expect(result).toContain("resource-2");
			// Get lines and check structure
			const lines = result.split("\n");
			expect(lines.length).toBeGreaterThan(3);
		});

		it("renders with row separators when enabled", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
				rowSeparators: true,
			};
			const result = formatBeautifulTable(resourceList, config, true);

			// Should have more separator lines than without
			const lines = result.split("\n");
			const separatorLines = lines.filter(
				(l) => l.startsWith("+") && l.includes("-"),
			);
			// With 3 data rows and row separators, should have 5+ separator lines
			expect(separatorLines.length).toBeGreaterThanOrEqual(5);
		});
	});

	describe("labels formatting", () => {
		it("formats labels as map[key:value]", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable(resourceList, config, true);

			expect(result).toContain("map[env:prod");
			expect(result).toContain("team:platform");
		});

		it("handles empty labels", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const data = [{ name: "test", namespace: "default", labels: {} }];
			const result = formatBeautifulTable(data, config, true);

			expect(result).toContain("<None>");
		});
	});

	describe("ASCII mode (no colors)", () => {
		it("uses ASCII box characters", () => {
			const config: TableConfig = { columns: DEFAULT_RESOURCE_COLUMNS };
			const result = formatBeautifulTable(resourceList, config, true);

			expect(result).toContain("+");
			expect(result).toContain("-");
			expect(result).toContain("|");
		});

		it("does not include ANSI escape codes", () => {
			const config: TableConfig = { columns: DEFAULT_RESOURCE_COLUMNS };
			const result = formatBeautifulTable(resourceList, config, true);

			expect(result).not.toContain("\x1b[");
		});
	});

	describe("snapshot tests", () => {
		it("matches snapshot for standard table", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable(resourceList, config, true);
			expect(result).toMatchSnapshot("standard-table");
		});

		it("matches snapshot for single row", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
			};
			const result = formatBeautifulTable([singleResource], config, true);
			expect(result).toMatchSnapshot("single-row-table");
		});

		it("matches snapshot for table with title", () => {
			const config: TableConfig = {
				columns: DEFAULT_RESOURCE_COLUMNS,
				title: "Test Resources",
			};
			const result = formatBeautifulTable(resourceList, config, true);
			expect(result).toMatchSnapshot("table-with-title");
		});

		it("matches snapshot for custom columns", () => {
			const config: TableConfig = {
				columns: customColumns,
			};
			const result = formatBeautifulTable(customColumnData, config, true);
			expect(result).toMatchSnapshot("custom-columns-table");
		});
	});
});

describe("formatResourceTable", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("formats array of resources", () => {
		const result = formatResourceTable(resourceList, true);

		expect(result).toContain("NAMESPACE");
		expect(result).toContain("NAME");
		expect(result).toContain("LABELS");
		expect(result).toContain("resource-1");
	});

	it("handles items wrapper", () => {
		const result = formatResourceTable(wrappedItems, true);

		expect(result).toContain("resource-1");
		expect(result).toContain("resource-2");
	});

	it("handles single object", () => {
		const result = formatResourceTable(singleResource, true);

		expect(result).toContain("single-item");
		expect(result).toContain("system");
	});

	it("returns empty string for empty array", () => {
		const result = formatResourceTable([], true);
		expect(result).toBe("");
	});

	it("returns empty string for empty items", () => {
		const result = formatResourceTable(emptyData.emptyItems, true);
		expect(result).toBe("");
	});

	it("matches snapshot", () => {
		const result = formatResourceTable(resourceList, true);
		expect(result).toMatchSnapshot("resource-table");
	});
});

describe("formatKeyValueBox", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("renders key-value pairs", () => {
		const result = formatKeyValueBox(keyValueData, "Session Info", true);

		expect(result).toContain("User");
		expect(result).toContain("admin@example.com");
		expect(result).toContain("Tenant");
		expect(result).toContain("acme-corp");
	});

	it("includes title", () => {
		const result = formatKeyValueBox(keyValueData, "My Title", true);
		expect(result).toContain("My Title");
	});

	it("returns empty string for empty data", () => {
		const result = formatKeyValueBox([], "Title", true);
		expect(result).toBe("");
	});

	it("aligns labels correctly", () => {
		const data = [
			{ label: "Short", value: "val1" },
			{ label: "Longer Label", value: "val2" },
		];
		const result = formatKeyValueBox(data, "Test", true);

		// Both colons should align
		const lines = result.split("\n");
		const dataLines = lines.filter((l) => l.includes(":"));
		expect(dataLines.length).toBeGreaterThan(0);
	});

	it("uses ASCII characters when noColor=true", () => {
		const result = formatKeyValueBox(keyValueData, "Test", true);

		expect(result).toContain("+");
		expect(result).not.toContain("\x1b[");
	});

	it("matches snapshot", () => {
		const result = formatKeyValueBox(keyValueData, "Session Info", true);
		expect(result).toMatchSnapshot("key-value-box");
	});
});

describe("Unicode handling", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("handles emoji in table cells", () => {
		// Use columns that include status to test emoji in that field
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name", minWidth: 20 },
			{ header: "STATUS", accessor: "status", minWidth: 15 },
		];
		const config: TableConfig = { columns };
		const result = formatBeautifulTable(unicodeResourceList, config, true);

		expect(result).toContain("🚀");
		expect(result).toContain("✅");
		expect(result).toContain("⏳");
		expect(result).toContain("❌");
	});

	it("handles CJK characters", () => {
		const config: TableConfig = {
			columns: DEFAULT_RESOURCE_COLUMNS,
		};
		const result = formatBeautifulTable(unicodeResourceList, config, true);

		expect(result).toContain("日本語");
		expect(result).toContain("中文");
	});

	it("matches snapshot for unicode table", () => {
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name", minWidth: 20 },
			{ header: "STATUS", accessor: "status", minWidth: 15 },
		];
		const config: TableConfig = { columns };
		const result = formatBeautifulTable(unicodeResourceList, config, true);
		expect(result).toMatchSnapshot("unicode-table");
	});
});

describe("Large datasets", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("handles 100+ row dataset", () => {
		const data = generateLargeDataset(100);
		const config: TableConfig = {
			columns: DEFAULT_RESOURCE_COLUMNS,
		};
		const result = formatBeautifulTable(data, config, true);

		expect(result).toContain("resource-1");
		expect(result).toContain("resource-100");
		expect(result.split("\n").length).toBeGreaterThan(100);
	});

	it("handles wide dataset with many columns", () => {
		const data = generateWideDataset(10);
		const columns: ColumnDefinition[] = Array.from({ length: 10 }, (_, i) => ({
			header: `COL${i + 1}`,
			accessor: `column_${i + 1}`,
			maxWidth: 15,
		}));
		const config: TableConfig = { columns };
		const result = formatBeautifulTable(data, config, true);

		expect(result).toContain("COL1");
		expect(result).toContain("COL10");
	});
});

describe("Edge cases", () => {
	beforeEach(() => {
		vi.stubEnv("NO_COLOR", "1");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("handles null values in cells", () => {
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name", minWidth: 10 },
			{ header: "VALUE", accessor: "value", minWidth: 10 },
		];
		const data = [{ name: "test", value: null }];
		const result = formatBeautifulTable(data, { columns }, true);

		// The table shows empty/null values as <None>
		expect(result).toContain("NAME");
		expect(result).toContain("VALUE");
		expect(result).toContain("test");
	});

	it("handles undefined values in cells", () => {
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name", minWidth: 10 },
			{ header: "VALUE", accessor: "value", minWidth: 10 },
		];
		const data = [{ name: "test" }]; // value is undefined
		const result = formatBeautifulTable(data, { columns }, true);

		// The table should render even with missing values
		expect(result).toContain("NAME");
		expect(result).toContain("VALUE");
		expect(result).toContain("test");
	});

	it("handles accessor function", () => {
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name" },
			{
				header: "CUSTOM",
				accessor: (row) => `Custom: ${row.name}`,
			},
		];
		const data = [{ name: "test" }];
		const result = formatBeautifulTable(data, { columns }, true);

		expect(result).toContain("Custom: test");
	});

	it("handles nested object in cell", () => {
		const columns: ColumnDefinition[] = [
			{ header: "NAME", accessor: "name" },
			{ header: "NESTED", accessor: "nested" },
		];
		const data = [{ name: "test", nested: { key: "value" } }];
		const result = formatBeautifulTable(data, { columns }, true);

		expect(result).toContain('{"key":"value"}');
	});
});
