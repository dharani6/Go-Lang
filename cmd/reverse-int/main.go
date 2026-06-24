package main

import (
	"fmt"
	"golearning/go-test-FW/reverse-int"
)

func main() {
	input := 112345
	fmt.Println("Input:", input)
	fmt.Println("Output:", reverseint.ReverseInteger(input))
}
