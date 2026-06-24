package main

import "fmt"

func reversestrings(items []string) {
	left := 0
	right := len(items) - 1

	for left < right {
		items[left], items[right] = items[right], items[left]
		left++
		right--
	}
}

func main() {
	devices := []string{"laptop", "phone", "tablet"}

	devices1 := []string{"router", "switch", "firewall", "server"}

	fmt.Println("original list:", devices)

	fmt.Println("original list:", devices1)

	reversestrings(devices)

	reversestrings(devices1)

	fmt.Println("reversed list:", devices)

	fmt.Println("reversed list:", devices1)

	// fmt.Println(reversestrings([]string{devices, devices1}...))

}
