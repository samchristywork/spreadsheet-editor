package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func alternateScreen() {
	fmt.Printf("\033[?1049h")
	fmt.Printf("\033[2J") // Clear the screen
}

func normalScreen() {
	fmt.Printf("\033[?1049l")
}

func color(r int, g int, b int) {
	fmt.Printf("\033[38;2;%d;%d;%dm", r, g, b)
}

func colorString(r int, g int, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func backgroundColor(r int, g int, b int) {
	fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
}

func backgroundColorString(r int, g int, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

func resetColor() {
	fmt.Printf("\033[0m")
}

func invert() {
	fmt.Printf("\033[7m")
}

func setCursorPosition(x int, y int) {
	fmt.Printf("\033[%d;%df", y, x)
}

func makeCursorInvisible() {
	fmt.Printf("\033[?25l")
}

func makeCursorVisible() {
	fmt.Printf("\033[?25h")
}

func screenDimensions() (int, int) {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		return 0, 0
	}
	return width, height
}
