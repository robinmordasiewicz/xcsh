import { describe, it, expect, beforeEach } from "vitest";
import { Completer } from "../../src/repl/completion/completer.js";

describe("Completer with trailing spaces", () => {
	let completer: Completer;

	beforeEach(() => {
		completer = new Completer();
	});

	it("should return profile subcommand for '/login ' (trailing space)", async () => {
		const suggestions = await completer.complete("/login ");
		const texts = suggestions.map((s) => s.text);
		expect(texts).toContain("profile");
	});

	it("should return profile subcommand for '/login p' (partial match)", async () => {
		const suggestions = await completer.complete("/login p");
		const texts = suggestions.map((s) => s.text);
		expect(texts).toContain("profile");
		expect(texts).toHaveLength(1); // Only profile matches 'p'
	});

	it("should return profile commands for '/login profile ' (trailing space)", async () => {
		const suggestions = await completer.complete("/login profile ");
		const texts = suggestions.map((s) => s.text);
		expect(texts).toContain("list");
		expect(texts).toContain("show");
		expect(texts).toContain("create");
		expect(texts).toContain("use");
		expect(texts).toContain("delete");
	});
});
