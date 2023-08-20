package main

import (
	"fmt"
)

func drawRectangle(x int, y int, width int, height int) {
	setCursorPosition(x, y)
	fmt.Printf("┌")
	for i := 0; i < width-2; i++ {
		fmt.Printf("─")
	}
	fmt.Printf("┐")

	for i := 0; i < height-2; i++ {
		setCursorPosition(x, y+i+1)
		fmt.Printf("│")
		setCursorPosition(x+width-1, y+i+1)
		fmt.Printf("│")
	}

	setCursorPosition(x, y+height-1)
	fmt.Printf("└")
	for i := 0; i < width-2; i++ {
		fmt.Printf("─")
	}
	fmt.Printf("┘")
}

func clearRectangle(x int, y int, width int, height int) {
	for i := 0; i < height; i++ {
		setCursorPosition(x, y+i)
		for j := 0; j < width; j++ {
			fmt.Printf(" ")
		}
	}
}
