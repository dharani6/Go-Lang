package main

import "fmt"

/* Reverse an integer without converting it to a string.
   For example, given the integer 12345, the function should return 54321.*/

/* pseducode
1. Initialize a variable to store the reversed integer.
2. While the input integer is not zero:
   a. Get the last digit of the input integer.
   b. Append the last digit to the reversed integer.
   c. Remove the last digit from the input integer.
3. Return the reversed integer.
*/

func reverseInteger(n int) int {
	reversed := 0
	if n < 0 {
		n = -n
	}

	for n > 0 {
		LastDigit := n % 10
		reversed = reversed*10 + LastDigit
		n = n / 10
		fmt.Printf("LastDigit is %d, reversed is %d, n is %d \n", LastDigit, reversed, n)
	}

	return reversed
}

func main() {
	input := 12345
	result := reverseInteger(input)

	fmt.Printf("The revsered integer of %d is %d \n", input, result)

	Neg_input := -3123
	Neg_result := reverseInteger(Neg_input)
	fmt.Printf("The revsered integer of %d is %d \n", Neg_input, Neg_result)
}
