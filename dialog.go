package main

import (
	"fmt"
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
