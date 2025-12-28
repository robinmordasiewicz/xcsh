/**
 * Banner component - F5 logo with connection info
 * Fixed 80-character width layout with logo, help text, and connection info
 */

import React from "react";
import { Box, Text } from "ink";
import stringWidth from "string-width";
import { F5_LOGO, CLI_FULL_NAME } from "../../branding/index.js";
import { toDisplayTier } from "../../domains/login/whoami/types.js";

interface ConnectionInfo {
	tenant?: string | undefined;
	username?: string | undefined;
	tier?: string | undefined;
	namespace: string;
	isConnected: boolean;
}

interface BannerProps {
	version: string;
	connectionInfo: ConnectionInfo;
}

// Box drawing characters
const BOX = {
	topLeft: "\u256D",
	topRight: "\u256E",
	bottomLeft: "\u2570",
	bottomRight: "\u256F",
	horizontal: "\u2500",
	vertical: "\u2502",
	leftT: "\u251C",
	rightT: "\u2524",
};

// Fixed total width of 80 characters (standard terminal width)
const TOTAL_WIDTH = 80;
const INNER_WIDTH = TOTAL_WIDTH - 2; // Subtract left and right borders

// Calculate logo dimensions
const logoLines = F5_LOGO.split("\n");
const LOGO_WIDTH = Math.max(...logoLines.map((l) => stringWidth(l)));

// Help text - positioned on specific logo rows (vertically centered)
const HELP_LINES = [
	"Type 'help' for commands",
	"Run 'namespace <ns>' to set",
	"Press Ctrl+C twice to exit",
];

// Rows where help text appears (0-indexed, centered in logo)
const HELP_START_ROW = 8; // Start at row 8 of the logo

/**
 * Pad a string to a specific display width
 */
function padToWidth(str: string, targetWidth: number): string {
	const currentWidth = stringWidth(str);
	if (currentWidth >= targetWidth) return str;
	return str + " ".repeat(targetWidth - currentWidth);
}

/**
 * Colorize a logo line - applies red to circle chars, white to F5 block chars
 */
function colorizeLogoLine(line: string): React.ReactElement[] {
	const elements: React.ReactElement[] = [];
	let currentColor: string | null = null;
	let buffer = "";
	let keyIndex = 0;

	const flush = () => {
		if (buffer) {
			if (currentColor) {
				elements.push(
					<Text key={keyIndex++} color={currentColor}>
						{buffer}
					</Text>,
				);
			} else {
				elements.push(<Text key={keyIndex++}>{buffer}</Text>);
			}
			buffer = "";
		}
	};

	for (const char of line) {
		let newColor: string | null;

		switch (char) {
			case "\u2593": // Dark shade - red
			case "\u2592": // Medium shade - red
			case "(":
			case ")":
			case "|":
			case "_":
				newColor = "#E4002B"; // F5 Red
				break;
			case "\u2588": // Full block - white F5 text
				newColor = "#ffffff";
				break;
			default:
				newColor = null;
		}

		if (newColor !== currentColor) {
			flush();
			currentColor = newColor;
		}

		buffer += char === "\u2593" ? "\u2588" : char;
	}

	flush();
	return elements;
}

/**
 * Build a single content row with logo and optional help text
 */
function buildContentRow(
	logoLine: string,
	helpText: string,
	innerWidth: number,
): React.ReactElement {
	// Pad logo to its fixed width
	const paddedLogo = padToWidth(logoLine, LOGO_WIDTH);

	// Calculate remaining space for help column
	// Format: "│ {logo} {help} │"
	// Inner = logo + 1 space + help
	const helpColumnWidth = innerWidth - LOGO_WIDTH - 1;
	const paddedHelp = padToWidth(helpText, helpColumnWidth);

	return (
		<Text>
			<Text color="#E4002B">{BOX.vertical}</Text>
			{colorizeLogoLine(paddedLogo)}
			<Text> </Text>
			<Text bold color="#ffffff">
				{paddedHelp}
			</Text>
			<Text color="#E4002B">{BOX.vertical}</Text>
		</Text>
	);
}

/**
 * Main Banner component
 */
export function Banner({
	version,
	connectionInfo,
}: BannerProps): React.ReactElement {
	const title = ` ${CLI_FULL_NAME} v${version} `;
	const titleWidth = stringWidth(title);

	// Calculate top border: "╭───{title}───...───╮" = 80 chars total
	// Total = ╭(1) + leftDashes(3) + title + rightDashes + ╮(1) = 80
	const leftDashes = 3;
	const rightDashes = TOTAL_WIDTH - 1 - leftDashes - titleWidth - 1;

	// Build connection info lines
	const connectionLines: string[] = [];
	if (!connectionInfo.isConnected) {
		connectionLines.push("Not connected - run: export F5XC_API_URL=...");
	} else {
		if (connectionInfo.tenant) {
			connectionLines.push(`Tenant: ${connectionInfo.tenant}`);
		}
		if (connectionInfo.username) {
			connectionLines.push(`User: ${connectionInfo.username}`);
		}
		if (connectionInfo.tier) {
			const displayTier = toDisplayTier(connectionInfo.tier);
			if (displayTier) {
				connectionLines.push(`Tier: ${displayTier}`);
			}
		}
	}
	connectionLines.push(`Namespace: ${connectionInfo.namespace}`);

	return (
		<Box flexDirection="column" marginY={1}>
			{/* Top border with title */}
			<Text>
				<Text color="#E4002B">
					{BOX.topLeft}
					{BOX.horizontal.repeat(leftDashes)}
				</Text>
				<Text bold color="#ffffff">
					{title}
				</Text>
				<Text color="#E4002B">
					{BOX.horizontal.repeat(Math.max(0, rightDashes))}
					{BOX.topRight}
				</Text>
			</Text>

			{/* Main content: logo with help text on specific rows */}
			{logoLines.map((logoLine, index) => {
				// Get help text for this row (if any)
				const helpIndex = index - HELP_START_ROW;
				const helpText =
					helpIndex >= 0 && helpIndex < HELP_LINES.length
						? HELP_LINES[helpIndex]
						: "";

				return (
					<Box key={index}>
						{buildContentRow(logoLine, helpText ?? "", INNER_WIDTH)}
					</Box>
				);
			})}

			{/* Separator */}
			<Text color="#E4002B">
				{BOX.leftT}
				{BOX.horizontal.repeat(INNER_WIDTH)}
				{BOX.rightT}
			</Text>

			{/* Connection info */}
			{connectionLines.map((line, index) => {
				const paddedLine = padToWidth(`  ${line}`, INNER_WIDTH);
				return (
					<Text key={index}>
						<Text color="#E4002B">{BOX.vertical}</Text>
						<Text bold color="#ffffff">
							{paddedLine}
						</Text>
						<Text color="#E4002B">{BOX.vertical}</Text>
					</Text>
				);
			})}

			{/* Bottom border */}
			<Text color="#E4002B">
				{BOX.bottomLeft}
				{BOX.horizontal.repeat(INNER_WIDTH)}
				{BOX.bottomRight}
			</Text>
		</Box>
	);
}

export default Banner;
