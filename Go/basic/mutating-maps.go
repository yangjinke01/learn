package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["num"] = 24
	fmt.Println("The num is ", m["num"])

	m["num"] = 25
	fmt.Println("The num is ", m["num"])

	delete(m, "num")
	fmt.Println("The num is ", m["num"])

	v, ok := m["num"]
	fmt.Println("The num is ", v, "present? ", ok)
}
