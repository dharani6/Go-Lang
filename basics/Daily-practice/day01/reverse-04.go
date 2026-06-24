// 4. Return a reversed copy of a list without modifying the original list.
package main

import "fmt"

func ReversedCopy[T any](items []T) []T {

	result := make([]T, len(items))

	for i := 0; i < len(items); i++ {
		result[i] = items[len(items)-1-i]
	}

	return result
}

func main() {
	ips := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"}

	reveresed := ReversedCopy(ips)

	fmt.Println("Original list:", ips)
	fmt.Println("Reversed copy:", reveresed)
}
