package main

import (
	"fmt"
	"math"
)

func sqrt(num float64) string {
	if num < 0 {
		return sqrt(-num) + "i"
	}
	return fmt.Sprint(math.Sqrt(num))
}

func main() {
	fmt.Println(sqrt(4), sqrt(-9))
}
