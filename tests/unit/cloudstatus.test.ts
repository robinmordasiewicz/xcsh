/**
 * Cloudstatus Domain Tests
 *
 * Comprehensive verification tests for the cloudstatus custom command group.
 * Tests output format handling, flag behavior, exit codes, and error handling.
 */

import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import type { REPLSession } from "../../src/repl/session.js";
import {
	mockStatusResponse,
	mockStatusResponseMinor,
	mockStatusResponseMajor,
	mockStatusResponseCritical,
	mockStatusResponseMaintenance,
	mockSummaryResponse,
	mockComponentsResponse,
	mockIncidentsResponse,
	mockIncidentsActiveOnly,
	mockMaintenancesResponse,
	mockMaintenancesUpcomingOnly,
} from "./cloudstatus-fixtures.js";

// Mock the cloudstatus client module before importing the domain
const mockGetStatus = vi.fn();
const mockGetSummary = vi.fn();
const mockGetComponents = vi.fn();
const mockGetIncidents = vi.fn();
const mockGetUnresolvedIncidents = vi.fn();
const mockGetMaintenances = vi.fn();
const mockGetUpcomingMaintenances = vi.fn();

vi.mock("../../src/cloudstatus/index.js", () => ({
	CloudstatusClient: vi.fn(() => ({
		getStatus: mockGetStatus,
		getSummary: mockGetSummary,
		getComponents: mockGetComponents,
		getIncidents: mockGetIncidents,
		getUnresolvedIncidents: mockGetUnresolvedIncidents,
		getMaintenances: mockGetMaintenances,
		getUpcomingMaintenances: mockGetUpcomingMaintenances,
	})),
	isIncidentActive: (incident: { status: string }) =>
		incident.status !== "resolved" && incident.status !== "postmortem",
	isMaintenanceActive: (maint: { status: string }) =>
		maint.status === "in_progress" || maint.status === "verifying",
	isMaintenanceUpcoming: (maint: { status: string }) =>
		maint.status === "scheduled",
	isMaintenanceCompleted: (maint: { status: string }) =>
		maint.status === "completed",
	isComponentDegraded: (comp: { status: string }) =>
		comp.status !== "operational",
	isComponentOperational: (comp: { status: string }) =>
		comp.status === "operational",
	statusIndicatorToExitCode: (indicator: string) => {
		switch (indicator) {
			case "none":
				return 0;
			case "minor":
				return 1;
			case "major":
				return 2;
			case "critical":
				return 3;
			case "maintenance":
				return 4;
			default:
				return 0;
		}
	},
}));

// Import after mocking
import { cloudstatusDomain } from "../../src/domains/cloudstatus/index.js";

// Create mock session with configurable output format
function createMockSession(outputFormat: string = "table"): REPLSession {
	return {
		getOutputFormat: () => outputFormat,
	} as unknown as REPLSession;
}

// Get command from domain
function getCommand(name: string) {
	return cloudstatusDomain.commands.get(name);
}

