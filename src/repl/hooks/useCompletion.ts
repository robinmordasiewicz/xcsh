/**
 * useCompletion hook - Tab completion state management
 * Manages suggestion visibility, selection, and filtering
 */

import { useState, useCallback, useEffect } from "react";
import type { CompletionSuggestion } from "../completion/types.js";
import { Completer } from "../completion/completer.js";
import type { REPLSession } from "../session.js";

interface UseCompletionOptions {
	session: REPLSession | null;
}

interface UseCompletionResult {
	suggestions: CompletionSuggestion[];
	selectedIndex: number;
	isShowing: boolean;
	triggerCompletion: (input: string) => Promise<CompletionSuggestion | null>; // Returns suggestion if single match (auto-complete)
	navigateUp: () => void;
	navigateDown: () => void;
	selectCurrent: () => CompletionSuggestion | null;
	hide: () => void;
	filterSuggestions: (input: string) => Promise<void>;
}

/**
 * Hook for managing tab completion state
 */
export function useCompletion(
	options: UseCompletionOptions,
): UseCompletionResult {
	const { session } = options;

	const [completer] = useState(() => {
		const c = new Completer();
		if (session) {
			c.setSession(session);
		}
		return c;
	});

	const [suggestions, setSuggestions] = useState<CompletionSuggestion[]>([]);
	const [selectedIndex, setSelectedIndex] = useState(0);
	const [isShowing, setIsShowing] = useState(false);

	// Update completer when session changes
	useEffect(() => {
		if (session) {
			completer.setSession(session);
		}
	}, [session, completer]);

	const hide = useCallback(() => {
		setIsShowing(false);
		setSuggestions([]);
		setSelectedIndex(0);
	}, []);

	const triggerCompletion = useCallback(
		async (input: string): Promise<CompletionSuggestion | null> => {
			const newSuggestions = await completer.complete(input);

			if (newSuggestions.length === 1) {
				// Single match - return it for auto-complete
				hide();
				return newSuggestions[0] ?? null;
			}

			if (newSuggestions.length > 0) {
				// Multiple matches - show suggestions
				setSuggestions(newSuggestions);
				setSelectedIndex(0);
				setIsShowing(true);
			} else {
				hide();
			}

			return null;
		},
		[completer, hide],
	);

	const filterSuggestions = useCallback(
		async (input: string): Promise<void> => {
			if (!isShowing) return;

			const newSuggestions = await completer.complete(input);
			if (newSuggestions.length === 0) {
				hide();
			} else {
				setSuggestions(newSuggestions);
				setSelectedIndex(0);
			}
		},
		[completer, isShowing, hide],
	);

	const navigateUp = useCallback(() => {
		if (!isShowing || suggestions.length === 0) return;

		setSelectedIndex((current) =>
			current > 0 ? current - 1 : suggestions.length - 1,
		);
	}, [isShowing, suggestions.length]);

	const navigateDown = useCallback(() => {
		if (!isShowing || suggestions.length === 0) return;

		setSelectedIndex((current) =>
			current < suggestions.length - 1 ? current + 1 : 0,
		);
	}, [isShowing, suggestions.length]);

	const selectCurrent = useCallback((): CompletionSuggestion | null => {
		if (!isShowing || suggestions.length === 0) return null;

		const selected = suggestions.at(selectedIndex);
		hide();
		return selected ?? null;
	}, [isShowing, suggestions, selectedIndex, hide]);

	return {
		suggestions,
		selectedIndex,
		isShowing,
		triggerCompletion,
		navigateUp,
		navigateDown,
		selectCurrent,
		hide,
		filterSuggestions,
	};
}

export default useCompletion;
