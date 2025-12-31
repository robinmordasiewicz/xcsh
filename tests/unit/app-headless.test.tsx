/**
 * Headless tests for the REPL App component using ink-testing-library
 *
 * These tests verify the App component renders correctly without requiring a TTY.
 * Note: The token warning is displayed in src/index.tsx BEFORE the App renders,
 * so this test focuses on verifying the App component behavior with pre-initialized sessions.
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import React from "react";
import { render } from "ink-testing-library";
import { App } from "../../src/repl/App.js";
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

// Mock executor to avoid actual command execution
vi.mock("../../src/repl/executor.js", () => ({
	executeCommand: vi.fn().mockResolvedValue({
		output: [],
		shouldExit: false,
		shouldClear: false,
		error: false,
	}),
}));

// Mock git info
vi.mock("../../src/repl/components/StatusBar.js", async () => {
	const actual = await vi.importActual("../../src/repl/components/StatusBar.js");
	return {
		...actual,
		getGitInfo: vi.fn().mockReturnValue(undefined),
	};
});

describe("App Headless Rendering", () => {
	beforeEach(() => {
		// Set up environment variables
		vi.stubEnv("F5XC_API_URL", "https://test.volterra.io");
		vi.stubEnv("F5XC_API_TOKEN", "test-token");
	});

	afterEach(() => {
		vi.unstubAllEnvs();
		vi.restoreAllMocks();
	});

	describe("Pre-initialized Session", () => {
		it("should render App with pre-initialized session", async () => {
			// Create and initialize session
			const session = new REPLSession();

			// Mock valid token
			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// Render App with pre-initialized session
			const { lastFrame } = render(<App initialSession={session} />);

			// Should render something (App should initialize)
			expect(lastFrame()).toBeDefined();
		});

		it("should show prompt when session is initialized", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			const { lastFrame } = render(<App initialSession={session} />);

			// Should contain prompt character
			const frame = lastFrame() ?? "";
			expect(frame.includes(">") || frame.includes("Initializing")).toBe(true);
		});
	});

	describe("Session State Verification", () => {
		it("should have correct session state when token is valid", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// Verify session state
			expect(session.isAuthenticated()).toBe(true);
			expect(session.isTokenValidated()).toBe(true);
			expect(session.getValidationError()).toBeNull();

			// This is the warning condition - should be false
			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();
			expect(showsWarning).toBeFalsy();

			// Now render App
			const { lastFrame } = render(<App initialSession={session} />);
			expect(lastFrame()).toBeDefined();
		});

		it("should have correct session state when token is invalid", async () => {
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: false,
					error: "Invalid token",
				});
			}

			await session.initialize();

			// Verify session state
			expect(session.isAuthenticated()).toBe(true);
			expect(session.isTokenValidated()).toBe(false);
			expect(session.getValidationError()).toBe("Invalid token");

			// This is the warning condition - should be true
			const showsWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();
			expect(showsWarning).toBeTruthy();

			// Render App (note: warning is displayed in index.tsx, not in App)
			const { lastFrame } = render(<App initialSession={session} />);
			expect(lastFrame()).toBeDefined();
		});
	});

	describe("Entry Point Warning Logic Simulation", () => {
		/**
		 * This test simulates the exact logic from src/index.tsx:115-125
		 * to verify the warning condition works correctly.
		 */
		it("should simulate REPL mode warning logic correctly", async () => {
			// Simulate what src/index.tsx does
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// This is the exact code from src/index.tsx:116-119
			const shouldShowWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			// With valid token, warning should NOT show
			expect(shouldShowWarning).toBeFalsy();
		});

		it("should simulate non-interactive mode warning logic correctly", async () => {
			// Simulate what executeNonInteractive does in src/index.tsx:176-189
			const session = new REPLSession();

			const apiClient = session.getAPIClient();
			if (apiClient) {
				vi.spyOn(apiClient, "validateToken").mockResolvedValue({
					valid: true,
					error: undefined,
				});
			}

			await session.initialize();

			// This is the exact code from src/index.tsx:181-184
			const shouldShowWarning =
				session.isAuthenticated() &&
				!session.isTokenValidated() &&
				session.getValidationError();

			// With valid token, warning should NOT show
			expect(shouldShowWarning).toBeFalsy();
		});

		it("REPL and non-interactive should have identical warning behavior", async () => {
			// Create two sessions to simulate both modes
			const replSession = new REPLSession();
			const nonIntSession = new REPLSession();

			// Mock both with same valid response
			[replSession, nonIntSession].forEach((session) => {
				const apiClient = session.getAPIClient();
				if (apiClient) {
					vi.spyOn(apiClient, "validateToken").mockResolvedValue({
						valid: true,
						error: undefined,
					});
				}
			});

			// Initialize both (same as both code paths do)
			await replSession.initialize();
			await nonIntSession.initialize();

			// Both should have identical state
			expect(replSession.isAuthenticated()).toBe(nonIntSession.isAuthenticated());
			expect(replSession.isTokenValidated()).toBe(nonIntSession.isTokenValidated());
			expect(replSession.getValidationError()).toBe(nonIntSession.getValidationError());

			// Both should have identical warning condition
			const replWarning =
				replSession.isAuthenticated() &&
				!replSession.isTokenValidated() &&
				replSession.getValidationError();

			const nonIntWarning =
				nonIntSession.isAuthenticated() &&
				!nonIntSession.isTokenValidated() &&
				nonIntSession.getValidationError();

			expect(replWarning).toBe(nonIntWarning);
			expect(replWarning).toBeFalsy();
		});
	});
});
