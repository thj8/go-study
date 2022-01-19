package engine

import (
	"context"
	"fmt"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	SubmitTask(Request)
	ConfigureWorkChan(chan Request)
}

func (c *ConcurrentEngine) Run(seeds []Request) {

	in := make(chan Request, 5)
	out := make(chan Respone)

	c.Scheduler.ConfigureWorkChan(in)

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < c.WorkerCount; i++ {
		fmt.Println("start")
		createWorker(in, out, ctx)
	}

	for _, r := range seeds {
		c.Scheduler.SubmitTask(r)
	}

	n := 1
	for {
		result := <-out
		fmt.Println(result.Back)
		n++
		if n == 99 {
			cancel()
			break
		}
	}

	//time.Sleep(time.Minute*1)
}

func createWorker(in chan Request, out chan Respone, ctx context.Context) {
	go func() {
		for {
			select {
			case request := <-in:
				result := doWork(request)
				out <- result
			case <-ctx.Done():
				fmt.Println("exit")
				return
			}
		}
	}()
}
