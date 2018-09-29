package api

import (
	"gofs/config"
	"github.com/golang/glog"
	"fmt"
	"strconv"
)


func (this *Api) Auth() {
	content := fmt.Sprintf(`{"action":"auth","id":%s,"key":"%s"}`, config.Config.Api.AppId, config.Config.Api.Key)
	this.sendData(content)
}

func (this *Api) Close(){
	this.ws.Close()
}

func (this *Api) sendData(str string)  {
	//this.lockWs.Lock()
	//defer this.lockWs.Unlock()
	_, err := this.ws.Write([]byte(str))
	if err != nil {
		glog.Fatal(err)
	}
}

func (this *Api) update(what, content string) {
	str := fmt.Sprintf(`{"action":"%s","content":"%s"}`, what, content)
	this.sendData(str)
}

func (this *Api) UpdateSipThread(user string) {
	this.update("sip_thread", user)
}

func (this *Api) UpdateWorkTime(user string) {
	this.update("worktime", user)
}

func (this *Api) UpdateTpl(id int) {
	this.update("tpl", strconv.Itoa(id))
}

func (this *Api) UpdateSim(id int) {
	this.update("sim", strconv.Itoa(id))
}

// request task by sim , the id is sim id
func (this *Api) SimTasks(id int) {
	this.update("sim_tasks", strconv.Itoa(id))
}

// request task by sip , the id is sip id
func (this *Api) SipTasks(id int) {
	this.update("sip_tasks", strconv.Itoa(id))
}

func (this *Api) SimTaskUser(simId,taskId int) {
	glog.V(2).Infoln("get sim task_user,sim id",simId,"task id",taskId)
	str:=fmt.Sprintf(`{"action":"sim_task_user","task_id":%d, "type":"sim","id":%d}`,taskId,simId)
	this.sendData(str)
}

func (this *Api) SipTaskUser(sipId,taskId int) {
	this.app = this.app
	str:=fmt.Sprintf(`{"action":"sip_task_user","task_id":%d, "type":"sip","id":%d}`,taskId,sipId)
	this.sendData(str)
}

func (this *Api) Tasks() map[int]*Task {
	this.app.lockTask.Lock()
	defer this.app.lockTask.Unlock()
	return this.app.tasks
}

func (this *Api)GetTaskInfo()*TaskInfo{
	this.app.lockTaskInfo.Lock()
	defer this.app.lockTaskInfo.Unlock()
	return this.app.taskInfo
}

//
//func (this *Api) GetGateways() (gateways []Gateway) {
//	this.app.lockGateway.Lock()
//	defer this.app.lockGateway.Unlock()
//
//	for _, g := range this.app.Gateways {
//		gateways = append(gateways, g)
//	}
//	return
//}
//
//// 根据id获取单个
//func (this *Api) GetGateway(gatewayId int) (gateway Gateway, ok bool) {
//	this.app.lockGateway.Lock()
//	gateway, ok = this.app.Gateways[gatewayId]
//	this.app.lockGateway.Unlock()
//
//	if !ok {
//		this.getaway()
//		this.app.lockGateway.Lock()
//		gateway, ok = this.app.Gateways[gatewayId]
//		this.app.lockGateway.Unlock()
//	}
//	return
//}
//
//// 根据网关获取电话卡列表
//func (this *Api) GetSims(gatewayId int) (sims []Sim) {
//	this.app.lockSimFree.Lock()
//	defer this.app.lockSimFree.Unlock()
//
//	for _, s := range this.app.Sims {
//		sims = append(sims, s)
//	}
//	return
//}
//
//func (this *Api) GetAllSims() (sims []Sim) {
//	this.app.lockSimFree.Lock()
//	defer this.app.lockSimFree.Unlock()
//	for _, s := range this.app.Sims {
//		sims = append(sims, s)
//	}
//	return
//}
//
//// 根据卡编号获取卡信息
//func (this *Api) GetSim(simId int) (sim Sim, ok bool) {
//	this.app.lockSimFree.Lock()
//	defer this.app.lockSimFree.Unlock()
//	sim, ok = this.app.Sims[simId]
//	return
//}
//
//// 根据卡编号获取任务列表
//func (this *Api) GetTasks(simId int) (tasks []Task) {
//	this.app.lockTask.Lock()
//	defer this.app.lockTask.Unlock()
//	tasks = this.app.Tasks[simId]
//	return
//}
//
//// 根据任务编号获取客户信息列表
//func (this *Api) GetTaskUsers(taskId int) (users []TaskUser, err error, over bool) {
//	req, err := http.NewRequest("GET", fmt.Sprintf(this.ApiUrl+"/task/%d/users", taskId), nil)
//	if err != nil {
//		return
//	}
//
//	res, err := this.Client.Do(req)
//	if err != nil {
//		return
//	}
//
//	ret, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return
//	}
//
//	glog.V(2).Infoln(string(ret))
//
//	var _users _TaskUser
//
//	err = json.Unmarshal(ret, &_users)
//	if err != nil {
//		return
//	}
//
//	// 如果当前任务客户信息为空，就删除当前任务
//	if _users.Meta.Code != 0 {
//		func() {
//			this.app.lockTask.Lock()
//			defer this.app.lockTask.Unlock()
//			delete(this.app.Tasks, taskId)
//			err = errors.New("当前任务已经没有更多客户了")
//			over = true
//		}()
//		return
//	}
//
//	func() {
//		for _, s := range _users.Data["users"] {
//			users = append(users, s)
//		}
//	}()
//
//	return
//}
//
//func (this *Api) GetSchedule(userId string) (workTime WorkTime, err error) {
//	str, err := this.getUserConfig(userId, "schedule")
//	if err != nil {
//		return
//	}
//
//	var _workTime _WorkTime
//	err = json.Unmarshal([]byte(str), &_workTime)
//	if err != nil {
//		return
//	}
//
//	if !_workTime.Meta.Success {
//		err = errors.New("用户没有配置工作时间")
//		return
//	}
//
//	err = json.Unmarshal([]byte(_workTime.Data["config"]["value"]), &workTime)
//	if err != nil {
//		return
//	}
//
//	return
//}
