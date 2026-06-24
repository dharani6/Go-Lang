package main

import (
	"fmt"
)

func main() {
	var test int = 25
	var test1 int = 26
	if test > test1 {
		fmt.Printf("test %v is greater than test1 %v", test, test1)
		fmt.Println("test is greater than test1")
		fmt.Print("test is greater than test1")
	} else {
		fmt.Printf("test1 %v is greater than test %v \n", test1, test)
		fmt.Println("test1 is greater than test")
		fmt.Print("test1 is greater than test")
	}
}
