import { defineConfig } from "tsup";
import { readFileSync } from "fs";

// Read package.json version
const packageJson = JSON.parse(readFileSync("./package.json", "utf-8"));

// Generate build version
// Priority: XCSH_VERSION env var > package.json version
function getBuildVersion(): string {
	if (process.env.XCSH_VERSION) {
		return process.env.XCSH_VERSION;
	}

	// Use package version with build timestamp for dev builds
	const now = new Date();
	const timestamp = [
		now.getFullYear(),
		String(now.getMonth() + 1).padStart(2, "0"),
		String(now.getDate()).padStart(2, "0"),
		String(now.getHours()).padStart(2, "0"),
		String(now.getMinutes()).padStart(2, "0"),
	].join("");

	return `${packageJson.version}-${timestamp}`;
}

export default defineConfig({
	entry: ["src/index.tsx"],
	format: ["esm"],
	dts: true,
	clean: true,
	define: {
		BUILD_VERSION: JSON.stringify(getBuildVersion()),
	},
});
