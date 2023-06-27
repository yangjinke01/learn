package main

import "fmt"

func main() {
	prime := [...]int{1, 2, 3, 4, 5, 6}
	var s []int = prime[1:4]
	fmt.Println(s)
}
