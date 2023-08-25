package main

import (
	"fmt"
	"github.com/pkg/term"
	"log"
	"os"
	"strconv"
)

var cellColorMap = map[int]map[int][]int{}
var clipboard string
var columnWidthMap = map[int]int{}
var contentMap = map[int]map[int]string{}
var currentCell = []int{0, 0}
var frame = 0
var modified = false
var running = true
var scrollOffset = []int{0, 0}
var shortcuts = map[string][]helpItem{}
var showGrid = false

func quit() bool {
	if modified {
		return promptBox("Unsaved changes.", "Really quit? (y/n)")
	} else {
		return true
	}
}

func tokenize(content string) ([]string, error) {
	var tokens []string
	var currentToken string

	for _, c := range content {
		if c >= '0' && c <= '9' {
			currentToken += string(c)
		} else if c >= 'a' && c <= 'z' {
			currentToken += string(c)
		} else if c >= 'A' && c <= 'Z' {
			currentToken += string(c)
		} else {
			if len(currentToken) > 0 {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}

			tokens = append(tokens, string(c))
		}
	}

	if len(currentToken) > 0 {
		tokens = append(tokens, currentToken)
	}

	return tokens, nil
}

func isCellIdentifier(token string) bool {
	row := row(token)
	col := column(token)

	if token == getColumnName(col)+strconv.Itoa(row) {
		return true
	}

	return false
}

func incrementToken(token string, colDelta int, rowDelta int) string {
	r := row(token)
	c := column(token)

	return getColumnName(c+colDelta) + strconv.Itoa(r+rowDelta)
}

func handleIncrement(content string, row int, col int, rowDelta int, colDelta int) {
	if len(content) > 0 && content[0] == '=' {
		tokens, err := tokenize(content)
		if err != nil {
			return
		}

		newContent := ""
		for i := 0; i < len(tokens); i++ {
			if !isCellIdentifier(tokens[i]) {
				newContent += tokens[i]
				continue
			}

			newContent += incrementToken(tokens[i], colDelta, rowDelta)
		}

		setCellContent(currentCell[0]+rowDelta, currentCell[1]+colDelta, newContent)
		return
	}

	setCellContent(currentCell[0]+rowDelta, currentCell[1]+colDelta, content)
}

func handleKeypress(bytes []byte) {
	for i := range shortcuts {
		for c := range shortcuts[i] {
			if keyPressed(shortcuts[i][c].key[0], shortcuts[i][c].key[1], shortcuts[i][c].key[2], bytes) {
				shortcuts[i][c].function()
			}
		}
	}
}

func main() {
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	defer f.Close()

	log.SetOutput(f)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		return
	}

	loadFile()

	t, err := term.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
		return
	}
	defer t.Restore()

	alternateScreen()
	makeCursorInvisible()
	t.SetRaw()

	registerFunctions(t)

	bytes := make([]byte, 3)
	for running {
		frame++

		if currentCell[0] < 0 {
			currentCell[0] = 0
		}
		if currentCell[1] < 0 {
			currentCell[1] = 0
		}

		render(bytes)

		bytes, err = nextKeyPress()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading keypress: %v\n", err)
			return
		}

		handleKeypress(bytes)

	}

	normalScreen()
	makeCursorVisible()
}
