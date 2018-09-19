package phone

import "sync"

type CallList struct {
	calls map[string]*Call
	lock sync.Mutex
}

func NewCallList()(*CallList) {
	return &CallList{
		calls: make(map[string]*Call),
	}
}

func (this *CallList) Set(key string, call *Call) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.calls[key] = call
}

func (this *CallList) Get(key string) (call *Call,ok bool) {
	if len(key)==0 {
		return
	}
	call,ok = this.calls[key]
	return
}
