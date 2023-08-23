package main

import (
	"fmt"
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

func renderStatusLine(bytes []byte) {
	width, height := screenDimensions()

	content := fmt.Sprintf("Keypress: %d %d %d, Frame: %d, Position: %d %d, Modified: %v", bytes[0], bytes[1], bytes[2], frame, currentCell[0], currentCell[1], modified)

	setCursorPosition(1, height)
	color(100, 100, 200)
	invert()
	fmt.Printf("%s", content)
	for i := len(content); i < width; i++ {
		fmt.Printf(" ")
	}
	resetColor()
}

func applyCellColors(row int, column int, evaluated bool) {
	c, err := getCellColor(row, column)

	if showGrid {
		backgroundColor(0, 0, 0)
		if row%2+column%2 == 1 {
			backgroundColor(20, 20, 20)
		}
	}

	// Apply color
	if err == nil && len(c) == 3 {
		color(c[0], c[1], c[2])
	}

	// Show currently selected cell
	if row == currentCell[0] && column == currentCell[1] {
		invert()
	}

	// Highlight cells that have been evaluated
	if evaluated && !(row == currentCell[0] && column == currentCell[1]) {
		color(100, 200, 100)
	}

	// Highlight cells that are out of bounds
	if row < 0 || column < 0 {
		color(100, 100, 100)
	}
}

func renderCell(row int, column int, width int) {
	row = row + scrollOffset[1]
	column = column + scrollOffset[0]

	content, _ := getCellContent(row, column)
	value, _ := getCellValue(row, column)

	applyCellColors(row, column, content != value)

	// Render cell
	if row < 0 || column < 0 {
		fmt.Printf("%s", fixedWidth("-----", width))
	} else if row == currentCell[0] && column == currentCell[1] {
		fmt.Printf("%s", fixedWidth(content, width))
	} else {
		fmt.Printf("%s", fixedWidth(value, width))
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

func renderTitle() {
	width, _ := screenDimensions()
	setCursorPosition(1, 1)

	content, _ := getCellContent(currentCell[0], currentCell[1])
	value, _ := getCellValue(currentCell[0], currentCell[1])

	s := fmt.Sprintf("Cell: %s%d, Content: '%s' Display: '%s'", getColumnName(currentCell[1]), currentCell[0], content, value)
	fmt.Printf("%s", fixedWidth(s, width))
}

func render(bytes []byte) {
	renderTitle()
	renderHeadings(1, 3)
	renderStatusLine(bytes)
	renderGrid()
}
