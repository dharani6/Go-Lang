package main

import "fmt"

func fizzBuzz(n int) string {
	if n%15 == 0 {
		return "FizzBuzz"
	} else if n%3 == 0 {
		return "Fizz"
	} else if n%5 == 0 {
		return "Buzz"
	}
	return fmt.Sprintf("%d", n)

}

func main() {
	for i := 1; i <= 20; i++ {
		fmt.Println(fizzBuzz(i))

	}
}
