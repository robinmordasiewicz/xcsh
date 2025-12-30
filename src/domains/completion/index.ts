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
	description:
		"Generate a bash shell completion script for xcsh. Output the script to stdout for manual installation or pipe to a file. Enables tab-completion for commands, domains, actions, and option flags.",
	descriptionShort: "Generate bash completion script",
	descriptionMedium:
		"Output bash completion script for tab-completion of commands, domains, and flags.",

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
	description:
		"Generate a zsh shell completion script for xcsh. Output the script to stdout for manual installation or add to your fpath. Provides rich tab-completion with descriptions for commands, domains, and options.",
	descriptionShort: "Generate zsh completion script",
	descriptionMedium:
		"Output zsh completion script with rich descriptions for commands and options.",

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
	description:
		"Generate a fish shell completion script for xcsh. Output the script to stdout for manual installation or save to your fish completions directory. Enables intelligent tab-completion with inline descriptions.",
	descriptionShort: "Generate fish completion script",
	descriptionMedium:
		"Output fish completion script with inline descriptions for commands and options.",

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
