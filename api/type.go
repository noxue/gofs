package api

import (
	"sync"
	"time"
	"golang.org/x/net/websocket"
)

type Api struct {
	appId  string
	key    string
	apiUrl string
	origin string
	lockWs sync.Mutex
	ws     *websocket.Conn
	app    App
}

type App struct {
	gateways     map[int]*Gateway
	lockGateway  sync.Mutex
	sims         map[int]*Sim
	lockSim      sync.Mutex
	templates    map[int]*Template
	lockTemplate sync.Mutex
	tasks        map[int]*Task
	lockTask     sync.Mutex
	users        map[string]*User
	lockUser     sync.Mutex
	taskInfo     *TaskInfo
}

type TaskInfo struct {
	TaskSim map[int][]int // map[simId][]taskId
	TaskSip map[int][]int
}

type Result struct {
	Action string `json:"action"`
	Code   int    `json:"code"`
	Data   string `json:"data"`
}

type User struct {
	Id        string `json:"id"`
	SipThread int    `json:"thread"`
	workTime  *WorkTime
}

// 网关信息
type Gateway struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	User string `json:"userId"`
}

// sim卡信息
type Sim struct {
	Id     int    `json:"id"`
	Gid    int    `json:"gatewayId"` // 卡属于哪个网关
	Number string `json:"number"`
	User   string `json:"userId"`
}

type Template struct {
	Id     int    `json:"id"`
	User   string `json:"userId"`
	Status int    `json:"status"`
	Tpl    string `json:"content"`
}

type Task struct {
	Id          int    `json:"id"`
	User        string `json:"userId"`
	Template    int    `json:"templateId"`
	Thread      int    `json:"thread"`
	ThreadCount int
	Total       int    `json:"total"`
	Status      int    `json:"status"`
	Break       int    `json:"interrupt"`
	IsTest      bool   `json:"test"`
}

type TaskUser struct {
	Id     int    `json:"id"`
	TaskId int    `json:"taskId"`
	Mobile string `json:"mobile"`
}

type Schedule struct {
	Times []time.Time `json:"workTime"`
}

type WorkTime struct {
	Repeat     []int      `json:"repeat"`
	Schedule   []Schedule `json:"schedule"`
	UpdateTime time.Time
}
