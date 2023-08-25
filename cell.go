package main

import (
	"fmt"
	"github.com/pkg/term"
)

func setCellColor(row int, column int) {
	if cellColorMap[row] == nil {
		cellColorMap[row] = make(map[int][]int)
	}

	cellColorMap[row][column] = []int{255, 0, 0}
}

func getCellColor(row int, column int) ([]int, error) {
	if cellColorMap[row] == nil {
		return nil, fmt.Errorf("No color for cell %d, %d", row, column)
	}

	color := cellColorMap[row][column]
	return color, nil
}

func waitForKeypress() ([]byte, error) {
	makeCursorVisible()
	bytes, err := nextKeyPress()
	makeCursorInvisible()

	return bytes, err
}

func editCell(t *term.Term) *string {
	width, _ := screenDimensions()
	entry, _ := getCellContent(currentCell[0], currentCell[1])

	cursorOffset := len(entry)
	for {
		setCursorPosition(1, 2)
		fmt.Printf("%s", fixedWidth(entry, width))

		setCursorPosition(cursorOffset+1, 2)
		bytes, err := waitForKeypress()
		if err != nil {
			break
		}

		if keyPressed(27, 0, 0, bytes) { // Escape
			return nil
		} else if keyPressed(27, 91, 68, bytes) { // Left
			cursorOffset--
		} else if keyPressed(27, 91, 67, bytes) { // Right
			cursorOffset++
		} else if keyPressed(21, 0, 0, bytes) { // Ctrl-u
			entry = ""
		} else if keyPressed(127, 0, 0, bytes) { // Backspace
			if cursorOffset > 0 {
				entry = entry[:cursorOffset-1] + entry[cursorOffset:]
				cursorOffset--
			}
		} else if keyPressed(13, 0, 0, bytes) { // Enter
			break
		} else if isPrintable(bytes) {
			entry = entry[:cursorOffset] + string(bytes[0]) + entry[cursorOffset:]
			cursorOffset++
		}

		cursorOffset = max(cursorOffset, 0)
		cursorOffset = min(cursorOffset, len(entry)+0)
	}

	setCursorPosition(1, 2)
	fmt.Printf("%s", fixedWidth("", width))

	return &entry
}

func copyCell() {
	clipboard, _ = getCellContent(currentCell[0], currentCell[1])
}

func pasteCell() {
	setCellContent(currentCell[0], currentCell[1], clipboard)
}

func getCellContent(row int, column int) (string, error) {
	if contentMap[row] == nil {
		return "", fmt.Errorf("Error getting cell: Cell %s%d is empty", getColumnName(column), row+1)
	}

	content := contentMap[row][column]
	return content, nil
}

func setCellContent(row int, column int, content string) {
	if contentMap[row] == nil {
		contentMap[row] = make(map[int]string)
	}

	contentMap[row][column] = content
	modified = true
}

func getCellValue(row int, column int) (string, error) {
	content, err := getCellContent(row, column)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", nil
	}

	if content[0] == '=' {
		val := eval(content[1:])
		return val, nil
	} else {
		return content, nil
	}
}
