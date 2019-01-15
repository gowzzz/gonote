package main

import "fmt"

type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job
	WorkerQueue chan chan Job
	// Worker     []*Worker
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}
func (wp *WorkerPool) Run() {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker(wp.WorkerQueue)
		// wp.Worker = append(wp.Worker, &worker)
		worker.Run()
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job
			}
		}
	}()
}
