package main

import "fmt"

func Reversed[T any](items []T) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}

func main() {

	names := []string{"Alice", "Bob", "Charlie", "David"}
	ports := []int{80, 443, 22, 8080}
	status := []bool{true, false, true, false}

	fmt.Printf("Original list: %v\n", names)
	fmt.Printf("Original list: %v\n", ports) 
	fmt.Printf("Original list: %v\n", status)

	fmt.Print("Helloworld\n")
	Reversed(names)
	Reversed(ports)
	Reversed(status)

	fmt.Printf("Reversed list: %v\n", names)
	fmt.Printf("Reversed list: %v\n", ports)
	fmt.Printf("Reversed list: %v\n", status)
}
