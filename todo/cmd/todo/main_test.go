package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {

	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}

	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		// run the program with the path, and then an array of the
		// arguments consisting of tokens delimited by a space
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		// build the command
		cmd := exec.Command(cmdPath, "-add")
		// connects a pipe to the command's STDIN when the
		// command starts
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		// write task2 to command's STDIN and close
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("CompleteTasks", func(t *testing.T) {
		// complete the task
		cmd := exec.Command(cmdPath, "-complete", "2")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		// list the tasks
		lsCmd := exec.Command(cmdPath, "-list")
		out, err := lsCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: %s\nX 2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	task3 := "test task number 3"
	t.Run("DeleteTaskFromList", func(t *testing.T) {
		// add task3
		addCmd := exec.Command(cmdPath, "-add", task3)
		if err := addCmd.Run(); err != nil {
			t.Fatal(err)
		}
		// delete task2
		delCmd := exec.Command(cmdPath, "-delete", "2")
		if err := delCmd.Run(); err != nil {
			t.Fatal(err)
		}
		// Run list
		lsCmd := exec.Command(cmdPath, "-list")
		out, err := lsCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		// expect task and task3
		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task3)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
