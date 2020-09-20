package memo5

import (
	"fmt"
)

type Func func(key string, done chan<- struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res       result
	ready     chan struct{}
	cancelled chan struct{}
}

type request struct {
	key      string
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{}), cancelled: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go func() {
			select {
			case <-e.cancelled:
				delete(cache, req.key)
				req.response <- result{nil, fmt.Errorf("cancelled: %s", req.key)}
			case <-e.ready:
				req.response <- e.res
			}
		}()
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key, e.cancelled)
	close(e.ready)
}
