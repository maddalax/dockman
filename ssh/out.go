package main

import (
	"bufio"
	"fmt"
	"io"
)

func streamOutput(pipe io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		writer.Write(scanner.Bytes())
		writer.Write([]byte("\n"))
	}
	if err := scanner.Err(); err != nil {
		_, err := writer.Write([]byte(fmt.Sprintf("Error reading output: %v", err)))
		if err != nil {
			return
		}
	}
}
