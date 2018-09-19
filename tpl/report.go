package tpl

import (
	"encoding/json"
	"time"
)

const (
	HumanReport = iota // human
	FlowReport
	KeywordReport
)

// 呼叫结果报告
type Report struct {
	Nodes          []*ReportNode `json:"nodes"`
	Type           int           `json:"type"`        // 客户类型
	Voice          string        `json:"voice"`       // 通话全程录音文件hash
	Status         int           `json:"status"`      // 电话拨打状态，0，空号。1，未打通。2，未接。3，客户挂机，4，机器人挂机
	TimeStart      time.Time     `json:"timeStart"`   // 开始拨打电话的时间
	TimeConnect    time.Time     `json:"timeConnect"` // 开始接通电话的时间
	TimeDisconnect time.Time     `json:"timeEnd"`     // 结束通话的时间
	Time           time.Duration `json:"time"`        // 通话时间
}

type ReportNode struct {
	Type     int       `json:"type"`  // 流程类型 0 human   1 flow   2  keyword
	Name     string    `json:"name"`  // 流程名称/或者关键字流程名称
	Word     string    `json:"word"`  // 所说的文字
	Voice    string    `json:"voice"` //  文字对应的wav录音hash
	UserType int       `json:"userType"`
	Time     time.Time `json:"time"`
}

// 添加人说的报告
func (this *Report) AddHuman(word, voice string) {
	this.Nodes = append(this.Nodes, &ReportNode{
		Type:  HumanReport,
		Word:  word,
		Voice: voice,
		Time:  time.Now(),
	})
}

// 添加机器人执行的流程
func (this *Report) AddFlow(name, word, voice string, userType int) {
	this.Nodes = append(this.Nodes, &ReportNode{
		Type:     FlowReport,
		Name:     name,
		Word:     word,
		Voice:    voice,
		UserType: userType,
		Time:     time.Now(),
	})
}

// 添加机器人执行的关键词
// name 流程名称
func (this *Report) AddKeyword(name, word, voice string, userType int) {
	this.Nodes = append(this.Nodes, &ReportNode{
		Type:     KeywordReport,
		Name:     name,
		Word:     word,
		Voice:    voice,
		UserType: userType,
		Time:     time.Now(),
	})
}

func (this *Report) ToJson() (str string, err error) {

	ret, err := json.Marshal(this)
	if err != nil {
		return
	}
	str = string(ret)
	return
}
