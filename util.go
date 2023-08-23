package main

import (
	"strconv"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func isPrintable(bytes []byte) bool {
	return bytes[0] >= 32 && bytes[0] <= 126
}

func isCapitalLetter(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
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
	if column < 26 {
		return string(rune('A' + column))
	}

	return getColumnName(column/26-1) + getColumnName(column%26)
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
			if len(content) > 0 && content[0] == '=' {
				content = eval(content[1:])
			}
			columnWidthMap[column] = max(columnWidthMap[column], len(content)+1)
		}
	}
}

func splitColumnRow(s string) (string, string) {
	if len(s) < 1 {
		return "", ""
	}

	pivot := 0
	for pivot < len(s) && isCapitalLetter(rune(s[pivot])) {
		pivot++
	}

	return s[0:pivot], s[pivot:]
}


func column(s string) int {
	column, _ := splitColumnRow(s)

	if isCapitalLetter(rune(s[0])) {
		return 0
	}

	for i := 0; i < len(s); i++ {
		if isCapitalLetter(rune(s[i])) {
			return col - 1
		}

		col = col*26 + int(s[i]-'a'+1)
	}

	return col - 1
}
