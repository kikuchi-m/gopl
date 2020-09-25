package memo5

import (
	"fmt"
)

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	response chan<- result
	done     <-chan struct{}
}

type Memo struct {
	requests         chan request
	cache            map[string]*entry
	newEntryRequest  chan string
	newEntryResponse chan newEntry
}

func New(f Func) *Memo {
	memo := &Memo{make(chan request), make(map[string]*entry), make(chan string), make(chan newEntry)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response, done}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	// TODO: prevent writing closed channels
	close(memo.requests)
	close(memo.newEntryRequest)
	close(memo.newEntryResponse)
}

type newEntry struct {
	e             *entry
	alreadyExists bool
}

func (memo *Memo) handleNewEntry(f Func) {
	for k := range memo.newEntryRequest {
		e := memo.cache[k]
		absent := e == nil
		if absent {
			e = &entry{ready: make(chan struct{})}
			memo.cache[k] = e
			go e.call(f, k)
		}
		memo.newEntryResponse <- newEntry{e, !absent}
	}
}

func (memo *Memo) server(f Func) {
	go memo.handleNewEntry(f)
	for req := range memo.requests {
		memo.newEntryRequest <- req.key
		newEntRes := <-memo.newEntryResponse
		if newEntRes.alreadyExists {
			go func(r request, e *entry) {
				select {
				case <-r.done:
					close(r.response)
				case <-e.ready:
					r.response <- e.res
				}
			}(req, newEntRes.e)
		} else {
			// exactly new
			go func(r request, e *entry) {
				select {
				case <-r.done:
					// delete(cache, r.key)
					r.response <- result{nil, fmt.Errorf("cancelled: %s", r.key)}
				case <-e.ready:
					req.response <- e.res
				}
			}(req, newEntRes.e)
		}
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}
