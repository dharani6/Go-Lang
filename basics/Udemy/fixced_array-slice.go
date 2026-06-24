package main

import (
	"fmt"
)

func main() {
	test := [4]int{1, 2}
	fmt.Println(test)
	test[2] = 6
	fmt.Println(len(test), test, test[4])
	fmt.Errorf(err, test[3])
}
