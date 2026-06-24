package common

import (
	"fmt"
	"strings"
)

// Assert provides fluent assertion helpers.
type Assert struct {
	steps *[]TestStep
}

// NewAssert creates an Assert that records steps into the provided slice.
func NewAssert(steps *[]TestStep) *Assert {
	return &Assert{steps: steps}
}

func (a *Assert) record(name string, status TestStatus, detail string) {
	if a.steps != nil {
		*a.steps = append(*a.steps, TestStep{
			Name:   name,
			Status: status,
			Detail: detail,
		})
	}
}

// Equal asserts that got == expected.
func (a *Assert) Equal(name string, expected, got interface{}) error {
	if fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", got) {
		a.record(name, StatusPass, fmt.Sprintf("expected=%v ✓", expected))
		return nil
	}
	detail := fmt.Sprintf("expected=%v got=%v", expected, got)
	a.record(name, StatusFail, detail)
	return &ErrTestFailed{Assertion: name, Expected: expected, Got: got}
}

// NotEmpty asserts that s is not an empty string.
func (a *Assert) NotEmpty(name, s string) error {
	if strings.TrimSpace(s) != "" {
		a.record(name, StatusPass, "value is non-empty ✓")
		return nil
	}
	a.record(name, StatusFail, "value is empty")
	return &ErrTestFailed{Assertion: name, Expected: "non-empty string", Got: s}
}

// Contains asserts that haystack contains needle.
func (a *Assert) Contains(name, haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		a.record(name, StatusPass, fmt.Sprintf("contains %q ✓", needle))
		return nil
	}
	a.record(name, StatusFail, fmt.Sprintf("%q not found in response", needle))
	return &ErrTestFailed{Assertion: name, Expected: needle, Got: haystack}
}

// True asserts that cond is true.
func (a *Assert) True(name string, cond bool, msg string) error {
	if cond {
		a.record(name, StatusPass, msg+" ✓")
		return nil
	}
	a.record(name, StatusFail, msg)
	return &ErrTestFailed{Assertion: name, Expected: true, Got: false}
}

// StatusCode asserts an HTTP status code.
func (a *Assert) StatusCode(got, expected int) error {
	name := fmt.Sprintf("HTTP status %d", expected)
	if got == expected {
		a.record(name, StatusPass, fmt.Sprintf("got %d ✓", got))
		return nil
	}
	a.record(name, StatusFail, fmt.Sprintf("expected %d got %d", expected, got))
	return &ErrTestFailed{Assertion: "HTTP StatusCode", Expected: expected, Got: got}
}
