/**
 * Unit tests for HistoryManager class
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { existsSync, unlinkSync, readFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { HistoryManager, redactSensitive } from '../../src/repl/history.js';

describe('HistoryManager', () => {
  let history: HistoryManager;

  beforeEach(() => {
    // Constructor takes (path: string, maxSize?: number)
    history = new HistoryManager(':memory:', 100);
  });

  describe('initial state', () => {
    it('should start empty', () => {
      expect(history.length).toBe(0);
      expect(history.getHistory()).toEqual([]);
    });
  });

  describe('add', () => {
    it('should add commands to history', () => {
      history.add('help');
      history.add('list');
      expect(history.length).toBe(2);
      expect(history.getHistory()).toEqual(['help', 'list']);
    });

    it('should not add empty commands', () => {
      history.add('');
      history.add('   ');
      expect(history.length).toBe(0);
    });

    it('should not add duplicate consecutive commands', () => {
      history.add('help');
      history.add('help');
      expect(history.length).toBe(1);
    });

    it('should add same command if not consecutive', () => {
      history.add('help');
      history.add('list');
      history.add('help');
      expect(history.length).toBe(3);
    });

    it('should respect maxSize limit by dropping oldest entries', () => {
      const smallHistory = new HistoryManager(':memory:', 3);
      smallHistory.add('cmd1');
      smallHistory.add('cmd2');
      smallHistory.add('cmd3');
      smallHistory.add('cmd4');
      expect(smallHistory.length).toBe(3);
      expect(smallHistory.getHistory()).toEqual(['cmd2', 'cmd3', 'cmd4']);
    });
  });

  describe('get', () => {
    beforeEach(() => {
      history.add('first');
      history.add('second');
      history.add('third');
    });

    it('should get command by index', () => {
      expect(history.get(0)).toBe('first');
      expect(history.get(1)).toBe('second');
      expect(history.get(2)).toBe('third');
    });

    it('should return undefined for out of bounds', () => {
      expect(history.get(10)).toBeUndefined();
      expect(history.get(-1)).toBe('third'); // .at() supports negative indices
    });
  });

  describe('getLast', () => {
    it('should return undefined when empty', () => {
      expect(history.getLast()).toBeUndefined();
    });

    it('should return last command', () => {
      history.add('first');
      history.add('second');
      expect(history.getLast()).toBe('second');
    });
  });

  describe('search', () => {
    beforeEach(() => {
      history.add('http_loadbalancer list');
      history.add('origin_pool create');
      history.add('http_loadbalancer get my-lb');
    });

    it('should find matching commands', () => {
      const results = history.search('http');
      expect(results).toHaveLength(2);
      expect(results).toContain('http_loadbalancer list');
      expect(results).toContain('http_loadbalancer get my-lb');
    });

    it('should be case insensitive', () => {
      const results = history.search('HTTP');
      expect(results).toHaveLength(2);
    });

    it('should return empty array when no matches', () => {
      const results = history.search('nonexistent');
      expect(results).toHaveLength(0);
    });
  });

  describe('clear', () => {
    it('should clear all history', () => {
      history.add('cmd1');
      history.add('cmd2');
      history.clear();
      expect(history.length).toBe(0);
      expect(history.getHistory()).toEqual([]);
    });
  });
});

/**
 * File Persistence Tests
 */
describe('HistoryManager - File Persistence', () => {
  const testPath = join(tmpdir(), `xcsh-test-history-${Date.now()}-${Math.random().toString(36).slice(2)}`);

  afterEach(() => {
    if (existsSync(testPath)) {
      unlinkSync(testPath);
    }
  });

  it('should save history to disk', async () => {
    const manager = new HistoryManager(testPath, 100);
    manager.add('command1');
    manager.add('command2');
    await manager.save();

    expect(existsSync(testPath)).toBe(true);
    const content = readFileSync(testPath, 'utf-8');
    expect(content).toBe('command1\ncommand2\n');
  });

  it('should load history from disk', async () => {
    // Create file first
    const manager1 = new HistoryManager(testPath, 100);
    manager1.add('saved1');
    manager1.add('saved2');
    await manager1.save();

    // Load in new instance
    const manager2 = await HistoryManager.create(testPath, 100);
    expect(manager2.length).toBe(2);
    expect(manager2.getHistory()).toEqual(['saved1', 'saved2']);
  });

  it('should round-trip correctly', async () => {
    const manager1 = new HistoryManager(testPath, 100);
    manager1.add('round1');
    manager1.add('round2');
    manager1.add('round3');
    await manager1.save();

    const manager2 = await HistoryManager.create(testPath, 100);
    manager2.add('round4');
    await manager2.save();

    const manager3 = await HistoryManager.create(testPath, 100);
    expect(manager3.getHistory()).toEqual(['round1', 'round2', 'round3', 'round4']);
  });

  it('should handle non-existent history file gracefully', async () => {
    const nonExistentPath = join(tmpdir(), `xcsh-nonexistent-${Date.now()}`);
    const manager = await HistoryManager.create(nonExistentPath, 100);
    expect(manager.length).toBe(0);
  });
});

