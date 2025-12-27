/**
 * HistoryManager handles command history persistence.
 * Stores commands to disk and provides navigation through history.
 */

import { readFileSync, writeFileSync, mkdirSync, existsSync } from "node:fs";
import { homedir } from "node:os";
import { join, dirname } from "node:path";

/**
 * Get the default history file path
 */
export function getHistoryFilePath(): string {
	try {
		return join(homedir(), ".xcsh_history");
	} catch {
		return ".xcsh_history";
	}
}

/**
 * HistoryManager handles command history persistence
 */
export class HistoryManager {
	private path: string;
	private maxSize: number;
	private history: string[];

	constructor(path: string, maxSize: number = 1000) {
		this.path = path;
		this.maxSize = maxSize;
		this.history = [];
	}

	/**
	 * Create a new history manager and load existing history
	 */
	static async create(
		path?: string,
		maxSize: number = 1000,
	): Promise<HistoryManager> {
		const historyPath = path ?? getHistoryFilePath();
		const manager = new HistoryManager(historyPath, maxSize);
		await manager.load();
		return manager;
	}

	/**
	 * Load reads history from the history file
	 */
	async load(): Promise<void> {
		try {
			if (!existsSync(this.path)) {
				return;
			}

			const content = readFileSync(this.path, "utf-8");
			const lines = content
				.split("\n")
				.filter((line) => line.trim() !== "");
			this.history = lines;

			// Trim to max size
			if (this.history.length > this.maxSize) {
				this.history = this.history.slice(
					this.history.length - this.maxSize,
				);
			}
		} catch (error) {
			// Ignore file not found errors, they're expected on first run
			if ((error as NodeJS.ErrnoException).code !== "ENOENT") {
				console.error("Warning: could not load history:", error);
			}
		}
	}

	/**
	 * Save writes history to the history file
	 */
	async save(): Promise<void> {
		try {
			// Ensure directory exists
			const dir = dirname(this.path);
			if (dir !== "." && dir !== "") {
				mkdirSync(dir, { recursive: true });
			}

			const content = this.history.join("\n") + "\n";
			writeFileSync(this.path, content, "utf-8");
		} catch (error) {
			console.error("Warning: could not save history:", error);
		}
	}

	/**
	 * Add a command to history
	 */
	add(cmd: string): void {
		// Don't add empty commands
		if (cmd.trim() === "") {
			return;
		}

		// Don't add duplicates of the last command
		if (
			this.history.length > 0 &&
			this.history[this.history.length - 1] === cmd
		) {
			return;
		}

		this.history.push(cmd);

		// Trim if necessary
		if (this.history.length > this.maxSize) {
			this.history.shift();
		}
	}

	/**
	 * Get all history entries
	 */
	getHistory(): string[] {
		return [...this.history];
	}

	/**
	 * Get the number of history entries
	 */
	get length(): number {
		return this.history.length;
	}

	/**
	 * Get a specific history entry by index (0 = oldest)
	 */
	get(index: number): string | undefined {
		return this.history.at(index);
	}

	/**
	 * Get the most recent history entry
	 */
	getLast(): string | undefined {
		return this.history[this.history.length - 1];
	}

	/**
	 * Search history for entries containing the query
	 */
	search(query: string): string[] {
		const lowerQuery = query.toLowerCase();
		return this.history.filter((entry) =>
			entry.toLowerCase().includes(lowerQuery),
		);
	}

	/**
	 * Clear all history
	 */
	clear(): void {
		this.history = [];
	}
}
