package main

import (
	"fmt"
	"math/cmplx"
)

var (
	toBe   bool       = false
	maxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
	bt     byte       = 65
	char   rune       = 'B'
)

func main() {
	fmt.Printf("Type: %T Value: %v\n", toBe, toBe)
	fmt.Printf("Type: %T Value: %v\n", maxInt, maxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)
	fmt.Printf("Type: %T Value: %v\n", bt, bt)
	fmt.Printf("Type: %T Value: %v String %s\n", bt, bt, string(bt))
	fmt.Printf("Type: %T Value: %v int: %d\n", char, char, int(char))
}
