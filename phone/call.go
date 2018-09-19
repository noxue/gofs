package phone

import (
	"errors"
	"fmt"
	"time"
	"github.com/0x19/goesl"
	"os"
	"sync"
)

type CallInterface interface {
	Create(call *Call)
	Answer(call *Call)
	Hangup(call *Call)
	Destroy(call *Call)
	SpeakStart(call *Call)
	SpeakEnd(call *Call)
}

type Call struct {
	number        string
	uuid          string
	createAt      time.Time
	answerAt      time.Time
	hangupAt      time.Time
	destroyAt     time.Time
	callInterface CallInterface
	client        *goesl.Client

	// for translate the user custom data
	dataMap     map[string]interface{}
	dataMapLock sync.Mutex

	// store the record wav file path
	records []string
	isRecording bool

	// channel alive max time
	timeout time.Duration
}

func NewCall(number string,callInterface CallInterface,timeout time.Duration) *Call{
	return &Call{
		number:number,
		dataMap :make(map[string]interface{}),
		callInterface:callInterface,
		timeout:timeout,
	}
}

func (this *Phone) MakeCall(gateway string, call *Call) (err error) {
	if len(gateway) == 0 || len(call.GetNumber()) == 0 {
		err = errors.New("gateway or number is empty")
		return
	}
	call.client = this.client
	this.calls.Set(call.GetNumber(), call)
	this.client.BgApi(fmt.Sprintf("originate {ignore_early_media=false,absolute_codec_string=pcma,origination_caller_id_number=" + gateway + "}sofia/gateway/" + gateway + "/" + call.GetNumber() + " 'ai:asdfaf' inline"))
	return
}

func (this *Call) SetData(key string, val interface{}) {
	this.dataMapLock.Lock()
	defer this.dataMapLock.Unlock()
	this.dataMap[key] = val
}

func (this *Call) GetData(key string) (val interface{}, ok bool) {
	this.dataMapLock.Lock()
	defer this.dataMapLock.Unlock()
	val, ok = this.dataMap[key]
	return
}

func (this *Call) GetDataString(key string) string {
	val, ok := this.GetData(key)
	if !ok {
		return ""
	}
	v, ok := val.(string)
	if !ok {
		return ""
	}
	return v
}

func (this *Call) GetNumber() string {
	return this.number
}

func (this *Call) GetUuid() string {
	return this.uuid
}

func (this *Call) SetUuid(uuid string) {
	this.uuid = uuid
}

func (this *Call) Play(wav string) (err error) {

	_, err = os.Stat(wav)
	if err != nil {
		err = errors.New("the wav file is not exists：" + wav)
		return
	}

	if len(this.GetUuid()) == 0 {
		err = errors.New("uuid is empty")
		return
	}

	this.client.BgApi("uuid_play " + this.GetUuid() + " " + wav)
	return
}

func (this *Call) Pause(on bool) (err error) {
	if len(this.GetUuid()) == 0 {
		err = errors.New("uuid is empty")
		return
	}

	if on {
		this.client.BgApi("uuid_pause " + this.GetUuid() + " on")
	} else {
		this.client.BgApi("uuid_pause " + this.GetUuid() + " off")
	}

	return
}

func (this *Call) Stop() (err error) {
	if len(this.GetUuid()) == 0 {
		err = errors.New("uuid is empty")
		return
	}

	this.client.BgApi("uuid_stop " + this.GetUuid())
	return
}

func (this *Call)record(wav string,flag bool) (err error) {

	if flag {
		if len(this.GetUuid()) == 0 {
			err = errors.New("uuid is empty")
			return
		}
		this.isRecording = true
		this.records = append(this.records, wav)
		this.client.BgApi("uuid_record " + this.GetUuid() + " start "+wav)
	} else {
		this.isRecording = false
		if len(this.records) == 0 {
			err = errors.New("no recording now")
			return
		}
		wav = this.records[len(this.records)-1]
		this.client.Api("uuid_record " + this.GetUuid() + " stop "+wav)
	}
	return
}

func (this *Call)Record(wav string) error {
	return this.record(wav,true)
}

func (this *Call)RecordStop() (string,error) {
	err:=this.record("", false)
	if err != nil {
		return "",err
	}
	return this.records[len(this.records)-1],nil
}

func (this *Call)Asr(on bool)(err error) {
	if len(this.GetUuid()) == 0 {
		err = errors.New("uuid is empty")
		return
	}
	if on {
		this.client.BgApi("uuid_asr " + this.GetUuid() + " on")
	} else {
		this.client.BgApi("uuid_asr " + this.GetUuid() + " off")
	}
	return
}
