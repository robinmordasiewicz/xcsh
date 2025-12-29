/**
 * Terminal capability detection for inline image support.
 * Detects iTerm2, Kitty, WezTerm, and other terminals with image support.
 */

/**
 * Terminal capabilities detected from environment.
 */
export interface TerminalCapabilities {
	/** Whether the terminal supports inline images */
	supportsInlineImages: boolean;
	/** Terminal program name if detected */
	terminalProgram: string | null;
	/** True if running in iTerm2 */
	isITerm2: boolean;
	/** True if running in Kitty */
	isKitty: boolean;
	/** True if running in WezTerm */
	isWezTerm: boolean;
	/** True if running in Mintty */
	isMintty: boolean;
}

/**
 * Detect terminal capabilities for inline image support.
 * Checks TERM_PROGRAM, LC_TERMINAL, and other environment variables.
 */
export function detectTerminalCapabilities(): TerminalCapabilities {
	const termProgram = process.env.TERM_PROGRAM ?? "";
	const lcTerminal = process.env.LC_TERMINAL ?? "";
	const term = process.env.TERM ?? "";

	// iTerm2 detection
	const isITerm2 =
		termProgram === "iTerm.app" ||
		lcTerminal === "iTerm2" ||
		!!process.env.ITERM_SESSION_ID;

	// Kitty detection
	const isKitty = termProgram === "kitty" || term.includes("kitty");

	// WezTerm detection
	const isWezTerm = termProgram === "WezTerm" || !!process.env.WEZTERM_PANE;

	// Mintty detection (Windows)
	const isMintty = !!process.env.MINTTY;

	// Any of these terminals support inline images
	const supportsInlineImages = isITerm2 || isKitty || isWezTerm || isMintty;

	return {
		supportsInlineImages,
		terminalProgram: termProgram || lcTerminal || null,
		isITerm2,
		isKitty,
		isWezTerm,
		isMintty,
	};
}

/**
 * Options for iTerm2 inline image display.
 */
export interface ITerm2ImageOptions {
	/** Width in character cells, pixels (Npx), or percentage (N%) */
	width?: string | number;
	/** Height in character cells, pixels (Npx), or percentage (N%) */
	height?: string | number;
	/** Whether to preserve aspect ratio (default: true) */
	preserveAspectRatio?: boolean;
	/** Display inline (true) or as download (false). Default: true */
	inline?: boolean;
	/** Optional filename for the image */
	name?: string;
}

/**
 * Generate iTerm2 inline image escape sequence.
 *
 * Protocol: ESC ] 1337 ; File = [args] : base64data BEL
 *
 * @param base64Data - Base64-encoded image data
 * @param options - Display options
 * @returns Escape sequence string to write to terminal
 */
export function generateITerm2ImageSequence(
	base64Data: string,
	options: ITerm2ImageOptions = {},
): string {
	const args: string[] = [];

	// Add width if specified
	if (options.width !== undefined) {
		args.push(`width=${options.width}`);
	}

	// Add height if specified
	if (options.height !== undefined) {
		args.push(`height=${options.height}`);
	}

	// Add aspect ratio preservation (default: on)
	if (options.preserveAspectRatio !== undefined) {
		args.push(`preserveAspectRatio=${options.preserveAspectRatio ? 1 : 0}`);
	}

	// Add filename if specified (base64 encoded)
	if (options.name) {
		const nameBase64 = Buffer.from(options.name, "utf-8").toString(
			"base64",
		);
		args.push(`name=${nameBase64}`);
	}

	// Always display inline by default
	args.push(`inline=${options.inline !== false ? 1 : 0}`);

	// Build the escape sequence
	// ESC ] 1337 ; File = <args> : <base64> BEL
	const argsString = args.join(";");
	return `\x1b]1337;File=${argsString}:${base64Data}\x07`;
}

/**
 * Generate Kitty terminal image sequence.
 * Kitty uses a different protocol than iTerm2.
 *
 * @param base64Data - Base64-encoded image data
 * @returns Escape sequence string for Kitty
 */
export function generateKittyImageSequence(base64Data: string): string {
	// Kitty protocol: ESC _ G <payload> ESC \
	// For simple display: a=T (transmit), f=100 (PNG), t=d (direct)
	// This is a simplified version - Kitty has a more complex chunked protocol
	return `\x1b_Ga=T,f=100,t=d;${base64Data}\x1b\\`;
}

/**
 * Get the appropriate image sequence for the detected terminal.
 * Falls back to iTerm2 format which has widest compatibility.
 *
 * @param base64Data - Base64-encoded PNG image data
 * @param capabilities - Detected terminal capabilities
 * @param options - Display options
 * @returns Escape sequence string or null if not supported
 */
export function getTerminalImageSequence(
	base64Data: string,
	capabilities: TerminalCapabilities,
	options: ITerm2ImageOptions = {},
): string | null {
	if (!capabilities.supportsInlineImages) {
		return null;
	}

	// Kitty has its own protocol
	if (capabilities.isKitty) {
		return generateKittyImageSequence(base64Data);
	}

	// iTerm2, WezTerm, and Mintty all support iTerm2 protocol
	return generateITerm2ImageSequence(base64Data, options);
}
