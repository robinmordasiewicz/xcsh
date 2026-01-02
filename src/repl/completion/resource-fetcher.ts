/**
 * Resource Fetcher - Live API integration for resource name completions
 * Fetches resource names from F5 XC API with caching
 */

import type { APIClient } from "../../api/index.js";

/**
 * Cache entry for resource names
 */
interface CacheEntry {
	names: string[];
	timestamp: number;
	loading: boolean;
}

/**
 * Default cache TTL (30 seconds for fresh data during completion)
 */
const CACHE_TTL = 30 * 1000;

/**
 * ResourceFetcher handles live fetching of resource names from the API
 */
export class ResourceFetcher {
	private cache: Map<string, CacheEntry> = new Map();
	private pendingRequests: Map<string, Promise<string[]>> = new Map();

	/**
	 * Convert resource type (snake_case) to API path (kebab-case with plural)
	 * Examples:
	 *   http_loadbalancer -> http-loadbalancers
	 *   origin_pool -> origin-pools
	 *   app_firewall -> app-firewalls
	 */
	private resourceTypeToApiPath(resourceType: string): string {
		// Convert snake_case to kebab-case
		const kebabCase = resourceType.replace(/_/g, "-");
		// Add 's' for plural form (most F5 XC resources are plural in API)
		return kebabCase.endsWith("s") ? kebabCase : `${kebabCase}s`;
	}

	/**
	 * Generate cache key from namespace and resource type
	 */
	private getCacheKey(namespace: string, resourceType: string): string {
		return `${namespace}:${resourceType}`;
	}

	/**
	 * Check if cache entry is still valid
	 */
	private isCacheValid(entry: CacheEntry | undefined): boolean {
		if (!entry) return false;
		return Date.now() - entry.timestamp < CACHE_TTL;
	}

	/**
	 * Check if a resource type is currently being loaded
	 */
	isLoading(namespace: string, resourceType: string): boolean {
		const key = this.getCacheKey(namespace, resourceType);
		const entry = this.cache.get(key);
		return entry?.loading ?? false;
	}

	/**
	 * Fetch resource names from API with caching and deduplication
	 */
	async fetchResourceNames(
		client: APIClient | null,
		namespace: string,
		resourceType: string,
		partial: string = "",
	): Promise<string[]> {
		const cacheKey = this.getCacheKey(namespace, resourceType);

		// Check cache first
		const cached = this.cache.get(cacheKey);
		if (cached && this.isCacheValid(cached)) {
			return this.filterByPrefix(cached.names, partial);
		}

		// Check if there's already a pending request for this resource type
		const pending = this.pendingRequests.get(cacheKey);
		if (pending) {
			const names = await pending;
			return this.filterByPrefix(names, partial);
		}

		// No client or not authenticated - return empty
		if (!client?.isAuthenticated()) {
			return [];
		}

		// Mark as loading (preserve existing data if available)
		const existingNames = cached?.names ?? [];
		const existingTimestamp = cached?.timestamp ?? 0;
		this.cache.set(cacheKey, {
			names: existingNames,
			timestamp: existingTimestamp,
			loading: true,
		});

		// Create the fetch promise
		const fetchPromise = this.doFetch(client, namespace, resourceType);
		this.pendingRequests.set(cacheKey, fetchPromise);

		try {
			const names = await fetchPromise;

			// Cache the result
			this.cache.set(cacheKey, {
				names,
				timestamp: Date.now(),
				loading: false,
			});

			return this.filterByPrefix(names, partial);
		} catch {
			// On error, clear loading state and return empty
			const existing = this.cache.get(cacheKey);
			if (existing) {
				existing.loading = false;
			}
			return [];
		} finally {
			// Clean up pending request
			this.pendingRequests.delete(cacheKey);
		}
	}

	/**
	 * Internal fetch implementation
	 */
	private async doFetch(
		client: APIClient,
		namespace: string,
		resourceType: string,
	): Promise<string[]> {
		const apiPath = this.resourceTypeToApiPath(resourceType);
		const endpoint = `/api/config/namespaces/${namespace}/${apiPath}`;

		try {
			const response = await client.get<{
				items?: Array<{ name?: string; metadata?: { name?: string } }>;
			}>(endpoint);

			if (response.ok && response.data?.items) {
				return response.data.items
					.map((item) => item.name ?? item.metadata?.name)
					.filter((name): name is string => !!name)
					.sort();
			}
		} catch {
			// Silently fail - resource may not exist or access denied
		}

		return [];
	}

	/**
	 * Filter names by prefix (case-insensitive)
	 */
	private filterByPrefix(names: string[], prefix: string): string[] {
		if (!prefix) return names;
		const lowerPrefix = prefix.toLowerCase();
		return names.filter((name) =>
			name.toLowerCase().startsWith(lowerPrefix),
		);
	}

	/**
	 * Prefetch resource names for a domain's primary resources
	 * Call this when entering a domain context to warm the cache
	 */
	async prefetch(
		client: APIClient | null,
		namespace: string,
		resourceTypes: string[],
	): Promise<void> {
		if (!client?.isAuthenticated()) return;

		// Prefetch in parallel (fire and forget)
		await Promise.allSettled(
			resourceTypes.map((rt) =>
				this.fetchResourceNames(client, namespace, rt),
			),
		);
	}

	/**
	 * Clear cache for a specific namespace (e.g., when switching namespaces)
	 */
	clearNamespaceCache(namespace: string): void {
		for (const key of this.cache.keys()) {
			if (key.startsWith(`${namespace}:`)) {
				this.cache.delete(key);
			}
		}
	}

	/**
	 * Clear all cached data
	 */
	clearAll(): void {
		this.cache.clear();
		this.pendingRequests.clear();
	}
}

/**
 * Singleton instance for shared caching across completions
 */
export const resourceFetcher = new ResourceFetcher();
