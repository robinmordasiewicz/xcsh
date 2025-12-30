/**
 * Unit tests for API client
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { APIClient, createClientFromEnv, buildResourcePath, APIError } from '../../src/api/client.js';

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('APIClient', () => {
	beforeEach(() => {
		mockFetch.mockReset();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	describe('constructor', () => {
		it('should normalize server URL by removing trailing slash', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com/',
			});
			expect(client.getServerUrl()).toBe('https://api.example.com');
		});

		it('should handle server URL without trailing slash', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com',
			});
			expect(client.getServerUrl()).toBe('https://api.example.com');
		});

		it('should remove multiple trailing slashes', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com///',
			});
			expect(client.getServerUrl()).toBe('https://api.example.com');
		});
	});

	describe('isAuthenticated', () => {
		it('should return true when API token is provided', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com',
				apiToken: 'test-token',
			});
			expect(client.isAuthenticated()).toBe(true);
		});

		it('should return false when no API token', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com',
			});
			expect(client.isAuthenticated()).toBe(false);
		});

		it('should return false for empty API token', () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com',
				apiToken: '',
			});
			expect(client.isAuthenticated()).toBe(false);
		});
	});

	describe('HTTP methods', () => {
		let client: APIClient;

		beforeEach(() => {
			client = new APIClient({
				serverUrl: 'https://api.example.com',
				apiToken: 'test-token',
			});
		});

		const mockSuccessResponse = (data: unknown) => {
			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 200,
				text: () => Promise.resolve(JSON.stringify(data)),
				headers: new Headers({ 'content-type': 'application/json' }),
			});
		};

		describe('GET', () => {
			it('should make GET request with correct URL', async () => {
				mockSuccessResponse({ items: [] });

				await client.get('/api/web/namespaces');

				expect(mockFetch).toHaveBeenCalledWith(
					'https://api.example.com/api/web/namespaces',
					expect.objectContaining({
						method: 'GET',
					})
				);
			});

			it('should include Authorization header', async () => {
				mockSuccessResponse({ items: [] });

				await client.get('/api/web/namespaces');

				expect(mockFetch).toHaveBeenCalledWith(
					expect.any(String),
					expect.objectContaining({
						headers: expect.objectContaining({
							'Authorization': 'APIToken test-token',
						}),
					})
				);
			});

			it('should include query parameters', async () => {
				mockSuccessResponse({ items: [] });

				await client.get('/api/web/namespaces', { limit: '10', offset: '0' });

				expect(mockFetch).toHaveBeenCalledWith(
					'https://api.example.com/api/web/namespaces?limit=10&offset=0',
					expect.any(Object)
				);
			});

			it('should return parsed JSON response', async () => {
				const responseData = { items: [{ name: 'test' }] };
				mockSuccessResponse(responseData);

				const response = await client.get('/api/web/namespaces');

				expect(response.statusCode).toBe(200);
				expect(response.ok).toBe(true);
				expect(response.data).toEqual(responseData);
			});
		});

		describe('POST', () => {
			it('should make POST request with body', async () => {
				mockSuccessResponse({ status: 'created' });

				const body = { metadata: { name: 'test' }, spec: {} };
				await client.post('/api/config/namespaces/shared/healthchecks', body);

				expect(mockFetch).toHaveBeenCalledWith(
					'https://api.example.com/api/config/namespaces/shared/healthchecks',
					expect.objectContaining({
						method: 'POST',
						body: JSON.stringify(body),
					})
				);
			});

			it('should set Content-Type header', async () => {
				mockSuccessResponse({ status: 'created' });

				await client.post('/api/test', { data: 'test' });

				expect(mockFetch).toHaveBeenCalledWith(
					expect.any(String),
					expect.objectContaining({
						headers: expect.objectContaining({
							'Content-Type': 'application/json',
						}),
					})
				);
			});
		});

		describe('PUT', () => {
			it('should make PUT request with body', async () => {
				mockSuccessResponse({ status: 'updated' });

				const body = { metadata: { name: 'test' }, spec: { updated: true } };
				await client.put('/api/config/namespaces/shared/healthchecks/test', body);

				expect(mockFetch).toHaveBeenCalledWith(
					expect.any(String),
					expect.objectContaining({
						method: 'PUT',
						body: JSON.stringify(body),
					})
				);
			});
		});

		describe('DELETE', () => {
			it('should make DELETE request', async () => {
				mockSuccessResponse({ status: 'deleted' });

				await client.delete('/api/config/namespaces/shared/healthchecks/test');

				expect(mockFetch).toHaveBeenCalledWith(
					'https://api.example.com/api/config/namespaces/shared/healthchecks/test',
					expect.objectContaining({
						method: 'DELETE',
					})
				);
			});

			it('should not include body', async () => {
				mockSuccessResponse({ status: 'deleted' });

				await client.delete('/api/test');

				expect(mockFetch).toHaveBeenCalledWith(
					expect.any(String),
					expect.objectContaining({
						body: null,
					})
				);
			});
		});

		describe('PATCH', () => {
			it('should make PATCH request with body', async () => {
				mockSuccessResponse({ status: 'patched' });

				const body = { spec: { partial: 'update' } };
				await client.patch('/api/config/namespaces/shared/test/item', body);

				expect(mockFetch).toHaveBeenCalledWith(
					expect.any(String),
					expect.objectContaining({
						method: 'PATCH',
						body: JSON.stringify(body),
					})
				);
			});
		});
	});

	describe('error handling', () => {
		let client: APIClient;

		beforeEach(() => {
			client = new APIClient({
				serverUrl: 'https://api.example.com',
				apiToken: 'test-token',
				timeout: 1000,
			});
		});

		it('should throw APIError for 401 Unauthorized', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 401,
				text: () => Promise.resolve(JSON.stringify({ message: 'Unauthorized' })),
				headers: new Headers(),
			});

			try {
				await client.get('/api/protected');
				expect.fail('Should have thrown');
			} catch (error) {
				expect(error).toBeInstanceOf(APIError);
				if (error instanceof APIError) {
					expect(error.statusCode).toBe(401);
					expect(error.message).toBe('Unauthorized');
				}
			}
		});

		it('should throw APIError for 403 Forbidden', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 403,
				text: () => Promise.resolve(JSON.stringify({ message: 'Forbidden' })),
				headers: new Headers(),
			});

			await expect(client.get('/api/admin')).rejects.toThrow(APIError);
		});

		it('should throw APIError for 404 Not Found', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 404,
				text: () => Promise.resolve(JSON.stringify({ message: 'Not found' })),
				headers: new Headers(),
			});

			await expect(client.get('/api/unknown')).rejects.toThrow(APIError);
		});

		it('should throw APIError for 500 Server Error', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 500,
				text: () => Promise.resolve(JSON.stringify({ message: 'Internal server error' })),
				headers: new Headers(),
			});

			await expect(client.get('/api/broken')).rejects.toThrow(APIError);
		});

		it('should handle network errors', async () => {
			mockFetch.mockRejectedValueOnce(new Error('Network failure'));

			await expect(client.get('/api/test')).rejects.toThrow(APIError);

			try {
				await client.get('/api/test');
			} catch (error) {
				if (error instanceof APIError) {
					expect(error.message).toContain('Network error');
				}
			}
		});

		it('should handle timeout', async () => {
			// Create a client with short timeout and no retries (to test timeout handling directly)
			const shortTimeoutClient = new APIClient({
				serverUrl: 'https://api.example.com',
				timeout: 10, // 10ms timeout
				retry: { maxRetries: 0 }, // Disable retries for this test
			});

			// Simulate AbortError from AbortController
			const abortError = new Error('The operation was aborted');
			abortError.name = 'AbortError';
			mockFetch.mockRejectedValueOnce(abortError);

			try {
				await shortTimeoutClient.get('/api/slow');
				expect.fail('Should have thrown');
			} catch (error) {
				expect(error).toBeInstanceOf(APIError);
				if (error instanceof APIError) {
					expect(error.message).toContain('timed out');
				}
			}
		});

		it('should include request context in APIError', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 400,
				text: () => Promise.resolve(JSON.stringify({ message: 'Bad request' })),
				headers: new Headers(),
			});

			try {
				await client.post('/api/create', { invalid: 'data' });
				expect.fail('Should have thrown');
			} catch (error) {
				expect(error).toBeInstanceOf(APIError);
				if (error instanceof APIError) {
					expect(error.operation).toBe('POST /api/create');
				}
			}
		});
	});

	describe('URL building edge cases', () => {
		it('should handle path without leading slash', async () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com',
			});

			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 200,
				text: () => Promise.resolve('{}'),
				headers: new Headers(),
			});

			await client.get('api/test');

			expect(mockFetch).toHaveBeenCalledWith(
				'https://api.example.com/api/test',
				expect.any(Object)
			);
		});

		it('should avoid duplicate /api in path', async () => {
			const client = new APIClient({
				serverUrl: 'https://api.example.com/api',
			});

			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 200,
				text: () => Promise.resolve('{}'),
				headers: new Headers(),
			});

			await client.get('/api/web/namespaces');

			expect(mockFetch).toHaveBeenCalledWith(
				'https://api.example.com/api/web/namespaces',
				expect.any(Object)
			);
		});
	});

	describe('response parsing', () => {
		let client: APIClient;

		beforeEach(() => {
			client = new APIClient({
				serverUrl: 'https://api.example.com',
			});
		});

		it('should parse valid JSON response', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 200,
				text: () => Promise.resolve('{"key": "value"}'),
				headers: new Headers({ 'content-type': 'application/json' }),
			});

			const response = await client.get('/api/test');
			expect(response.data).toEqual({ key: 'value' });
		});

		it('should handle empty response body', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 204,
				text: () => Promise.resolve(''),
				headers: new Headers(),
			});

			const response = await client.get('/api/test');
			expect(response.data).toEqual({});
		});

		it('should handle non-JSON response', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: true,
				status: 200,
				text: () => Promise.resolve('Plain text response'),
				headers: new Headers({ 'content-type': 'text/plain' }),
			});

			const response = await client.get('/api/test');
			expect(response.data).toBe('Plain text response');
		});
	});
});

describe('createClientFromEnv', () => {
	const originalEnv = process.env;

	beforeEach(() => {
		vi.resetModules();
		process.env = { ...originalEnv };
	});

	afterEach(() => {
		process.env = originalEnv;
	});

	it('should create client when F5XC_API_URL is set', () => {
		process.env.F5XC_API_URL = 'https://api.f5xc.com';
		process.env.F5XC_API_TOKEN = 'test-token';

		const client = createClientFromEnv();

		expect(client).not.toBeNull();
		expect(client?.getServerUrl()).toBe('https://api.f5xc.com');
		expect(client?.isAuthenticated()).toBe(true);
	});

	it('should return null when F5XC_API_URL is not set', () => {
		delete process.env.F5XC_API_URL;

		const client = createClientFromEnv();

		expect(client).toBeNull();
	});

	it('should create client without token', () => {
		process.env.F5XC_API_URL = 'https://api.f5xc.com';
		delete process.env.F5XC_API_TOKEN;

		const client = createClientFromEnv();

		expect(client).not.toBeNull();
		expect(client?.isAuthenticated()).toBe(false);
	});

	it('should use custom env prefix', () => {
		process.env.CUSTOM_API_URL = 'https://custom.api.com';
		process.env.CUSTOM_API_TOKEN = 'custom-token';

		const client = createClientFromEnv('CUSTOM');

		expect(client).not.toBeNull();
		expect(client?.getServerUrl()).toBe('https://custom.api.com');
	});

	it('should enable debug mode from env', () => {
		process.env.F5XC_API_URL = 'https://api.f5xc.com';
		process.env.F5XC_DEBUG = 'true';

		const client = createClientFromEnv();

		expect(client).not.toBeNull();
		// Debug mode is internal, but client should be created successfully
	});
});

describe('buildResourcePath', () => {
	it('should build path for list action with namespace', () => {
		const path = buildResourcePath('config', 'http_loadbalancers', 'list', 'shared');
		expect(path).toBe('/api/config/namespaces/shared/http_loadbalancers');
	});

	it('should build path for get action with name', () => {
		const path = buildResourcePath('config', 'http_loadbalancers', 'get', 'shared', 'my-lb');
		expect(path).toBe('/api/config/namespaces/shared/http_loadbalancers/my-lb');
	});

	it('should build path without namespace', () => {
		const path = buildResourcePath('web', 'namespaces', 'list');
		expect(path).toBe('/api/web/namespaces');
	});

	it('should build path for system namespace', () => {
		const path = buildResourcePath('config', 'sites', 'list', 'system');
		expect(path).toBe('/api/config/namespaces/system/sites');
	});

	it('should handle various resource types', () => {
		expect(buildResourcePath('config', 'origin_pools', 'list', 'shared'))
			.toBe('/api/config/namespaces/shared/origin_pools');

		expect(buildResourcePath('config', 'healthchecks', 'list', 'shared'))
			.toBe('/api/config/namespaces/shared/healthchecks');

		expect(buildResourcePath('config', 'app_firewalls', 'list', 'shared'))
			.toBe('/api/config/namespaces/shared/app_firewalls');
	});
});

describe('APIError', () => {
	it('should store all error properties', () => {
		const error = new APIError('Test error', 400, { code: 'ERR_001' }, 'GET /api/test');

		expect(error.message).toBe('Test error');
		expect(error.statusCode).toBe(400);
		expect(error.response).toEqual({ code: 'ERR_001' });
		expect(error.operation).toBe('GET /api/test');
	});

	it('should be an instance of Error', () => {
		const error = new APIError('Test', 500);
		expect(error).toBeInstanceOf(Error);
	});

	it('should have correct name', () => {
		const error = new APIError('Test', 500);
		expect(error.name).toBe('APIError');
	});

	it('should provide helpful hints for common status codes', () => {
		const error401 = new APIError('Unauthorized', 401);
		expect(error401.getHint()).toContain('Authentication');

		const error403 = new APIError('Forbidden', 403);
		expect(error403.getHint()).toContain('Permission denied');

		const error404 = new APIError('Not found', 404);
		expect(error404.getHint()).toContain('not found');

		const error500 = new APIError('Server error', 500);
		expect(error500.getHint()).toContain('Server error');
	});
});
