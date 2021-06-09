package multicas

import (
	"sync"
	"sync/atomic"
)

type MultipleCAS interface {
	Set(interface{}) bool
	Unset(interface{})
}

func NewMultipleCAS() MultipleCAS {
	return &multiplecas{
		pool: &sync.Pool{
			New: func() interface{} {
				return &casunit{}
			},
		},
	}
}

type casunit struct {
	v int32
}

func (c *casunit) set() bool {
	return atomic.CompareAndSwapInt32(&c.v, 0, 1)
}

func (c *casunit) unset() bool {
	return atomic.CompareAndSwapInt32(&c.v, 1, 0)
}

type casunitCounter struct {
	counter int64
	unit    *casunit
}

type multiplecas struct {
	inUse sync.Map
	pool  *sync.Pool
}

// Set returns true if a goroutine successfully set the value of key as 1.
//
// Then that goroutine is responsible for calling `Unset` to set the value of key back to 0.
func (m *multiplecas) Set(key interface{}) bool {
	cas := m.getCAS(key)
	ok := cas.unit.set()
	if ok {
		atomic.AddInt64(&cas.counter, 1)
	}
	return ok
}

func (m *multiplecas) Unset(key interface{}) {
	cas := m.getCAS(key)
	cas.unit.unset()
	m.putBack(key, cas)
}

func (m *multiplecas) getCAS(key interface{}) *casunitCounter {
	res, _ := m.inUse.LoadOrStore(key, &casunitCounter{
		counter: 0,
		unit:    m.pool.Get().(*casunit),
	})
	return res.(*casunitCounter)
}

func (m *multiplecas) putBack(key interface{}, c *casunitCounter) {
	atomic.AddInt64(&c.counter, -1)
	if c.counter <= 0 {
		m.pool.Put(c.unit)
		m.inUse.Delete(key)
	}
}
