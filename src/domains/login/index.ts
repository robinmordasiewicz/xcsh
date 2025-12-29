/**
 * Login Domain - Authentication, identity, and session management for F5 XC
 */

import type { DomainDefinition, SubcommandGroup } from "../registry.js";
import { listCommand } from "./profile/list.js";
import { showCommand } from "./profile/show.js";
import { createCommand } from "./profile/create.js";
import { useCommand } from "./profile/use.js";
import { deleteCommand } from "./profile/delete.js";
import { contextSubcommands } from "./context/index.js";
import { bannerCommand } from "./banner/index.js";
import { whoamiCommand } from "./whoami/index.js";

/**
 * Profile subcommand group
 */
const profileSubcommands: SubcommandGroup = {
	name: "profile",
	description: "Manage saved connection profiles",
	commands: new Map([
		["list", listCommand],
		["show", showCommand],
		["create", createCommand],
		["use", useCommand],
		["delete", deleteCommand],
	]),
};

/**
 * Login domain definition
 */
export const loginDomain: DomainDefinition = {
	name: "login",
	description:
		"Authentication, identity, and session management for F5 XC. Manage connection profiles to save and switch between tenants, handle context for namespace targeting, and verify current authentication status with whoami.",
	descriptionShort: "Authentication and session management",
	descriptionMedium:
		"Manage connection profiles, authentication contexts, and session identity for F5 Distributed Cloud.",
	commands: new Map([
		["show", whoamiCommand],
		["banner", bannerCommand],
	]),
	subcommands: new Map([
		["profile", profileSubcommands],
		["context", contextSubcommands],
	]),
};
