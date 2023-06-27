package main

import "fmt"

func main() {
	num := 0
	for i := 0; i <= 10; i++ {
		num += i
	}
	fmt.Println(num)
}
