package main

import "fmt"

func main() {
	names := [4]string{
		"John",
		"Paul",
		"Jack",
		"Rose",
	}
	fmt.Println(names)
	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)
	a[1] = "Alice"
	fmt.Println(a, b)
}
