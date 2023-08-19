package main

import (
	"fmt"
	"strings"
)

func renderHeadings(x int, y int) {
	width, height := screenDimensions()

	// Render top headings
	xoff := 3
	for i := 0; xoff < width; i++ {
		setCursorPosition(x+xoff, y)
		columnWidth := getColumnWidth(i)

		i := i + scrollOffset[0]

		color(100, 100, 200)
		invert()
		fmt.Printf("%s", fixedWidth(getColumnName(i), columnWidth))
		resetColor()

		xoff += columnWidth
	}

	// Render left headings
	for i := 0; i < height; i++ {
		setCursorPosition(x, y+i+1)
		i := i + scrollOffset[1]

		color(100, 100, 200)
		invert()
		fmt.Printf("%s", fixedWidth(fmt.Sprintf("%d", i), 3))
		resetColor()
	}
}

func renderStatusLine() {
	width, height := screenDimensions()

	content := fmt.Sprintf("Keypress: %d %d %d, Frame: %d, Position: %d %d", bytes[0], bytes[1], bytes[2], frame, currentCell[0], currentCell[1])

	setCursorPosition(1, height)
	color(100, 100, 200)
	invert()
	fmt.Printf("%s", content)
	for i := len(content); i < width; i++ {
		fmt.Printf(" ")
	}
	resetColor()
}

func renderCell(row int, column int, width int) {
	row = row + scrollOffset[1]
	column = column + scrollOffset[0]
	content := ""
	contentReference := getCellContent(row, column)
	if contentReference != nil {
		content = *contentReference
	}

	c := getCellColor(row, column)
	if c != nil {
		if len(*c) == 3 {
			color((*c)[0], (*c)[1], (*c)[2])
		}
	}

	if row == currentCell[0] && column == currentCell[1] {
		if !showGrid {
			invert()
		} else {
			if (currentCell[0]+currentCell[1])%2 == 0 {
			} else {
				invert()
			}
		}
		fmt.Printf("%s", fixedWidth(content, width))
	} else if row < 0 || column < 0 {
		color(100, 100, 100)
		fmt.Printf("%s", fixedWidth("-----", width))
	} else {
		if showGrid {
			if (row+column)%2 == 0 {
				invert()
			}
		}
		fmt.Printf("%s", fixedWidth(content, width))
	}

	resetColor()
}

func renderRow(row int, width int) {
	xoff := 4
	for column := 0; xoff < width; column++ {
		setCursorPosition(xoff, row+4)

		columnWidth := getColumnWidth(column)

		if xoff+columnWidth > width {
			columnWidth := width - xoff
			renderCell(row, column, columnWidth)
		} else {
			renderCell(row, column, columnWidth)
		}

		xoff += columnWidth
	}
}

func renderGrid() {
	width, height := screenDimensions()

	for i := 0; i < height-4; i++ {
		renderRow(i, width)
	}
}

func render() {
	setCursorPosition(1, 1)
	fmt.Printf("Hello, World!\n")
	frame++

	renderHeadings(1, 3)

	renderStatusLine()

	renderGrid()
}

func drawRectangle(x int, y int, width int, height int) {
	setCursorPosition(x, y)
	fmt.Printf("┌")
	for i := 0; i < width-2; i++ {
		fmt.Printf("─")
	}
	fmt.Printf("┐")

	for i := 0; i < height-2; i++ {
		setCursorPosition(x, y+i+1)
		fmt.Printf("│")
		setCursorPosition(x+width-1, y+i+1)
		fmt.Printf("│")
	}

	setCursorPosition(x, y+height-1)
	fmt.Printf("└")
	for i := 0; i < width-2; i++ {
		fmt.Printf("─")
	}
	fmt.Printf("┘")
}

func clearRectangle(x int, y int, width int, height int) {
	for i := 0; i < height; i++ {
		setCursorPosition(x, y+i)
		for j := 0; j < width; j++ {
			fmt.Printf(" ")
		}
	}
}

func messageBox(title string, message string) {
	width, height := screenDimensions()

	messages := strings.Split(message, "\n")

	stringWidth := len(title)
	for i := 0; i < len(messages); i++ {
		if len(messages[i]) > stringWidth {
			stringWidth = len(messages[i])
		}
	}

	boxWidth := stringWidth
	boxHeight := len(messages) + 2

	x := width/2 - boxWidth/2
	y := height/2 - boxHeight/2

	clearRectangle(x-2, y-2, boxWidth+4, boxHeight+2)
	drawRectangle(x-2, y-2, boxWidth+4, boxHeight+2)

	setCursorPosition(x, y-1)
	invert()
	fmt.Printf("%s", fixedWidth(title, boxWidth))
	resetColor()

	for i := 0; i < len(messages); i++ {
		setCursorPosition(x, y+1+i)
		fmt.Printf("%s", fixedWidth(messages[i], boxWidth))
	}

	nextKeyPress()
}
