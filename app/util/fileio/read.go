package fileio

import (
	"bufio"
	"os"
)

func ReadLines(filePath string, cb func(line string)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// all good
		cb(line)
	}
	return nil
}
