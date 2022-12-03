package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
</html>`
)

// content type represents the HTML content to add into the template
type content struct {
	Title string
	Body  template.HTML
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	// if user did not provide a tFname and did provide an env var
	if *tFname == "" && os.Getenv("DEFAULT_TEMPLATE") != "" {
		*tFname = os.Getenv("DEFAULT_TEMPLATE")
	}

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// reads md file, converts to html, writes to tmp file,
func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// convert markdown > HTML
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	// Create temp file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)

	// saveHTML returns an error if it fails
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

// Accepts markdown and returns html
func parseContent(input []byte, tFname string) ([]byte, error) {
	// Parse the md file through blackfriday to generate HTML
	output := blackfriday.Run(input)
	// send the md content through blue monday for security
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content type, adding the title and body
	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	// compose the page using a buffer of bytes to write to a file
	var buffer bytes.Buffer

	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outName string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outName, data, 0644)
}

func preview(fname string) error {
	cName := ""
	// slice literal that creates empty slice of strings
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS Not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}

	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Give browser time to open before deleting it
	time.Sleep(2 * time.Second)
	return err
}

// TODO
func getFile(r io.Reader, arg string) (string, error) {
	var file string
	// if flag arg present
	if arg != "" {
		file = arg
		return file, nil
	}
	// scan filename from STDIN
	s := bufio.NewScanner(r)
	// scan by word, not line
	s.Split(bufio.ScanWords)

	for s.Scan() {
		// if non-EOF error
		if err := s.Err(); err != nil {
			return "", err
		}
		file = s.Text()

		if len(s.Text()) == 0 {
			return "", fmt.Errorf("Must provide file")
		}
	}
	return file, nil
}
