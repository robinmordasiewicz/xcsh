/**
 * Unit tests for HistoryManager class
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { HistoryManager } from '../../src/repl/history.js';

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
