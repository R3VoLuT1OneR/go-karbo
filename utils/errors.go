package utils

import "fmt"

type assertionError struct {
	msg string
}

func AssertionError(msg string) error {
	return &assertionError{msg}
}

func (a assertionError) Error() string {
	return fmt.Sprintf("assertion failed: %s", a.msg)
}
