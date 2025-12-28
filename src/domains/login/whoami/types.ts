/**
 * Whoami Types
 * Type definitions for connection info and identity display
 */

import type {
	AddonServiceInfo,
	QuotaSummary,
} from "../../../subscription/types.js";

/**
 * Connection and identity information
 */
export interface WhoamiInfo {
	// Always available
	serverUrl: string;
	namespace: string;
	isAuthenticated: boolean;

	// Fetched (undefined if fetch fails - omit from display)
	tenant?: string;
	username?: string;
	email?: string;
	tier?: "Standard" | "Advanced";

	// Optional (fetched on demand with flags)
	quotas?: QuotaSummary;
	addons?: AddonServiceInfo[];
}

/**
 * Options for whoami display
 */
export interface WhoamiOptions {
	includeQuotas?: boolean;
	includeAddons?: boolean;
	verbose?: boolean;
	json?: boolean;
}

/**
 * Display tier type for user-facing output
 */
export type DisplayTier = "Standard" | "Advanced";

/**
 * Convert API tier value to display tier
 */
export function toDisplayTier(tier: string): DisplayTier | undefined {
	const normalized = tier.toUpperCase();
	switch (normalized) {
		case "STANDARD":
		case "BASIC": // Legacy mapping
			return "Standard";
		case "ADVANCED":
		case "PREMIUM": // Legacy mapping
		case "ENTERPRISE":
			return "Advanced";
		default:
			return undefined;
	}
}
