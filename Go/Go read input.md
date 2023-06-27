Go read input tutorial shows how to read input from a user. Standard input, often abbreviated *stdin*, is a stream from
which a program reads its input data.

To read input from users in Go, we use the `fmt`, `bufio`, and `os` packages.

## Go read input with Scanf

The `Scanf` function scans text read from standard input, storing successive space-separated values into successive
arguments as determined by the format. It returns the number of items successfully scanned.

read_input.go

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	var name string
	rand.Seed(time.Now().UnixNano())
	fmt.Print("Enter your name: ")
	fmt.Scanf("%s", &name)
	fmt.Println("Hello", name)
}
```

read_input2.go

```go
package main

import "fmt"

func main() {

	var name string
	var age int

	fmt.Print("Enter your name & age: ")
	fmt.Scanf("%s %d", &name, &age)
	fmt.Printf("%s is %d years old\n", name, age)
}
```

## Go read input with NewReader

The `bufio` package implements buffered I/O. Buffered I/O has much better performance than non-buffered. The package
wraps an `io.Reader` or `io.Writer` object, creating another object (Reader or Writer) that also implements the
interface but provides buffering and some help for textual I/O.

read_input3.go

```go
package main

import (
	"os"
	"bufio"
	"fmt"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")

	name, _ := reader.ReadString('\n')
	fmt.Printf("Hello %s\n", name)
}
```

## Go read input with NewScanner

The `Scanner` provides a convenient interface for reading data such as a file of newline-delimited lines of text.

read_input4.go

```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	names := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter name: ")

		scanner.Scan()

		text := scanner.Text()

		if len(text) != 0 {

			fmt.Println(text)
			names = append(names, text)
		} else {
			break
		}
	}

	fmt.Println(names)
}
```
