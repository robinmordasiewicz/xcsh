/**
 * Unit tests for executor module
 */

import { describe, it, expect } from 'vitest';
import { parseCommand } from '../../src/repl/executor.js';

describe('parseCommand', () => {
  describe('empty input', () => {
    it('should handle empty string', () => {
      const result = parseCommand('');
      expect(result.raw).toBe('');
      expect(result.isBuiltin).toBe(false);
      expect(result.isDirectNavigation).toBe(false);
    });

    it('should handle whitespace only', () => {
      const result = parseCommand('   ');
      expect(result.raw).toBe('');
    });
  });

  describe('built-in commands', () => {
    it('should recognize help', () => {
      const result = parseCommand('help');
      expect(result.isBuiltin).toBe(true);
      expect(result.raw).toBe('help');
    });

    it('should recognize clear', () => {
      const result = parseCommand('clear');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize quit', () => {
      const result = parseCommand('quit');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize exit', () => {
      const result = parseCommand('exit');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize back', () => {
      const result = parseCommand('back');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize ..', () => {
      const result = parseCommand('..');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize /', () => {
      const result = parseCommand('/');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize context', () => {
      const result = parseCommand('context');
      expect(result.isBuiltin).toBe(true);
    });

    it('should recognize history', () => {
      const result = parseCommand('history');
      expect(result.isBuiltin).toBe(true);
    });
  });

  describe('direct navigation', () => {
    it('should not recognize /unknown as navigation', () => {
      const result = parseCommand('/unknown_domain');
      // Not a valid domain, so not direct navigation
      expect(result.isDirectNavigation).toBe(false);
    });

    it('should handle / alone as builtin', () => {
      const result = parseCommand('/');
      expect(result.isBuiltin).toBe(true);
      expect(result.isDirectNavigation).toBe(false);
    });
  });

  describe('regular commands', () => {
    it('should parse regular command with args', () => {
      const result = parseCommand('list --namespace default');
      expect(result.isBuiltin).toBe(false);
      expect(result.isDirectNavigation).toBe(false);
      expect(result.args).toEqual(['list', '--namespace', 'default']);
    });

    it('should handle single word command', () => {
      const result = parseCommand('something');
      expect(result.isBuiltin).toBe(false);
      expect(result.args).toEqual(['something']);
    });
  });
});
