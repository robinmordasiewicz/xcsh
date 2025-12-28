/**
 * Whoami Service
 * Core service for fetching connection and identity information
 */

import type { REPLSession } from "../../../repl/session.js";
import type { WhoamiInfo, WhoamiOptions } from "./types.js";
import { toDisplayTier } from "./types.js";
import { SubscriptionClient } from "../../../subscription/client.js";
import { isAddonActive } from "../../../subscription/types.js";

/**
 * Get whoami information from session
 * Returns connection and identity info based on options
 */
export async function getWhoamiInfo(
	session: REPLSession,
	options: WhoamiOptions = {},
): Promise<WhoamiInfo> {
	const info: WhoamiInfo = {
		serverUrl: session.getServerUrl(),
		namespace: session.getNamespace(),
		isAuthenticated: session.isAuthenticated(),
	};

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

	// Fetch tier from subscription API
	const apiClient = session.getAPIClient();
	if (apiClient) {
		const subscriptionClient = new SubscriptionClient(apiClient);

		// Get tier
		try {
			const tierValue = await subscriptionClient.getTierFromCurrentPlan();
			const displayTier = toDisplayTier(tierValue);
			if (displayTier) {
				info.tier = displayTier;
			}
		} catch {
			// Tier fetch failed - just omit from display
		}

		// Get quotas if requested
		if (options.includeQuotas || options.verbose) {
			try {
				const subscriptionInfo =
					await subscriptionClient.getSubscriptionInfo();
				info.quotas = subscriptionInfo.quotaSummary;
			} catch {
				// Quota fetch failed - omit from display
			}
		}

		// Get addons if requested
		if (options.includeAddons || options.verbose) {
			try {
				const addons = await subscriptionClient.getAddonServices();
				// Only include active addons
				info.addons = addons.filter(isAddonActive);
			} catch {
				// Addon fetch failed - omit from display
			}
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
	return getWhoamiInfo(session, {
		includeQuotas: false,
		includeAddons: false,
		verbose: false,
	});
}
