/**
 * REPL module exports.
 * Re-exports all REPL-related types and functions.
 */

// Core modules
export { ContextPath, ContextValidator } from "./context.js";
export { HistoryManager, getHistoryFilePath } from "./history.js";
export { REPLSession, createSession, type SessionConfig } from "./session.js";
export {
	buildPlainPrompt,
	buildColoredPrompt,
	getPromptParts,
	type PromptParts,
} from "./prompt.js";

// Main application
export { App } from "./App.js";

// Hooks
export { useDoubleCtrlC, useHistory, useCompletion } from "./hooks/index.js";

// Components
export * from "./components/index.js";

// Completion
export * from "./completion/index.js";
