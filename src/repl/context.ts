/**
 * ContextPath represents the current navigation context in the REPL.
 * Supports hierarchical navigation through domains and actions.
 */

import {
	isValidDomain as checkDomain,
	isValidAction as checkAction,
	resolveDomain,
	aliasRegistry,
	domainRegistry,
} from "../types/domains.js";

/**
 * ContextPath represents the current navigation context in the REPL
 */
export class ContextPath {
	domain: string = ""; // e.g., "load_balancer" - empty at root
	action: string = ""; // e.g., "list" - empty if only domain context

	/**
	 * Returns true if at root context (no domain selected)
	 */
	isRoot(): boolean {
		return this.domain === "";
	}

	/**
	 * Returns true if in a domain context (but no action)
	 */
	isDomain(): boolean {
		return this.domain !== "" && this.action === "";
	}

	/**
	 * Returns true if in an action context
	 */
	isAction(): boolean {
		return this.domain !== "" && this.action !== "";
	}

	/**
	 * Returns the path as "domain/action" or "domain" or ""
	 */
	toString(): string {
		if (this.domain === "") {
			return "";
		}
		if (this.action === "") {
			return this.domain;
		}
		return `${this.domain}/${this.action}`;
	}

	/**
	 * Reset clears the context to root
	 */
	reset(): void {
		this.domain = "";
		this.action = "";
	}

	/**
	 * Navigate up one level in the context hierarchy
	 * Returns true if navigation occurred, false if already at root
	 */
	navigateUp(): boolean {
		if (this.action !== "") {
			this.action = "";
			return true;
		}
		if (this.domain !== "") {
			this.domain = "";
			return true;
		}
		return false; // Already at root
	}

	/**
	 * Enter a domain context
	 */
	setDomain(domain: string): void {
		this.domain = domain;
		this.action = "";
	}

	/**
	 * Enter an action context within current domain
	 */
	setAction(action: string): void {
		this.action = action;
	}

	/**
	 * Clone this context path
	 */
	clone(): ContextPath {
		const copy = new ContextPath();
		copy.domain = this.domain;
		copy.action = this.action;
		return copy;
	}
}

/**
 * ContextValidator provides validation for navigation commands
 */
export class ContextValidator {
	private domains: Set<string>;

	constructor() {
		this.domains = new Set<string>();
		this.refresh();
	}

	/**
	 * Refresh the cached domain list from registries
	 */
	refresh(): void {
		this.domains.clear();

		// Add all canonical domain names
		for (const domain of domainRegistry.keys()) {
			this.domains.add(domain);
		}

		// Also include aliases
		for (const alias of aliasRegistry.keys()) {
			this.domains.add(alias);
		}
	}

	/**
	 * Check if input is a valid domain name or alias
	 */
	isValidDomain(name: string): boolean {
		return this.domains.has(name) || checkDomain(name);
	}

	/**
	 * Check if input is a valid action command
	 */
	isValidAction(name: string): boolean {
		return checkAction(name);
	}

	/**
	 * Resolve an alias to its canonical domain name
	 */
	resolveDomain(name: string): string | undefined {
		return resolveDomain(name);
	}
}
