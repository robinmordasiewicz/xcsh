import { describe, it, expect } from "vitest";
import { parseInput } from "../../src/repl/completion/completer.js";

describe("parseInput with trailing spaces", () => {
	it("should set empty currentWord when input ends with space", () => {
		const result = parseInput("/login ");
		expect(result.currentWord).toBe("");
		expect(result.args).toEqual(["login"]);
	});

	it("should set partial word when no trailing space", () => {
		const result = parseInput("/login p");
		expect(result.currentWord).toBe("p");
		expect(result.args).toEqual(["login", "p"]);
	});

	it("should handle subcommand with trailing space", () => {
		const result = parseInput("/login profile ");
		expect(result.currentWord).toBe("");
		expect(result.args).toEqual(["login", "profile"]);
	});
});
