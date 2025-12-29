/**
 * Patched ansi-escapes that preserves scrollback buffer.
 *
 * The original clearTerminal uses \x1b[3J which erases the scrollback buffer.
 * This patch removes that escape sequence so terminal history is preserved
 * when Ink clears the screen.
 *
 * Original: \x1b[2J\x1b[3J\x1b[H (erase screen + erase scrollback + cursor home)
 * Patched:  \x1b[2J\x1b[H       (erase screen + cursor home, scrollback preserved)
 *
 * Note: We import directly from node_modules to avoid circular alias reference.
 */

// Import directly from the node_modules path (not the aliased name)
import originalDefault from "../../node_modules/ansi-escapes/index.js";
import * as originalNamed from "../../node_modules/ansi-escapes/index.js";

// ESC sequence building blocks
const ESC = "\u001B[";
const eraseScreen = `${ESC}2J`;
const cursorHome = `${ESC}H`;

// Patched clearTerminal without scrollback erase (\x1b[3J)
const patchedClearTerminal = `${eraseScreen}${cursorHome}`;

// Re-export all named exports from original
export const {
	beep,
	clearScreen,
	clearViewport,
	cursorBackward,
	cursorDown,
	cursorForward,
	cursorGetPosition,
	cursorHide,
	cursorLeft,
	cursorMove,
	cursorNextLine,
	cursorPrevLine,
	cursorRestorePosition,
	cursorSavePosition,
	cursorShow,
	cursorTo,
	cursorUp,
	enterAlternativeScreen,
	eraseDown,
	eraseEndLine,
	eraseLine,
	eraseLines,
	eraseScreen: eraseScreenOriginal,
	eraseStartLine,
	eraseUp,
	exitAlternativeScreen,
	image,
	link,
	scrollDown,
	scrollUp,
	setCwd,
	iTerm,
	ConEmu,
} = originalNamed;

// Override clearTerminal as a named export
export const clearTerminal = patchedClearTerminal;

// Default export with all original functions plus patched clearTerminal
export default {
	...originalDefault,
	clearTerminal: patchedClearTerminal,
};
