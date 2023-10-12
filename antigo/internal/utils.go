package internal

import (
	"os"
	"path/filepath"
)

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}

func isDirectoryEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			isEmpty, err := isDirectoryEmpty(filepath.Join(path, entry.Name()))
			if err != nil {
				return false, err
			}
			if !isEmpty {
				return false, nil
			}
		} else {
			return false, nil
		}
	}

	return true, nil
}
