package main

import (
	"fmt"
	"github.com/pkg/term"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var bytes []byte

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

	xoff := 3
	for i := 0; xoff < width; i++ {
		setCursorPosition(x+xoff, y)

		columnWidth := getColumnWidth(i)
		fmt.Printf("\033[7m")
		fmt.Printf("%s", fixedWidth(getColumnName(i), columnWidth))
		fmt.Printf("\033[0m")

		xoff += columnWidth
	}

	for i := 0; i < height; i++ {
		setCursorPosition(x, y+i+1)
		fmt.Printf("\033[7m")
		fmt.Printf("%s", fixedWidth(fmt.Sprintf("%d", i), 3))
		fmt.Printf("\033[0m")
	}
}

func render() {
	setCursorPosition(1, 1)
	fmt.Printf("Hello, World!\n")
	frame++

	renderHeadings(1, 3)
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
		}
	}

	normalScreen()
	makeCursorVisible()
}
