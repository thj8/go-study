package engine

import "fmt"

type QueueScheduler struct {
	requestChan chan Request
}

func (s *QueueScheduler) SubmitTask(r Request) {
	s.requestChan <- r
}

func (s *QueueScheduler) Run() {
	s.requestChan = make(chan Request)

	go func() {
		for {
			select {
			case t := <-s.requestChan:
				fmt.Println(t.Url)
			}
		}
	}()
}
