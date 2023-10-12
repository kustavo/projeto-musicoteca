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

func isDiretorioVazio(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			isEmpty, err := isDiretorioVazio(filepath.Join(path, entry.Name()))
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

func ObterDiretorios(path string) ([]string, error) {
	diretorios := []string{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return diretorios, err
	}

	for _, e := range entries {
		if e.IsDir() {
			nomeDir := e.Name()
			if nomeDir[:5] != "_000_" {
				diretorios = append(diretorios, nomeDir)
			}
		}
	}
	return diretorios, nil
}
