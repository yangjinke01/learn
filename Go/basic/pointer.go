package main

import "fmt"

func main() {
	i, j := 27, 2701
	p := &i
	fmt.Println(*p)
	*p = 42
	fmt.Println(i)

	p = &j
	*p = *p / 37
	fmt.Println(j)
}
