/**
 * login profile delete - Delete a profile
 */

import type { CommandDefinition } from "../../registry.js";
import { successResult, errorResult } from "../../registry.js";
import { getProfileManager } from "../../../profile/index.js";

export const deleteCommand: CommandDefinition = {
	name: "delete",
	description:
		"Delete a saved connection profile permanently. Requires --force flag to delete the currently active profile to prevent accidental disconnection from active tenant.",
	descriptionShort: "Delete a saved profile",
	descriptionMedium:
		"Remove a saved profile permanently. Use --force to delete active profile.",
	usage: "<name> [--force]",
	aliases: ["rm", "remove"],

	async execute(args, _session) {
		const manager = getProfileManager();

		// Parse arguments
		const name = args.find((arg) => !arg.startsWith("-"));
		const force = args.includes("--force") || args.includes("-f");

		if (!name) {
			return errorResult("Usage: login profile delete <name> [--force]");
		}

		try {
			// Check if profile exists
			const profile = await manager.get(name);
			if (!profile) {
				return errorResult(`Profile '${name}' not found.`);
			}

			// Check if it's the active profile
			const activeProfile = await manager.getActive();
			if (profile.name === activeProfile) {
				if (!force) {
					return errorResult(
						[
							`Cannot delete active profile '${name}'.`,
							``,
							`Either:`,
							`  - Switch to another profile first: login profile use <other>`,
							`  - Force delete with: login profile delete ${name} --force`,
						].join("\n"),
					);
				}
				// Force delete - clear active profile
				// Note: This would need to be implemented in the manager
			}

			// Delete the profile
			const result = await manager.delete(name);

			if (!result.success) {
				return errorResult(result.message);
			}

			return successResult([`Profile '${name}' deleted successfully.`]);
		} catch (error) {
			return errorResult(
				`Failed to delete profile: ${error instanceof Error ? error.message : "Unknown error"}`,
			);
		}
	},

	async completion(partial, _args, _session) {
		const manager = getProfileManager();
		const profiles = await manager.list();
		return profiles
			.map((p) => p.name)
			.filter((name) =>
				name.toLowerCase().startsWith(partial.toLowerCase()),
			);
	},
};
