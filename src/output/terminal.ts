/**
 * Terminal dimensions context
 * Provides current terminal width to output formatters
 *
 * This module maintains a global terminal width state that can be accessed
 * by any output formatter without threading width through the entire
 * execution pipeline.
 */

const DEFAULT_WIDTH = 80;
const MIN_WIDTH = 40;

let terminalWidth = process.stdout.columns ?? DEFAULT_WIDTH;

/**
 * Get current terminal width
 * Returns the cached terminal width, falling back to default if unavailable
 */
export function getTerminalWidth(): number {
	return Math.max(terminalWidth, MIN_WIDTH);
}

/**
 * Set terminal width manually
 * Used by REPL to sync Ink's width state with this context
 */
export function setTerminalWidth(width: number): void {
	terminalWidth = Math.max(width, MIN_WIDTH);
}

/**
 * Initialize terminal resize listener
 * Call once at startup for non-REPL mode (REPL uses setTerminalWidth)
 */
export function initTerminalResize(): void {
	if (process.stdout.isTTY) {
		process.stdout.on("resize", () => {
			terminalWidth = process.stdout.columns ?? DEFAULT_WIDTH;
		});
	}
}

/**
 * Get raw terminal columns without minimum enforcement
 * Useful for checking actual terminal size
 */
export function getRawTerminalWidth(): number {
	return process.stdout.columns ?? DEFAULT_WIDTH;
}
