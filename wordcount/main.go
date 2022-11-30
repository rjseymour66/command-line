package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
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

	l, b := count(os.Stdin, *lines, *byteCount)
	// count lines
	// count bytes
	// count words and bytes
	// read from file
	// default is count words
	switch {
	case *lines:
		fmt.Println(l)
	case *byteCount:
		fmt.Println(b)
	case *lines == false && *byteCount:
		fmt.Println(count(os.Stdin, *lines, *byteCount))
	case *file != "":
		// create buffer
		var buffer bytes.Buffer
		// read file into var
		output, err := os.ReadFile(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// write var into buffer
		buffer.Write(output)
		// pass buffer to count()
		fmt.Println(count(&buffer, *lines, *byteCount))
	}

	if *file == "" {
		fmt.Println(count(os.Stdin, *lines, *byteCount))
	} else {
		// create buffer
		var buffer bytes.Buffer
		// read file into var
		output, err := os.ReadFile(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// write var into buffer
		buffer.Write(output)
		// pass buffer to count()
		fmt.Println(count(&buffer, *lines, *byteCount))
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
