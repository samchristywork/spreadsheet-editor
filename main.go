package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/term"
	"os"
	"strconv"
	"strings"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var bytes []byte
var clipboard string
var columnWidthMap = map[int]int{}
var contentMap = map[int]map[int]string{}
var cellColorMap = map[int]map[int][]int{}
var currentCell = []int{0, 0}
var frame = 0
var scrollOffset = []int{0, 0}
var showGrid = false

func save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	for row := 0; row <= maxRow(); row++ {
		for column := 0; column <= maxColumn(); column++ {
			contentReference := getCellContent(row, column)
			if contentReference != nil {
				content := *contentReference
				file.WriteString(content)
			}
			if column < maxColumn() {
				file.WriteString(",")
			}
		}
		file.WriteString("\n")
	}
}

func handleMovement() bool {
	if keyPressed(27, 91, 65, bytes) { // Up
		scrollOffset[1]--
	} else if keyPressed(27, 91, 66, bytes) { // Down
		scrollOffset[1]++
	} else if keyPressed(27, 91, 67, bytes) { // Right
		scrollOffset[0]++
	} else if keyPressed(27, 91, 68, bytes) { // Left
		scrollOffset[0]--
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
	} else {
		return false
	}

	return true
}

func handleClipboard() bool {
	if keyPressed(byte('y'), 0, 0, bytes) {
		if !nextKeyPress() {
			return true
		}

		if keyPressed(byte('y'), 0, 0, bytes) {
			copyCell()
			return true
		}
	} else if keyPressed(byte('p'), 0, 0, bytes) {
		pasteCell()
		return true
	}

	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	lineNumber := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSuffix(line, "\n")

		segments := strings.Split(line, ",")
		for i := 0; i < len(segments); i++ {
			setCellContent(lineNumber, i, segments[i])
		}

		lineNumber++
	}

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
		if currentCell[0] < 0 {
			currentCell[0] = 0
		}
		if currentCell[1] < 0 {
			currentCell[1] = 0
		}

		render()

		if !nextKeyPress() {
			break
		}

		if handleMovement() {
			continue
		}

		if handleClipboard() {
			continue
		}

		if keyPressed(27, 0, 0, bytes) { // Escape
			break
		} else if keyPressed(byte('q'), 0, 0, bytes) {
			break
		} else if keyPressed(byte('s'), 0, 0, bytes) {
			save(filename)
		} else if keyPressed(1, 0, 0, bytes) { // Ctrl-A
			content := ""
			contentReference := getCellContent(currentCell[0], currentCell[1])
			if contentReference != nil {
				content = *contentReference
			}

			contentInt, err := strconv.Atoi(content)
			if err != nil {
				continue
			}

			contentInt++
			content = strconv.Itoa(contentInt)
			setCellContent(currentCell[0], currentCell[1], content)
		} else if keyPressed(24, 0, 0, bytes) { // Ctrl-X
			content := ""
			contentReference := getCellContent(currentCell[0], currentCell[1])
			if contentReference != nil {
				content = *contentReference
			}

			contentInt, err := strconv.Atoi(content)
			if err != nil {
				continue
			}

			contentInt--
			content = strconv.Itoa(contentInt)
			setCellContent(currentCell[0], currentCell[1], content)
		} else if keyPressed(byte('c'), 0, 0, bytes) {
			c := getCellColor(currentCell[0], currentCell[1])
			if c != nil {
				if len(*c) == 3 {
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
		} else if keyPressed(byte('='), 0, 0, bytes) {
			// TODO: Doesn't work with scroll
			equalizeColumns()
		} else if keyPressed(13, 0, 0, bytes) { // Enter
			content := editCell(t)
			if content != nil {
				setCellContent(currentCell[0], currentCell[1], *content)
				currentCell[0]++
			}
		}
	}

	normalScreen()
	makeCursorVisible()
}
