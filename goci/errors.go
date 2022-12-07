package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("Validation failed")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

// unwrap the errror to see if an underlying error matches target
func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}
