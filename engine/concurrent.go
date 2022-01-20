package engine

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
	out := make(chan bool, c.WorkerCount)

	c.Scheduler.ConfigureWorkChan(in)

	for i := 0; i < c.WorkerCount; i++ {
		createWorker(in, out)
	}

	go c.submit(in, seeds)

	for i := 0; i < c.WorkerCount; i++ {
		<-out
	}
}

func (c *ConcurrentEngine) submit(in chan Request, s []Request) {
	for _, r := range s {
		c.Scheduler.SubmitTask(r)
	}
	close(in)
}

func createWorker(in chan Request, out chan bool) {
	go func() {
		for r := range in {
			doWork(r)
		}

		out <- true
	}()
}
