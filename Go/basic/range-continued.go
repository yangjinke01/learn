package main

import "fmt"

func main() {
	pow := make([]int, 10)
	for i := range pow {
		pow[i] = 1 << i
	}

	for _, v := range pow {
		fmt.Println(v)
	}
}
