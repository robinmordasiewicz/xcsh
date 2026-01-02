/**
 * GenAI API Client
 *
 * Client for interacting with F5 Distributed Cloud AI Assistant APIs
 */

import type { APIClient } from "../../api/client.js";
import type {
	GenAIQueryRequest,
	GenAIQueryResponse,
	GenAIFeedbackRequest,
} from "./types.js";

/**
 * GenAI API Client
 *
 * Provides methods for querying the AI assistant and submitting feedback
 */
export class GenAIClient {
	constructor(private apiClient: APIClient) {}

	/**
	 * Query the AI assistant
	 *
	 * @param namespace - The namespace context for the query
	 * @param query - The natural language query
	 * @returns The AI assistant response with query_id and response data
	 */
	async query(namespace: string, query: string): Promise<GenAIQueryResponse> {
		const request: GenAIQueryRequest = {
			current_query: query,
			namespace,
		};

		const response = await this.apiClient.post<GenAIQueryResponse>(
			`/api/gen-ai/namespaces/${namespace}/query`,
			request,
		);

		if (!response.ok) {
			throw new Error(
				`GenAI query failed: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Submit feedback for a query
	 *
	 * @param request - The feedback request with query_id and feedback type
	 */
	async feedback(request: GenAIFeedbackRequest): Promise<void> {
		const response = await this.apiClient.post<Record<string, unknown>>(
			`/api/gen-ai/namespaces/${request.namespace}/query_feedback`,
			request,
		);

		if (!response.ok) {
			throw new Error(
				`GenAI feedback failed: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}

	/**
	 * Query the AI assistant in eval mode (for RBAC testing)
	 *
	 * @param namespace - The namespace context for the query
	 * @param query - The natural language query
	 * @returns The AI assistant response
	 */
	async evalQuery(
		namespace: string,
		query: string,
	): Promise<GenAIQueryResponse> {
		const request: GenAIQueryRequest = {
			current_query: query,
			namespace,
		};

		const response = await this.apiClient.post<GenAIQueryResponse>(
			`/api/gen-ai/namespaces/${namespace}/eval_query`,
			request,
		);

		if (!response.ok) {
			throw new Error(
				`GenAI eval query failed: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Submit feedback for an eval query
	 *
	 * @param request - The feedback request with query_id and feedback type
	 */
	async evalFeedback(request: GenAIFeedbackRequest): Promise<void> {
		const response = await this.apiClient.post<Record<string, unknown>>(
			`/api/gen-ai/namespaces/${request.namespace}/eval_query_feedback`,
			request,
		);

		if (!response.ok) {
			throw new Error(
				`GenAI eval feedback failed: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}
}

/**
 * Cached client instance
 */
let cachedClient: GenAIClient | null = null;

/**
 * Get or create a GenAI client instance
 *
 * Uses lazy initialization pattern to avoid creating client until needed
 *
 * @param apiClient - The API client from the session
 * @returns The GenAI client instance
 */
export function getGenAIClient(apiClient: APIClient): GenAIClient {
	if (!cachedClient) {
		cachedClient = new GenAIClient(apiClient);
	}
	return cachedClient;
}

/**
 * Reset the cached client (for testing or session changes)
 */
export function resetGenAIClient(): void {
	cachedClient = null;
}
