package main

import "fmt"

func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

func main() {
	si := []int{1, 2, 3, 4, 5}
	fmt.Println(Index(si, 3))

	ss := []string{"a", "b", "c", "d", "e"}
	fmt.Println(Index(ss, "d"))
}
