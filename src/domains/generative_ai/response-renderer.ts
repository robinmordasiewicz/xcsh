/**
 * GenAI Response Renderer
 *
 * Format and display the various response types from the AI Assistant
 */

import type {
	GenAIQueryResponse,
	GenericResponse,
	ExplainLogRecordResponse,
	GenDashboardFilterResponse,
	ListResponse,
	SiteAnalysisResponse,
	WidgetResponse,
} from "./types.js";
import { getResponseType } from "./types.js";

/**
 * Render a GenAI response for display
 *
 * Routes to the appropriate renderer based on response type
 */
export function renderResponse(response: GenAIQueryResponse): string[] {
	const lines: string[] = [];

	const responseType = getResponseType(response);

	switch (responseType) {
		case "generic_response":
			lines.push(...renderGenericResponse(response.generic_response!));
			break;
		case "explain_log":
			lines.push(...renderExplainLog(response.explain_log!));
			break;
		case "gen_dashboard_filter":
			lines.push(
				...renderDashboardFilter(response.gen_dashboard_filter!),
			);
			break;
		case "list_response":
			lines.push(...renderListResponse(response.list_response!));
			break;
		case "site_analysis_response":
			lines.push(...renderSiteAnalysis(response.site_analysis_response!));
			break;
		case "widget_response":
			lines.push(...renderWidgetResponse(response.widget_response!));
			break;
		default:
			lines.push("No response content.");
	}

	// Add follow-up suggestions if available
	if (response.follow_up_queries && response.follow_up_queries.length > 0) {
		lines.push("");
		lines.push("Suggested follow-up questions:");
		response.follow_up_queries.forEach((q, i) => {
			lines.push(`  ${i + 1}. ${q}`);
		});
	}

	return lines;
}

/**
 * Render a generic text response
 */
function renderGenericResponse(response: GenericResponse): string[] {
	const lines: string[] = [];

	if (response.error) {
		lines.push(`Error: ${response.error}`);
		return lines;
	}

	if (response.text) {
		// Split text by newlines for proper rendering
		lines.push(...response.text.split("\n"));
	}

	if (response.links && response.links.length > 0) {
		lines.push("");
		lines.push("Related links:");
		for (const link of response.links) {
			lines.push(`  - ${link.title}: ${link.url}`);
		}
	}

	return lines;
}

/**
 * Render a security event explanation
 */
function renderExplainLog(response: ExplainLogRecordResponse): string[] {
	const lines: string[] = [];

	lines.push("=== Security Event Analysis ===");
	lines.push("");

	if (response.summary) {
		lines.push("Summary:");
		lines.push(`  ${response.summary}`);
		lines.push("");
	}

	if (response.action) {
		lines.push(`Action Taken: ${response.action}`);
	}

	if (response.accuracy) {
		lines.push(`Accuracy: ${response.accuracy}`);
	}

	if (response.violations && response.violations.length > 0) {
		lines.push("");
		lines.push("Violations Detected:");
		for (const v of response.violations) {
			lines.push(`  - ${v.name}: ${v.description}`);
		}
	}

	if (response.threat_campaigns && response.threat_campaigns.length > 0) {
		lines.push("");
		lines.push("Threat Campaigns:");
		for (const campaign of response.threat_campaigns) {
			lines.push(`  - ${campaign}`);
		}
	}

	if (response.request_details) {
		lines.push("");
		lines.push("Request Details:");
		lines.push(`  ${JSON.stringify(response.request_details, null, 2)}`);
	}

	return lines;
}

/**
 * Render a dashboard filter expression
 */
function renderDashboardFilter(response: GenDashboardFilterResponse): string[] {
	const lines: string[] = [];

	lines.push("=== Dashboard Filter ===");
	lines.push("");

	if (response.filter_expression) {
		lines.push("Filter Expression:");
		lines.push(`  ${response.filter_expression}`);
	}

	if (response.dashboard_context) {
		lines.push("");
		lines.push("Dashboard Context:");
		lines.push(`  ${response.dashboard_context}`);
	}

	return lines;
}

