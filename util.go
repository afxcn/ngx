package main

import (
	"os"
	"path/filepath"
)

func mkdirAll(dir string, perm os.FileMode) error {
	return os.MkdirAll(dir, perm)
}

func sameDir(filename string, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	return mkdirAll(dir, perm)
}
