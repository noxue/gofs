package main

import (
	"gofs/phone"
	"fmt"
	"math/rand"
	"time"
)

type EndPoint struct {
}

func (this *EndPoint) Create(call *phone.Call) {

}

func (this *EndPoint) Answer(call *phone.Call) {
	call.Play("d:/0012.wav")
}

func (this *EndPoint) Hangup(call *phone.Call) {

}

func (this *EndPoint) Destroy(call *phone.Call) {

}

func (this *EndPoint) SpeakStart(call *phone.Call) {
	call.Pause(true)
	call.Record(fmt.Sprintf("d:/%s.wav", rand.Intn(100000)))
}

func (this *EndPoint) SpeakEnd(call *phone.Call) {
	call.Stop()
	wav, _ := call.RecordStop()
	call.Play(wav)
	fmt.Println(call.GetDataString("word"), call.GetDataString("file"))
}

func main() {
	go Api.Handle()
	time.Sleep(time.Second*3)
	for {
		tasks := Api.Tasks()
		for id, _ := range tasks {
			//fmt.Println("----",id)
			Api.TaskUser(id)
		}
		time.Sleep(time.Millisecond * 1000)
	}
	for {
		time.Sleep(time.Second)
	}
}

//
//func main() {
//	p, err := phone.New("localhost", 8021, "ClueCon", 10)
//	if err != nil {
//		glog.Error(err)
//	}
//
//	call:=phone.NewCall("13758277505",&EndPoint{},time.Minute*5)
//	p.MakeCall("xigao", call)
//	p.Handle()
//}
