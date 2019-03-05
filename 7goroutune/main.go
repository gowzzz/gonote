package a

import (
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() {
	// fmt.Println("num:", s.Num)
	time.Sleep(10 * time.Millisecond)
}
func NewPool() *WorkerPool {
	workernum := 100 * 100 * 1
	jobnum := 100 * 100 * 10
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := NewWorkerPool(workernum, jobnum)
	p.Start()
	return p
}
func Runpool(p *WorkerPool, i int) {
	sc := &Score{Num: i}
	p.JobQueue <- sc
	// datanum := 100 * 100 * 100 * 100
	// go func() {
	// 	for i := 1; i <= datanum; i++ {
	// 		sc := &Score{Num: i}
	// 		p.JobQueue <- sc
	// 	}
	// }()

	// for {
	// 	fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
	// 	time.Sleep(2 * time.Second)
	// }
}
