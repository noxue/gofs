package fs

import (
	"github.com/0x19/goesl"
	"strings"
	"github.com/golang/glog"
	"time"
)

type Phone struct {
	client  *goesl.Client
	calls   *CallList
	host    string
	port    uint
	pass    string
	timeout int
}

func New(host string, port uint, pass string, timeout int) (client *Phone, err error) {
	c := &Phone{
		host:    host,
		port:    port,
		pass:    pass,
		timeout: timeout,
		calls:   NewCallList(),
	}
	c.client, err = goesl.NewClient(host, port, pass, timeout)
	if err != nil {
		return
	}
	client = c
	return
}

func (this *Phone) Handle() {
	client := this.client
	go client.Handle()

	//this.client.Send("events json CHANNEL_UUID CHANNEL_CREATE CHANNEL_PROGRESS CHANNEL_PROGRESS_MEDIA CHANNEL_ANSWER CHANNEL_DESTROY CHANNEL_HANGUP CHANNEL_HANGUP_COMPLETE PLAYBACK_START PLAYBACK_STOP CUSTOM asr::start_speak asr::end_speak")
	this.client.Send("event json CHANNEL_DESTROY")
	//this.client.Send("event json CHANNEL_CREATE")
	//this.client.Send("event json CHANNEL_ANSWER")
	//this.client.Send("event json CHANNEL_HANGUP")
	//this.client.Send("event json ALL")

	for {
		msg, err := client.ReadMessage()
		if err != nil {
			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				glog.Error("Error while reading Freeswitch message:", err)
			}
			break
		}
		if (msg == nil) {
			continue
		}
		glog.V(3).Infoln(msg)
		uuid := msg.GetHeader("Caller-Unique-ID")
		eventName := msg.GetHeader("Event-Name")
		numberArr := strings.Split(msg.GetHeader("Caller-Destination-Number"),"@")
		number:=""
		if len(numberArr)>=1 {
			number = numberArr[0]
		} else {
			continue
		}
		call, ok := this.calls.Get(number)
		glog.V(2).Infoln("Event Name:", eventName, "Number:", number)
		if !ok {
			continue
		}

		glog.V(2).Infoln("Call Number:", call.number)
		switch  eventName {
		case "CHANNEL_CREATE":
			call.SetUuid(uuid)
			go func() {
				time.Sleep(call.timeout)
				call.Hungup()
			}()
			go call.callInterface.Create(call)
		case "CHANNEL_PROGRESS":
			go call.callInterface.Progress(call)
		case "CHANNEL_PROGRESS_MEDIA":
			go call.callInterface.ProgressMedia(call)
		case "CHANNEL_ANSWER":
			go call.callInterface.Answer(call)
		case "CHANNEL_DESTROY":
			go call.callInterface.Destroy(call)
		case "CHANNEL_HANGUP":
			go call.callInterface.Hangup(call)
		case "CHANNEL_HANGUP_COMPLETE":
			go call.callInterface.HangupComplete(call)
		case "PLAYBACK_START":
			go call.callInterface.PlaybackStart(call)
		case "PLAYBACK_STOP":
			go call.callInterface.PlaybackStop(call)
		case "asr::start_speak":
			go call.callInterface.SpeakStart(call)
		case "asr::end_speak":
			word := msg.GetHeader("Word")
			file := msg.GetHeader("File")
			call.SetData("word", word)
			call.SetData("file", file)
			go call.callInterface.SpeakEnd(call)

		case "CHANNEL_UUID":
			call.SetUuid(uuid)
		}


	}
}
