package main

import (
	"fmt"
	"os"
)

func keyPressed(a byte, b byte, c byte, bytes []byte) bool {
	return bytes[0] == a && bytes[1] == b && bytes[2] == c
}

func nextKeyPress() ([]byte, error) {
	bytes := make([]byte, 3)
	bytes[0] = 0
	bytes[1] = 0
	bytes[2] = 0
	_, err := os.Stdin.Read(bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading keypress: %v\n", err)
		return bytes, err
	}

	return bytes, nil
}

func handleClipboard(bytes []byte) bool {
	if keyPressed(byte('y'), 0, 0, bytes) {
		bytes, err := nextKeyPress()
		if err != nil {
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

func handleMovement(bytes []byte) bool {
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
