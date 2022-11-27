package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>Markdown Preview Tool</title>
</head>
<body>
`
	footer = `
</body>
</html>`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// convert markdown > HTML
	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	// saveHTML returns an error if it fails
	return saveHTML(outName, htmlData)
}

// Accepts markdown and returns html
func parseContent(input []byte) []byte {
	// Parse the md file through blackfriday to generate HTML
	output := blackfriday.Run(input)
	// send the md content through blue monday for security
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// compose the page using a buffer of bytes to write to a file
	var buffer bytes.Buffer

	// write html to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outName string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outName, data, 0644)
}
