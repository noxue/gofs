package api

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
	this.lockTemplate.Lock()
	defer this.lockTemplate.Unlock()
	tpl, ok := this.templates[id]
	if !ok {
		return nil
	}
	return tpl
}
