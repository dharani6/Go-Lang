package main

import (
	"fmt"
)

func main() {

	for i := 1; i <= 5; i++ {
		if i == 4 {
			break
		}
		if i == 2 {
			continue
		}
		fmt.Print(i, " ")

	}
}
