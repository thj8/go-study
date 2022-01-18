package main

import (
	"gowork/engine"
	"strconv"
	"time"
)

func main() {
	rs := gReqs()

	qs := engine.QueueScheduler{}
	qs.Run()

	for _, r := range rs {
		qs.SubmitTask(r)
	}

	time.Sleep(100 * time.Second)
}

func gReqs() []engine.Request {
	var reqs []engine.Request
	for i := 0; i < 99; i++ {
		reqs = append(reqs, engine.Request{Url: "www.baidu.com" + strconv.Itoa(i)})
	}

	return reqs
}
