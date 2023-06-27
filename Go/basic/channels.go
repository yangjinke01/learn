package main

import "fmt"

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func main() {
	si := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	c := make(chan int)

	go sum(si[:len(si)/2], c)
	go sum(si[len(si)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}
