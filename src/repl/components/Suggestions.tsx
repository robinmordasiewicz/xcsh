/**
 * Suggestions component - Completion popup for tab completion
 * Displays filtered suggestions with keyboard navigation
 */

import React from "react";
import { Box, Text, useInput } from "ink";

/**
 * Suggestion item with display text and completion value
 */
export interface Suggestion {
	label: string; // Display text
	value: string; // Completion value
	description?: string; // Optional description
	category?: string; // e.g., "domain", "action", "flag"
}

interface SuggestionsProps {
	suggestions: Suggestion[];
	selectedIndex: number;
	onSelect: (suggestion: Suggestion) => void;
	onNavigate: (direction: "up" | "down") => void;
	onCancel: () => void;
	maxVisible?: number;
	isActive?: boolean;
}

/**
 * Get category color for visual distinction
 */
function getCategoryColor(category?: string): string {
	switch (category) {
		case "domain":
			return "#2196f3"; // Blue
		case "action":
			return "#4caf50"; // Green
		case "flag":
			return "#ffc107"; // Yellow
		case "value":
			return "#9c27b0"; // Purple
		default:
			return "#ffffff"; // White
	}
}

/**
 * Single suggestion item
 */
function SuggestionItem({
	suggestion,
	isSelected,
	maxLabelWidth,
}: {
	suggestion: Suggestion;
	isSelected: boolean;
	index: number;
	maxLabelWidth: number;
}): React.ReactElement {
	const categoryColor = getCategoryColor(suggestion.category);

	return (
		<Box>
			{/* Selection indicator - use 3 chars for alignment (▶ may render as wide char) */}
			<Text color={isSelected ? "#CA260A" : "#333333"}>
				{isSelected ? "▶ " : "   "}
			</Text>

			{/* Label with category color - padded to align descriptions */}
			<Text color={categoryColor} bold={isSelected} inverse={isSelected}>
				{suggestion.label.padEnd(maxLabelWidth)}
			</Text>

			{/* Description if present */}
			{suggestion.description && (
				<Text color="#666666"> - {suggestion.description}</Text>
			)}
		</Box>
	);
}

/**
 * Suggestions popup component
 * Displays a scrollable list of suggestions with keyboard navigation
 */
export function Suggestions({
	suggestions,
	selectedIndex,
	onSelect,
	onNavigate,
	onCancel,
	maxVisible = 20,
	isActive = true,
}: SuggestionsProps): React.ReactElement | null {
	// Handle keyboard navigation
	useInput(
		(_input, key) => {
			if (!isActive) return;

			if (key.upArrow) {
				onNavigate("up");
				return;
			}

			if (key.downArrow) {
				onNavigate("down");
				return;
			}

			if (key.return || key.tab) {
				const selected = suggestions.at(selectedIndex);
				if (selected) {
					onSelect(selected);
				}
				return;
			}

			if (key.escape) {
				onCancel();
				return;
			}
		},
		{ isActive },
	);

	// Don't render if no suggestions
	if (suggestions.length === 0) {
		return null;
	}

	// Calculate visible window
	const totalCount = suggestions.length;
	const startIndex = Math.max(
		0,
		Math.min(
			selectedIndex - Math.floor(maxVisible / 2),
			totalCount - maxVisible,
		),
	);
	const visibleSuggestions = suggestions.slice(
		startIndex,
		startIndex + maxVisible,
	);

	// Show scroll indicators
	const showScrollUp = startIndex > 0;
	const showScrollDown = startIndex + maxVisible < totalCount;

	// Calculate max label width for column alignment
	const maxLabelWidth = Math.max(
		...visibleSuggestions.map((s) => s.label.length),
	);

	return (
		<Box
			flexDirection="column"
			borderStyle="round"
			borderColor="#CA260A"
			paddingX={1}
		>
			{/* Scroll up indicator */}
			{showScrollUp && (
				<Text color="#666666" dimColor>
					{"\u25B2"} ({startIndex} more above)
				</Text>
			)}

			{/* Suggestion items */}
			{visibleSuggestions.map((suggestion, index) => (
				<SuggestionItem
					key={suggestion.value}
					suggestion={suggestion}
					isSelected={startIndex + index === selectedIndex}
					index={startIndex + index}
					maxLabelWidth={maxLabelWidth}
				/>
			))}

			{/* Scroll down indicator */}
			{showScrollDown && (
				<Text color="#666666" dimColor>
					{"\u25BC"} ({totalCount - startIndex - maxVisible} more
					below)
				</Text>
			)}

			{/* Help text */}
			<Box marginTop={1}>
				<Text color="#666666" dimColor>
					Tab: select | Up/Down: navigate | Esc: cancel
				</Text>
			</Box>
		</Box>
	);
}

/**
 * Hook for managing suggestion state
 */
export function useSuggestions(allSuggestions: Suggestion[]) {
	const [selectedIndex, setSelectedIndex] = React.useState(0);
	const [filteredSuggestions, setFilteredSuggestions] =
		React.useState(allSuggestions);

	// Reset selection when suggestions change
	React.useEffect(() => {
		setSelectedIndex(0);
		setFilteredSuggestions(allSuggestions);
	}, [allSuggestions]);

	const navigate = React.useCallback(
		(direction: "up" | "down") => {
			setSelectedIndex((current) => {
				if (direction === "up") {
					return current > 0
						? current - 1
						: filteredSuggestions.length - 1;
				} else {
					return current < filteredSuggestions.length - 1
						? current + 1
						: 0;
				}
			});
		},
		[filteredSuggestions.length],
	);

	const filter = React.useCallback(
		(query: string) => {
			const lowerQuery = query.toLowerCase();
			const filtered = allSuggestions.filter(
				(s) =>
					s.label.toLowerCase().includes(lowerQuery) ||
					s.value.toLowerCase().includes(lowerQuery),
			);
			setFilteredSuggestions(filtered);
			setSelectedIndex(0);
		},
		[allSuggestions],
	);

	const reset = React.useCallback(() => {
		setSelectedIndex(0);
		setFilteredSuggestions(allSuggestions);
	}, [allSuggestions]);

	return {
		suggestions: filteredSuggestions,
		selectedIndex,
		selectedSuggestion: filteredSuggestions.at(selectedIndex),
		navigate,
		filter,
		reset,
		setSelectedIndex,
	};
}

export default Suggestions;
