/**
 * Unit tests for REPLSession token validation
 *
 * Tests the session initialization and token validation logic
 * to ensure the warning condition works correctly.
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { REPLSession } from "../../src/repl/session.js";

// Mock the profile manager to avoid file system operations
vi.mock("../../src/profile/index.js", () => ({
	getProfileManager: () => ({
		getActive: vi.fn().mockResolvedValue(null),
		list: vi.fn().mockResolvedValue([]),
		get: vi.fn().mockResolvedValue(null),
		setActive: vi.fn().mockResolvedValue({ success: true }),
	}),
}));

// Mock the history manager to avoid file system operations
vi.mock("../../src/repl/history.js", () => ({
	HistoryManager: {
		create: vi.fn().mockResolvedValue({
			add: vi.fn(),
			getHistory: vi.fn().mockReturnValue([]),
			save: vi.fn().mockResolvedValue(undefined),
		}),
	},
	getHistoryFilePath: vi.fn().mockReturnValue("/tmp/test-history"),
}));

describe("REPLSession Token Validation", () => {
	beforeEach(() => {
		// Set up environment variables for testing
		vi.stubEnv("F5XC_API_URL", "https://test.volterra.io");
		vi.stubEnv("F5XC_API_TOKEN", "valid-test-token");
		vi.stubEnv("F5XC_NAMESPACE", "default");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
		vi.restoreAllMocks();
	});

	describe("Session Creation", () => {
		it("should create session with API client when URL and token are set", () => {
			const session = new REPLSession();

			expect(session.getServerUrl()).toBe("https://test.volterra.io");
			expect(session.isConnected()).toBe(true);
			expect(session.getAPIClient()).not.toBeNull();
		});

		it("should report authenticated when token is present", () => {
			const session = new REPLSession();

			expect(session.isAuthenticated()).toBe(true);
		});

		it("should not be token validated before initialize()", () => {
			const session = new REPLSession();

			// Before initialize(), token has not been validated yet
			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBeNull();
		});
	});

	describe("Token Validation on Initialize", () => {
		it("should set tokenValidated=true when API returns valid", async () => {
			const session = new REPLSession();

			// Mock the validateToken method to return valid
			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			expect(session.isTokenValidated()).toBe(true);
			expect(session.getValidationError()).toBeNull();
		});

		it("should set tokenValidated=false with error when API returns invalid", async () => {
			const session = new REPLSession();

			// Mock the validateToken method to return invalid
			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Invalid or expired API token",
				});
			}

			await session.initialize();

			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBe(
				"Invalid or expired API token",
			);
		});
	});

	describe("Warning Condition Logic", () => {
		it("should NOT trigger warning when token is valid", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// This is the exact warning condition from src/index.tsx:115-119
			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			expect(showsWarning).toBeFalsy();
		});

		it("should trigger warning when token validation fails with error", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Invalid or expired API token",
				});
			}

			await session.initialize();

			// This is the exact warning condition from src/index.tsx:115-119
			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			expect(showsWarning).toBeTruthy();
		});

		it("should NOT trigger warning when token validation fails but no error message", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: undefined, // No error message
				});
			}

			await session.initialize();

			// Warning requires error message to be truthy
			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			expect(showsWarning).toBeFalsy();
		});

		it("should NOT trigger warning when not authenticated", async () => {
			// Clear environment to test unauthenticated state
			vi.stubEnv("F5XC_API_URL", "");
			vi.stubEnv("F5XC_API_TOKEN", "");

			const session = new REPLSession();
			await session.initialize();

			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			expect(showsWarning).toBeFalsy();
		});
	});

	describe("Session State After Initialization", () => {
		it("should have consistent state after valid token initialization", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// Verify all state is consistent
			expect(session.isConnected()).toBe(true);
			expect(session.isAuthenticated()).toBe(true);
			expect(session.isTokenValidated()).toBe(true);
			expect(session.getValidationError()).toBeNull();
			expect(session.getServerUrl()).toBe("https://test.volterra.io");
		});

		it("should have consistent state after invalid token initialization", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Token expired",
				});
			}

			await session.initialize();

			// Verify all state is consistent
			expect(session.isConnected()).toBe(true);
			expect(session.isAuthenticated()).toBe(true); // Has token, just invalid
			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBe("Token expired");
		});
	});

	describe("HTTP Status Code Handling", () => {
		it("should handle HTTP 401 as invalid token", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				// Simulate what the API client does on 401
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Invalid or expired API token",
				});
			}

			await session.initialize();

			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBe(
				"Invalid or expired API token",
			);
		});

		it("should handle HTTP 403 as invalid token", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				// Simulate what the API client does on 403
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Access forbidden",
				});
			}

			await session.initialize();

			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBe("Access forbidden");
		});

		it("should handle HTTP 200 as valid token", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				// Simulate what the API client does on 200
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			expect(session.isTokenValidated()).toBe(true);
			expect(session.getValidationError()).toBeNull();
		});
	});
});
