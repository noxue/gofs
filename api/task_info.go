package api

import "github.com/golang/glog"

func (this *TaskInfo) GetSimId(taskId int) (id int) {
	this.lockSimTask.Lock()
	defer this.lockSimTask.Unlock()
	ids, ok := this.simTask[taskId]
	if !ok {
		return
	}
	for i, v := range ids {
		if v == taskId {
			id = i
			break
		}
	}
	return
}

func (this *TaskInfo) AddFreeSim(id int) {
	this.lockSimFree.Lock()
	defer this.lockSimFree.Unlock()
	this.simFree[id] = true
}

func (this *TaskInfo) RemoveFreeSim(id int) {
	this.lockSimFree.Lock()
	defer this.lockSimFree.Unlock()
	// if the sim exists, remove it
	if id != 0 {
		delete(this.simFree,id)
	}
}
func (this *TaskInfo) GetFreeSim() (id int) {
	this.lockSimFree.Lock()
	defer this.lockSimFree.Unlock()

	if len(this.simFree)==0 {
		return
	}

	for id,_=range this.simFree {
		break
	}



	return
}



func (this *TaskInfo)GetTaskIdBySimid(simId int) (id int) {
	this.lockSimTask.Lock()
	defer this.lockSimTask.Unlock()

	ids,ok:=this.simTask[simId]
	if !ok || len(ids)==0 {
		return
	}
	id=ids[0]

	// if the task exists, put the first number to end
	if _,ok:=TaskApi.app.tasks[id]; ok {
		this.simTask[simId]=append(this.simTask[simId],id)
	} else {
		id = 0
	}
	this.simTask[simId]=this.simTask[simId][1:]

	return
}

// set the sim free or not free, true means set the sim free
func (this *TaskInfo)SimFree(simId int,isFree bool){
	this.lockSimFree.Lock()
	defer this.lockSimFree.Unlock()

	if isFree {
		this.simFree[simId]=true
		glog.V(2).Infoln("sim is free，id:",simId)
	} else {
		glog.V(2).Infoln("sim is not free，id:",simId)
		delete(this.simFree,simId)
	}
}


