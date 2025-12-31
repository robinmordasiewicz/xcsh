import { defineConfig } from "tsup";
import { readFileSync } from "fs";

// Get upstream API version from specs
// Priority: .specs/index.json > .specs/.version > fallback
function getUpstreamApiVersion(): string {
	// Primary: Read from .specs/index.json
	try {
		const indexJson = JSON.parse(readFileSync("./.specs/index.json", "utf-8"));
		if (indexJson.version) {
			return indexJson.version; // e.g., "1.0.78"
		}
	} catch {
		// File not found or invalid JSON
	}

	// Fallback: Read from .specs/.version
	try {
		const version = readFileSync("./.specs/.version", "utf-8").trim();
		return version.replace(/^v/, ""); // Strip leading v if present
	} catch {
		// File not found
	}

	// Ultimate fallback
	return "0.0.0";
}

// Detect if running in CI environment
function isCI(): boolean {
	return Boolean(process.env.CI || process.env.GITHUB_ACTIONS);
}

// Generate build version
// Format: v{upstream}-YYMMDDHHMM[-BETA]
// Priority: XCSH_VERSION env var > generated version
function getBuildVersion(): string {
	// Environment override (for CI or explicit version)
	if (process.env.XCSH_VERSION) {
		return process.env.XCSH_VERSION;
	}

	const upstreamVersion = getUpstreamApiVersion();

	// Generate YYMMDDHHMM timestamp (10 digits, UTC)
	const now = new Date();
	const timestamp = [
		String(now.getUTCFullYear()).slice(-2), // YY
		String(now.getUTCMonth() + 1).padStart(2, "0"), // MM
		String(now.getUTCDate()).padStart(2, "0"), // DD
		String(now.getUTCHours()).padStart(2, "0"), // HH
		String(now.getUTCMinutes()).padStart(2, "0"), // MM
	].join("");

	// CI builds get no suffix, local builds get -BETA
	const suffix = isCI() ? "" : "-BETA";

	return `v${upstreamVersion}-${timestamp}${suffix}`;
}

export default defineConfig({
	entry: ["src/index.tsx"],
	// Use ESM format - required for ink v5+ which uses top-level await
	format: ["esm"],
	dts: true,
	clean: true,
	// Bundle all dependencies into single file for standalone binary
	noExternal: [/.*/],
	// Disable code splitting - single file output required for pkg
	splitting: false,
	// Shebang is added automatically by tsup when bin is defined in package.json
	define: {
		BUILD_VERSION: JSON.stringify(getBuildVersion()),
	},
	platform: "node",
	// Add banner to polyfill require for CJS dependencies in ESM bundle
	banner: {
		js: `import { createRequire as __createRequire } from 'module';
const require = __createRequire(import.meta.url);`,
	},
	esbuildOptions(options) {
		// Module aliasing for stubs and patches
		options.alias = {
			// Stub react-devtools-core (optional dependency not needed for production)
			"react-devtools-core": "./src/stubs/devtools.ts",
			// Patch ansi-escapes to preserve scrollback buffer
			// The original clearTerminal uses \x1b[3J which erases scrollback
			"ansi-escapes": "./src/stubs/ansi-escapes.ts",
		};
	},
});
