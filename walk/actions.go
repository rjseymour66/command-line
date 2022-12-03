package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Filters out the given path if it is a directory, is greater than the minSize,
// or the file extension does not match the ext
func filterOut(path string, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
