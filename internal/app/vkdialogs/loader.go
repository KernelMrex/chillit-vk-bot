package vkdialogs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// LoadFromFile loads dialogs from file
func LoadFromFile(path string) (Dialogs, error) {
	if !fileExists(path) {
		return nil, fmt.Errorf("could not open file: file '%s' does not exists", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()

	return Load(f)
}

// Load loads dialogs from reader
func Load(in io.Reader) (Dialogs, error) {
	dialogs := make(map[string]*Message)
	if err := json.NewDecoder(in).Decode(&dialogs); err != nil {
		return nil, fmt.Errorf("error while parsing config: %v", err)
	}
	return dialogs, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
