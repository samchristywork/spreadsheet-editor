package main

import (
	"fmt"
	"github.com/pkg/term"
	"os"
)

var bytes []byte

func alternateScreen() {
	fmt.Printf("\033[?1049h")
}

func normalScreen() {
	fmt.Printf("\033[?1049l")
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
		//		render()
		//
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
