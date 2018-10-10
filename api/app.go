package api

import (
	"time"
)

func (this *App) GetTask(id int) *Task {
	this.lockTask.Lock()
	defer this.lockTask.Unlock()
	task, ok := this.tasks[id]
	if !ok {
		return nil
	}
	return task
}

func (this *App) GetGateway(id int) *Gateway {
	this.lockGateway.Lock()
	defer this.lockGateway.Unlock()
	gateway, ok := this.gateways[id]
	if !ok {
		return nil
	}
	return gateway
}

func (this *App) GetSim(id int) *Sim {
	this.lockSim.Lock()
	defer this.lockSim.Unlock()
	sim, ok := this.sims[id]
	if !ok {
		return nil
	}
	return sim
}

func (this *App) GetTpl(id int) *Template {

	n := 0
	this.lockTemplate.Lock()
	tpl, ok := this.templates[id]
	this.lockTemplate.Unlock()

	if !ok {
		TaskApi.UpdateTpl(id)

	ReTry:
		time.Sleep(time.Second)
		tpl, ok = this.templates[id]
		n++
		if !ok{
			if n<3{
				goto ReTry
			}
			return nil
		}
	}
	return tpl
}
