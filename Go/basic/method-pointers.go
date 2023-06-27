package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	if v == nil {
		fmt.Println("nil")
		return 0
	}
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X, v.Y = v.X*f, v.Y*f
}

func main() {
	var v *Vertex
	fmt.Printf("%v,%T\n", v, v)
	fmt.Println(v.Abs())
}
