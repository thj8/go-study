package syncx

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-playground/assert"
)

func TestSingleCall(t *testing.T) {
	fun := func() interface{} {
		return 123
	}

	ns := NewShareCall()
	ret := ns.Do("thj", fun)

	assert.Equal(t, ret, 123)
}

func TestShareCall(t *testing.T) {
	var c int32
	fun := func() interface{} {
		atomic.AddInt32(&c, 1)
		time.Sleep(time.Millisecond * 10)
		return 123
	}

	sc := NewShareCall()
	for j := 0; j < 100; j++ {
		go func() {
			sc.Do("thj", fun)
		}()
	}

	time.Sleep(time.Millisecond * 100)
	get := atomic.LoadInt32(&c)
	fmt.Println(get)
	assert.Equal(t, get, int32(1))
}
