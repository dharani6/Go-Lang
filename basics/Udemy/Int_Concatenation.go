package main

import (
	"fmt"
)

func main() {
	var test int = 2
	var test1 int = 2
	finalvlaue := test + test1
	fmt.Println(finalvlaue)

	// Arthmetic operations
	var test12 int = 1
	var test123 int = 2
	final := test12 + test123
	fmt.Println(final)
	fmt.Println(test12)
	test12 += 1
	fmt.Println(test12)
	test12 *= 3
	fmt.Println(test12)
	fmt.Println(test12 < 8 && test123 > 20)
	fmt.Println(test12 < 8 || test123 > 20)
	fmt.Println(test12 & test123)
	fmt.Println(test12 | test123)
	fmt.Println(test12 << 2)
	fmt.Println(test123 >> 2)

}
