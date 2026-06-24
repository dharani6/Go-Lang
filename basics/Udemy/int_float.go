package main

import (
	"fmt"
)

func main() {
	var test int = 5
	test1 := float64(test)
	fmt.Printf("Output value of test1: is %v \n", test1)
	fmt.Printf("output value of test1: is %T \n", test1)

	var test11 int = 100
	var test12 int = -100
	fmt.Println(test11, test12)
	var test23 uint = 100
	fmt.Println(test23)

	var t float32 = 1.2345678
	var t_float float64 = 7.8954321
	fmt.Println(t, t_float)
	var str string = " hello there "
	var b bool
	fmt.Println(str, b)

	var testt string = "hello"
	var testt1 string = "there"

	fmt.Print(testt, " ", "\n")
	fmt.Print(testt1, " ", "\n")

	value1 := "vlaue1"
	value2 := "value2"

	fmt.Println(value1, value2)

	fmt.Printf("test has value: %v", value1)
	fmt.Printf("test has type: %T", value1)

	var testtt int = 10
	testttt1 := float64(testtt)

	fmt.Printf("value for test %v \n", testttt1)
	fmt.Printf("value for test %T \n", testttt1)
	var test123 float32 = 1.21

	fmt.Printf("value for test %v \n", test123)
	fmt.Printf("value for test %T \n", test123)

	test113 := int(test12)
	fmt.Printf("value for test %v \n", test113)
	fmt.Printf("value for test %T \n", test113)

	var ttest string = "helllo"
	var ttest1 int = len(ttest)

	fmt.Println(" ttest value", ttest, "ttest1 ", ttest1)

	var ttest2 string = " There "
	var ttest12 string = ttest + ttest2
	var ttest123 string = ttest + " " + ttest12
	fmt.Println(ttest123)

	iint := 100
	str_int := string(iint)
	fmt.Println("str_int", str_int)

}
