/**
 * InputBox component - Command input with F5-branded horizontal rules
 * Uses ink-text-input for proper keyboard handling including backspace
 */

import React, { useState, useCallback } from "react";
import { Box, Text } from "ink";
import TextInput from "ink-text-input";

interface InputBoxProps {
	prompt: string;
	value: string;
	onChange: (value: string) => void;
	onSubmit: (value: string) => void;
	width?: number;
	isActive?: boolean;
	inputKey?: number; // Change to force cursor reset to end
}

/**
 * Render a horizontal rule in F5 red
 */
function HorizontalRule({ width }: { width: number }): React.ReactElement {
	const rule = "\u2500".repeat(Math.max(width, 1));
	return <Text color="#E4002B">{rule}</Text>;
}

/**
 * InputBox component with F5 branding
 * Displays input between two horizontal rules with a customizable prompt
 */
export function InputBox({
	prompt,
	value,
	onChange,
	onSubmit,
	width = 80,
	isActive = true,
	inputKey = 0,
}: InputBoxProps): React.ReactElement {
	return (
		<Box flexDirection="column" width={width}>
			{/* Top horizontal rule */}
			<HorizontalRule width={width} />

			{/* Prompt and input line */}
			<Box>
				<Text bold color="#ffffff">
					{prompt}
				</Text>
				<TextInput
					key={inputKey}
					value={value}
					onChange={onChange}
					onSubmit={onSubmit}
					focus={isActive}
				/>
			</Box>

			{/* Bottom horizontal rule */}
			<HorizontalRule width={width} />
		</Box>
	);
}

/**
 * Hook for managing input state
 */
export function useInputState(initialValue: string = "") {
	const [value, setValue] = useState(initialValue);
	const [cursorPosition, setCursorPosition] = useState(initialValue.length);

	const handleChange = useCallback((newValue: string) => {
		setValue(newValue);
		setCursorPosition(newValue.length);
	}, []);

	const reset = useCallback(() => {
		setValue("");
		setCursorPosition(0);
	}, []);

	const moveCursorLeft = useCallback(() => {
		setCursorPosition((pos) => Math.max(0, pos - 1));
	}, []);

	const moveCursorRight = useCallback(() => {
		setCursorPosition((pos) => Math.min(value.length, pos + 1));
	}, [value.length]);

	const moveCursorHome = useCallback(() => {
		setCursorPosition(0);
	}, []);

	const moveCursorEnd = useCallback(() => {
		setCursorPosition(value.length);
	}, [value.length]);

	return {
		value,
		cursorPosition,
		onChange: handleChange,
		reset,
		setValue,
		setCursorPosition,
		moveCursorLeft,
		moveCursorRight,
		moveCursorHome,
		moveCursorEnd,
	};
}

export default InputBox;