/**
 * Render a list response as a table
 */
function renderListResponse(response: ListResponse): string[] {
	const lines: string[] = [];

	if (response.formatted_list) {
		lines.push(response.formatted_list);
		return lines;
	}

	if (response.items && response.items.length > 0) {
		if (response.total_count !== undefined) {
			lines.push(`Total: ${response.total_count} items`);
			lines.push("");
		}

		for (const item of response.items) {
			lines.push(JSON.stringify(item, null, 2));
			lines.push("");
		}
	} else {
		lines.push("No items found.");
	}

	return lines;
}

/**
 * Render a site analysis response
 */
function renderSiteAnalysis(response: SiteAnalysisResponse): string[] {
	const lines: string[] = [];

	lines.push("=== Site Analysis ===");
	lines.push("");

	if (response.site_name) {
		lines.push(`Site: ${response.site_name}`);
	}

	if (response.health_status) {
		lines.push(`Health: ${response.health_status}`);
	}

	if (response.metrics) {
		lines.push("");
		lines.push("Metrics:");
		for (const [key, value] of Object.entries(response.metrics)) {
			lines.push(`  ${key}: ${JSON.stringify(value)}`);
		}
	}

	if (response.recommendations && response.recommendations.length > 0) {
		lines.push("");
		lines.push("Recommendations:");
		for (const rec of response.recommendations) {
			lines.push(`  - ${rec}`);
		}
	}

	return lines;
}

/**
 * Render a widget response
 */
function renderWidgetResponse(response: WidgetResponse): string[] {
	const lines: string[] = [];

	if (response.display_type) {
		lines.push(`Widget Type: ${response.display_type}`);
		lines.push("");
	}

	if (response.table) {
		// Render table with rows and cells
		for (const row of response.table.rows) {
			const cells = row.cells.map((cell) => {
				let value = cell.value;
				if (cell.properties?.status_style) {
					value = `[${cell.properties.status_style}] ${value}`;
				}
				return value;
			});
			lines.push(cells.join(" | "));
		}
	}

	if (response.links && response.links.length > 0) {
		lines.push("");
		lines.push("Links:");
		for (const link of response.links) {
			lines.push(`  - ${link.title}: ${link.url}`);
		}
	}

	return lines;
}

/**
 * Format response as compact single-line summary
 */
export function renderResponseCompact(response: GenAIQueryResponse): string {
	const responseType = getResponseType(response);

	switch (responseType) {
		case "generic_response":
			if (response.generic_response?.text) {
				const text = response.generic_response.text;
				return text.length > 100 ? text.slice(0, 100) + "..." : text;
			}
			return response.generic_response?.error ?? "No content";

		case "explain_log":
			return (
				response.explain_log?.summary ??
				`Action: ${response.explain_log?.action}`
			);

		case "gen_dashboard_filter":
			return (
				response.gen_dashboard_filter?.filter_expression ??
				"Filter generated"
			);

		case "list_response":
			return `${response.list_response?.total_count ?? response.list_response?.items?.length ?? 0} items`;

		case "site_analysis_response":
			return `${response.site_analysis_response?.site_name}: ${response.site_analysis_response?.health_status}`;

		case "widget_response":
			return `Widget: ${response.widget_response?.display_type ?? "table"}`;

		default:
			return "No response";
	}
}

/**
 * Extract plain text from response for chat history
 */
export function extractResponseText(response: GenAIQueryResponse): string {
	const responseType = getResponseType(response);

	switch (responseType) {
		case "generic_response":
			return (
				response.generic_response?.text ??
				response.generic_response?.error ??
				""
			);

		case "explain_log":
			return response.explain_log?.summary ?? "";

		case "gen_dashboard_filter":
			return response.gen_dashboard_filter?.filter_expression ?? "";

		case "list_response":
			return response.list_response?.formatted_list ?? "";

		case "site_analysis_response": {
			const site = response.site_analysis_response;
			return site?.recommendations?.join("; ") ?? "";
		}

		case "widget_response":
			return "";

		default:
			return "";
	}
}
