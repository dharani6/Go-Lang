package main

import "testing"

func TestNumbers(t *testing.T) {
	if CompareNumbers(6, 5) {
		t.Errorf("Expectation 6 > 5 to be true")
	}
	if CheckEquality(6, 6) {
		t.Errorf("Expectation 6 == 6 to be true")
	}
	if CheckNotEqual(5, 5) {
		t.Errorf("Expectation 5 != 5 to be true")
	}
}
