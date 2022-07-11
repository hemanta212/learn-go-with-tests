package main

import "fmt"

// Add takes two integers and returns their sum.
func Add(x, y int) int {
	return x + y
}

func main() {
	fmt.Println(Add(2, 2))
}
