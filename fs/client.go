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

	this.client.Send("events json CHANNEL_CREATE CHANNEL_PROGRESS CHANNEL_PROGRESS_MEDIA CHANNEL_ANSWER CHANNEL_DESTROY CHANNEL_HANGUP CHANNEL_HANGUP_COMPLETE CUSTOM asr::start_speak asr::end_speak")
	for {
		msg, err := client.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				glog.Error("Error while reading Freeswitch message: %s", err)
			}
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
			call.callInterface.Create(call)
		case "CHANNEL_PROGRESS":
			call.callInterface.Progress(call)
		case "CHANNEL_PROGRESS_MEDIA":
			call.callInterface.ProgressMedia(call)
		case "CHANNEL_ANSWER":
			call.callInterface.Answer(call)
		case "CHANNEL_DESTROY":
			call.callInterface.Destroy(call)
		case "CHANNEL_HANGUP":
			call.callInterface.Hangup(call)
		case "CHANNEL_HANGUP_COMPLETE":
			call.callInterface.HangupComplete(call)
		case "asr::start_speak":
			call.callInterface.SpeakStart(call)
		case "asr::end_speak":
			word := msg.GetHeader("Word")
			file := msg.GetHeader("File")
			call.SetData("word", word)
			call.SetData("file", file)
			call.callInterface.SpeakEnd(call)
		}

	}
}
