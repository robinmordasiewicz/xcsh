/**
 * Login Domain - Profile management for F5 XC authentication
 */

import type { DomainDefinition, SubcommandGroup } from "../registry.js";
import { listCommand } from "./profile/list.js";
import { showCommand } from "./profile/show.js";
import { createCommand } from "./profile/create.js";
import { useCommand } from "./profile/use.js";
import { deleteCommand } from "./profile/delete.js";

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
	description: "Authentication and profile management",
	commands: new Map(),
	subcommands: new Map([["profile", profileSubcommands]]),
};
