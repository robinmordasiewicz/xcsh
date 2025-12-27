/**
 * StatusBar component - Displays git status and keyboard hints
 * Bottom bar with CLI name, git branch info, and help text
 */

import React from "react";
import { Box, Text, Spacer } from "ink";
import { CLI_NAME } from "../../branding/index.js";

/**
 * Git repository status information
 */
export interface GitInfo {
	inRepo: boolean;
	branch: string;
	isDirty: boolean;
	ahead: number;
	behind: number;
}

interface StatusBarProps {
	gitInfo?: GitInfo | undefined;
	width?: number;
	hint?: string;
}

/**
 * Get status icon based on git state
 */
function getStatusIcon(gitInfo: GitInfo): string {
	if (gitInfo.ahead > 0 && gitInfo.behind > 0) {
		return "\u21C5"; // up-down arrow
	}
	if (gitInfo.ahead > 0) {
		return "\u2191"; // up arrow
	}
	if (gitInfo.behind > 0) {
		return "\u2193"; // down arrow
	}
	if (gitInfo.isDirty) {
		return "*";
	}
	return "\u2713"; // checkmark
}

/**
 * Get color for status icon
 */
function getStatusColor(gitInfo: GitInfo): string {
	if (gitInfo.ahead > 0 || gitInfo.behind > 0) {
		return "#2196f3"; // Blue
	}
	if (gitInfo.isDirty) {
		return "#ffc107"; // Yellow
	}
	return "#00c853"; // Green
}

/**
 * StatusBar component
 * Displays CLI name, git status, and keyboard hints
 */
export function StatusBar({
	gitInfo,
	width = 80,
	hint = "Ctrl+C: quit",
}: StatusBarProps): React.ReactElement {
	// Left side: CLI name and git info
	const renderLeft = (): React.ReactElement => {
		const content = CLI_NAME;

		if (gitInfo?.inRepo) {
			const icon = getStatusIcon(gitInfo);
			const color = getStatusColor(gitInfo);

			return (
				<Text>
					<Text color="#ffffff">{CLI_NAME} (</Text>
					<Text color="#ffffff">{gitInfo.branch}</Text>
					<Text> </Text>
					<Text color={color}>{icon}</Text>
					<Text color="#ffffff">)</Text>
				</Text>
			);
		}

		return <Text color="#ffffff">{content}</Text>;
	};

	// Right side: keyboard hints
	const renderRight = (): React.ReactElement => {
		return <Text color="#666666">{hint}</Text>;
	};

	return (
		<Box width={width} paddingX={1} justifyContent="space-between">
			{renderLeft()}
			<Spacer />
			{renderRight()}
		</Box>
	);
}

/**
 * Get git repository info (stub - actual implementation will call git)
 */
export function getGitInfo(): GitInfo {
	// TODO: Implement actual git status detection
	// For now, return a default state
	return {
		inRepo: false,
		branch: "",
		isDirty: false,
		ahead: 0,
		behind: 0,
	};
}

export default StatusBar;
