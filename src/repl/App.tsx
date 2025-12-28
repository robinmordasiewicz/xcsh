/**
 * App - Main Ink application for the REPL
 * Orchestrates all UI components and handles keyboard input
 */

import React, { useState, useCallback, useEffect, useRef } from "react";
import { Box, Text, useApp, useInput, useStdout } from "ink";

import {
	Banner,
	InputBox,
	StatusBar,
	Suggestions,
} from "./components/index.js";
import type { Suggestion } from "./components/Suggestions.js";
import { getGitInfo, type GitInfo } from "./components/StatusBar.js";
import { REPLSession } from "./session.js";
import { buildPlainPrompt } from "./prompt.js";
import { useDoubleCtrlC } from "./hooks/useDoubleCtrlC.js";
import { useHistory } from "./hooks/useHistory.js";
import { useCompletion } from "./hooks/useCompletion.js";
import { executeCommand } from "./executor.js";
import { CLI_VERSION } from "../branding/index.js";

/**
 * Convert completion suggestions to UI suggestions
 */
function toUISuggestions(
	suggestions: Array<{
		text: string;
		description: string;
		category?: string;
	}>,
): Suggestion[] {
	return suggestions.map((s) => ({
		label: s.text,
		value: s.text,
		description: s.description,
		category: s.category ?? "builtin", // Provide default to avoid undefined
	}));
}

/**
 * Main REPL Application
 */