/**
 * Sensitive Data Redaction Tests
 */
describe('Sensitive Data Redaction', () => {
  describe('redactSensitive function', () => {
    it('should redact --token values', () => {
      expect(redactSensitive('login --token abc123xyz')).toBe('login --token ******');
    });

    it('should redact --token= values', () => {
      expect(redactSensitive('login --token=abc123xyz')).toBe('login --token=******');
    });

    it('should redact --api-token values', () => {
      expect(redactSensitive('configure --api-token mysecret')).toBe('configure --api-token ******');
    });

    it('should redact --api-token= values', () => {
      expect(redactSensitive('configure --api-token=mysecret')).toBe('configure --api-token=******');
    });

    it('should redact --password values', () => {
      expect(redactSensitive('login --password supersecret')).toBe('login --password ******');
    });

    it('should redact --secret values', () => {
      expect(redactSensitive('configure --secret mykey')).toBe('configure --secret ******');
    });

    it('should redact --certificate values', () => {
      expect(redactSensitive('configure --certificate /path/to/cert.pem')).toBe('configure --certificate ******');
    });

    it('should redact --cert values', () => {
      expect(redactSensitive('configure --cert /path/to/cert.pem')).toBe('configure --cert ******');
    });

    it('should redact --private-key values', () => {
      expect(redactSensitive('configure --private-key keydata')).toBe('configure --private-key ******');
    });

    it('should redact --api-key values', () => {
      expect(redactSensitive('configure --api-key myapikey')).toBe('configure --api-key ******');
    });

    it('should redact short -t flag values', () => {
      expect(redactSensitive('login -t abc123')).toBe('login -t ******');
    });

    it('should redact short -p flag values', () => {
      expect(redactSensitive('login -p mypassword')).toBe('login -p ******');
    });

    it('should redact Authorization: Bearer headers', () => {
      expect(redactSensitive('curl -H "Authorization: Bearer xyz789"')).toBe('curl -H "Authorization: Bearer ******"');
    });

    it('should redact Authorization: APIToken headers', () => {
      expect(redactSensitive('curl -H "Authorization: APIToken abc123"')).toBe('curl -H "Authorization: APIToken ******"');
    });

    it('should redact Authorization: Basic headers', () => {
      expect(redactSensitive('curl -H "Authorization: Basic dXNlcjpwYXNz"')).toBe('curl -H "Authorization: Basic ******"');
    });

    it('should redact F5 XC APIToken', () => {
      expect(redactSensitive('APIToken abc123xyz')).toBe('APIToken ******');
    });

    it('should redact APIToken with quotes', () => {
      expect(redactSensitive('APIToken "abc123xyz"')).toBe('APIToken ******');
    });

    it('should redact environment variable assignments', () => {
      expect(redactSensitive('export API_TOKEN=abc123')).toBe('export API_TOKEN=******');
      expect(redactSensitive('export F5XC_API_TOKEN=xyz789')).toBe('export F5XC_API_TOKEN=******');
      expect(redactSensitive('export PASSWORD=secret')).toBe('export PASSWORD=******');
    });

    it('should redact multiple sensitive values', () => {
      expect(redactSensitive('cmd --token abc --password xyz')).toBe('cmd --token ****** --password ******');
    });

    it('should not redact non-sensitive flags', () => {
      expect(redactSensitive('list --namespace default --format json')).toBe('list --namespace default --format json');
    });

    it('should not redact non-sensitive content', () => {
      expect(redactSensitive('help')).toBe('help');
      expect(redactSensitive('http_loadbalancer list')).toBe('http_loadbalancer list');
    });

    it('should be case insensitive for flag names', () => {
      expect(redactSensitive('login --TOKEN abc123')).toBe('login --TOKEN ******');
      expect(redactSensitive('login --Password secret')).toBe('login --Password ******');
    });
  });

  describe('HistoryManager integration with redaction', () => {
    let history: HistoryManager;

    beforeEach(() => {
      history = new HistoryManager(':memory:', 100);
    });

    it('should store redacted commands', () => {
      history.add('login --token abc123xyz');
      expect(history.getLast()).toBe('login --token ******');
    });

    it('should redact sensitive data before storing', () => {
      history.add('curl -H "Authorization: Bearer mysecrettoken"');
      expect(history.getLast()).toBe('curl -H "Authorization: Bearer ******"');
    });

    it('should deduplicate based on redacted value', () => {
      history.add('login --token abc123');
      history.add('login --token xyz789'); // Same command with different token
      // Both should be stored as same redacted value, so only one entry
      expect(history.length).toBe(1);
      expect(history.getLast()).toBe('login --token ******');
    });

    it('should allow different commands with different redacted flags', () => {
      history.add('login --token abc123');
      history.add('configure --api-key xyz789');
      expect(history.length).toBe(2);
      expect(history.getHistory()).toEqual(['login --token ******', 'configure --api-key ******']);
    });
  });
});
