/**
 * Completion system type definitions
 */

/**
 * A completion suggestion
 */
export interface CompletionSuggestion {
	text: string;
	description: string;
	category?:
		| "domain"
		| "action"
		| "flag"
		| "builtin"
		| "navigation"
		| "value"
		| "argument"
		| "subcommand"
		| "command"
		| "resource"
		| "resource-name";
}

/**
 * Completion context information
 */
export interface CompletionContext {
	input: string;
	cursorPosition: number;
	contextDomain: string;
	contextAction: string;
}

/**
 * Result of parsing input for completion
 */
export interface ParsedInput {
	args: string[];
	currentWord: string;
	isEscapedToRoot: boolean;
	isCompletingFlag: boolean;
	isCompletingFlagValue: boolean;
	currentFlag: string | null;
}
