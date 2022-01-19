package engine

type QueueScheduler struct {
	requestChan chan Request
	workChan    chan chan Request
}

func (s *QueueScheduler) SubmitTask(r Request) {
	s.requestChan <- r
}

func (s *QueueScheduler) Run() {
	s.requestChan = make(chan Request)
	s.workChan = make(chan chan Request)

	go func() {
		var requestQ []Request

		for {
			select {
			case t := <-s.requestChan:
				requestQ = append(requestQ, t)
			}
		}
	}()
}
