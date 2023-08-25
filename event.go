package main

import (
	"fmt"
	"os"
)

type helpItem struct {
	name        string
	description string
	key         []byte
	function    func()
}

func shortcut(a byte, b byte, c byte, category string, shortcut string, description string, function func()) {
	if _, ok := shortcuts[category]; !ok {
		shortcuts[category] = []helpItem{}
	}

	shortcuts[category] = append(shortcuts[category], helpItem{shortcut, description, []byte{a, b, c}, function})
}

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
