package phone

import (
	"runtime"
	"github.com/0x19/goesl"
	"strings"
	"fmt"
	"github.com/golang/glog"
)

type Phone struct {
	client *goesl.Client
	calls *CallList
	host string
	port uint
	pass string
	timeout int
}

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func New(host string, port uint, pass string, timeout int) (client *Phone,err error){
	c := &Phone{
		host:host,
		port:port,
		pass:pass,
		timeout:timeout,
		calls:NewCallList(),
	}
	c.client,err=goesl.NewClient(host, port, pass, timeout)
	if err != nil {
		return
	}
	client = c
	return
}

func (this *Phone)Handle(){
	client := this.client
	go client.Handle()
	this.client.Send("events json CHANNEL_CREATE CHANNEL_ANSWER CHANNEL_DESTROY CHANNEL_HANGUP CHANNEL_HANGUP_COMPLETE CUSTOM asr::start_speak asr::end_speak")
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
		fmt.Println(msg)
		//msg.Parse()
		uuid := msg.GetHeader("Caller-Unique-ID")
		eventName := msg.GetHeader("Event-Name")
		number := msg.GetHeader("Caller-Destination-Number")
		call,ok := this.calls.Get(number)
		if !ok{
			continue
		}
		switch  eventName {
		case "CHANNEL_CREATE":
			call.SetUuid(uuid)
			call.callInterface.Create(call)
			break;
		case "CHANNEL_ANSWER":
			call.callInterface.Answer(call)
			break;
		case "CHANNEL_DESTROY":
			call.callInterface.Destroy(call)
			break;
		case "CHANNEL_HANGUP":
			call.callInterface.Hangup(call)
			break;
		case "CHANNEL_HANGUP_COMPLETE":
			break;
		case "asr::start_speak":
			call.callInterface.SpeakStart(call)
			break;
		case "asr::end_speak":
			word := msg.GetHeader("Word")
			file := msg.GetHeader("File")
			call.SetData("word", word)
			call.SetData("file",file)
			call.callInterface.SpeakEnd(call)
			break;
		}

	}
}
