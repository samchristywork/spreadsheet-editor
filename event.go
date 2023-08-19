package main

import (
	"fmt"
	"os"
)

func keyPressed(a byte, b byte, c byte, bytes []byte) bool {
	return bytes[0] == a && bytes[1] == b && bytes[2] == c
}

func nextKeyPress() bool {
	bytes[0] = 0
	bytes[1] = 0
	bytes[2] = 0
	_, err := os.Stdin.Read(bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		return false
	}

	return true
}
