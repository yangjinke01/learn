package main

import "fmt"

func fibnacci() func() int {
	x, y := 0, 1
	return func() int {
		x, y = y, x+y
		return y - x
	}
}

func main() {
	fib := fibnacci()
	for i := 0; i < 10; i++ {
		fmt.Println(fib())
	}
}
