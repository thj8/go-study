package engine

type SimpleScheduler struct {
	workChan chan Request
}

func (s *SimpleScheduler) SubmitTask(r Request) {
	go func() { s.workChan <- r }()
}

func (s *SimpleScheduler) ConfigureWorkChan(c chan Request) {
	s.workChan = c
}
