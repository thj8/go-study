package main

import (
	"gowork/engine"
	"strconv"
)

func main() {
	rs := gReqs()

	e := engine.ConcurrentEngine{
		Scheduler:   &engine.SimpleScheduler{},
		WorkerCount: 5,
	}

	e.Run(rs)
}

func gReqs() []engine.Request {
	var reqs []engine.Request
	for i := 0; i < 99; i++ {
		reqs = append(reqs, engine.Request{Url: "www.baidu.com" + strconv.Itoa(i)})
	}

	return reqs
}
