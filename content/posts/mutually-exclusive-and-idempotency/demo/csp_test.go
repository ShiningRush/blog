package demo

import (
	"sync"
	"testing"
)

var (
	_defaultGuard = NewGuard()
	_cspGlobalMap = map[string]string{}
)

type ModifyCommand struct {
	key   string
	value string
}

type ReadCommand struct {
	key        string
	receiveCnl chan string
}

type Guard struct {
	modifyChl chan *ModifyCommand
	readChl   chan *ReadCommand
}

func NewGuard() *Guard {
	g := &Guard{}
	g.modifyChl = make(chan *ModifyCommand)
	g.readChl = make(chan *ReadCommand)
	go func() {
		for {
			select {
			case mc := <-g.modifyChl:
				_cspGlobalMap[mc.key] = mc.value
			case rc := <-g.readChl:
				rc.receiveCnl <- _cspGlobalMap[rc.key]
			}
		}
	}()
	return g
}
func (g *Guard) Get(key string) string {
	rnl := make(chan string)
	g.readChl <- &ReadCommand{
		key:        key,
		receiveCnl: rnl,
	}
	rs := ""
	for rs = range rnl {
		close(rnl)
	}
	return rs
}
func (g *Guard) Set(key, value string) {
	g.modifyChl <- &ModifyCommand{
		key:   key,
		value: value,
	}
}

func cspReadValue(key string) string {
	return _defaultGuard.Get(key)
}
func cspWriteValue(key, value string) {
	_defaultGuard.Set(key, value)
}

func TestCsp(t *testing.T) {
	workers := 1000000
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			cspReadValue("key")
			// do something
			cspWriteValue("key", "value")
			wg.Done()
		}()
	}
	wg.Wait()
}
