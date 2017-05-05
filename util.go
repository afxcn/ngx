package main

import "os"

func createDir(dir string, perm os.FileMode) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		if err := os.MkdirAll(dir, perm); err != nil {
			fatalf("create dir failure: %v", err)
		}
	}
}
