package main

import "fmt"

func add1(a int, b int) int {
	return a + b
}

func main() {
	result := add1(10, 20)
	fmt.Printf("sum of %d and %d given as input is = %d", 10, 20, result)
}
