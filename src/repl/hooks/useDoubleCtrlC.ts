/**
 * useDoubleCtrlC hook - Double Ctrl+C detection for exit confirmation
 * Requires two Ctrl+C presses within 500ms window to exit
 */

import { useState, useCallback, useRef } from "react";

interface UseDoubleCtrlCOptions {
	windowMs?: number; // Time window for double press (default: 500ms)
	onFirstPress?: () => void; // Callback on first press
	onDoublePress?: () => void; // Callback on confirmed double press
}

interface UseDoubleCtrlCResult {
	handleCtrlC: () => boolean; // Returns true if this is a double press
	reset: () => void; // Reset the timer
	isWaiting: boolean; // Whether waiting for second press
}

/**
 * Hook for detecting double Ctrl+C for exit confirmation
 */
export function useDoubleCtrlC(
	options: UseDoubleCtrlCOptions = {},
): UseDoubleCtrlCResult {
	const { windowMs = 500, onFirstPress, onDoublePress } = options;

	const lastPressRef = useRef<number>(0);
	const [isWaiting, setIsWaiting] = useState(false);
	const timeoutRef = useRef<NodeJS.Timeout | null>(null);

	const reset = useCallback(() => {
		lastPressRef.current = 0;
		setIsWaiting(false);
		if (timeoutRef.current) {
			clearTimeout(timeoutRef.current);
			timeoutRef.current = null;
		}
	}, []);

	const handleCtrlC = useCallback((): boolean => {
		const now = Date.now();
		const elapsed = now - lastPressRef.current;

		if (elapsed < windowMs && lastPressRef.current !== 0) {
			// Double press detected
			reset();
			onDoublePress?.();
			return true;
		}

		// First press
		lastPressRef.current = now;
		setIsWaiting(true);
		onFirstPress?.();

		// Clear waiting state after window expires
		if (timeoutRef.current) {
			clearTimeout(timeoutRef.current);
		}
		timeoutRef.current = setTimeout(() => {
			setIsWaiting(false);
		}, windowMs);

		return false;
	}, [windowMs, onFirstPress, onDoublePress, reset]);

	return {
		handleCtrlC,
		reset,
		isWaiting,
	};
}

export default useDoubleCtrlC;
