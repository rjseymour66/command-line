package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}
		})
	}
}

func TestRunDeleteExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		// no match occurs, all files match, mixed matches
		{name: "DeleteExtentionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 0, nNoDelete: 10,
			expected: ""},
		{name: "DeleteExtentionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "", nDelete: 10, nNoDelete: 0,
			expected: ""},
		{name: "DeleteExtentionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 5, nNoDelete: 5,
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			// set the logger with testCases.cfg
			tc.cfg.wLog = &logBuffer

			// files are deleted with the ext value used in the command
			// files that are not deleted are created but not acted upon
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			// run the tool from the tempDir, write to &buffer, and use the
			// config in the test cases. This deletes the files associated with the
			// cfg.ext and nDelete
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}
			// the buffer should be empty each time
			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}

			// the number of remaining files should equal the number that was
			// created but not deleted
			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(nil)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n",
					tc.nNoDelete, len(filesLeft))
			}

			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("expected %d log lines, got %d instead\n",
					expLogLines, len(lines))
			}
		})
	}
}

// creates a temp directory, then creates files that will be deleted, and
// files that will not be deleted according to the cfg setup.
// files is the number of files this function creates for each extension
// returns the dir name that it creates, and a func that cleans up test artifacts
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	// make a temp directory in the default directory for temp files
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		// for each number of files, write 'dummy' to "file1.log"
		// while j is less or equal to the number of files to delete or not
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k) // fname = file1.gz
			fpath := filepath.Join(tempDir, fname) // fpath = /tmp/walktest
			// write "dummy" to a file
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	// func removes path and its children
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtensionNoMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionMixed",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			// Create temp dirs for RunArchive test
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.nArchive,
				tc.extNoArchive: tc.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// match the dynamically named archive test files
			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}
			expOut := strings.Join(expFiles, "\n")
			res := strings.TrimSpace(buffer.String())

			if expOut != res {
				t.Errorf("Expected %q, got %q instead\n", expOut, res)
			}

			// get number of files that were archived
			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			// compare the number of files archived with the expected number
			if len(filesArchived) != tc.nArchive {
				t.Errorf("Expected %d files archived, got %d instead\n",
					tc.nArchive, len(filesArchived))
			}
		})
	}
}
