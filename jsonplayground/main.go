package main

import (
	"fmt"
	"os"
)

func readFromFileToConsole(filename string) ([]byte, error) {
	// open the file
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		err := fmt.Errorf("Error reading from file: %s", filename)
		return nil, err
	}

	return fileContents, nil
}

func writeJSONtoFile(inFile string, outFile string) error {
	data, err := readFromFileToConsole(inFile)
	if err != nil {
		return fmt.Errorf("Did not read from file: %s", err)
	}
	os.WriteFile(outFile, data, 0644)
	return nil
}

func main() {
	fileContent, _ := readFromFileToConsole("person.json")
	err := writeJSONtoFile(string(fileContent), "write_test")
	if err != nil {
		fmt.Printf("%s", err)
	}
	// fmt.Printf("File contents:\n%s\n", fileContent)
}
