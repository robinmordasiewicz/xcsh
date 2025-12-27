/**
 * Unit tests for ContextPath class
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { ContextPath } from '../../src/repl/context.js';

describe('ContextPath', () => {
  let ctx: ContextPath;

  beforeEach(() => {
    ctx = new ContextPath();
  });

  describe('initial state', () => {
    it('should start at root', () => {
      expect(ctx.isRoot()).toBe(true);
      expect(ctx.isDomain()).toBe(false);
      expect(ctx.isAction()).toBe(false);
    });

    it('should have empty domain and action', () => {
      expect(ctx.domain).toBe('');
      expect(ctx.action).toBe('');
    });
  });

  describe('setDomain', () => {
    it('should set domain and leave root state', () => {
      ctx.setDomain('http_loadbalancer');
      expect(ctx.isRoot()).toBe(false);
      expect(ctx.isDomain()).toBe(true);
      expect(ctx.domain).toBe('http_loadbalancer');
    });

    it('should clear action when setting new domain', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.setAction('list');
      ctx.setDomain('origin_pool');
      expect(ctx.domain).toBe('origin_pool');
      expect(ctx.action).toBe('');
    });
  });

  describe('setAction', () => {
    it('should set action when in domain context', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.setAction('list');
      expect(ctx.isAction()).toBe(true);
      expect(ctx.action).toBe('list');
    });

    it('should set action even at root (implementation allows it)', () => {
      // Note: Implementation allows setting action at root, though isAction() returns false
      // because isAction checks both domain and action are set
      ctx.setAction('list');
      expect(ctx.action).toBe('list');
      expect(ctx.isAction()).toBe(false); // domain is empty, so not "in action"
    });
  });

  describe('navigateUp', () => {
    it('should do nothing at root', () => {
      ctx.navigateUp();
      expect(ctx.isRoot()).toBe(true);
    });

    it('should clear action when in action context', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.setAction('list');
      ctx.navigateUp();
      expect(ctx.isDomain()).toBe(true);
      expect(ctx.isAction()).toBe(false);
      expect(ctx.domain).toBe('http_loadbalancer');
    });

    it('should return to root when in domain context', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.navigateUp();
      expect(ctx.isRoot()).toBe(true);
    });
  });

  describe('reset', () => {
    it('should return to root from any state', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.setAction('list');
      ctx.reset();
      expect(ctx.isRoot()).toBe(true);
      expect(ctx.domain).toBe('');
      expect(ctx.action).toBe('');
    });
  });

  describe('toString', () => {
    it('should return empty string at root', () => {
      expect(ctx.toString()).toBe('');
    });

    it('should return domain in domain context', () => {
      ctx.setDomain('http_loadbalancer');
      expect(ctx.toString()).toBe('http_loadbalancer');
    });

    it('should return domain/action in action context', () => {
      ctx.setDomain('http_loadbalancer');
      ctx.setAction('list');
      expect(ctx.toString()).toBe('http_loadbalancer/list');
    });
  });
});
