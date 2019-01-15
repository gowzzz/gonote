package main

import (
	"fmt"
	"runtime"
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)
}

func main() {
	num := 100 * 100 * 20
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := NewWorkerPool(num)
	p.Run()
	datanum := 100 * 100 * 100 * 100
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}

}
