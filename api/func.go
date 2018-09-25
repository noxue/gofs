package api

import (
	"encoding/json"
	"github.com/golang/glog"
	"log"
	"golang.org/x/net/websocket"
	"time"
	"strconv"
	"fmt"
)

func New(apiUrl, origin, appid, key string) *Api {

	if appid == "" || key == "" {
		glog.Error("robot id or key empty")
		return nil
	}

	api := &Api{
		appId:  appid,
		key:    key,
		apiUrl: apiUrl,
		origin: origin,
		app: *&App{
			gateways:  make(map[int]*Gateway),
			sims:      make(map[int]*Sim),
			templates: make(map[int]*Template),
			tasks:     make(map[int]*Task),
			users:     make(map[string]*User),
			taskInfo: &TaskInfo{
				TaskSim:make(map[int][]int),
				TaskSip:make(map[int][]int),
			},
		},
	}

	var err error
	api.ws, err = websocket.Dial(apiUrl, "", origin)
	if err != nil {
		log.Panic(err)
	}
	glog.Info("connect server success")

	return api
}

func (this *Api) Handle() {

	var data [1024000]byte
	var len = 0
	var msg [4096]byte
	for {
		n, err := this.GetWs().Read(msg[:])

		if n == 0 {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		copy(data[len:], msg[:n])
		len += n
		if n < 6|| "\r\n\r\n" != string(msg[n-4:n]) {
			continue
		}
		glog.V(3).Infoln(this.GetWs().Len(), n, err, string(data[:len-4]))

		var result Result
		err = json.Unmarshal(data[0:len-4], &result)
		len = 0

		if err != nil {
			glog.Error(err)
			continue
		}

		if (result.Code != 0){
			glog.Error(result.Data)
			return
		}

		switch result.Action {
		case "NotAuth":
			this.notAuth(&result)
		case "AuthFinish":
			this.auth(&result)
		case "TplUpdate":
			this.tplUpdate(&result)
		case "TplDelete":
			this.tplDelete(&result)
		case "GatewayUpdate":
			this.gatewayUpdate(&result)
		case "GatewaysUpdate":
			this.gatewaysUpdate(&result)
		case "GatewayDelete":
			this.gatewayDelete(&result)
		case "SimUpdate":
			this.simUpdate(&result)
		case "SimDelete":
			this.simDelete(&result)
		case "TaskUpdate":
			this.taskUpdate(&result)
		case "TasksUpdate":
			this.tasksUpdate(&result)
		case "TaskUserUpdate":
			this.taskUserUpdate(&result)
		case "TaskDelete":
			this.taskDelete(&result)
		case "WorkTimeUpdate":
			this.workTimeUpdate(&result)
		case "SipThreadUpdate":
			this.sipThreadUpdate(&result)
		}
	}
	return
}

func (this *Api) notAuth(result *Result) {
	glog.Infoln("start authentication")
	this.Auth()
}

func (this *Api) auth(result *Result) {
	if result.Code != 0 {
		glog.Error(result.Data)
		panic(result.Data)
	}
	glog.Infoln("robot authenticationOk", result.Data)
}

func (this *Api) tplUpdate(result *Result) {
	var tpl Template
	err := json.Unmarshal([]byte(result.Data), &tpl)
	if err != nil {
		glog.Error(err)
		return
	}
	this.app.lockTemplate.Lock()
	defer this.app.lockTemplate.Unlock()
	this.app.templates[tpl.Id] = &tpl

	glog.V(3).Infoln(*(this.app.templates[tpl.Id]))
	glog.V(1).Infoln("template update success, template id:", tpl.Id)
}

func (this *Api) tplDelete(result *Result) {
	id, err := strconv.Atoi(result.Data)
	if err != nil {
		glog.Error("delete template failed, id:", result.Data)
		return
	}
	this.app.lockTemplate.Lock()
	this.app.lockTemplate.Unlock()
	delete(this.app.templates, id)
	glog.V(1).Infoln("delete template id:" + result.Data)
}

func (this *Api) gatewaysUpdate(result *Result) {

	var gateways []*Gateway
	err:=json.Unmarshal([]byte(result.Data), &gateways)
	if err != nil {
		glog.Error(err)
		return
	}
	this.app.lockGateway.Lock()
	defer this.app.lockGateway.Unlock()

	for _, gateway := range gateways {
		this.app.gateways[gateway.Id] = gateway
		glog.V(3).Infoln(*(this.app.gateways[gateway.Id]))
		glog.Infoln("update gateways success, id:", gateway.Id)
	}
}

func (this *Api) gatewayUpdate(result *Result) {

	var gateway Gateway
	err := json.Unmarshal([]byte(result.Data), &gateway)
	if err != nil {
		glog.Error(err)
		return
	}
	this.app.lockGateway.Lock()
	defer this.app.lockGateway.Unlock()

	this.app.gateways[gateway.Id] = &gateway
	glog.V(3).Infoln(*(this.app.gateways[gateway.Id]))
	glog.Infoln("update gateways success, id:", gateway.Id)

}

func (this *Api) gatewayDelete(result *Result) {
	id, err := strconv.Atoi(result.Data)
	if err != nil {
		glog.Error("delete template failed, id:", result.Data)
		return
	}
	this.app.lockGateway.Lock()
	this.app.lockGateway.Unlock()
	delete(this.app.gateways, id)
	glog.V(1).Infoln("delete template id:" + result.Data)
}

func (this *Api) simUpdate(result *Result) {
	var sim Sim
	json.Unmarshal([]byte(result.Data), &sim)

	this.app.lockSim.Lock()
	defer this.app.lockSim.Unlock()
	this.app.sims[sim.Id] = &sim

	glog.V(3).Infoln(*(this.app.sims[sim.Id]))
	glog.V(1).Infoln("sim update success, sim id:", sim.Id)

	this.SimTasks(sim.Id)
}

func (this *Api) simDelete(result *Result) {
	id, err := strconv.Atoi(result.Data)
	if err != nil {
		glog.Error("delete sim failed, id:", result.Data)
		return
	}
	this.app.lockSim.Lock()
	this.app.lockSim.Unlock()
	delete(this.app.sims, id)
	glog.V(1).Infoln("delete sim id:" + result.Data)
}

func (this *Api) taskUpdate(result *Result) {
	var task Task
	json.Unmarshal([]byte(result.Data), &task)

	this.app.lockTask.Lock()
	defer this.app.lockTask.Unlock()
	this.app.tasks[task.Id] = &task

	glog.V(3).Infoln(*(this.app.tasks[task.Id]))
	glog.V(1).Infoln("task update success, task id:", task.Id)
}

func (this *Api) tasksUpdate(result *Result) {

	m := make(map[string]interface{})

	err:=json.Unmarshal([]byte(result.Data), &m)
	if err != nil {
		glog.Error("tasks update failed:", err)
		return
	}

	var tasks []Task
	err = json.Unmarshal([]byte(m["tasks"].(string)), &tasks)
	if err != nil {
		glog.Error("tasks update failed:", err)
		return
	}

	sim_id:=0
	if _,ok:=m["sim_id"]; ok{
		sim_id=int(m["sim_id"].(float64))
	}
	sip_id:=0
	if _,ok:=m["sip_id"]; ok{
		sip_id=int(m["sip_id"].(float64))
	}
	if sim_id==0 && sip_id==0 {
		glog.Error("sim_id and sip_id is invalid")
		return
	}

	this.app.lockTask.Lock()
	defer this.app.lockTask.Unlock()

	for _,task:= range tasks {
		if sim_id>0 {
			this.app.taskInfo.TaskSim[task.Id] = append(this.app.taskInfo.TaskSim[task.Id],sim_id)
		} else if sip_id>0{
			this.app.taskInfo.TaskSip[task.Id] = append(this.app.taskInfo.TaskSip[task.Id],sip_id)
		}
		this.app.tasks[task.Id] = &task
		glog.V(3).Infoln(*(this.app.tasks[task.Id]))
		glog.V(1).Infoln("task update success, task id:", task.Id)
	}
}

func (this *Api) taskDelete(result *Result) {
	id, err := strconv.Atoi(result.Data)
	if err != nil {
		glog.Error("delete task failed, id:", result.Data)
		return
	}
	this.app.lockTask.Lock()
	this.app.lockTask.Unlock()
	delete(this.app.tasks, id)
	glog.V(1).Infoln("delete task id:" + result.Data)
}

func (this *Api) taskUserUpdate(result *Result) {
	var taskUser TaskUser
	err:= json.Unmarshal([]byte(result.Data), &taskUser)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infoln(fmt.Sprintf("execute task:%d\ttask user:%d", taskUser.TaskId,taskUser.Id))

	
}


func (this *Api) workTimeUpdate(result *Result) {

	m := make(map[string]string)

	err := json.Unmarshal([]byte(result.Data), &m)
	if err != nil {
		glog.Error(err)
		return
	}

	this.app.lockUser.Lock()
	defer this.app.lockUser.Unlock()
	if v, ok := m["uid"]; !ok || v == "" {
		glog.Error("parse uid from json failed")
		return
	}

	var workTime WorkTime
	err = json.Unmarshal([]byte(m["worktime"]), &workTime)
	if err != nil {
		glog.Error(err)
		return
	}

	var user User
	user.workTime = &workTime
	this.app.users[m["uid"]] = &user

	glog.V(3).Infoln(*(this.app.users[m["uid"]].workTime))
	glog.V(1).Infoln("workTime update success, user id:", m["uid"])
}

func (this *Api) sipThreadUpdate(result *Result) {
	var task Task
	json.Unmarshal([]byte(result.Data), &task)

	this.app.lockTask.Lock()
	defer this.app.lockTask.Unlock()
	this.app.tasks[task.Id] = &task

	glog.V(3).Infoln(*(this.app.tasks[task.Id]))
	glog.V(1).Infoln("task update success, task id:", task.Id)
}

//func (this *Api) update() {
//	this.getaway()
//	for _, g := range this.app.Gateways {
//		this.sim(g.Id)
//	}
//
//	for _, s := range this.app.Sims {
//		this.tasks(s.Id)
//	}
//}
//
//// 获取当前app管理的网关信息
//func (this *Api) getaway() (gateways []Gateway, err error) {
//
//	req, err := http.NewRequest("GET", this.ApiUrl+"/gateways", nil)
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
//	var gateway _Gateway
//	err = json.Unmarshal(ret, &gateway)
//	if err != nil {
//		return
//	}
//
//	if gateway.Meta.Code != 0 {
//		return nil, errors.New(fmt.Sprintf("获取网关信息错误码：%d", gateway.Meta.Code))
//	}
//
//	func() {
//		this.app.lockGateway.Lock()
//		defer this.app.lockGateway.Unlock()
//		for _, g := range gateway.Data["gateways"] {
//			this.app.Gateways[g.Id] = g
//		}
//	}()
//
//	return
//}
//
//// 根据网关获取该网关管理的sim卡
//func (this *Api) sim(gatewayId int) (sims []Sim, err error) {
//	req, err := http.NewRequest("GET", fmt.Sprintf(this.ApiUrl+"/gateway/%d/sims", gatewayId), nil)
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
//	var _sim _Sim
//
//	err = json.Unmarshal(ret, &_sim)
//	if err != nil {
//		return
//	}
//
//	if _sim.Meta.Code != 0 {
//		err = errors.New(fmt.Sprintf("获取sim卡信息错误码：%d", _sim.Meta.Code))
//		return
//	}
//
//	func() {
//		this.app.lockSim.Lock()
//		defer this.app.lockSim.Unlock()
//		for _, s := range _sim.Data["sims"] {
//			this.app.Sims[s.Id] = s
//		}
//	}()
//
//	return
//}
//
//func (this *Api) GetTpl(id int) (tpl Template, err error) {
//
//	// 先从map上面找，如果找到了就不用从远程获取了
//	this.app.lockTemplate.Lock()
//	tpl, ok := this.app.Templates[id]
//	this.app.lockTemplate.Unlock()
//	if ok {
//		return
//	}
//
//	req, err := http.NewRequest("GET", fmt.Sprintf(this.ApiUrl+"/tpl/%d", id), nil)
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
//	var _tpl _Template
//
//	err = json.Unmarshal(ret, &_tpl)
//	if err != nil {
//		return
//	}
//
//	if _tpl.Meta.Code != 0 {
//		err = errors.New(fmt.Sprintf("获取模板信息错误码：%d", _tpl.Meta.Code))
//		return
//	}
//
//	func() {
//		this.app.lockTemplate.Lock()
//		defer this.app.lockTemplate.Unlock()
//		tpl = _tpl.Data["tpl"]
//		this.app.Templates[tpl.Id] = tpl
//	}()
//	return
//}
//
//func (this *Api) tasks(sim_id int) (tasks []Task, err error) {
//	req, err := http.NewRequest("GET", fmt.Sprintf(this.ApiUrl+"/sim/%d/tasks", sim_id), nil)
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
//	var _task _Task
//
//	err = json.Unmarshal(ret, &_task)
//	if err != nil {
//		return
//	}
//
//	if _task.Meta.Code != 0 {
//		err = errors.New(fmt.Sprintf("获取任务信息错误码：%d", _task.Meta.Code))
//		return
//	}
//
//	func() {
//		this.app.lockTask.Lock()
//		defer this.app.lockTask.Unlock()
//		var tasks []Task
//		for _, t := range _task.Data["tasks"] {
//			tasks = append(tasks, t)
//		}
//		this.app.Tasks[sim_id] = tasks
//	}()
//	return
//}
//
//func (this *Api) Poll() {
//
//	for {
//		this.update()
//		time.Sleep(time.Second * 10)
//	}
//}
//
//// 执行api请求
//func (this *HttpClient) Do(req *http.Request) (*http.Response, error) {
//	req.Header.Set("appid", this.AppId)
//	req.Header.Set("key", this.Key)
//
//	client := &http.Client{}
//	return client.Do(req)
//}
//
//func (this *Api) DownloadVoice(localPath string, voices map[string]*tpl.Voice) error {
//	for k, _ := range voices {
//		pcm := this.ApiUrl + "/voice/file/pcm/" + k
//
//		localPcm := localPath + "/" + k + ".pcm"
//
//		// 如果文件存在，无需下载
//		if _, err := os.Stat(localPcm); err == nil {
//			voices[k].LocalPcm = localPcm
//			continue
//		}
//
//		fmt.Println(pcm)
//		req, err := http.NewRequest("GET", pcm, nil)
//		if err != nil {
//			return err
//		}
//
//		res, err := this.Client.Do(req)
//		if err != nil {
//			return err
//		}
//
//		if res.StatusCode == http.StatusOK {
//			f, err := os.Create(localPcm)
//			if err != nil {
//				return err
//			}
//			io.Copy(f, res.Body)
//			voices[k].LocalPcm = localPcm
//		} else {
//			return errors.New(fmt.Sprint("获取pcm文件:[ ", pcm, "] 失败，http状态码:", res.StatusCode))
//		}
//	}
//	return nil
//}
//
//func (this *Api) UploadVoice(wav string) (hash string, err error) {
//
//	file, err := os.Open(wav)
//	if err != nil {
//		return
//	}
//	defer file.Close()
//
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("voice", wav)
//	if err != nil {
//		return
//	}
//	_, err = io.Copy(part, file)
//
//	err = writer.Close()
//	if err != nil {
//		return
//	}
//	request, err := http.NewRequest("POST", this.ApiUrl+"/voice/upload", body)
//	request.Header.Set("Data-Type", writer.FormDataContentType())
//
//	res, err := this.Client.Do(request)
//	if err != nil {
//		return
//	}
//	ret, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return
//	}
//
//	var _voice _Voice
//	err = json.Unmarshal(ret, &_voice)
//	if err != nil {
//		return
//	}
//
//	if _voice.Meta.Code != 0 {
//		err = errors.New(fmt.Sprintf("获取sim卡信息错误码：%d", _voice.Meta.Code))
//		return
//	}
//
//	hash = _voice.Data["voice"].Hash
//
//	return
//}
//
//func (this *Api) getUserConfig(userId, name string) (config string, err error) {
//	req, err := http.NewRequest("GET", fmt.Sprintf(this.ApiUrl+"/user/%s/config/%s", userId, name), nil)
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
//	config = string(ret)
//	return
//}
//
//func (this *Api) SendReport(template *tpl.Tpl, taskUserId int) (err error) {
//
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//
//	for i, t := range template.Report.Nodes {
//		if t.Type != tpl.HumanReport {
//			continue
//		}
//		template.Report.Nodes[i].Voice, err = this.UploadVoice(t.Voice)
//		if err != nil {
//			return
//		}
//	}
//
//	// 客户有接电话
//	if template.Report.Time > 0 && len(template.Report.Nodes) > 0 {
//		var str string
//		str, err = template.Report.ToJson()
//		if err != nil {
//			return
//		}
//		writer.WriteField("report", str)
//	} else {
//		writer.WriteField("report", "")
//	}
//
//	writer.WriteField("time", fmt.Sprint(int(template.Report.Time.Seconds())))
//	writer.WriteField("type", fmt.Sprint(template.Report.Type))
//
//	err = writer.Close()
//	if err != nil {
//		return
//	}
//	request, err := http.NewRequest("POST", fmt.Sprintf("%s/task/user/%d/report", this.ApiUrl, taskUserId), body)
//	request.Header.Set("Data-Type", writer.FormDataContentType())
//
//	res, err := this.Client.Do(request)
//	if err != nil {
//		return
//	}
//
//	ret, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return
//	}
//
//	var _voice _Voice
//	err = json.Unmarshal(ret, &_voice)
//	if err != nil {
//		return
//	}
//
//	if _voice.Meta.Code != 0 {
//		err = errors.New(fmt.Sprintf("发送拨打报告失败错误码：%d", _voice.Meta.Code))
//		return
//	}
//
//	return
//}
