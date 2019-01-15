package main

type Worker struct {
	WorkerQueue chan chan Job
	Jobchan     chan Job
}

func NewWorker(wc chan chan Job) Worker {
	return Worker{WorkerQueue: wc, Jobchan: make(chan Job)}
}
func (w Worker) Run() {
	go func() {
		for {
			w.WorkerQueue <- w.Jobchan
			select {
			case job := <-w.Jobchan:
				job.Do()
			}
		}
	}()
}
