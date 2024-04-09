package main

import (
	"github.com/pkg/term"
	"os"
	"strconv"
)

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

func registerIncrement() {
	category := "Increment"

	shortcut(59, 50, 65, category, "Shift-Up", "Increment cell upwards", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])
		handleIncrement(content, currentCell[0], currentCell[1], -1, 0)
		currentCell[0]--
	})

	shortcut(59, 50, 66, category, "Shift-Down", "Increment cell downwards", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])
		handleIncrement(content, currentCell[0], currentCell[1], 1, 0)
		currentCell[0]++
	})

	shortcut(59, 50, 67, category, "Shift-Right", "Increment cell rightwards", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])
		handleIncrement(content, currentCell[0], currentCell[1], 0, 1)
		currentCell[1]++
	})

	shortcut(59, 50, 68, category, "Shift-Left", "Increment cell leftwards", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])
		handleIncrement(content, currentCell[0], currentCell[1], 0, -1)
		currentCell[1]--
	})
}

func registerGeneral(t *term.Term) {
	category := "General"

	shortcut(27, 79, 80, category, "F1", "Show help", func() {
		content := ""

		keys := make([]string, 0, len(shortcuts))
		for k := range shortcuts {
			keys = append(keys, k)
		}

		for i := 0; i < len(keys); i++ {
			//content += colorString(0, 0, 255)
			//content += backgroundColorString(20, 20, 20)
			content += keys[i] + "\n"
			//content += "\033[0m"

			for j := 0; j < len(shortcuts[keys[i]]); j++ {
				content += colorString(0, 255, 255)
				content += fixedWidth(shortcuts[keys[i]][j].name, 12)
				content += "\033[0m"
				content += shortcuts[keys[i]][j].description + "\n"
			}

			if i != len(keys)-1 {
				content += "\n"
			}
		}

		//for i := range shortcuts {
		//	content += fixedWidth(shortcuts[i].name, 12) + shortcuts[i].description + "\n"
		//}
		messageBox("Help", content)
		nextKeyPress()
	})

	shortcut(27, 0, 0, category, "Escape", "Quit out of the program", func() {
		if quit() {
			running = false
		}
	})

	shortcut(byte('q'), 0, 0, category, "q", "Quit out of the program", func() {
		if quit() {
			running = false
		}
	})

	shortcut(byte('s'), 0, 0, category, "s", "Save the tsv file", func() {
		filename := os.Args[1]
		save(filename)
	})

	shortcut(13, 0, 0, category, "Enter", "Edit cell", func() {
		content := editCell(t)
		if content != nil {
			setCellContent(currentCell[0], currentCell[1], *content)
			currentCell[0]++
			modified = true
		}
	})
}

func registerMiscellaneous(t *term.Term) {
	category := "Miscellaneous"

	shortcut(byte('y'), 0, 0, category, "y", "Copy current cell", func() {
		bytes, _ := nextKeyPress()

		if keyPressed(byte('y'), 0, 0, bytes) {
			copyCell()
		}
	})

	shortcut(byte('p'), 0, 0, category, "p", "Paste current cell", func() {
		pasteCell()
	})

	shortcut(byte('e'), 0, 0, category, "e", "Edit the tsv file with vim", func() {
		if modified {
			messageBox("Unsaved changes", "Cannot edit the file unless it is saved.")
			nextKeyPress()
		}

		normalScreen()
		t.Restore()
		editFile()

		alternateScreen()
		makeCursorInvisible()
		t.SetRaw()
		loadFile()
	})

	shortcut(1, 0, 0, category, "Ctrl-A", "Increment cell content", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])

		contentInt, _ := strconv.Atoi(content)

		contentInt++
		content = strconv.Itoa(contentInt)
		setCellContent(currentCell[0], currentCell[1], content)
	})

	shortcut(24, 0, 0, category, "Ctrl-X", "Decrement cell content", func() {
		content, _ := getCellContent(currentCell[0], currentCell[1])

		contentInt, _ := strconv.Atoi(content)

		contentInt--
		content = strconv.Itoa(contentInt)
		setCellContent(currentCell[0], currentCell[1], content)
	})

	shortcut(byte('c'), 0, 0, category, "c", "Mark cell with color", func() {
		c, _ := getCellColor(currentCell[0], currentCell[1])
		if c != nil {
			if len(c) == 3 {
				delete(cellColorMap[currentCell[0]], currentCell[1])
			}
		}

		setCellColor(currentCell[0], currentCell[1])
	})

	shortcut(byte('g'), 0, 0, category, "g", "Show grid", func() {
		showGrid = !showGrid
	})

	shortcut(byte('x'), 0, 0, category, "x", "Delete cell contents", func() {
		delete(contentMap[currentCell[0]], currentCell[1])
	})

	shortcut(byte('+'), 0, 0, category, "+", "Reset column widths", func() {
		columnWidthMap = make(map[int]int)
	})

	shortcut(byte('='), 0, 0, category, "=", "Equalize column widths", func() {
		// TODO: Doesn't work with scroll
		equalizeColumns()
	})
}

func registerFunctions(t *term.Term) {

	registerGeneral(t)

	registerMovement()

	registerIncrement()

	registerMiscellaneous(t)
}
