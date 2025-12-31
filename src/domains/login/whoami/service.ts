/**
 * Whoami Service
 * Core service for fetching connection and identity information
 */

import type { REPLSession } from "../../../repl/session.js";
import type { WhoamiInfo, WhoamiOptions } from "./types.js";

/**
 * Get whoami information from session
 * Returns connection and identity info based on options
 */
export async function getWhoamiInfo(
	session: REPLSession,
	_options: WhoamiOptions = {},
): Promise<WhoamiInfo> {
	const info: WhoamiInfo = {
		serverUrl: session.getServerUrl(),
		namespace: session.getNamespace(),
		isAuthenticated: session.isAuthenticated(),
	};

	// Add validation fields only if they have values
	const isValidated = session.isTokenValidated();
	const validationError = session.getValidationError();
	if (isValidated !== undefined) {
		info.isValidated = isValidated;
	}
	if (validationError) {
		info.validationError = validationError;
	}

	// If not authenticated, return minimal info
	if (!info.isAuthenticated) {
		return info;
	}

	// Get tenant from session
	const tenant = session.getTenant();
	if (tenant) {
		info.tenant = tenant;
	}

	// Get username from session
	const username = session.getUsername();
	if (username) {
		// Username might be email
		if (username.includes("@")) {
			info.email = username;
		} else {
			info.username = username;
		}
	}

	return info;
}

/**
 * Get basic whoami info for banner display (minimal API calls)
 */
export async function getWhoamiInfoBasic(
	session: REPLSession,
): Promise<WhoamiInfo> {
	return getWhoamiInfo(session, { verbose: false });
}
