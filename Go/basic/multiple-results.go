package main

import "fmt"

func swap(a, b string) (string, string) {
	return b, a
}

func main() {
	fmt.Println(swap("jack", "hello"))
}
