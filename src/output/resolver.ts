/**
 * Output Format Resolver
 * Resolves output format based on precedence: CLI flags > env vars > config > defaults
 */

import { ENV_PREFIX } from "../branding/index.js";
import type { OutputContext, OutputFormat } from "./types.js";
import { isValidOutputFormat } from "./types.js";

/**
 * Environment variable name for output format
 */
export const OUTPUT_FORMAT_ENV_VAR = `${ENV_PREFIX}_OUTPUT_FORMAT`;

/**
 * Resolve the effective output format based on precedence:
 * 1. CLI flags (--output, -o)
 * 2. Environment variable (F5XC_OUTPUT_FORMAT)
 * 3. Config file (~/.xcshconfig)
 * 4. Default (table)
 */
export function resolveOutputFormat(context: OutputContext): OutputFormat {
	// 1. CLI flag takes highest precedence
	if (context.cliFormat) {
		return context.cliFormat;
	}

	// 2. Environment variable
	if (context.envFormat) {
		return context.envFormat;
	}

	// 3. Config file
	if (context.configFormat) {
		return context.configFormat;
	}

	// 4. Default: table for all contexts
	return "table";
}

/**
 * Parse output format from environment variable
 */
export function getOutputFormatFromEnv(): OutputFormat | undefined {
	const envValue = process.env[OUTPUT_FORMAT_ENV_VAR]?.toLowerCase().trim();

	if (envValue && isValidOutputFormat(envValue)) {
		return envValue as OutputFormat;
	}

	return undefined;
}

/**
 * Parse --output/-o flag from command arguments
 * Returns the format and remaining arguments with the flag removed
 */
export function parseOutputFlag(args: string[]): {
	format: OutputFormat | undefined;
	remainingArgs: string[];
} {
	const remainingArgs: string[] = [];
	let format: OutputFormat | undefined;

	for (let i = 0; i < args.length; i++) {
		const arg = args[i];

		if (arg === "-o" || arg === "--output") {
			// Value is next argument
			const nextArg = args[i + 1];
			if (nextArg && !nextArg.startsWith("-")) {
				if (isValidOutputFormat(nextArg)) {
					format = nextArg.toLowerCase() as OutputFormat;
				}
				i++; // Skip the value
			}
		} else if (arg?.startsWith("--output=")) {
			// Value is part of the argument
			const value = arg.slice("--output=".length);
			if (isValidOutputFormat(value)) {
				format = value.toLowerCase() as OutputFormat;
			}
		} else if (arg?.startsWith("-o=")) {
			// Short form with value
			const value = arg.slice("-o=".length);
			if (isValidOutputFormat(value)) {
				format = value.toLowerCase() as OutputFormat;
			}
		} else if (arg) {
			remainingArgs.push(arg);
		}
	}

	return { format, remainingArgs };
}

/**
 * Parse --spec flag from command arguments
 * Returns whether spec was requested and remaining arguments
 */
export function parseSpecFlag(args: string[]): {
	spec: boolean;
	remainingArgs: string[];
} {
	const remainingArgs: string[] = [];
	let spec = false;

	for (const arg of args) {
		if (arg === "--spec") {
			spec = true;
		} else {
			remainingArgs.push(arg);
		}
	}

	return { spec, remainingArgs };
}

/**
 * Determine if colors should be used based on context
 * Auto-detects TTY and respects --no-color flag
 */
export function shouldUseColors(
	isTTY: boolean = process.stdout.isTTY ?? false,
	noColorFlag: boolean = false,
): boolean {
	// Explicit --no-color flag
	if (noColorFlag) {
		return false;
	}

	// NO_COLOR environment variable (https://no-color.org/)
	if (process.env.NO_COLOR !== undefined) {
		return false;
	}

	// FORCE_COLOR forces colors even in non-TTY
	if (process.env.FORCE_COLOR !== undefined) {
		return true;
	}

	// Auto-detect based on TTY
	return isTTY;
}

/**
 * Build output context from current environment
 */
export function buildOutputContext(options?: {
	cliFormat?: OutputFormat;
	configFormat?: OutputFormat;
	isInteractive?: boolean;
	noColor?: boolean;
}): OutputContext {
	const context: OutputContext = {
		isInteractive: options?.isInteractive ?? false,
		isTTY: process.stdout.isTTY ?? false,
	};

	// Only add optional properties if they have values (exactOptionalPropertyTypes)
	if (options?.cliFormat !== undefined) {
		context.cliFormat = options.cliFormat;
	}
	const envFormat = getOutputFormatFromEnv();
	if (envFormat !== undefined) {
		context.envFormat = envFormat;
	}
	if (options?.configFormat !== undefined) {
		context.configFormat = options.configFormat;
	}
	if (options?.noColor !== undefined) {
		context.noColor = options.noColor;
	}

	return context;
}

/**
 * Parse all output-related flags from arguments
 * Combines --output and --spec parsing
 */
export function parseOutputFlags(args: string[]): {
	format: OutputFormat | undefined;
	spec: boolean;
	remainingArgs: string[];
} {
	// First parse --output
	const { format, remainingArgs: afterOutput } = parseOutputFlag(args);

	// Then parse --spec from remaining
	const { spec, remainingArgs } = parseSpecFlag(afterOutput);

	return { format, spec, remainingArgs };
}
