/**
 * Headless Mode Protocol Types
 *
 * Defines the JSON protocol for stdin/stdout communication
 * between Claude Code (or other AI agents) and xcsh.
 */

/**
 * Input message types that xcsh can receive via stdin
 */
export type HeadlessInputType =
	| "command" // Execute a command
	| "completion_request" // Request tab completions
	| "interrupt" // Cancel current operation
	| "exit"; // Exit the session

/**
 * Output message types that xcsh sends via stdout
 */
export type HeadlessOutputType =
	| "output" // Command output
	| "prompt" // Ready for input
	| "completion_response" // Tab completion suggestions
	| "error" // Error message
	| "event" // Debug/lifecycle event
	| "exit"; // Session ended

/**
 * Input message structure (stdin → xcsh)
 */
export interface HeadlessInput {
	/** Message type */
	type: HeadlessInputType;
	/** Command string (for "command" type) */
	value?: string;
	/** Partial input for completion (for "completion_request" type) */
	partial?: string;
}

/**
 * Output format types
 */
export type HeadlessOutputFormat = "table" | "json" | "yaml" | "text";

/**
 * Output message structure (xcsh → stdout)
 */
export interface HeadlessOutput {
	/** Message type */
	type: HeadlessOutputType;
	/** Output content (for "output", "error" types) */
	content?: string;
	/** Output format hint (for "output" type) */
	format?: HeadlessOutputFormat;
	/** Prompt string (for "prompt" type) */
	prompt?: string;
	/** Completion suggestions (for "completion_response" type) */
	suggestions?: CompletionSuggestion[];
	/** Error message (for "error" type) */
	message?: string;
	/** Exit/error code */
	code?: number;
	/** Event name (for "event" type) */
	event?: string;
	/** Event data (for "event" type) */
	data?: Record<string, unknown>;
	/** Timestamp */
	timestamp?: string;
}

/**
 * Completion suggestion structure
 */
export interface CompletionSuggestion {
	/** The suggestion text */
	text: string;
	/** Description of the suggestion */
	description?: string;
	/** Category (domain, action, builtin, etc.) */
	category?: string;
}

/**
 * Authentication source type for headless output
 */
export type HeadlessAuthSource =
	| "env"
	| "profile"
	| "mixed"
	| "profile-fallback"
	| "none";

/**
 * Session state for headless output
 */
export interface HeadlessSessionState {
	/** Whether authenticated to F5 XC */
	authenticated: boolean;
	/** Whether token has been validated */
	tokenValidated: boolean;
	/** Authentication source (env, profile, mixed, profile-fallback, none) */
	authSource: HeadlessAuthSource;
	/** Current namespace */
	namespace: string;
	/** Server URL */
	serverUrl: string;
	/** Active profile name */
	activeProfile: string | null;
	/** Current context path */
	context: {
		domain: string | null;
		action: string | null;
	};
}

/**
 * Parse a JSON input line into HeadlessInput
 */
export function parseInput(line: string): HeadlessInput | null {
	try {
		const parsed = JSON.parse(line) as unknown;

		// Validate structure
		if (typeof parsed !== "object" || parsed === null) {
			return null;
		}

		const obj = parsed as Record<string, unknown>;

		if (typeof obj.type !== "string") {
			return null;
		}

		// Validate type
		const validTypes: HeadlessInputType[] = [
			"command",
			"completion_request",
			"interrupt",
			"exit",
		];
		if (!validTypes.includes(obj.type as HeadlessInputType)) {
			return null;
		}

		const result: HeadlessInput = {
			type: obj.type as HeadlessInputType,
		};
		if (typeof obj.value === "string") {
			result.value = obj.value;
		}
		if (typeof obj.partial === "string") {
			result.partial = obj.partial;
		}
		return result;
	} catch {
		return null;
	}
}

/**
 * Format an output message as JSON
 */
export function formatOutput(output: HeadlessOutput): string {
	return JSON.stringify({
		...output,
		timestamp: output.timestamp ?? new Date().toISOString(),
	});
}

/**
 * Create a command output message
 */
export function createOutputMessage(
	content: string,
	format: HeadlessOutputFormat = "text",
): HeadlessOutput {
	return {
		type: "output",
		content,
		format,
		timestamp: new Date().toISOString(),
	};
}

/**
 * Create a prompt message
 */
export function createPromptMessage(prompt: string): HeadlessOutput {
	return {
		type: "prompt",
		prompt,
		timestamp: new Date().toISOString(),
	};
}

/**
 * Create a completion response message
 */
export function createCompletionResponse(
	suggestions: CompletionSuggestion[],
): HeadlessOutput {
	return {
		type: "completion_response",
		suggestions,
		timestamp: new Date().toISOString(),
	};
}

/**
 * Create an error message
 */
export function createErrorMessage(
	message: string,
	code: number = 1,
): HeadlessOutput {
	return {
		type: "error",
		message,
		code,
		timestamp: new Date().toISOString(),
	};
}

/**
 * Create an event message
 */
export function createEventMessage(
	event: string,
	data: Record<string, unknown> = {},
): HeadlessOutput {
	return {
		type: "event",
		event,
		data,
		timestamp: new Date().toISOString(),
	};
}

/**
 * Create an exit message
 */
export function createExitMessage(code: number = 0): HeadlessOutput {
	return {
		type: "exit",
		code,
		timestamp: new Date().toISOString(),
	};
}
