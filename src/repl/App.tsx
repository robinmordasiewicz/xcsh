/**
 * App - Main Ink application for the REPL
 * Orchestrates all UI components and handles keyboard input
 */

import React, { useState, useCallback, useEffect, useRef } from "react";
import { Box, Text, useApp, useInput, useStdout, Static } from "ink";

import { InputBox, StatusBar, Suggestions } from "./components/index.js";
import type { Suggestion } from "./components/Suggestions.js";
import { getGitInfo, type GitInfo } from "./components/StatusBar.js";
import { REPLSession } from "./session.js";
import { buildPlainPrompt } from "./prompt.js";
import { useDoubleCtrlC } from "./hooks/useDoubleCtrlC.js";
import { useHistory } from "./hooks/useHistory.js";
import { useCompletion } from "./hooks/useCompletion.js";
import { executeCommand } from "./executor.js";
import { isCustomDomain } from "../domains/index.js";
import { domainRegistry } from "../types/domains.js";
import { extensionRegistry } from "../extensions/index.js";

/**
 * Props for the App component
 */
export interface AppProps {
	/** Pre-initialized session (optional - will create new if not provided) */
	initialSession?: REPLSession;
}

/**
 * Check if a word is a valid domain (custom, API, or extension)
 */
function isValidDomain(word: string): boolean {
	const lowerWord = word.toLowerCase();
	// Check custom domains
	if (isCustomDomain(lowerWord)) return true;
	// Check API domains
	if (domainRegistry.has(lowerWord)) return true;
	// Check extension-only domains
	if (extensionRegistry.getExtension(lowerWord)) return true;
	return false;
}

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
export function App({ initialSession }: AppProps = {}): React.ReactElement {
	const { exit } = useApp();
	const { stdout } = useStdout();

	// Session state - use provided session or create new one
	const [session] = useState(() => initialSession ?? new REPLSession());
	// If session was pre-initialized, start as initialized
	const [isInitialized, setIsInitialized] = useState(!!initialSession);

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

	// Output items with unique IDs for Static component (goes to scrollback)
	const [outputItems, setOutputItems] = useState<
		Array<{ id: number; content: string }>
	>([]);
	const outputIdRef = useRef(0);
	const [prompt, setPrompt] = useState("> ");
	const [width, setWidth] = useState(stdout?.columns ?? 80);
	const [gitInfo, setGitInfo] = useState<GitInfo | undefined>(undefined);
	const [statusHint, setStatusHint] = useState("Ctrl+C twice to exit");
	const [historyArray, setHistoryArray] = useState<string[]>([]);
	const [inputKey, setInputKey] = useState(0); // Key to reset cursor position
	// Raw stdout handling - hide status bar while writing direct stdout content
	const [hideStatusBar, setHideStatusBar] = useState(false);
	const [pendingRawStdout, setPendingRawStdout] = useState<string | null>(
		null,
	);

	// Effect to handle raw stdout writing when status bar is hidden
	// This ensures the status bar is removed from render BEFORE we write content
	// that may cause terminal scrolling, preventing it from being captured in scrollback
	// CRITICAL: Do not modify without testing - prevents status bar in scrollback
	useEffect(() => {
		if (hideStatusBar && pendingRawStdout) {
			// Status bar is now hidden from Ink's render, safe to write raw stdout
			process.stdout.write(pendingRawStdout);
			// CRITICAL: Newlines to prevent Ink from truncating banner bottom border
			// - 3 newlines is minimum needed to prevent truncation
			// - DO NOT REDUCE below 3 - will truncate bottom of banner frame
			// - DO NOT INCREASE above 3 - causes excessive deadspace below banner
			process.stdout.write("\n\n\n");
			// Restore state - status bar will reappear in next render
			setPendingRawStdout(null);
			setHideStatusBar(false);
		}
	}, [hideStatusBar, pendingRawStdout]);

	// Completion state
	const completion = useCompletion({
		session: isInitialized ? session : null,
	});

	// History state
	const history = useHistory({
		history: historyArray,
		onSelect: (cmd) => {
			setInput(cmd);
			// Force TextInput remount to reflect new value
			setInputKey((k) => k + 1);
		},
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

	// Initialize session (or set up UI state from pre-initialized session)
	useEffect(() => {
		const init = async () => {
			// Only initialize if not pre-initialized
			if (!isInitialized) {
				await session.initialize();
				setIsInitialized(true);
			}

			// Always set up prompt, history, and git info
			setPrompt(buildPlainPrompt(session));

			// Get initial history array
			const histMgr = session.getHistory();
			if (histMgr) {
				setHistoryArray(histMgr.getHistory());
			}

			// Get git repository info
			setGitInfo(getGitInfo());
		};
		init();
		// eslint-disable-next-line react-hooks/exhaustive-deps
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

	// Add output line(s) - each line gets unique ID for Static component
	const addOutput = useCallback((line: string) => {
		const lines = line.split("\n");
		const newItems = lines.map((content) => ({
			id: outputIdRef.current++,
			content,
		}));
		setOutputItems((prev) => {
			const combined = [...prev, ...newItems];
			// Keep buffer under 1000 items
			if (combined.length > 1000) {
				return combined.slice(combined.length - 1000);
			}
			return combined;
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

			// Execute via executor module
			const result = await executeCommand(trimmed, session);

			// Handle raw stdout content (e.g., image banner with cursor positioning)
			// This bypasses Ink's rendering to avoid status bar in scrollback
			if (result.rawStdout) {
				// Hide status bar before writing raw content
				// The useEffect will handle the actual writing after render
				setHideStatusBar(true);
				setPendingRawStdout(result.rawStdout);
				// Don't process normal output - rawStdout is the complete output
			} else if (result.shouldClear) {
				// Handle clear
				setOutputItems([]);
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
				return;
			}

			// Check if "/" was typed - trigger completion
			if (newValue !== oldValue && newValue.endsWith("/")) {
				const beforeSlash = newValue.slice(0, -1);
				if (beforeSlash === "" || beforeSlash.endsWith(" ")) {
					completion.triggerCompletion(newValue);
					return;
				}
			}

			// Check if space was typed after a known domain - trigger contextual completion
			// This works at any depth: "login ", "login profile ", "login profile use ", etc.
			if (
				newValue !== oldValue &&
				newValue.endsWith(" ") &&
				!oldValue.endsWith(" ")
			) {
				const trimmed = newValue.trimEnd();
				const words = trimmed.split(/\s+/);
				if (words.length > 0 && words[0]) {
					const firstWord = words[0].startsWith("/")
						? words[0].slice(1)
						: words[0];
					if (isValidDomain(firstWord)) {
						completion.triggerCompletion(newValue);
					}
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

			// Clear input immediately to prevent duplicate display in scrollback
			// (InputBox would show command while addOutput also adds it to Static)
			setInput("");
			history.reset();

			// Execute command
			await runCommand(value);
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

		// Tab - trigger completion or select current suggestion
		if (key.tab) {
			// Use ref to get current input value (avoids stale closure)
			const currentInput = inputRef.current;
			if (completion.isShowing) {
				// Tab selects the current suggestion (same as Enter)
				const selected = completion.suggestions.at(
					completion.selectedIndex,
				);
				if (selected) {
					applyCompletion({
						label: selected.text,
						value: selected.text,
						description: selected.description,
						category: selected.category ?? "builtin",
					});
					completion.hide();
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

	// Always mount Static from start to prevent tree restructure issues
	// This keeps the component tree stable and prevents Ink's screen clearing
	return (
		<Box flexDirection="column" width={width}>
			{/* Static content - goes to scrollback, never re-rendered */}
			<Static items={outputItems}>
				{(item) => <Text key={item.id}>{item.content}</Text>}
			</Static>

			{/* Conditionally render active UI or loading state */}
			{/* Hide entire active UI when writing raw stdout to prevent it appearing in scrollback */}
			{isInitialized && !hideStatusBar ? (
				<>
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

					{/* Suggestions popup OR Status bar - mutually exclusive to conserve space */}
					{completion.isShowing &&
					completion.suggestions.length > 0 ? (
						<Suggestions
							suggestions={toUISuggestions(
								completion.suggestions,
							)}
							selectedIndex={completion.selectedIndex}
							onSelect={handleSuggestionSelect}
							onNavigate={handleSuggestionNavigate}
							onCancel={completion.hide}
							maxVisible={20}
							isActive={false} // Let App handle keyboard
						/>
					) : (
						<StatusBar
							gitInfo={gitInfo}
							width={width}
							hint={statusHint}
						/>
					)}
				</>
			) : !isInitialized ? (
				<Text>Initializing...</Text>
			) : null}
		</Box>
	);
}

export default App;
