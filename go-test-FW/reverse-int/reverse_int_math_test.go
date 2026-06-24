package reverseint

import "testing"

/*
Step 1 — your first Go test (same idea as pytest's single assert).

Rules:
  - File name must end with _test.go
  - Function name must start with Test
  - Receive *testing.T as the only parameter
  - Call t.Errorf(...) when something fails (test keeps running)
  - Run: go test -v <v flag for verbose output>
*/

func TestReverseIntegerBasic(t *testing.T) {
	result := ReverseInteger(12345)
	expected := 54321

	t.Log("input is 12345, expected output is 54321, got ", result)

	if result != expected {
		t.Errorf("ReverseInteger(12345) = %d, want %d", result, expected)
	}
}

/*
Step 2 — table-driven tests (next lesson; same pattern as pytest parametrize).

Uncomment after we walk through subtests together.
*/

// func TestReverseIntegerCases(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    int
// 		expected int
// 	}{
// 		{"positive number", 12345, 54321},
// 		{"single digit", 5, 5},
// 		{"zero", 0, 0},
// 		{"number with zeros", 1000, 1},
// 		{"large number", 123456789, 987654321},
// 		{"number ending in zero", 1200, 21},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := ReverseInteger(tt.input)
// 			if got != tt.expected {
// 				t.Errorf("ReverseInteger(%d) = %d, want %d", tt.input, got, tt.expected)
// 			}
// 		})
// 	}
// }
