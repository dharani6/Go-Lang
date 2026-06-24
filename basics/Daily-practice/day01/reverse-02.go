package main

import "fmt"

func reversestrings(num []int) {

	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {

		num[i], num[j] = num[j], num[i]
	}
}

func main() {

	numbers := []int{10, 20, 30, 40, 50}

	fmt.Println("original list:", numbers)
	reversestrings(numbers)
	fmt.Println("reversed list:", numbers)
}
