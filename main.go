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
var scrollOffset = []int{0, 0}
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

var running = true

func registerMovement() {
	category := "Movement"

	shortcut(27, 91, 65, category, "Up", "Move up one cell", func() {
		scrollOffset[1]--
	})

	shortcut(27, 91, 66, category, "Down", "Move down one cell", func() {
		scrollOffset[1]++
	})

	shortcut(27, 91, 67, category, "Right", "Move right one cell", func() {
		scrollOffset[0]++
	})

	shortcut(27, 91, 68, category, "Left", "Move left one cell", func() {
		scrollOffset[0]--
	})

	shortcut(byte('h'), 0, 0, category, "h", "Move left one cell", func() {
		currentCell[1]--
	})

	shortcut(byte('j'), 0, 0, category, "j", "Move down one cell", func() {
		currentCell[0]++
	})

	shortcut(byte('k'), 0, 0, category, "k", "Move up one cell", func() {
		currentCell[0]--
	})

	shortcut(byte('l'), 0, 0, category, "l", "Move right one cell", func() {
		currentCell[1]++
	})

	shortcut(byte('H'), 0, 0, category, "H", "Move left five cells", func() {
		currentCell[1] -= 5
	})

	shortcut(byte('J'), 0, 0, category, "J", "Move down five cells", func() {
		currentCell[0] += 5
	})

	shortcut(byte('K'), 0, 0, category, "K", "Move up five cells", func() {
		currentCell[0] -= 5
	})

	shortcut(byte('L'), 0, 0, category, "L", "Move right five cells", func() {
		currentCell[1] += 5
	})

	shortcut(byte('0'), 0, 0, category, "0", "Move cursor to origin", func() {
		scrollOffset[0] = 0
		scrollOffset[1] = 0
		currentCell[0] = 0
		currentCell[1] = 0
	})
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

	bytes := make([]byte, 3)
	for {
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

		if handleMovement(bytes) {
			continue
		}

		if handleClipboard(bytes) {
			continue
		}

		if keyPressed(27, 0, 0, bytes) { // Escape
			if quit() {
				break
			}
		} else if keyPressed(byte('q'), 0, 0, bytes) {
			if quit() {
				break
			}
		} else if keyPressed(59, 50, 65, bytes) { // Up
			content, _ := getCellContent(currentCell[0], currentCell[1])
			setCellContent(currentCell[0]-1, currentCell[1], content)
			currentCell[0]--
		} else if keyPressed(59, 50, 66, bytes) { // Down
			content, _ := getCellContent(currentCell[0], currentCell[1])
			setCellContent(currentCell[0]+1, currentCell[1], content)
			currentCell[0]++
		} else if keyPressed(59, 50, 67, bytes) { // Right
			content, _ := getCellContent(currentCell[0], currentCell[1])
			setCellContent(currentCell[0], currentCell[1]+1, content)
			currentCell[1]++
		} else if keyPressed(59, 50, 68, bytes) { // Left
			content, _ := getCellContent(currentCell[0], currentCell[1])
			setCellContent(currentCell[0], currentCell[1]-1, content)
			currentCell[1]--
		} else if keyPressed(byte('e'), 0, 0, bytes) {
			if modified {
				messageBox("Unsaved changes", "Cannot edit the file unless it is saved.")
				nextKeyPress()
				continue
			}
			normalScreen()
			t.Restore()
			editFile()

			alternateScreen()
			makeCursorInvisible()
			t.SetRaw()
			loadFile()
		} else if keyPressed(byte('s'), 0, 0, bytes) {
			filename := os.Args[1]
			save(filename)
		} else if keyPressed(1, 0, 0, bytes) { // Ctrl-A
			content, _ := getCellContent(currentCell[0], currentCell[1])

			contentInt, err := strconv.Atoi(content)
			if err != nil {
				continue
			}

			contentInt++
			content = strconv.Itoa(contentInt)
			setCellContent(currentCell[0], currentCell[1], content)
		} else if keyPressed(24, 0, 0, bytes) { // Ctrl-X
			content, _ := getCellContent(currentCell[0], currentCell[1])

			contentInt, err := strconv.Atoi(content)
			if err != nil {
				continue
			}

			contentInt--
			content = strconv.Itoa(contentInt)
			setCellContent(currentCell[0], currentCell[1], content)
		} else if keyPressed(byte('c'), 0, 0, bytes) {
			c, _ := getCellColor(currentCell[0], currentCell[1])
			if c != nil {
				if len(c) == 3 {
					delete(cellColorMap[currentCell[0]], currentCell[1])
					continue
				}
			}

			setCellColor(currentCell[0], currentCell[1])
		} else if keyPressed(byte('0'), 0, 0, bytes) {
			scrollOffset[0] = 0
			scrollOffset[1] = 0
			currentCell[0] = 0
			currentCell[1] = 0
		} else if keyPressed(byte('g'), 0, 0, bytes) {
			showGrid = !showGrid
		} else if keyPressed(byte('x'), 0, 0, bytes) {
			delete(contentMap[currentCell[0]], currentCell[1])
		} else if keyPressed(byte('+'), 0, 0, bytes) {
			columnWidthMap = make(map[int]int)
		} else if keyPressed(byte('='), 0, 0, bytes) {
			// TODO: Doesn't work with scroll
			equalizeColumns()
		} else if keyPressed(13, 0, 0, bytes) { // Enter
			content := editCell(t)
			if content != nil {
				setCellContent(currentCell[0], currentCell[1], *content)
				currentCell[0]++
				modified = true
			}
		}
	}

	normalScreen()
	makeCursorVisible()
}
