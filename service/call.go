package service

import (
	"gofs/fs"
	"fmt"
	"math/rand"
)

type EndPoint struct {
}

func (this *EndPoint) Create(call *fs.Call) {
	fmt.Println("-------------------------------------------create========================")

}

func (this *EndPoint) Answer(call *fs.Call) {
	call.Play("d:/0012.wav")
	fmt.Println("-------------------------------------------answer========================")
}

func (this *EndPoint) Hangup(call *fs.Call) {

}

func (this *EndPoint) Destroy(call *fs.Call) {

}

func (this *EndPoint) SpeakStart(call *fs.Call) {
	call.Pause(true)
	call.Record(fmt.Sprintf("d:/%s.wav", rand.Intn(100000)))
}

func (this *EndPoint) SpeakEnd(call *fs.Call) {
	call.Stop()
	wav, _ := call.RecordStop()
	call.Play(wav)
	fmt.Println(call.GetDataString("word"), call.GetDataString("file"))
}

