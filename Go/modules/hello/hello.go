package main

import (
	"example.com/greetings"
	"fmt"
	"log"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	names := []string{"jack", "tingting", "chenchen"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	for name, greeting := range messages {
		fmt.Printf("%v's greeting: %v\n", name, greeting)
	}
}
