/**
 * Debug Protocol for F5 XC Shell
 *
 * Provides structured event logging in JSONL format for AI/PTY debugging.
 * Events are emitted to stderr so they don't interfere with normal output.
 *
 * Usage:
 *   F5XC_DEBUG_EVENTS=jsonl xcsh
 *
 * Output format (one JSON object per line):
 *   {"timestamp":"2024-01-01T00:00:00.000Z","type":"session","event":"init","data":{...}}
 */

import { ENV_PREFIX } from "../branding/index.js";

/**
 * Debug event types
 */
export type DebugEventType =
	| "session" // Session lifecycle events
	| "api" // API call events
	| "ui" // UI rendering events
	| "error" // Error events
	| "profile" // Profile loading events
	| "auth"; // Authentication/validation events

/**
 * Debug event structure
 */
export interface DebugEvent {
	timestamp: string;
	type: DebugEventType;
	event: string;
	data: Record<string, unknown>;
}

/**
 * Debug event format options
 */
export type DebugFormat = "jsonl" | "human" | "none";

/**
 * Get debug format from environment
 */
function getDebugFormat(): DebugFormat {
	const envValue = process.env[`${ENV_PREFIX}_DEBUG_EVENTS`];
	if (envValue === "jsonl" || envValue === "human") {
		return envValue;
	}
	// Legacy F5XC_DEBUG=true maps to human format
	if (process.env[`${ENV_PREFIX}_DEBUG`] === "true") {
		return "human";
	}
	return "none";
}

/**
 * Debug Protocol singleton for structured event logging
 */
class DebugProtocolImpl {
	private readonly format: DebugFormat;
	private readonly events: DebugEvent[] = [];
	private readonly startTime: number;

	constructor() {
		this.format = getDebugFormat();
		this.startTime = Date.now();
	}

	/**
	 * Check if debug events are enabled
	 */
	isEnabled(): boolean {
		return this.format !== "none";
	}

	/**
	 * Check if JSONL format is enabled
	 */
	isJsonl(): boolean {
		return this.format === "jsonl";
	}

	/**
	 * Emit a debug event
	 */
	emit(
		type: DebugEventType,
		event: string,
		data: Record<string, unknown> = {},
	): void {
		if (!this.isEnabled()) return;

		const entry: DebugEvent = {
			timestamp: new Date().toISOString(),
			type,
			event,
			data: {
				...data,
				elapsedMs: Date.now() - this.startTime,
			},
		};

		this.events.push(entry);

		if (this.format === "jsonl") {
			// JSONL format - one JSON object per line
			console.error(JSON.stringify(entry));
		} else if (this.format === "human") {
			// Human-readable format
			const prefix = `DEBUG [${type}:${event}]`;
			const dataStr =
				Object.keys(data).length > 0 ? ` ${JSON.stringify(data)}` : "";
			console.error(`${prefix}${dataStr}`);
		}
	}

	/**
	 * Emit session-related event
	 */
	session(event: string, data: Record<string, unknown> = {}): void {
		this.emit("session", event, data);
	}

	/**
	 * Emit API-related event
	 */
	api(event: string, data: Record<string, unknown> = {}): void {
		this.emit("api", event, data);
	}

	/**
	 * Emit authentication-related event
	 */
	auth(event: string, data: Record<string, unknown> = {}): void {
		this.emit("auth", event, data);
	}

	/**
	 * Emit profile-related event
	 */
	profile(event: string, data: Record<string, unknown> = {}): void {
		this.emit("profile", event, data);
	}

	/**
	 * Emit UI-related event
	 */
	ui(event: string, data: Record<string, unknown> = {}): void {
		this.emit("ui", event, data);
	}

	/**
	 * Emit error event
	 */
	error(
		event: string,
		error: unknown,
		data: Record<string, unknown> = {},
	): void {
		this.emit("error", event, {
			...data,
			error: error instanceof Error ? error.message : String(error),
			stack: error instanceof Error ? error.stack : undefined,
		});
	}

	/**
	 * Get all captured events
	 */
	getEvents(): DebugEvent[] {
		return [...this.events];
	}

	/**
	 * Dump session state on exit (for --dump-state-on-exit)
	 */
	dumpState(state: Record<string, unknown>): void {
		if (!this.isEnabled()) return;

		console.error("\n=== Session State Dump ===");
		console.error(JSON.stringify(state, null, 2));
		console.error("=== End Session State ===\n");
	}
}

/**
 * Global debug protocol instance
 */
export const debugProtocol = new DebugProtocolImpl();

/**
 * Convenience function to emit session state after initialization
 */
export function emitSessionState(session: {
	isAuthenticated: () => boolean;
	isTokenValidated: () => boolean;
	getValidationError: () => string | null;
	getServerUrl: () => string;
	getActiveProfileName: () => string | null;
	getAPIClient: () => unknown | null;
}): void {
	debugProtocol.auth("session_state", {
		isAuthenticated: session.isAuthenticated(),
		isTokenValidated: session.isTokenValidated(),
		validationError: session.getValidationError(),
		serverUrl: session.getServerUrl(),
		activeProfile: session.getActiveProfileName(),
		hasApiClient: session.getAPIClient() !== null,
	});

	// Also emit the warning condition check
	const showsWarning =
		session.isAuthenticated() &&
		!session.isTokenValidated() &&
		session.getValidationError();

	debugProtocol.auth("warning_check", {
		shouldShowWarning: !!showsWarning,
		isAuthenticated: session.isAuthenticated(),
		isTokenValidated: session.isTokenValidated(),
		hasValidationError: !!session.getValidationError(),
	});
}
