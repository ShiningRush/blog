package demo

import (
	"sync"
	"testing"
)

var (
	_globalMap = map[string]string{}
	_mux       = sync.Mutex{}
)

func mutexReadValue(key string) string {
	_mux.Lock()
	defer _mux.Unlock()
	return _globalMap[key]
}
func mutexWriteValue(key, value string) {
	_mux.Lock()
	defer _mux.Unlock()
	_globalMap[key] = value
}

func TestMutex(t *testing.T) {
	workers := 1000000
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			mutexReadValue("key")
			// do something
			mutexWriteValue("key", "value")
			wg.Done()
		}()
	}
	wg.Wait()
}
