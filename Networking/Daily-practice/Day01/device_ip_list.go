package main

import "fmt"

// main is the entry point for the Go program.
func main() {
	// devices holds a list of device IP addresses.
	devices := []string{
		"10.1.1.1",
		"10.1.1.2",
		"10.1.1.3",
	}

	// Loop over each device IP and print it to the console.
	for _, device := range devices {
		fmt.Println(device)
	}
}
