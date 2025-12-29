import { defineConfig } from "tsup";
import { readFileSync } from "fs";
import { execSync } from "child_process";

// Read package.json version
const packageJson = JSON.parse(readFileSync("./package.json", "utf-8"));

// Fetch latest release version from GitHub
function fetchLatestRelease(): string | null {
	try {
		const result = execSync(
			'curl -s --connect-timeout 3 "https://api.github.com/repos/robinmordasiewicz/xcsh/releases/latest" | grep \'"tag_name":\' | sed -E \'s/.*"v?([^"]+)".*/\\1/\'',
			{ encoding: "utf-8", timeout: 5000 },
		).trim();
		// Validate we got a version-like string (e.g., "6.9.0")
		if (result && /^\d+\.\d+\.\d+/.test(result)) {
			return result;
		}
		return null;
	} catch {
		// Network unavailable or timeout
		return null;
	}
}

// Generate build version
// Priority: XCSH_VERSION env var > latest release + timestamp > DEV + timestamp
function getBuildVersion(): string {
	if (process.env.XCSH_VERSION) {
		return process.env.XCSH_VERSION;
	}

	// Generate YYYYMMDDHHMM timestamp for dev builds
	const now = new Date();
	const timestamp = [
		now.getFullYear(),
		String(now.getMonth() + 1).padStart(2, "0"),
		String(now.getDate()).padStart(2, "0"),
		String(now.getHours()).padStart(2, "0"),
		String(now.getMinutes()).padStart(2, "0"),
	].join("");

	// Try to fetch latest release version
	const latestRelease = fetchLatestRelease();
	if (latestRelease) {
		return `${latestRelease}-${timestamp}`;
	}

	// Offline fallback
	return `DEV-${timestamp}`;
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
