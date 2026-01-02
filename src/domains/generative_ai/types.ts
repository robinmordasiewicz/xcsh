/**
 * TypeScript types for the Generative AI domain
 *
 * Based on F5 Distributed Cloud OpenAPI spec for AI Assistant APIs
 */

/**
 * Negative feedback type enum
 */
export type NegativeFeedbackType =
	| "OTHER"
	| "INACCURATE_DATA"
	| "IRRELEVANT_CONTENT"
	| "POOR_FORMAT"
	| "SLOW_RESPONSE";

/**
 * Negative feedback details
 */
export interface NegativeFeedbackDetails {
	remarks: NegativeFeedbackType[];
}

/**
 * AI Assistant Query Request
 */
export interface GenAIQueryRequest {
	current_query: string;
	namespace: string;
	[key: string]: unknown;
}

/**
 * Generic text response
 */
export interface GenericResponse {
	text?: string;
	links?: Array<{
		title: string;
		url: string;
	}>;
	error?: string;
}

/**
 * Security event explanation response
 */
export interface ExplainLogRecordResponse {
	summary?: string;
	request_details?: Record<string, unknown>;
	violations?: Array<{
		name: string;
		description: string;
	}>;
	action?: "ACTION_NONE" | "ALLOW" | "BLOCK" | "REDIRECT";
	accuracy?: string;
	waf_event?: Record<string, unknown>;
	bot_defense_event?: Record<string, unknown>;
	service_policy_event?: Record<string, unknown>;
	ip_reputation?: Record<string, unknown>;
	threat_campaigns?: string[];
}

/**
 * Dashboard filter response
 */
export interface GenDashboardFilterResponse {
	filter_expression?: string;
	dashboard_context?: string;
}

/**
 * List response with items
 */
export interface ListResponse {
	items?: Array<Record<string, unknown>>;
	formatted_list?: string;
	total_count?: number;
}

/**
 * Site analysis response
 */
export interface SiteAnalysisResponse {
	site_name?: string;
	health_status?: string;
	metrics?: Record<string, unknown>;
	recommendations?: string[];
}

/**
 * Widget response for dashboard visualization
 */
export interface WidgetResponse {
	table?: {
		rows: Array<{
			cells: Array<{
				value: string;
				properties?: {
					status_style?: string;
				};
			}>;
		}>;
	};
	display_type?: string;
	links?: Array<{
		title: string;
		url: string;
	}>;
}

/**
 * Response type discriminator
 */
export type GenAIResponseType =
	| "generic_response"
	| "explain_log"
	| "gen_dashboard_filter"
	| "list_response"
	| "site_analysis_response"
	| "widget_response";

/**
 * AI Assistant Query Response
 */
export interface GenAIQueryResponse {
	query_id: string;
	current_query: string;
	follow_up_queries?: string[];
	// One of the following response types will be present
	generic_response?: GenericResponse;
	explain_log?: ExplainLogRecordResponse;
	gen_dashboard_filter?: GenDashboardFilterResponse;
	list_response?: ListResponse;
	site_analysis_response?: SiteAnalysisResponse;
	widget_response?: WidgetResponse;
}

/**
 * AI Assistant Query Feedback Request
 */
export interface GenAIFeedbackRequest {
	query: string;
	query_id: string;
	namespace: string;
	comment?: string | undefined;
	// One of positive_feedback or negative_feedback
	positive_feedback?: Record<string, never>; // Empty object for positive
	negative_feedback?: NegativeFeedbackDetails;
	[key: string]: unknown;
}

/**
 * Helper to determine which response type is present
 */
export function getResponseType(
	response: GenAIQueryResponse,
): GenAIResponseType | null {
	if (response.generic_response) return "generic_response";
	if (response.explain_log) return "explain_log";
	if (response.gen_dashboard_filter) return "gen_dashboard_filter";
	if (response.list_response) return "list_response";
	if (response.site_analysis_response) return "site_analysis_response";
	if (response.widget_response) return "widget_response";
	return null;
}

/**
 * Chat session state for tracking conversation
 */
export interface GenAIChatSession {
	namespace: string;
	lastQueryId: string | null;
	lastQuery: string | null;
	followUpQueries: string[];
	isActive: boolean;
}

/**
 * Mapping of user-friendly feedback type names to API values
 */
export const FEEDBACK_TYPE_MAP: Record<string, NegativeFeedbackType> = {
	other: "OTHER",
	inaccurate: "INACCURATE_DATA",
	irrelevant: "IRRELEVANT_CONTENT",
	poor_format: "POOR_FORMAT",
	slow: "SLOW_RESPONSE",
};

/**
 * Get all valid negative feedback type names (for help text)
 */
export function getValidFeedbackTypes(): string[] {
	return Object.keys(FEEDBACK_TYPE_MAP);
}
