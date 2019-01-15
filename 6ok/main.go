package main

import (
	"fmt"
	"sync"
	"time"
)

var EventConfLock = new(sync.RWMutex)

func trans(in interface{}) {
	EventConfLock.Lock()
	var str string
	if v, ok := in.(string); ok {
		str = v
		fmt.Println("str:", str)
	} else {
		fmt.Println("not string")
	}
	EventConfLock.Unlock()

}

var m = make(map[int]int)
var wg sync.WaitGroup

func test() {
	EventConfLock.Lock()
	m[1] = 1
	EventConfLock.Unlock()
	// fmt.Println("ok")
	// wg.Done()
}
func main() {
	for i := 0; i < 100; i++ {
		go test()
		go test()
	}
	time.Sleep(5 * time.Second)
}

/*
v, ok = m[key] // map lookup
v, ok = x.(T) // type assertion
v, ok = <-ch // channel receive
*/