export function App(): React.ReactElement {
	const { exit } = useApp();
	const { stdout } = useStdout();

	// Session state
	const [session] = useState(() => new REPLSession());
	const [isInitialized, setIsInitialized] = useState(false);

	// UI state
	const [input, setInputState] = useState("");
	const inputRef = useRef(""); // Ref to track current input for async callbacks

	// Custom setInput that also updates the ref synchronously
	const setInput = useCallback(
		(value: string | ((prev: string) => string)) => {
			setInputState((prev) => {
				const newValue =
					typeof value === "function" ? value(prev) : value;
				inputRef.current = newValue; // Update ref synchronously
				return newValue;
			});
		},
		[],
	);

	const [output, setOutput] = useState<string[]>([]);
	const [showBanner, setShowBanner] = useState(true);
	const [prompt, setPrompt] = useState("<xc> ");
	const [width, setWidth] = useState(stdout?.columns ?? 80);
	const [gitInfo, setGitInfo] = useState<GitInfo | undefined>(undefined);
	const [statusHint, setStatusHint] = useState("Ctrl+C twice to exit");
	const [historyArray, setHistoryArray] = useState<string[]>([]);
	const [inputKey, setInputKey] = useState(0); // Key to reset cursor position

	// Completion state
	const completion = useCompletion({
		session: isInitialized ? session : null,
	});

	// History state
	const history = useHistory({
		history: historyArray,
		onSelect: (cmd) => setInput(cmd),
	});

	// Double Ctrl+C detection
	const ctrlC = useDoubleCtrlC({
		windowMs: 500,
		onFirstPress: () => {
			addOutput("\nPress Ctrl+C again to exit, or continue typing");
			setStatusHint("Press Ctrl+C again to exit");
		},
		onDoublePress: () => {
			// Save history before exiting
			session.saveHistory().finally(() => exit());
		},
	});

	// Initialize session
	useEffect(() => {
		const init = async () => {
			await session.initialize();
			setPrompt(buildPlainPrompt(session));
			// Get initial history array
			const histMgr = session.getHistory();
			if (histMgr) {
				setHistoryArray(histMgr.getHistory());
			}
			setIsInitialized(true);
			// Get git repository info
			setGitInfo(getGitInfo());
		};
		init();
	}, [session]);

	// Handle terminal resize
	useEffect(() => {
		const handleResize = () => {
			if (stdout) {
				setWidth(stdout.columns ?? 80);
			}
		};

		stdout?.on("resize", handleResize);
		return () => {
			stdout?.off("resize", handleResize);
		};
	}, [stdout]);

	// Add output line(s)
	const addOutput = useCallback((line: string) => {
		setOutput((prev) => {
			const newLines = [...prev, ...line.split("\n")];
			// Keep buffer under 1000 lines
			if (newLines.length > 1000) {
				return newLines.slice(newLines.length - 1000);
			}
			return newLines;
		});
	}, []);

	// Apply completion to input - uses inputRef for current value
	const applyCompletion = useCallback(
		(suggestion: Suggestion) => {
			// Use ref to get current input (avoids stale closure)
			const currentInput = inputRef.current;
			let newValue: string;
			const words = currentInput.split(/\s+/);

			if (words.length === 0 || currentInput === "") {
				newValue = suggestion.value + " ";
			} else {
				const lastWord = words[words.length - 1] ?? "";

				// Handle "/" escape prefix - strip it for matching but preserve in output
				const hasEscapePrefix =
					words.length === 1 && lastWord.startsWith("/");
				const matchWord = hasEscapePrefix
					? lastWord.slice(1)
					: lastWord;

				if (
					matchWord === "" ||
					suggestion.value
						.toLowerCase()
						.startsWith(matchWord.toLowerCase())
				) {
					// Replace the partial word with the full suggestion
					const prefix = currentInput.slice(
						0,
						currentInput.length - lastWord.length,
					);
					// Preserve the "/" prefix if it was there
					const completedValue = hasEscapePrefix
						? "/" + suggestion.value
						: suggestion.value;
					newValue = prefix + completedValue + " ";
				} else {
					// No match - just append with space (shouldn't normally happen)
					newValue = currentInput + " " + suggestion.value + " ";
				}
			}

			setInput(newValue);
			// Increment key to force TextInput remount, moving cursor to end
			setInputKey((k) => k + 1);
		},
		[], // No dependencies needed since we use inputRef
	);

	// Refresh history array from session
	const refreshHistory = useCallback(() => {
		const histMgr = session.getHistory();
		if (histMgr) {
			setHistoryArray(histMgr.getHistory());
		}
	}, [session]);

	// Execute command using executor module
	const runCommand = useCallback(
		async (cmd: string) => {
			const trimmed = cmd.trim();
			if (!trimmed) return;

			// Show command in output
			addOutput(prompt + trimmed);

			// Hide banner after first command
			setShowBanner(false);

			// Execute via executor module
			const result = await executeCommand(trimmed, session);

			// Handle clear
			if (result.shouldClear) {
				setOutput([]);
			} else {
				// Add output lines
				result.output.forEach((line) => addOutput(line));
			}

			// Handle exit
			if (result.shouldExit) {
				await session.saveHistory();
				exit();
				return;
			}

			// Update prompt if context changed
			if (result.contextChanged) {
				setPrompt(buildPlainPrompt(session));
			}

			// Refresh history
			refreshHistory();
		},
		[session, prompt, addOutput, exit, refreshHistory],
	);

	// Handle input change
	const handleInputChange = useCallback(
		(newValue: string) => {
			const oldValue = input;
			setInput(newValue);

			// Live filtering when suggestions are showing
			if (completion.isShowing) {
				completion.filterSuggestions(newValue);
			}

			// Check if "/" was typed - trigger completion
			if (newValue !== oldValue && newValue.endsWith("/")) {
				const beforeSlash = newValue.slice(0, -1);
				if (beforeSlash === "" || beforeSlash.endsWith(" ")) {
					completion.triggerCompletion(newValue);
				}
			}
		},
		[input, completion],
	);

	// Handle input submission
	const handleSubmit = useCallback(
		async (value: string) => {
			// If showing suggestions, apply selected
			if (completion.isShowing && completion.suggestions.length > 0) {
				const selected = completion.selectCurrent();
				if (selected) {
					applyCompletion({
						label: selected.text,
						value: selected.text,
						description: selected.description,
						category: selected.category ?? "builtin",
					});
				}
				return;
			}

			// Execute command
			await runCommand(value);
			setInput("");
			history.reset();
		},
		[completion, applyCompletion, runCommand, history],
	);

	// Keyboard input handling
	useInput((char, key) => {
		// Ctrl+C - double press to exit
		if (key.ctrl && char === "c") {
			ctrlC.handleCtrlC();
			return;
		}

		// Ctrl+D - immediate exit
		if (key.ctrl && char === "d") {
			session.saveHistory().finally(() => exit());
			return;
		}

		// Tab - trigger or cycle completion
		if (key.tab) {
			// Use ref to get current input value (avoids stale closure)
			const currentInput = inputRef.current;
			if (completion.isShowing) {
				if (key.shift) {
					completion.navigateUp();
				} else {
					completion.navigateDown();
				}
			} else {
				// triggerCompletion is async - returns suggestion for single match
				completion
					.triggerCompletion(currentInput)
					.then((suggestion) => {
						if (suggestion) {
							// Single match - auto-complete
							applyCompletion({
								label: suggestion.text,
								value: suggestion.text,
								description: suggestion.description,
								category: suggestion.category ?? "builtin",
							});
						}
						// Multiple matches - suggestions are shown by the hook
					})
					.catch(() => {
						// Ignore completion errors
					});
			}
			return;
		}

		// Up/Down arrows
		if (key.upArrow) {
			if (completion.isShowing) {
				completion.navigateUp();
			} else {
				history.navigateUp();
			}
			return;
		}

		if (key.downArrow) {
			if (completion.isShowing) {
				completion.navigateDown();
			} else {
				history.navigateDown();
			}
			return;
		}

		// Escape - cancel suggestions
		if (key.escape) {
			completion.hide();
			return;
		}

		// Right arrow - apply suggestion
		if (
			key.rightArrow &&
			completion.isShowing &&
			completion.suggestions.length > 0
		) {
			const selected = completion.selectCurrent();
			if (selected) {
				applyCompletion({
					label: selected.text,
					value: selected.text,
					description: selected.description,
					category: selected.category ?? "builtin",
				});
			}
			return;
		}
	});

	// Suggestion navigation handlers
	const handleSuggestionNavigate = useCallback(
		(direction: "up" | "down") => {
			if (direction === "up") {
				completion.navigateUp();
			} else {
				completion.navigateDown();
			}
		},
		[completion],
	);

	const handleSuggestionSelect = useCallback(
		(suggestion: Suggestion) => {
			applyCompletion(suggestion);
			completion.hide();
		},
		[applyCompletion, completion],
	);

	// Render loading state
	if (!isInitialized) {
		return (
			<Box>
				<Text>Initializing...</Text>
			</Box>
		);
	}

	// Calculate output height
	const maxOutputLines = Math.max(5, 20); // Show reasonable output

	// Build connection info for banner
	const connectionInfo = {
		tenant: session.getTenant() || undefined,
		username: session.getUsername() || undefined,
		tier: session.getTier() || undefined,
		namespace: session.getNamespace(),
		isConnected: session.isConnected(),
	};

	return (
		<Box flexDirection="column" width={width}>
			{/* Banner - shown until first command */}
			{showBanner && (
				<Banner version={CLI_VERSION} connectionInfo={connectionInfo} />
			)}

			{/* Output area */}
			<Box flexDirection="column">
				{output.slice(-maxOutputLines).map((line, i) => (
					<Text key={i}>{line}</Text>
				))}
			</Box>

			{/* Suggestions popup */}
			{completion.isShowing && completion.suggestions.length > 0 && (
				<Suggestions
					suggestions={toUISuggestions(completion.suggestions)}
					selectedIndex={completion.selectedIndex}
					onSelect={handleSuggestionSelect}
					onNavigate={handleSuggestionNavigate}
					onCancel={completion.hide}
					maxVisible={10}
					isActive={false} // Let App handle keyboard
				/>
			)}

			{/* Input box */}
			<InputBox
				prompt={prompt}
				value={input}
				onChange={handleInputChange}
				onSubmit={handleSubmit}
				width={width}
				isActive={true}
				inputKey={inputKey}
			/>

			{/* Status bar */}
			<StatusBar gitInfo={gitInfo} width={width} hint={statusHint} />
		</Box>
	);
}

export default App;
