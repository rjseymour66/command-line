package pomodoro_test

import (
	"pomo/pomodoro"
	"pomo/pomodoro/repository"
	"testing"
)

// Test .Helper() to get instance of the repo
func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
