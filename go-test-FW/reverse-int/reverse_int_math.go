package reverseint

//package main

/*
Reverse an integer without converting it to a string.
Example: 12345 -> 54321

Pseudocode:
  1. reversed = 0
  2. while n > 0:
       lastDigit = n % 10
       reversed = reversed*10 + lastDigit
       n = n / 10
  3. return reversed

Concepts:
  - function: reverseInteger encapsulates the logic
  - loop: for n > 0 { ... } runs once per digit
  - variables: reversed, lastDigit, n change each iteration

Time:  O(d)  — d = number of digits
Space: O(1)  — constant extra memory
*/

func ReverseInteger(n int) int {
	reversed := 0
	for n > 0 {
		lastDigit := n % 10
		reversed = reversed*10 + lastDigit
		n = n / 10
	}
	return reversed
}

// Demo main — run with: go run .
/*
func main() {
	input := 12345
	result := ReverseInteger(input)
	fmt.Printf("The reversed integer of %d is %d\n", input, result)
}
*/
