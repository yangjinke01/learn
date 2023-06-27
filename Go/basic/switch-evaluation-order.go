package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("When is Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("Aftertomorrow.")
	default:
		fmt.Println("Too far away.")
	}
}
