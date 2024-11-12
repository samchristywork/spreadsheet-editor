package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	for row := 0; row <= maxRow(); row++ {
		for column := 0; column <= maxColumn(); column++ {
			content, _ := getCellContent(row, column)
			file.WriteString(content)

			if column < maxColumn() {
				file.WriteString("	")
			}
		}
		file.WriteString("\n")
	}

	modified = false
}

func loadFile() {
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

		line = strings.TrimSuffix(line, "\n")

		if len(line) > 0 {
			segments := strings.Split(line, "	")
			for i := 0; i < len(segments); i++ {
				setCellContent(lineNumber, i, segments[i])
			}
			lineNumber++
		}

		if err != nil {
			break
		}
	}
}

func editFile() {
	filename := os.Args[1]

	cmd := exec.Command("vim", filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
