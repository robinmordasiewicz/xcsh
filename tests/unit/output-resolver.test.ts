/**
 * Output Resolver Tests
 * Tests for format resolution and flag parsing
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import {
	resolveOutputFormat,
	getOutputFormatFromEnv,
	parseOutputFlag,
	parseSpecFlag,
	parseOutputFlags,
	shouldUseColors,
	buildOutputContext,
	OUTPUT_FORMAT_ENV_VAR,
} from "../../src/output/resolver.js";
import type { OutputContext } from "../../src/output/types.js";

describe("resolveOutputFormat", () => {
	describe("precedence order", () => {
		it("CLI flag takes highest precedence", () => {
			const context: OutputContext = {
				cliFormat: "json",
				envFormat: "yaml",
				configFormat: "tsv",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("json");
		});

		it("env format is second priority", () => {
			const context: OutputContext = {
				envFormat: "yaml",
				configFormat: "tsv",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("yaml");
		});

		it("config format is third priority", () => {
			const context: OutputContext = {
				configFormat: "tsv",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("tsv");
		});

		it("defaults to table when no format specified", () => {
			const context: OutputContext = {
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("table");
		});
	});

	describe("format types", () => {
		it("resolves json format", () => {
			const context: OutputContext = {
				cliFormat: "json",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("json");
		});

		it("resolves yaml format", () => {
			const context: OutputContext = {
				cliFormat: "yaml",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("yaml");
		});

		it("resolves table format", () => {
			const context: OutputContext = {
				cliFormat: "table",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("table");
		});

		it("resolves tsv format", () => {
			const context: OutputContext = {
				cliFormat: "tsv",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("tsv");
		});

		it("resolves none format", () => {
			const context: OutputContext = {
				cliFormat: "none",
				isInteractive: false,
				isTTY: true,
			};
			expect(resolveOutputFormat(context)).toBe("none");
		});
	});
});

describe("getOutputFormatFromEnv", () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("returns undefined when env var not set", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "");
		expect(getOutputFormatFromEnv()).toBeUndefined();
	});

	it("parses json format", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "json");
		expect(getOutputFormatFromEnv()).toBe("json");
	});

	it("parses yaml format", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "yaml");
		expect(getOutputFormatFromEnv()).toBe("yaml");
	});

	it("handles uppercase format", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "JSON");
		expect(getOutputFormatFromEnv()).toBe("json");
	});

	it("handles format with whitespace", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "  json  ");
		expect(getOutputFormatFromEnv()).toBe("json");
	});

	it("returns undefined for invalid format", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "invalid");
		expect(getOutputFormatFromEnv()).toBeUndefined();
	});
});

describe("parseOutputFlag", () => {
	describe("long form --output", () => {
		it("parses --output json", () => {
			const result = parseOutputFlag(["--output", "json"]);
			expect(result.format).toBe("json");
			expect(result.remainingArgs).toEqual([]);
		});

		it("parses --output yaml", () => {
			const result = parseOutputFlag(["--output", "yaml"]);
			expect(result.format).toBe("yaml");
		});

		it("parses --output=json (equals form)", () => {
			const result = parseOutputFlag(["--output=json"]);
			expect(result.format).toBe("json");
			expect(result.remainingArgs).toEqual([]);
		});

		it("preserves other arguments", () => {
			const result = parseOutputFlag([
				"--name",
				"test",
				"--output",
				"json",
				"--namespace",
				"default",
			]);
			expect(result.format).toBe("json");
			expect(result.remainingArgs).toEqual([
				"--name",
				"test",
				"--namespace",
				"default",
			]);
		});
	});

	describe("short form -o", () => {
		it("parses -o json", () => {
			const result = parseOutputFlag(["-o", "json"]);
			expect(result.format).toBe("json");
			expect(result.remainingArgs).toEqual([]);
		});

		it("parses -o=yaml (equals form)", () => {
			const result = parseOutputFlag(["-o=yaml"]);
			expect(result.format).toBe("yaml");
		});
	});

	describe("edge cases", () => {
		it("returns undefined when no output flag", () => {
			const result = parseOutputFlag(["--name", "test"]);
			expect(result.format).toBeUndefined();
			expect(result.remainingArgs).toEqual(["--name", "test"]);
		});

		it("handles --output without value", () => {
			const result = parseOutputFlag(["--output"]);
			expect(result.format).toBeUndefined();
		});

		it("handles --output followed by another flag", () => {
			const result = parseOutputFlag(["--output", "--name", "test"]);
			expect(result.format).toBeUndefined();
			expect(result.remainingArgs).toContain("--name");
		});

		it("ignores invalid format values", () => {
			const result = parseOutputFlag(["--output", "invalid"]);
			expect(result.format).toBeUndefined();
		});

		it("handles empty array", () => {
			const result = parseOutputFlag([]);
			expect(result.format).toBeUndefined();
			expect(result.remainingArgs).toEqual([]);
		});

		it("handles case insensitivity", () => {
			const result = parseOutputFlag(["--output", "JSON"]);
			expect(result.format).toBe("json");
		});
	});
});

describe("parseSpecFlag", () => {
	it("detects --spec flag", () => {
		const result = parseSpecFlag(["--spec"]);
		expect(result.spec).toBe(true);
		expect(result.remainingArgs).toEqual([]);
	});

	it("removes --spec from remaining args", () => {
		const result = parseSpecFlag(["--name", "test", "--spec", "--namespace", "default"]);
		expect(result.spec).toBe(true);
		expect(result.remainingArgs).toEqual(["--name", "test", "--namespace", "default"]);
	});

	it("returns false when no --spec", () => {
		const result = parseSpecFlag(["--name", "test"]);
		expect(result.spec).toBe(false);
		expect(result.remainingArgs).toEqual(["--name", "test"]);
	});

	it("handles empty array", () => {
		const result = parseSpecFlag([]);
		expect(result.spec).toBe(false);
		expect(result.remainingArgs).toEqual([]);
	});
});

describe("parseOutputFlags", () => {
	it("parses both --output and --spec", () => {
		const result = parseOutputFlags(["--output", "json", "--spec"]);
		expect(result.format).toBe("json");
		expect(result.spec).toBe(true);
		expect(result.remainingArgs).toEqual([]);
	});

	it("handles --output only", () => {
		const result = parseOutputFlags(["--output", "yaml"]);
		expect(result.format).toBe("yaml");
		expect(result.spec).toBe(false);
	});

	it("handles --spec only", () => {
		const result = parseOutputFlags(["--spec"]);
		expect(result.format).toBeUndefined();
		expect(result.spec).toBe(true);
	});

	it("preserves other arguments", () => {
		const result = parseOutputFlags([
			"--name",
			"test",
			"--output",
			"json",
			"--spec",
			"--ns",
			"default",
		]);
		expect(result.format).toBe("json");
		expect(result.spec).toBe(true);
		expect(result.remainingArgs).toEqual(["--name", "test", "--ns", "default"]);
	});
});

describe("shouldUseColors", () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	describe("noColor flag", () => {
		it("returns false when noColorFlag is true", () => {
			expect(shouldUseColors(true, true)).toBe(false);
		});

		it("ignores TTY when noColorFlag is true", () => {
			expect(shouldUseColors(true, true)).toBe(false);
		});
	});

	describe("NO_COLOR environment variable", () => {
		it("returns false when NO_COLOR is set", () => {
			vi.stubEnv("NO_COLOR", "1");
			expect(shouldUseColors(true, false)).toBe(false);
		});

		it("returns false when NO_COLOR is empty string", () => {
			vi.stubEnv("NO_COLOR", "");
			expect(shouldUseColors(true, false)).toBe(false);
		});
	});

	describe("FORCE_COLOR environment variable", () => {
		it("returns true when FORCE_COLOR is set", () => {
			vi.stubEnv("FORCE_COLOR", "1");
			expect(shouldUseColors(false, false)).toBe(true);
		});

		it("FORCE_COLOR overrides non-TTY", () => {
			vi.stubEnv("FORCE_COLOR", "1");
			expect(shouldUseColors(false, false)).toBe(true);
		});
	});

	describe("TTY detection", () => {
		it("returns true for TTY without flags", () => {
			expect(shouldUseColors(true, false)).toBe(true);
		});

		it("returns false for non-TTY without FORCE_COLOR", () => {
			expect(shouldUseColors(false, false)).toBe(false);
		});
	});

	describe("precedence", () => {
		it("noColorFlag takes precedence over FORCE_COLOR", () => {
			vi.stubEnv("FORCE_COLOR", "1");
			expect(shouldUseColors(true, true)).toBe(false);
		});

		it("NO_COLOR takes precedence over FORCE_COLOR", () => {
			vi.stubEnv("NO_COLOR", "1");
			vi.stubEnv("FORCE_COLOR", "1");
			expect(shouldUseColors(true, false)).toBe(false);
		});
	});
});

describe("buildOutputContext", () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("builds context with CLI format", () => {
		const context = buildOutputContext({ cliFormat: "json" });
		expect(context.cliFormat).toBe("json");
	});

	it("builds context with config format", () => {
		const context = buildOutputContext({ configFormat: "yaml" });
		expect(context.configFormat).toBe("yaml");
	});

	it("includes env format when set", () => {
		vi.stubEnv(OUTPUT_FORMAT_ENV_VAR, "tsv");
		const context = buildOutputContext();
		expect(context.envFormat).toBe("tsv");
	});

	it("sets isInteractive", () => {
		const context = buildOutputContext({ isInteractive: true });
		expect(context.isInteractive).toBe(true);
	});

	it("sets noColor", () => {
		const context = buildOutputContext({ noColor: true });
		expect(context.noColor).toBe(true);
	});

	it("defaults isInteractive to false", () => {
		const context = buildOutputContext();
		expect(context.isInteractive).toBe(false);
	});

	it("does not include undefined optional properties", () => {
		const context = buildOutputContext();
		expect("cliFormat" in context).toBe(false);
		expect("configFormat" in context).toBe(false);
	});

	it("includes isTTY from process.stdout", () => {
		const context = buildOutputContext();
		expect(typeof context.isTTY).toBe("boolean");
	});
});
