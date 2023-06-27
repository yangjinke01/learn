package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	lock sync.Mutex
	v    map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.lock.Lock()
	c.v[key]++
	c.lock.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.v[key]
}

func main() {
	c := &SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("someKey")
	}
	time.Sleep(time.Second)
	fmt.Println(c.Value("someKey"))
}
