package main

import (
	"fmt"
	"sort"
	"strings"
)

func messageBox(title string, message string) {
	width, height := screenDimensions()

	messages := strings.Split(message, "\n")

	stringWidth := len(title)
	for i := 0; i < len(messages); i++ {
		if len(messages[i]) > stringWidth {
			stringWidth = len(messages[i])
		}
	}

	boxWidth := stringWidth
	boxHeight := len(messages) + 2

	x := width/2 - boxWidth/2
	y := height/2 - boxHeight/2

	clearRectangle(x-2, y-2, boxWidth+4, boxHeight+2)
	drawRectangle(x-2, y-2, boxWidth+4, boxHeight+2)

	setCursorPosition(x, y-1)
	invert()
	fmt.Printf("%s", fixedWidth(title, boxWidth))
	resetColor()

	for i := 0; i < len(messages); i++ {
		setCursorPosition(x, y+1+i)
		fmt.Printf("%s", fixedWidth(messages[i], boxWidth))
	}
}

func helpBox(title string) {
	width, height := screenDimensions()

	keys := make([]string, 0, len(shortcuts))
	for k := range shortcuts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	type block struct {
		header string
		items  []helpItem
	}

	blocks := make([]block, len(keys))
	for i, k := range keys {
		blocks[i] = block{k, shortcuts[k]}
	}

	blockWidth := func(b block) int {
		w := len(b.header)
		for _, item := range b.items {
			if lineWidth := 12 + len(item.description); lineWidth > w {
				w = lineWidth
			}
		}
		return w
	}

	numCols := 2
	colSize := (len(blocks) + numCols - 1) / numCols

	colWidths := make([]int, numCols)
	colHeights := make([]int, numCols)
	for col := 0; col < numCols; col++ {
		start := col * colSize
		end := min(start+colSize, len(blocks))
		h := 0
		for i := start; i < end; i++ {
			if w := blockWidth(blocks[i]); w > colWidths[col] {
				colWidths[col] = w
			}
			h += 1 + len(blocks[i].items)
			if i < end-1 {
				h++ // blank line between blocks
			}
		}
		colHeights[col] = h
	}

	colSep := 4
	totalWidth := colWidths[0] + colSep + colWidths[1]
	if len(title) > totalWidth {
		totalWidth = len(title)
	}

	boxHeight := max(colHeights[0], colHeights[1]) + 2

	bx := width/2 - totalWidth/2
	by := height/2 - boxHeight/2

	clearRectangle(bx-2, by-2, totalWidth+4, boxHeight+2)
	drawRectangle(bx-2, by-2, totalWidth+4, boxHeight+2)

	setCursorPosition(bx, by-1)
	invert()
	fmt.Printf("%s", fixedWidth(title, totalWidth))
	resetColor()

	colX := [2]int{bx, bx + colWidths[0] + colSep}

	for col := 0; col < numCols; col++ {
		start := col * colSize
		end := min(start+colSize, len(blocks))
		row := 0
		for i := start; i < end; i++ {
			b := blocks[i]
			setCursorPosition(colX[col], by+1+row)
			fmt.Printf("%s", b.header)
			row++
			for _, item := range b.items {
				setCursorPosition(colX[col], by+1+row)
				fmt.Printf("%s%s\033[0m%s", colorString(0, 255, 255), fixedWidth(item.name, 12), item.description)
				row++
			}
			if i < end-1 {
				row++
			}
		}
	}
}

func promptBox(title string, message string) bool {
	messageBox(title, message)
	for {
		bytes, err := nextKeyPress()
		if err != nil {
			panic(err)
		}

		if keyPressed('y', 0, 0, bytes) {
			return true
		} else if keyPressed('n', 0, 0, bytes) {
			return false
		}
	}
}
