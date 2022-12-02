package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
	wordcount counts the number of words or lines from STDIN.
*/

func main() {

	lines := flag.Bool("l", false, "Count lines")
	// add file flag
	file := flag.String("file", "", "File to read from")
	byteCount := flag.Bool("b", false, "Count bytes")
	flag.Parse()

	if *file == "" {
		fmt.Println(count(os.Stdin, *lines, *byteCount))
	} else {
		fileContents, err := getFile(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// create reader from string
		output := strings.NewReader(string(fileContents))
		// print output
		fmt.Println(count(output, *lines, *byteCount))
	}

}

func count(r io.Reader, countLines bool, countBytes bool) (int, int) {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0
	blength := 0

	if !countBytes {
		for scanner.Scan() {
			wc++
		}
	} else {
		for scanner.Scan() {
			wc++
			blength += len(scanner.Bytes())
		}
	}
	return wc, blength
}

func getFile(fileName string) ([]byte, error) {
	// create buffer
	var buffer bytes.Buffer
	// read file into var
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// write var into buffer
	buffer.Write(fileContents)
	return buffer.Bytes(), nil
}
