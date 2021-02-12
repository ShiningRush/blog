package demo

import (
	"sync"
	"testing"
)

var (
	_rwglobalMap = map[string]string{}
	_rwmux       = sync.RWMutex{}
)

func rwmutexReadValue(key string) string {
	_rwmux.RLock()
	defer _rwmux.RUnlock()
	return _rwglobalMap[key]
}
func rwmutexWriteValue(key, value string) {
	_rwmux.Lock()
	defer _rwmux.Unlock()
	_rwglobalMap[key] = value
}

func TestRwMutex(t *testing.T) {
	workers := 1000000
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			rwmutexReadValue("key")
			// do something
			rwmutexWriteValue("key", "value")
			wg.Done()
		}()
	}
	wg.Wait()
}
