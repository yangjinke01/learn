package main

import "fmt"

func main() {
	var i interface{}
	i = 42
	fmt.Printf("%v,%T\n", i, i)
	i = "jack"
	fmt.Printf("%v,%T\n", i, i)
}
