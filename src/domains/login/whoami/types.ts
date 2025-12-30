/**
 * Whoami Types
 * Type definitions for connection info and identity display
 */

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
}

/**
 * Options for whoami display
 */
export interface WhoamiOptions {
	verbose?: boolean;
	json?: boolean;
}
