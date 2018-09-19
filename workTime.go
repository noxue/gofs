package main

import "sync"
import noApi "gofs/api"

type WorkTime struct {
	Lock     sync.Mutex
	WorkTime map[string]noApi.WorkTime
}

func (this *WorkTime) Get(name string) (noApi.WorkTime, bool) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	w, o := this.WorkTime[name]
	return w, o
}