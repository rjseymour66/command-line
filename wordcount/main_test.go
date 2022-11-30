package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")

	exp := 4

	res, _ := count(b, false, false)

	if res != exp {
		t.Errorf("Expected %d, got %d instead", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\nline2\nline3 word1")

	exp := 3

	res, _ := count(b, true, false)

	if res != exp {
		t.Errorf("Expected %d, got %d instead", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("1234567890")

	exp := 10

	_, res := count(b, false, true)

	if res != exp {
		t.Errorf("Expected %d, got %d instead", exp, res)
	}
}

func TestCountWordsAndBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")

	expWords := 4
	expBytes := 2

	resWords, resBytes := count(b, false, true)

	if resWords != expWords {
		t.Errorf("Expected %d, got %d instead", expWords, resWords)
	}

	if resBytes != expBytes {
		t.Errorf("Expected %d, got %d instead", expBytes, resBytes)
	}
}

func TestCountLinesAndBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\nline2\nline3 word1")

	expLines := 3
	expBytes := 39

	resLines, resBytes := count(b, true, true)

	if resLines != expLines {
		t.Errorf("Expected %d, got %d instead", expLines, resLines)
	}

	if resBytes != expBytes {
		t.Errorf("Expected %d, got %d instead", expBytes, resBytes)
	}
}
