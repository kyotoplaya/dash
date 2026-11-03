package main

import (
	"os"
)

func isDirExists(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}
