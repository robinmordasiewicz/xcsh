#!/usr/bin/env tsx
/**
 * Generate shell completion files
 *
 * This script generates shell completion files for bash, zsh, and fish.
 * It is called during the build process to ensure completions are always up-to-date.
 *
 * Usage: tsx scripts/generate-completions.ts
 */

import { writeFileSync, mkdirSync, existsSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

// Get project root
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const projectRoot = join(__dirname, "..");

// Import generators - these will be available after the TypeScript is compiled
// For script execution, we dynamically import from the source
async function main() {
	console.log("🔧 Generating shell completions...");

	// Ensure completions directory exists
	const completionsDir = join(projectRoot, "completions");
	if (!existsSync(completionsDir)) {
		mkdirSync(completionsDir, { recursive: true });
	}

	try {
		// Import domains module first to register custom domains (login, completion, etc.)
		// This populates the customDomains registry before generators access it
		await import("../src/domains/index.js");

		// Dynamic import to load the generators
		const { generateBashCompletion, generateZshCompletion, generateFishCompletion } =
			await import("../src/domains/completion/generators.js");

		// Generate bash completion
		const bashScript = generateBashCompletion();
		const bashPath = join(completionsDir, "xcsh.bash");
		writeFileSync(bashPath, bashScript, "utf-8");
		console.log(`✓ Generated: ${bashPath}`);

		// Generate zsh completion
		const zshScript = generateZshCompletion();
		const zshPath = join(completionsDir, "_xcsh");
		writeFileSync(zshPath, zshScript, "utf-8");
		console.log(`✓ Generated: ${zshPath}`);

		// Generate fish completion
		const fishScript = generateFishCompletion();
		const fishPath = join(completionsDir, "xcsh.fish");
		writeFileSync(fishPath, fishScript, "utf-8");
		console.log(`✓ Generated: ${fishPath}`);

		console.log("✅ Shell completions generated successfully!");
	} catch (error) {
		console.error("❌ Failed to generate completions:", error);
		process.exit(1);
	}
}

main();
