package api

import (
	"github.com/golang/glog"
	"gofs/fs"
	"fmt"
	"gofs/tpl"
	"math/rand"
)

type EndPoint struct {
	IsPlay bool
	Type string
	Id int // simId or sipId
	TaskId int
	TaskUserId int
	Tpl *tpl.Tpl
}

func (this *EndPoint) Create(call *fs.Call) {
	fmt.Println("-------------------------------------------create========================")
	call.Asr(false)
}

func (this *EndPoint) Answer(call *fs.Call) {
	call.Asr(true)
	call.Play("d:/A1.wav")
	fmt.Println("-------------------------------------------answer========================")
}

func (this *EndPoint) Hangup(call *fs.Call) {
	fmt.Println("-------------------------------------------hungup========================")
}

func (this *EndPoint) Destroy(call *fs.Call) {
	fmt.Println("-------------------------------------------destroy========================")
	if call.GetDataString("type") == "sim" {
		t, ok := call.GetData("simId")
		if ok {
			simId := int(t.(float64))
			if simId > 0 {
				TaskApi.GetTaskInfo().SimFree(simId, true)
			}
		}
	}

	if this.Type=="sim"{
		TaskApi.GetTaskInfo().SimFree(this.Id,true)
	}
}

func (this *EndPoint) SpeakStart(call *fs.Call) {
	//call.Pause(true)
	//call.Record(fmt.Sprintf("d:/%s.wav", rand.Intn(100000)))
}

func (this *EndPoint) SpeakEnd(call *fs.Call) {
	//call.Stop()
	//wav, _ := call.RecordStop()
	//call.Play(wav)
	fmt.Println(call.GetDataString("word"), call.GetDataString("file"))
}

func (this *EndPoint) Progress(call *fs.Call) {

}

func (this *EndPoint) ProgressMedia(call *fs.Call) {
	call.Record(fmt.Sprintf("d:/records/%d.wav", rand.Intn(100000)))
}

func (this *EndPoint) HangupComplete(call *fs.Call) {
	fmt.Println(call.RecordStop())
}

func (this *EndPoint) PlaybackStart(call *fs.Call) {
	this.IsPlay = true
	glog.V(3).Info("================开始播放录音===================")
}

func (this *EndPoint) PlaybackStop(call *fs.Call) {
	this.IsPlay = false
	glog.V(3).Info("================停止播放录音===================")
}
