/* My-Version
Psedocode == define or declar zero valued varaible result
loop through the given string and get each charater as ch
append the charater with reuslt in order to revers the string
	result  = ch + result
time complexity O(n) for a for loop its O(n)
*/

/*
	AI-version

Pseudocode:
- define a zero-valued string variable result
- loop through the string and get each character as ch
- prepend each character to result to reverse the string
- time complexity is O(n)
*/
package main // declaring the main package

import "fmt" // importing the formate for print and scan module

func main() { // main function
	result := ""    // zero valued string varaible // empty result string
	s := "Learning" // declaring a string varaible s and assinging //input string
	//  string vlaue which epxected to return
	for _, cha := range s { // looping syntax for _(dummpy value), // loop over runes in s
		// declarted varaible for itreation which in range of s declared vlaue

		result = string(cha) + result // print the reversed string
		// result value equating to charater of the string pre-append the result var
	} // closing paranthises of for loop

	fmt.Println(result) // printing the result
	fmt.Println("****************")

	fmt.Println("****************")
	fmt.Println(**** loop 2 ***********)
	//loop 2 : print each charaters
	for _, cha := range result {
		fmt.Println(string(cha))
	} // closing paranthises of for loop
	fmt.Println("****************")
} // closing paranthisees of main function
