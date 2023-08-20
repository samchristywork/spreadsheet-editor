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

func getCellColor(row int, column int) *[]int {
	if cellColorMap[row] == nil {
		return nil
	}

	color := cellColorMap[row][column]
	return &color
}

func editCell(t *term.Term) *string {
	width, _ := screenDimensions()

	entry, _ := getCellContent(currentCell[0], currentCell[1])

	x := len(entry) + 1

	for {
		setCursorPosition(1, 2)
		fmt.Printf("%s", fixedWidth(entry, width))

		setCursorPosition(x, 2)
		makeCursorVisible()

		if !nextKeyPress() {
			break
		}

		makeCursorInvisible()

		if keyPressed(27, 0, 0, bytes) { // Escape
			entry = ""
			break
		} else if keyPressed(27, 91, 68, bytes) { // Left
			if x > 1 {
				x--
			}
		} else if keyPressed(27, 91, 67, bytes) { // Right
			if x < len(entry)+1 {
				x++
			}
		} else if keyPressed(21, 0, 0, bytes) { // Ctrl-u
			entry = ""
			x = 1
		} else if keyPressed(127, 0, 0, bytes) { // Backspace
			if x > 1 {
				entry = entry[:x-2] + entry[x-1:]
				x--
			}
		} else if keyPressed(13, 0, 0, bytes) { // Enter
			break
		} else if isPrintable(bytes) {
			entry = entry[:x-1] + string(bytes[0]) + entry[x-1:]
			x++
		}
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
		return "", fmt.Errorf("Error getting cell: Cell %s%s is empty", getColumnName(column), row+1)
	}

	content := contentMap[row][column]
	return content, nil
}

func setCellContent(row int, column int, content string) {
	if contentMap[row] == nil {
		contentMap[row] = make(map[int]string)
	}

	contentMap[row][column] = content
}
