/**
 * CompletionCache - Caches API responses for tab completion
 * Implements TTL-based caching to reduce API calls while keeping completions fresh
 */

/**
 * Default TTL for cached completion data (5 minutes)
 */
const DEFAULT_TTL_MS = 5 * 60 * 1000;

/**
 * Cache entry with expiration
 */
interface CacheEntry<T> {
	data: T;
	expiresAt: number;
}

/**
 * CompletionCache provides TTL-based caching for completion API lookups
 */
export class CompletionCache {
	private namespaceCache: CacheEntry<string[]> | null = null;
	private resourceCache: Map<string, CacheEntry<string[]>> = new Map();
	private readonly ttlMs: number;

	constructor(ttlMs: number = DEFAULT_TTL_MS) {
		this.ttlMs = ttlMs;
	}

	/**
	 * Check if a cache entry is still valid
	 */
	private isValid<T>(
		entry: CacheEntry<T> | null | undefined,
	): entry is CacheEntry<T> {
		return (
			entry !== null &&
			entry !== undefined &&
			Date.now() < entry.expiresAt
		);
	}

	/**
	 * Get namespaces with caching
	 * @param fetcher - Async function to fetch namespaces if cache miss
	 */
	async getNamespaces(fetcher: () => Promise<string[]>): Promise<string[]> {
		if (this.isValid(this.namespaceCache)) {
			return this.namespaceCache.data;
		}

		try {
			const namespaces = await fetcher();
			this.namespaceCache = {
				data: namespaces,
				expiresAt: Date.now() + this.ttlMs,
			};
			return namespaces;
		} catch {
			// Fall back to common defaults on error
			return ["default", "system"];
		}
	}

	/**
	 * Get resource names with caching
	 * @param key - Cache key (e.g., "domain:resourceType")
	 * @param fetcher - Async function to fetch resource names if cache miss
	 */
	async getResourceNames(
		key: string,
		fetcher: () => Promise<string[]>,
	): Promise<string[]> {
		const cached = this.resourceCache.get(key);
		if (this.isValid(cached)) {
			return cached.data;
		}

		try {
			const names = await fetcher();
			this.resourceCache.set(key, {
				data: names,
				expiresAt: Date.now() + this.ttlMs,
			});
			return names;
		} catch {
			return [];
		}
	}

	/**
	 * Clear all cached data
	 */
	clear(): void {
		this.namespaceCache = null;
		this.resourceCache.clear();
	}

	/**
	 * Clear namespace cache only
	 */
	clearNamespaces(): void {
		this.namespaceCache = null;
	}

	/**
	 * Clear resource cache for a specific key
	 */
	clearResources(key?: string): void {
		if (key) {
			this.resourceCache.delete(key);
		} else {
			this.resourceCache.clear();
		}
	}

	/**
	 * Get cache statistics (for debugging)
	 */
	getStats(): {
		namespaceCached: boolean;
		resourceKeys: string[];
	} {
		return {
			namespaceCached: this.isValid(this.namespaceCache),
			resourceKeys: Array.from(this.resourceCache.keys()).filter((key) =>
				this.isValid(this.resourceCache.get(key)),
			),
		};
	}
}

/**
 * Global completion cache instance
 */
let completionCache: CompletionCache | null = null;

/**
 * Get the global completion cache instance
 */
export function getCompletionCache(): CompletionCache {
	if (!completionCache) {
		completionCache = new CompletionCache();
	}
	return completionCache;
}

export default CompletionCache;
