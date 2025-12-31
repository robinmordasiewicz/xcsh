/**
 * Headless Controller
 *
 * Provides a JSON-based stdin/stdout interface for xcsh,
 * enabling AI agents like Claude Code to interact with the CLI
 * without requiring a TTY.
 */

import * as readline from "node:readline";
import { REPLSession } from "../repl/session.js";
import { executeCommand, getCommandSuggestions } from "../repl/executor.js";
import { debugProtocol, emitSessionState } from "../debug/protocol.js";
import { CLI_NAME } from "../branding/index.js";
import {
	parseInput,
	formatOutput,
	createOutputMessage,
	createPromptMessage,
	createCompletionResponse,
	createErrorMessage,
	createEventMessage,
	createExitMessage,
	type HeadlessInput,
	type HeadlessSessionState,
	type CompletionSuggestion,
} from "./protocol.js";

/**
 * HeadlessController handles JSON protocol communication
 * for AI agent interaction with xcsh
 */
export class HeadlessController {
	private session: REPLSession;
	private rl: readline.Interface | null = null;
	private running = false;

	constructor() {
		this.session = new REPLSession();
	}

	/**
	 * Initialize the headless session
	 */
	async initialize(): Promise<void> {
		await this.session.initialize();

		// Emit debug events
		debugProtocol.session("init", { mode: "headless" });
		emitSessionState(this.session);

		// Emit session state as event
		this.emitEvent(
			"session_initialized",
			this.getSessionState() as unknown as Record<string, unknown>,
		);

		// Emit explicit warning event if token validation failed
		// This provides clear feedback to AI agents about authentication issues
		if (
			this.session.isAuthenticated() &&
			!this.session.isTokenValidated() &&
			this.session.getValidationError()
		) {
			this.emitEvent("warning", {
				message: this.session.getValidationError(),
				type: "token_validation",
			});
		}
	}

	/**
	 * Get current session state for output
	 */
	private getSessionState(): HeadlessSessionState {
		const ctx = this.session.getContextPath();
		return {
			authenticated: this.session.isAuthenticated(),
			tokenValidated: this.session.isTokenValidated(),
			authSource: this.session.getAuthSource(),
			namespace: this.session.getNamespace(),
			serverUrl: this.session.getServerUrl(),
			activeProfile: this.session.getActiveProfileName(),
			context: {
				domain: ctx.domain,
				action: ctx.action,
			},
		};
	}

	/**
	 * Build prompt string based on current context
	 */
	private buildPrompt(): string {
		const ctx = this.session.getContextPath();

		if (ctx.isRoot()) {
			return `${CLI_NAME}> `;
		}

		if (ctx.isAction()) {
			return `${CLI_NAME}:${ctx.domain}/${ctx.action}> `;
		}

		return `${CLI_NAME}:${ctx.domain}> `;
	}

	/**
	 * Write output to stdout
	 */
	private write(output: ReturnType<typeof createOutputMessage>): void {
		console.log(formatOutput(output));
	}

	/**
	 * Emit an event
	 */
	private emitEvent(event: string, data: Record<string, unknown> = {}): void {
		this.write(createEventMessage(event, data));
	}

	/**
	 * Send prompt
	 */
	private sendPrompt(): void {
		this.write(createPromptMessage(this.buildPrompt()));
	}

	/**
	 * Handle a command input
	 */
	private async handleCommand(value: string): Promise<void> {
		debugProtocol.session("command_start", { command: value });

		try {
			const result = await executeCommand(value, this.session);

			// Emit command result
			if (result.output.length > 0) {
				const content = result.output.join("\n");
				// Detect format from output
				let format: "json" | "yaml" | "table" | "text" = "text";
				if (content.startsWith("{") || content.startsWith("[")) {
					format = "json";
				} else if (content.includes(": ") && content.includes("\n")) {
					format = "yaml";
				}
				this.write(createOutputMessage(content, format));
			}

			// Check for exit
			if (result.shouldExit) {
				this.running = false;
				this.write(createExitMessage(0));
				return;
			}

			// Emit context change event
			if (result.contextChanged) {
				this.emitEvent("context_changed", {
					context: {
						domain: this.session.getContextPath().domain,
						action: this.session.getContextPath().action,
					},
				});
			}

			// Emit error if present
			if (result.error) {
				this.write(createErrorMessage(result.error, 1));
			}

			debugProtocol.session("command_complete", {
				command: value,
				success: !result.error,
			});
		} catch (error) {
			const message =
				error instanceof Error ? error.message : String(error);
			this.write(createErrorMessage(message, 1));
			debugProtocol.error("command_failed", error, { command: value });
		}

		// Send prompt for next command
		this.sendPrompt();
	}

	/**
	 * Handle completion request
	 */
	private handleCompletionRequest(partial: string): void {
		debugProtocol.session("completion_request", { partial });

		const rawSuggestions = getCommandSuggestions(partial, this.session);

		const suggestions: CompletionSuggestion[] = rawSuggestions.map((s) => ({
			text: s.text,
			description: s.description,
			category: s.category,
		}));

		this.write(createCompletionResponse(suggestions));
	}

	/**
	 * Handle interrupt signal
	 */
	private handleInterrupt(): void {
		debugProtocol.session("interrupt", {});
		this.emitEvent("interrupted", {});
		this.sendPrompt();
	}

	/**
	 * Handle exit request
	 */
	private handleExit(): void {
		this.running = false;
		this.write(createExitMessage(0));
	}

	/**
	 * Process a single input message
	 */
	private async processInput(input: HeadlessInput): Promise<void> {
		switch (input.type) {
			case "command":
				if (input.value !== undefined) {
					await this.handleCommand(input.value);
				} else {
					this.write(
						createErrorMessage(
							'Missing "value" for command input',
							1,
						),
					);
					this.sendPrompt();
				}
				break;

			case "completion_request":
				this.handleCompletionRequest(input.partial ?? "");
				break;

			case "interrupt":
				this.handleInterrupt();
				break;

			case "exit":
				this.handleExit();
				break;

			default:
				this.write(
					createErrorMessage(`Unknown input type: ${input.type}`, 1),
				);
				this.sendPrompt();
		}
	}

	/**
	 * Run the headless controller
	 * Reads JSON messages from stdin, processes them, and writes responses to stdout
	 */
	async run(): Promise<void> {
		// Initialize session
		await this.initialize();

		// Set up readline interface
		this.rl = readline.createInterface({
			input: process.stdin,
			output: process.stdout,
			terminal: false,
		});

		this.running = true;

		// Send initial prompt
		this.sendPrompt();

		// Process lines
		for await (const line of this.rl) {
			if (!this.running) break;

			const trimmed = line.trim();
			if (!trimmed) continue;

			const input = parseInput(trimmed);
			if (!input) {
				this.write(
					createErrorMessage(`Invalid JSON input: ${trimmed}`, 1),
				);
				this.sendPrompt();
				continue;
			}

			await this.processInput(input);
		}

		// Clean up
		this.rl?.close();
	}

	/**
	 * Stop the controller
	 */
	stop(): void {
		this.running = false;
		this.rl?.close();
	}
}
