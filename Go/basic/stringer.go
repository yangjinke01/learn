package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v, %v years.", p.Name, p.Age)
}

func main() {
	jack := Person{"jack", 28}
	tingting := Person{"tingting", 26}
	fmt.Println(jack, tingting)
}
