package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}

	// create the pipeline
	pipeline := make([]step, 2)

	// create a step
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)

	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)
	// loop through steps and execute
	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		// if executes without an error, print
		// success or failure msg
		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
