package syncx

import "sync"

type (
	SCFun func() interface{}

	ShareCall interface {
		Do(key string, fn SCFun) interface{}
	}

	call struct {
		val interface{}
		wg  sync.WaitGroup
	}

	shareGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

func NewShareCall() ShareCall {
	return &shareGroup{
		calls: make(map[string]*call),
	}
}

func (s *shareGroup) Do(key string, fn SCFun) interface{} {
	c, done := s.createCall(key, fn)
	if done {
		return c.val
	}

	s.makeCall(c, fn, key)
	return c.val
}

func (s *shareGroup) createCall(key string, fn SCFun) (*call, bool) {
	s.lock.Lock()
	if c, ok := s.calls[key]; ok {
		s.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c := new(call)
	s.calls[key] = c
	c.wg.Add(1)
	s.lock.Unlock()
	return c, false
}

func (s *shareGroup) makeCall(c *call, fn SCFun, key string) {
	defer func() {
		s.lock.Lock()
		delete(s.calls, key)
		s.lock.Unlock()
		c.wg.Done()
	}()

	c.val = fn()
}
