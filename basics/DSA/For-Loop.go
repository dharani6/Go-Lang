package main

import "fmt"

func main() {
	N := 10
	sum := 0
	for i := 0; i <= N; i++ {
		sum = sum + i
		//fmt.Println("sum in 1st loop %v", sum)
	}
	//fmt.Println(sum)
	for i := 0; i <= N; i++ {
		sum = sum + i
	}
	//fmt.Println("Sum", sum)

	for j := 10; j >= 5; j-- {
		//fmt.Println(j)
	}

	for i := 0; i < 5; i++ {
		//fmt.Println(i)
	}

	i := 0
	for i < 5 {
		//fmt.Println(i)
		i++
	}
	ii := 0
	for ii > 10 {
		fmt.Println(ii)
		ii--
	}

	for {
		if i >= 5 {
			break
		}
		fmt.Println(i)
		i++
	}
}
