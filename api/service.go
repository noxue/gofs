package api

import (
	"gofs/fs"
	"fmt"
	"math/rand"
)

type EndPoint struct {
}

func (this *EndPoint) Create(call *fs.Call) {
	fmt.Println("-------------------------------------------create========================")
	call.Asr(false)
}

func (this *EndPoint) Answer(call *fs.Call) {
	call.Asr(true)
	call.Play("d:/0012.wav")
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
