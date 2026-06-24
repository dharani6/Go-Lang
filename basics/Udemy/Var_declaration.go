package main

import (
	"fmt"
)

func main() {
	var test int = 20
	fmt.Println(test)

	var test1 int
	test1 = 25

	fmt.Println(test1)

	ttvar := "hello there"
	fmt.Println(ttvar)

	var var_t, var_t1, var_t2 int = 1, 2, 3
	fmt.Println(var_t, var_t1, var_t2)

	var tt, tt1, tt2 = 1, "hello", 3
	fmt.Println(tt, tt1, tt2)

	var (
		x int    = 1
		y int    = 2
		z string = "hello"
	)
	fmt.Println(x, y, z)

	const testt int = 10
	fmt.Println(testt)
}