describe("cloudstatus domain", () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		vi.clearAllMocks();
	});

	describe("status command", () => {
		const command = getCommand("status");

		it("returns JSON output with --output json", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(parsed.status.indicator).toBe("none");
		});

		it("returns YAML output with --output yaml", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "yaml"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("indicator:");
			expect(yamlOutput).toContain("none");
		});

		it("returns TSV output with --output tsv", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "tsv"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns table output by default", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
			// Table format has box drawing characters
			const tableOutput = result.output.join("\n");
			expect(tableOutput).toContain("STATUS");
			expect(tableOutput).toContain("DESCRIPTION");
		});

		it("returns empty array with --output none", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "none"], session);

			expect(result.error).toBeUndefined();
			expect(result.output).toEqual([]);
		});

		it("returns exit code info with --quiet flag", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBe(1);
			expect(result.output[0]).toContain("Exit code:");
			expect(result.output[0]).toContain("0");
		});

		it("returns exit code 0 for healthy status", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponse);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output[0]).toContain("Exit code: 0");
		});

		it("returns exit code 1 for minor status", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponseMinor);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output[0]).toContain("Exit code: 1");
		});

		it("returns exit code 2 for major status", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponseMajor);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output[0]).toContain("Exit code: 2");
		});

		it("returns exit code 3 for critical status", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponseCritical);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output[0]).toContain("Exit code: 3");
		});

		it("returns exit code 4 for maintenance status", async () => {
			mockGetStatus.mockResolvedValueOnce(mockStatusResponseMaintenance);
			const session = createMockSession();
			const result = await command!.execute(["--quiet"], session);

			expect(result.error).toBeUndefined();
			expect(result.output[0]).toContain("Exit code: 4");
		});
	});

	describe("summary command", () => {
		const command = getCommand("summary");

		it("returns JSON output with --output json", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(parsed.status).toBeDefined();
			expect(parsed.components).toBeDefined();
		});

		it("returns YAML output with --output yaml", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "yaml"], session);

			expect(result.error).toBeUndefined();
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("status:");
			expect(yamlOutput).toContain("components:");
		});

		it("returns TSV output with --output tsv", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "tsv"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns full summary by default", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			expect(output).toContain("=== OVERALL STATUS ===");
			expect(output).toContain("=== COMPONENTS ===");
			expect(output).toContain("=== ACTIVE INCIDENTS ===");
		});

		it("returns empty array with --output none", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "none"], session);

			expect(result.error).toBeUndefined();
			expect(result.output).toEqual([]);
		});

		it("returns brief summary with --brief flag", async () => {
			mockGetSummary.mockResolvedValueOnce(mockSummaryResponse);
			const session = createMockSession();
			const result = await command!.execute(["--brief"], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			expect(output).toContain("Status:");
			expect(output).toContain("Components:");
			expect(output).toContain("Incidents:");
			expect(output).toContain("Maintenance:");
			// Brief should not have full section headers
			expect(output).not.toContain("=== OVERALL STATUS ===");
		});
	});

	describe("components command", () => {
		const command = getCommand("components");

		it("returns JSON output with --output json", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(Array.isArray(parsed)).toBe(true);
		});

		it("returns YAML output with --output yaml", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "yaml"], session);

			expect(result.error).toBeUndefined();
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("name:");
			expect(yamlOutput).toContain("status:");
		});

		it("returns TSV output with --output tsv", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "tsv"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns table output by default", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			expect(output).toContain("NAME");
			expect(output).toContain("STATUS");
		});

		it("returns empty array with --output none", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "none"], session);

			expect(result.error).toBeUndefined();
			expect(result.output).toEqual([]);
		});

		it("filters group components from output", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			const parsed = JSON.parse(jsonOutput);
			// Group components should be filtered out
			const groups = parsed.filter(
				(c: { group: boolean }) => c.group === true,
			);
			expect(groups.length).toBe(0);
		});

		it("shows only degraded with --degraded-only flag", async () => {
			mockGetComponents.mockResolvedValueOnce(mockComponentsResponse);
			const session = createMockSession();
			const result = await command!.execute(
				["--degraded-only", "--output", "json"],
				session,
			);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			const parsed = JSON.parse(jsonOutput);
			// All components should be non-operational
			for (const comp of parsed) {
				expect(comp.status).not.toBe("operational");
			}
		});
	});

	describe("incidents command", () => {
		const command = getCommand("incidents");

		it("returns JSON output with --output json", async () => {
			mockGetIncidents.mockResolvedValueOnce(mockIncidentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(Array.isArray(parsed)).toBe(true);
		});

		it("returns YAML output with --output yaml", async () => {
			mockGetIncidents.mockResolvedValueOnce(mockIncidentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "yaml"], session);

			expect(result.error).toBeUndefined();
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("name:");
			expect(yamlOutput).toContain("status:");
		});

		it("returns TSV output with --output tsv", async () => {
			mockGetIncidents.mockResolvedValueOnce(mockIncidentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "tsv"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns formatted text by default", async () => {
			mockGetIncidents.mockResolvedValueOnce(mockIncidentsResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			expect(output).toMatch(/\[(ACTIVE|RESOLVED)\]/);
			expect(output).toContain("Impact:");
		});

		it("returns empty array with --output none", async () => {
			mockGetIncidents.mockResolvedValueOnce(mockIncidentsResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "none"], session);

			expect(result.error).toBeUndefined();
			expect(result.output).toEqual([]);
		});

		it("shows only active with --active-only flag", async () => {
			mockGetUnresolvedIncidents.mockResolvedValueOnce(
				mockIncidentsActiveOnly,
			);
			const session = createMockSession();
			const result = await command!.execute(
				["--active-only", "--output", "json"],
				session,
			);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			const parsed = JSON.parse(jsonOutput);
			// All incidents should be active (not resolved or postmortem)
			for (const inc of parsed) {
				expect(inc.status).not.toBe("resolved");
				expect(inc.status).not.toBe("postmortem");
			}
		});
	});

	describe("maintenance command", () => {
		const command = getCommand("maintenance");

		it("returns JSON output with --output json", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "json"], session);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			expect(() => JSON.parse(jsonOutput)).not.toThrow();
			const parsed = JSON.parse(jsonOutput);
			expect(Array.isArray(parsed)).toBe(true);
		});

		it("returns YAML output with --output yaml", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "yaml"], session);

			expect(result.error).toBeUndefined();
			const yamlOutput = result.output.join("\n");
			expect(yamlOutput).toContain("name:");
			expect(yamlOutput).toContain("status:");
		});

		it("returns TSV output with --output tsv", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "tsv"], session);

			expect(result.error).toBeUndefined();
			expect(result.output.length).toBeGreaterThan(0);
		});

		it("returns formatted text by default", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			expect(output).toMatch(/\[(SCHEDULED|IN PROGRESS)\]/);
			expect(output).toContain("Scheduled:");
		});

		it("returns empty array with --output none", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute(["--output", "none"], session);

			expect(result.error).toBeUndefined();
			expect(result.output).toEqual([]);
		});

		it("hides completed maintenance by default", async () => {
			mockGetMaintenances.mockResolvedValueOnce(mockMaintenancesResponse);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeUndefined();
			const output = result.output.join("\n");
			// Completed maintenance should not appear in default table view
			expect(output).not.toContain("Security Patch Deployment");
		});

		it("shows only upcoming with --upcoming flag", async () => {
			mockGetUpcomingMaintenances.mockResolvedValueOnce(
				mockMaintenancesUpcomingOnly,
			);
			const session = createMockSession();
			const result = await command!.execute(
				["--upcoming", "--output", "json"],
				session,
			);

			expect(result.error).toBeUndefined();
			const jsonOutput = result.output.join("\n");
			const parsed = JSON.parse(jsonOutput);
			// All maintenance should be scheduled
			for (const maint of parsed) {
				expect(maint.status).toBe("scheduled");
			}
		});
	});

	describe("error handling", () => {
		it("handles API timeout gracefully", async () => {
			const command = getCommand("status");
			mockGetStatus.mockRejectedValueOnce(new Error("Request timeout"));
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeDefined();
			expect(result.output[0]).toContain("Failed to get status");
			expect(result.output[0]).toContain("Request timeout");
		});

		it("handles network errors gracefully", async () => {
			const command = getCommand("components");
			mockGetComponents.mockRejectedValueOnce(
				new Error("Network error: Unable to connect"),
			);
			const session = createMockSession();
			const result = await command!.execute([], session);

			expect(result.error).toBeDefined();
			expect(result.output[0]).toContain("Failed to get components");
		});
	});
});
