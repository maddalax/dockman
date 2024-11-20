package filekv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileLocker ensures safe concurrent access
type FileLocker struct {
	mu sync.Mutex
}

// ReadKeyValues reads all key-value pairs from a file into a map
func ReadKeyValues(filename string) (map[string]string, error) {
	kvMap := make(map[string]string)

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(filename)
			os.MkdirAll(dir, 0700)
			_, err = os.Create(filename)
			if err != nil {
				return nil, err
			}
			// Return an empty map if the file does not exist
			return kvMap, nil
		}
		return nil, err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			kvMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return kvMap, nil
}

// WriteKeyValue safely writes or updates a key-value pair in a file
func WriteKeyValue(filename, key, value string, locker *FileLocker) error {
	// Acquire lock
	locker.mu.Lock()
	defer locker.mu.Unlock()

	// Load existing key-value pairs
	kvMap, err := ReadKeyValues(filename)
	if err != nil {
		return err
	}

	// Update or add the key-value pair
	kvMap[key] = value

	// Write the updated map back to the file
	tempFile := filename + ".tmp"
	output, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer output.Close()

	for k, v := range kvMap {
		_, err := fmt.Fprintf(output, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}

	// Replace the original file with the updated one
	return os.Rename(tempFile, filename)
}
