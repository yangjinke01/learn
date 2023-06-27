package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	printSlice(s)
	s = s[:0]
	printSlice(s)
	s = s[:4]
	printSlice(s)
	s = s[2:]
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len: %d cap: %d %v\n", len(s), cap(s), s)
}
