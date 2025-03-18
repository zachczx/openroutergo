// Package assert provides a set of functions for asserting the results of tests.
package assert

import "testing"

// Equal asserts that two values are equal.
func Equal[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v to equal %v", expected, actual)
	}
}

// NotEqual asserts that two values are not equal.
func NotEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		t.Errorf("expected %v to not equal %v", expected, actual)
	}
}

// Nil asserts that a value is nil.
func Nil(t *testing.T, actual any) {
	t.Helper()
	if actual != nil {
		t.Errorf("expected %v to be nil", actual)
	}
}

// NotNil asserts that a value is not nil.
func NotNil(t *testing.T, actual any) {
	t.Helper()
	if actual == nil {
		t.Errorf("expected %v to not be nil", actual)
	}
}

// Error asserts that an error is not nil.
func Error(t *testing.T, expected error, actual error) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v to equal %v", expected, actual)
	}
}

// NoError asserts that an error is nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// True asserts that a boolean is true.
func True(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Errorf("expected true, got false")
	}
}

// False asserts that a boolean is false.
func False(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Errorf("expected false, got true")
	}
}
