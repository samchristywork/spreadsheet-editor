package main

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func isPrintable(bytes []byte) bool {
	return bytes[0] >= 32 && bytes[0] <= 126
}

func fixedWidth(s string, width int) string {
	if len(s) > width {
		if width < 1 {
			return ""
		}
		return s[0:width-1] + "…"
	}

	for len(s) < width {
		s += " "
	}

	return s
}

func getColumnWidth(column int) int {
	if width, ok := columnWidthMap[column]; ok {
		return max(width, 3)
	}
	return 8
}

func setColumnWidth(column int, width int) {
	columnWidthMap[column] = width
}

func getColumnName(column int) string {
	column = column % len(alphabet)
	if column < 0 {
		column = len(alphabet) + column
	}
	return alphabet[column : column+1]
}

func maxColumn() int {
	maxColumn := 0
	for row := range contentMap {
		for column := range contentMap[row] {
			maxColumn = max(maxColumn, column)
		}
	}
	return maxColumn
}

func maxRow() int {
	maxRow := 0
	for row := range contentMap {
		maxRow = max(maxRow, row)
	}
	return maxRow
}

func equalizeColumns() {
	columnWidthMap = make(map[int]int)
	for row := range contentMap {
		for column := range contentMap[row] {
			content := contentMap[row][column]
			columnWidthMap[column] = max(columnWidthMap[column], len(content)+1)
		}
	}
}
