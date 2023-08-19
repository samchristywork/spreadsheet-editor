package main

import (
	"fmt"
	"github.com/pkg/term"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var bytes []byte

var contentMap = map[int]map[int]string{}
var currentCell = []int{0, 0}
var frame = 0

func alternateScreen() {
	fmt.Printf("\033[?1049h")
}

func normalScreen() {
	fmt.Printf("\033[?1049l")
}

func screenDimensions() (int, int) {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		return 0, 0
	}
	return width, height
}

func getColumnWidth(column int) int {
	return 5
}

func getColumnName(column int) string {
	column = column % len(alphabet)
	return alphabet[column : column+1]
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

func renderHeadings(x int, y int) {
	width, height := screenDimensions()

	// Render top headings
	xoff := 3
	for i := 0; xoff < width; i++ {
		setCursorPosition(x+xoff, y)

		columnWidth := getColumnWidth(i)
		fmt.Printf("\033[7m")
		fmt.Printf("%s", fixedWidth(getColumnName(i), columnWidth))
		fmt.Printf("\033[0m")

		xoff += columnWidth
	}

	// Render left headings
	for i := 0; i < height; i++ {
		setCursorPosition(x, y+i+1)
		fmt.Printf("\033[7m")
		fmt.Printf("%s", fixedWidth(fmt.Sprintf("%d", i), 3))
		fmt.Printf("\033[0m")
	}
}

func renderCell(row int, column int, width int) {
	content := contentMap[row][column]
	if row == currentCell[0] && column == currentCell[1] {
		fmt.Printf("\033[7m")
		fmt.Printf("%s", fixedWidth(content, width))
		fmt.Printf("\033[0m")
	} else {
		fmt.Printf("%s", fixedWidth(content, width))
	}
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
	renderGrid()
}

func setCursorPosition(x, y int) {
	fmt.Printf("\033[%d;%df", y, x)
}

func makeCursorInvisible() {
	fmt.Printf("\033[?25l")
}

func makeCursorVisible() {
	fmt.Printf("\033[?25h")
}

func keyPressed(a byte, b byte, c byte, bytes []byte) bool {
	return bytes[0] == a && bytes[1] == b && bytes[2] == c
}

func isPrintable(bytes []byte) bool {
	return bytes[0] >= 32 && bytes[0] <= 126
}

func getCellContent(row int, column int) *string {
	if contentMap[row] == nil {
		return nil
	}

	content := contentMap[row][column]
	return &content
}

func editCell(t *term.Term) *string {
	width, _ := screenDimensions()

	entry := ""
	entryReference := getCellContent(currentCell[0], currentCell[1])
	if entryReference != nil {
		entry = *entryReference
	}

	x := len(entry) + 1

	for {
		setCursorPosition(1, 2)
		fmt.Printf("%s", fixedWidth(entry, width))

		setCursorPosition(x, 2)
		makeCursorVisible()
		bytes[0] = 0
		bytes[1] = 0
		bytes[2] = 0
		_, err := t.Read(bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
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


func main() {
	t, err := term.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
		return
	}
	defer t.Restore()

	alternateScreen()
	makeCursorInvisible()
	t.SetRaw()

	bytes = make([]byte, 3) // buffer to read escape sequences
	for {
		render()

		bytes[0] = 0
		bytes[1] = 0
		bytes[2] = 0
		_, err := t.Read(bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
			break
		}

		if keyPressed(27, 0, 0, bytes) { // Escape
			break
		} else if keyPressed(27, 91, 65, bytes) { // Up
			currentCell[0]--
		} else if keyPressed(27, 91, 66, bytes) { // Down
			currentCell[0]++
		} else if keyPressed(27, 91, 67, bytes) { // Right
			currentCell[1]++
		} else if keyPressed(27, 91, 68, bytes) { // Left
			currentCell[1]--
		} else if keyPressed(byte('h'), 0, 0, bytes) {
			currentCell[1]--
		} else if keyPressed(byte('j'), 0, 0, bytes) {
			currentCell[0]++
		} else if keyPressed(byte('k'), 0, 0, bytes) {
			currentCell[0]--
		} else if keyPressed(byte('l'), 0, 0, bytes) {
			currentCell[1]++
		} else if keyPressed(byte('H'), 0, 0, bytes) {
			currentCell[1] -= 5
		} else if keyPressed(byte('J'), 0, 0, bytes) {
			currentCell[0] += 5
		} else if keyPressed(byte('K'), 0, 0, bytes) {
			currentCell[0] -= 5
		} else if keyPressed(byte('L'), 0, 0, bytes) {
			currentCell[1] += 5
		} else if keyPressed(13, 0, 0, bytes) { // Enter
			content := editCell(t)
			if content != nil {
				if contentMap[currentCell[0]] == nil {
					contentMap[currentCell[0]] = make(map[int]string)
				}
				contentMap[currentCell[0]][currentCell[1]] = *content
			}
		}
	}

	normalScreen()
	makeCursorVisible()
}
