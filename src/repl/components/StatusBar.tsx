/**
 * StatusBar component - Displays git status and keyboard hints
 * Bottom bar with CLI name, git branch info, and help text
 */

import React from "react";
import { Box, Text, Spacer } from "ink";
import { execSync } from "child_process";
import { CLI_NAME } from "../../branding/index.js";

/**
 * Git repository status information
 */
export interface GitInfo {
	inRepo: boolean;
	repoName: string;
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
	// Left side: repo/branch or CLI name
	const renderLeft = (): React.ReactElement => {
		if (gitInfo?.inRepo) {
			const icon = getStatusIcon(gitInfo);
			const color = getStatusColor(gitInfo);

			return (
				<Text>
					<Text color="#ffffff">{gitInfo.repoName}</Text>
					<Text color="#666666">/</Text>
					<Text color={color}>{gitInfo.branch}</Text>
					<Text> </Text>
					<Text color={color}>{icon}</Text>
				</Text>
			);
		}

		return <Text color="#ffffff">{CLI_NAME}</Text>;
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
 * Get git repository info by executing git commands
 */
export function getGitInfo(): GitInfo {
	const defaultInfo: GitInfo = {
		inRepo: false,
		repoName: "",
		branch: "",
		isDirty: false,
		ahead: 0,
		behind: 0,
	};

	try {
		// Check if we're in a git repo
		try {
			execSync("git rev-parse --is-inside-work-tree", {
				encoding: "utf-8",
				stdio: ["pipe", "pipe", "pipe"],
			});
		} catch {
			return defaultInfo;
		}

		// Get repository name from top-level directory
		let repoName = "";
		try {
			const topLevel = execSync("git rev-parse --show-toplevel", {
				encoding: "utf-8",
				stdio: ["pipe", "pipe", "pipe"],
			}).trim();
			repoName = topLevel.split("/").pop() ?? "";
		} catch {
			repoName = "repo";
		}

		// Get current branch
		let branch = "";
		try {
			branch = execSync("git rev-parse --abbrev-ref HEAD", {
				encoding: "utf-8",
				stdio: ["pipe", "pipe", "pipe"],
			}).trim();
		} catch {
			branch = "unknown";
		}

		// Check for uncommitted changes
		let isDirty = false;
		try {
			const status = execSync("git status --porcelain", {
				encoding: "utf-8",
				stdio: ["pipe", "pipe", "pipe"],
			});
			isDirty = status.trim().length > 0;
		} catch {
			// Ignore errors
		}

		// Get ahead/behind counts
		let ahead = 0;
		let behind = 0;
		try {
			const counts = execSync(
				"git rev-list --left-right --count HEAD...@{upstream}",
				{
					encoding: "utf-8",
					stdio: ["pipe", "pipe", "pipe"],
				},
			).trim();
			const [aheadStr, behindStr] = counts.split(/\s+/);
			ahead = parseInt(aheadStr ?? "0", 10) || 0;
			behind = parseInt(behindStr ?? "0", 10) || 0;
		} catch {
			// No upstream or error - that's fine
		}

		return {
			inRepo: true,
			repoName,
			branch,
			isDirty,
			ahead,
			behind,
		};
	} catch {
		// If anything fails, return default
		return defaultInfo;
	}
}

export default StatusBar;
