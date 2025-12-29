/**
 * Stub modules for optional dependencies.
 * These provide minimal implementations for packages that are optional
 * or need patching for terminal compatibility.
 */

// Export patched ansi-escapes (preserves scrollback buffer)
export * as ansiEscapes from "./ansi-escapes.js";

// Export devtools stub (optional peer dependency)
export * as devtools from "./devtools.js";
