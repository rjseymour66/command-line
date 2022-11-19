package main

import (
	"bufio"
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
	byteCount := flag.Bool("b", false, "Count bytes")
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines, *byteCount))
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
