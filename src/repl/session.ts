/**
 * REPLSession holds state across the REPL lifetime.
 * Manages context, history, namespace, and user information.
 */

import { ContextPath, ContextValidator } from "./context.js";
import { HistoryManager, getHistoryFilePath } from "./history.js";
import { ENV_PREFIX } from "../branding/index.js";
import {
	getProfileManager,
	type Profile,
	type ProfileManager,
} from "../profile/index.js";
import { APIClient } from "../api/index.js";
import { SubscriptionClient } from "../subscription/client.js";
import type { OutputFormat } from "../output/index.js";

/**
 * Configuration for creating a REPL session
 */
export interface SessionConfig {
	namespace?: string;
	serverUrl?: string;
	apiToken?: string;
	outputFormat?: OutputFormat;
	debug?: boolean;
}

/**
 * REPLSession holds state across the REPL lifetime
 */
export class REPLSession {
	private _history: HistoryManager | null = null;
	private _namespace: string;
	private _lastExitCode: number = 0;
	private _contextPath: ContextPath;
	private _tenant: string = "";
	private _username: string = "";
	private _validator: ContextValidator;
	private _serverUrl: string = "";
	private _apiToken: string = "";
	private _apiClient: APIClient | null = null;
	private _outputFormat: OutputFormat = "yaml";
	private _debug: boolean = false;
	private _profileManager: ProfileManager;
	private _activeProfile: Profile | null = null;
	private _activeProfileName: string | null = null;
	private _tier: string = "";

	constructor(config: SessionConfig = {}) {
		this._namespace = config.namespace ?? this.getDefaultNamespace();
		this._contextPath = new ContextPath();
		this._validator = new ContextValidator();
		this._profileManager = getProfileManager();
		this._serverUrl =
			config.serverUrl ?? process.env[`${ENV_PREFIX}_API_URL`] ?? "";
		this._apiToken =
			config.apiToken ?? process.env[`${ENV_PREFIX}_API_TOKEN`] ?? "";
		this._outputFormat = config.outputFormat ?? "yaml";
		this._debug =
			config.debug ?? process.env[`${ENV_PREFIX}_DEBUG`] === "true";

		// Extract tenant from server URL if available
		if (this._serverUrl) {
			this._tenant = this.extractTenant(this._serverUrl);
		}

		// Create API client if we have server URL
		if (this._serverUrl) {
			this._apiClient = new APIClient({
				serverUrl: this._serverUrl,
				apiToken: this._apiToken,
				debug: this._debug,
			});
		}
	}

	/**
	 * Initialize the session (async operations)
	 */
	async initialize(): Promise<void> {
		// Initialize history manager
		try {
			this._history = await HistoryManager.create(
				getHistoryFilePath(),
				1000,
			);
		} catch (error) {
			console.error("Warning: could not initialize history:", error);
			this._history = new HistoryManager(getHistoryFilePath(), 1000);
		}

		// Load active profile if one is set
		await this.loadActiveProfile();

		// Fetch user info if connected and authenticated
		if (this._apiClient?.isAuthenticated()) {
			await this.fetchUserInfo();
			await this.fetchTier();
		}
	}

	/**
	 * Fetch subscription tier from the API
	 */
	private async fetchTier(): Promise<void> {
		if (!this._apiClient) return;

		try {
			const subscriptionClient = new SubscriptionClient(this._apiClient);
			this._tier = await subscriptionClient.getTierFromCurrentPlan();
		} catch {
			// Ignore tier fetch errors - not critical for session
		}
	}

	/**
	 * Fetch user info from the API
	 */
	private async fetchUserInfo(): Promise<void> {
		if (!this._apiClient) return;

		try {
			const response = await this._apiClient.get<{
				email?: string;
				name?: string;
				username?: string;
			}>("/api/web/custom/user/info");

			if (response.ok && response.data) {
				this._username =
					response.data.email ||
					response.data.name ||
					response.data.username ||
					"";
			}
		} catch {
			// Ignore user info fetch errors - not critical for session
		}
	}

	/**
	 * Load the active profile from profile manager
	 */
	async loadActiveProfile(): Promise<void> {
		try {
			const activeName = await this._profileManager.getActive();
			if (activeName) {
				const profile = await this._profileManager.get(activeName);
				if (profile) {
					this._activeProfileName = activeName;
					this._activeProfile = profile;

					// Apply profile settings to session
					if (profile.apiUrl) {
						this._serverUrl = profile.apiUrl;
						this._tenant = this.extractTenant(profile.apiUrl);
					}
					if (profile.apiToken) {
						this._apiToken = profile.apiToken;
					}
					if (profile.defaultNamespace) {
						this._namespace = profile.defaultNamespace;
					}

					// Recreate API client with profile settings
					if (this._serverUrl) {
						this._apiClient = new APIClient({
							serverUrl: this._serverUrl,
							apiToken: this._apiToken,
							debug: this._debug,
						});
					}
				}
			}
		} catch {
			// Ignore profile loading errors - session can work without profile
		}
	}

