package main

import "fmt"

func main() {
	// explicit type + value
	var x int = 10
	// inferred type 
	var y = 20
	// short declaration
	z := 30
	// Zero-value declaration
	var a int 
	var s stringvar 
	b bool
	//2 Grouped declarations
	var ( 
		i int = 1
		j = 2
	)

	// constant
	const pi = 3.14
	const MaxUsers = 100
	
	// Custom types 

	type Age int
	type Name string
	type Person struct {
		Name string
		Age int
	}
	
	name := "test" // short decl
	items := []int{1, 2, 3}
	d := map[string]int{"a": 1, "b": 2, "c": 3}
	if x > 5 {
		fmt.Println(name)
	}
	for _, v := range items {
		fmt.Println(v)
	}

	for k, v := range d {
		fmt.Println(k, v)
	}
}
