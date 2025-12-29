/**
 * Banner Command - Display the xcsh banner with optional logo modes
 *
 * Supports multiple logo display modes:
 * - auto: Use image if terminal supports it, otherwise ASCII
 * - image: Image only (falls back to ASCII if unsupported)
 * - ascii: ASCII art only (traditional banner)
 * - both: Image followed by ASCII art
 * - none: No logo at all
 */

export { bannerCommand, renderBanner } from "./display.js";
