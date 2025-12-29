/**
 * Completion Domain - Generate shell completion scripts
 *
 * Generate completion scripts for bash, zsh, and fish shells.
 * Used by the build process to create completion files for distribution.
 */

import type { DomainDefinition, CommandDefinition } from "../registry.js";
import { successResult, errorResult } from "../registry.js";
import {
	generateBashCompletion,
	generateZshCompletion,
	generateFishCompletion,
} from "./generators.js";

/**
 * Bash completion command
 */
const bashCommand: CommandDefinition = {
	name: "bash",
	description: "Generate bash completion script",

	async execute() {
		try {
			const script = generateBashCompletion();
			return successResult([script]);
		} catch (error) {
			return errorResult(
				`Failed to generate bash completion: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};

/**
 * Zsh completion command
 */
const zshCommand: CommandDefinition = {
	name: "zsh",
	description: "Generate zsh completion script",

	async execute() {
		try {
			const script = generateZshCompletion();
			return successResult([script]);
		} catch (error) {
			return errorResult(
				`Failed to generate zsh completion: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};

/**
 * Fish completion command
 */
const fishCommand: CommandDefinition = {
	name: "fish",
	description: "Generate fish completion script",

	async execute() {
		try {
			const script = generateFishCompletion();
			return successResult([script]);
		} catch (error) {
			return errorResult(
				`Failed to generate fish completion: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},
};

/**
 * Completion domain definition
 */
export const completionDomain: DomainDefinition = {
	name: "completion",
	description:
		"Generate shell completion scripts for bash, zsh, and fish shells. Enables tab-completion for xcsh commands, domains, actions, and flags in your preferred shell environment.",
	descriptionShort: "Shell completion script generation",
	descriptionMedium:
		"Generate tab-completion scripts for bash, zsh, and fish shells to enhance the xcsh command-line experience.",
	commands: new Map([
		["bash", bashCommand],
		["zsh", zshCommand],
		["fish", fishCommand],
	]),
	subcommands: new Map(),
};