	/**
	 * Get the default namespace from environment or config
	 */
	private getDefaultNamespace(): string {
		return process.env[`${ENV_PREFIX}_NAMESPACE`] ?? "default";
	}

	/**
	 * Extract tenant name from server URL
	 */
	private extractTenant(url: string): string {
		try {
			const parsed = new URL(url);
			const hostname = parsed.hostname;

			// Extract subdomain as tenant (e.g., "mycompany" from "mycompany.console.ves.volterra.io")
			const parts = hostname.split(".");
			if (parts.length > 0 && parts[0]) {
				return parts[0];
			}
			return hostname;
		} catch {
			return "";
		}
	}

	/**
	 * Set the default namespace for the session
	 */
	setNamespace(ns: string): void {
		this._namespace = ns;
	}

	/**
	 * Get the current default namespace
	 */
	getNamespace(): string {
		return this._namespace;
	}

	/**
	 * Get the exit code of the last command
	 */
	getLastExitCode(): number {
		return this._lastExitCode;
	}

	/**
	 * Set the exit code of the last command
	 */
	setLastExitCode(code: number): void {
		this._lastExitCode = code;
	}

	/**
	 * Get the current navigation context
	 */
	getContextPath(): ContextPath {
		return this._contextPath;
	}

	/**
	 * Get the current tenant name
	 */
	getTenant(): string {
		return this._tenant;
	}

	/**
	 * Get the logged-in user's name/email
	 */
	getUsername(): string {
		return this._username;
	}

	/**
	 * Set the username (used when fetched from API)
	 */
	setUsername(username: string): void {
		this._username = username;
	}

	/**
	 * Get the subscription tier (Standard/Advanced)
	 */
	getTier(): string {
		return this._tier;
	}

	/**
	 * Get the context validator
	 */
	getValidator(): ContextValidator {
		return this._validator;
	}

	/**
	 * Get the history manager
	 */
	getHistory(): HistoryManager | null {
		return this._history;
	}

	/**
	 * Get the server URL
	 */
	getServerUrl(): string {
		return this._serverUrl;
	}

	/**
	 * Check if connected to an API server
	 */
	isConnected(): boolean {
		return this._serverUrl !== "" && this._apiClient !== null;
	}

	/**
	 * Check if authenticated with API
	 */
	isAuthenticated(): boolean {
		return this._apiClient?.isAuthenticated() ?? false;
	}

	/**
	 * Get the API client
	 */
	getAPIClient(): APIClient | null {
		return this._apiClient;
	}

	/**
	 * Get the current output format
	 */
	getOutputFormat(): OutputFormat {
		return this._outputFormat;
	}

	/**
	 * Set the output format
	 */
	setOutputFormat(format: OutputFormat): void {
		this._outputFormat = format;
	}

	/**
	 * Get debug mode status
	 */
	isDebug(): boolean {
		return this._debug;
	}

	/**
	 * Get the profile manager
	 */
	getProfileManager(): ProfileManager {
		return this._profileManager;
	}

	/**
	 * Get the active profile
	 */
	getActiveProfile(): Profile | null {
		return this._activeProfile;
	}

	/**
	 * Get the active profile name
	 */
	getActiveProfileName(): string | null {
		return this._activeProfileName;
	}

	/**
	 * Switch to a different profile
	 */
	async switchProfile(profileName: string): Promise<boolean> {
		const profile = await this._profileManager.get(profileName);
		if (!profile) {
			return false;
		}

		const result = await this._profileManager.setActive(profileName);
		if (!result.success) {
			return false;
		}

		// Update session with new profile settings
		this._activeProfileName = profileName;
		this._activeProfile = profile;

		if (profile.apiUrl) {
			this._serverUrl = profile.apiUrl;
			this._tenant = this.extractTenant(profile.apiUrl);
		}
		if (profile.apiToken) {
			this._apiToken = profile.apiToken;
		}
		if (profile.defaultNamespace) {
			this._namespace = profile.defaultNamespace;
		}

		// Recreate API client with new profile settings
		if (this._serverUrl) {
			this._apiClient = new APIClient({
				serverUrl: this._serverUrl,
				apiToken: this._apiToken,
				debug: this._debug,
			});
		} else {
			this._apiClient = null;
		}

		return true;
	}

	/**
	 * Clear the active profile
	 */
	clearActiveProfile(): void {
		this._activeProfileName = null;
		this._activeProfile = null;
	}

	/**
	 * Add a command to history
	 */
	addToHistory(cmd: string): void {
		this._history?.add(cmd);
	}

	/**
	 * Save history to disk
	 */
	async saveHistory(): Promise<void> {
		await this._history?.save();
	}

	/**
	 * Clean up session resources
	 */
	async cleanup(): Promise<void> {
		await this.saveHistory();
	}
}

/**
 * Create and initialize a new REPL session
 */
export async function createSession(
	config: SessionConfig = {},
): Promise<REPLSession> {
	const session = new REPLSession(config);
	await session.initialize();
	return session;
}
