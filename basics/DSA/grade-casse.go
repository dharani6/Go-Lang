package main

import "fmt"

func grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

func main() {
	var array = []int{99, 77, 88, 44, 55, 66}
	result := []int{}

	for i := 1; i < len(array); i++ {
		//result = (grade(array[i]))
		result = append(result, array[i])
		fmt.Printf("%v \n", result)
	}
	fmt.Println("Master", result)
}

/*	for _, s := range array {
			fmt.Printf("Score %d -> Grade %s\n", s, grade(s))
	}
} */
