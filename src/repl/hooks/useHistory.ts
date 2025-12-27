/**
 * useHistory hook - Command history navigation for the REPL
 * Handles up/down arrow navigation through command history
 */

import { useState, useCallback } from "react";

interface UseHistoryOptions {
	history: string[]; // Full command history array
	onSelect?: (command: string) => void; // Callback when history item selected
}

interface UseHistoryResult {
	navigateUp: () => string | null; // Navigate to older command
	navigateDown: () => string | null; // Navigate to newer command
	reset: () => void; // Reset to not navigating
	isNavigating: boolean; // Whether currently navigating history
	currentIndex: number; // Current position (-1 = not navigating)
}

/**
 * Hook for navigating command history with up/down arrows
 */
export function useHistory(options: UseHistoryOptions): UseHistoryResult {
	const { history, onSelect } = options;

	// historyIndex: -1 means not navigating, 0 = most recent, 1 = second most recent, etc.
	const [historyIndex, setHistoryIndex] = useState(-1);

	const reset = useCallback(() => {
		setHistoryIndex(-1);
	}, []);

	const navigateUp = useCallback((): string | null => {
		if (history.length === 0) {
			return null;
		}

		const newIndex = Math.min(historyIndex + 1, history.length - 1);
		setHistoryIndex(newIndex);

		// History is stored oldest first, so we access from the end
		const command = history[history.length - 1 - newIndex];
		if (command !== undefined) {
			onSelect?.(command);
			return command;
		}
		return null;
	}, [history, historyIndex, onSelect]);

	const navigateDown = useCallback((): string | null => {
		if (historyIndex <= 0) {
			// Already at most recent or not navigating
			if (historyIndex === 0) {
				setHistoryIndex(-1);
				onSelect?.("");
				return "";
			}
			return null;
		}

		const newIndex = historyIndex - 1;
		setHistoryIndex(newIndex);

		const command = history[history.length - 1 - newIndex];
		if (command !== undefined) {
			onSelect?.(command);
			return command;
		}
		return null;
	}, [history, historyIndex, onSelect]);

	return {
		navigateUp,
		navigateDown,
		reset,
		isNavigating: historyIndex >= 0,
		currentIndex: historyIndex,
	};
}

export default useHistory;
