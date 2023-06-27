package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, v := range strings.Split(s, " ") {
		if _, ok := m[v]; ok {
			m[v] = m[v] + 1
		} else {
			m[v] = 1
		}
	}
	return m
}

func main() {
	fmt.Println(WordCount("big big girl, big big word"))
}
