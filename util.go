package main

import (
	"os"
	"path/filepath"
)

func createDir(dir string, perm os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		if err := os.MkdirAll(dir, perm); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func createFileDir(filename string, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	return createDir(dir, perm)
}
