package tpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)


func New(tplStr, localBasePath, voiceRemoteBaseUrl string) (*Tpl, error) {
	var tpl1 Tpl
	err := json.Unmarshal([]byte(tplStr), &tpl1)
	if err != nil {
		return nil, err
	}

	tpl1.LocalBasePath = localBasePath
	tpl1.VoiceRemoteBaseUrl = voiceRemoteBaseUrl

	tpl1.Ai = &Ai{}
	tpl1.Person = &Person{}

	if err := tpl1.initFlowKeyword(); err != nil {
		return nil, err
	}

	if err := tpl1.initKeyword(); err != nil {
		return nil, err
	}

	if err := tpl1.initKeywordCondKeyword(); err != nil {
		return nil, err
	}

	tpl1.SetStartTime()
	tpl1.BreakChan = make(chan bool)

	return &tpl1, nil
}

func (this *Tpl) Start() (err error) {

	// 主流程设置为当前流程
	this.Cur = this.Main

	this.SetConnectTime()
Loop:

	var opt *Operation
	var voice *Voice
	var match bool
	var next string
	var noMatchTimes int

	flow, err := this.getFlow(this.Cur)
	if err != nil {
		goto End
	}

	// 机器人说，就是让人听到
	{
		voice, err := this.choiceVoice(flow.Section.Voice, flow.Section.Choice)
		if err != nil {
			goto End
		}
		opt := &Operation{
			Opt:   flow.Section.Type,
			Text:  voice.Text,
			Sound: voice.LocalPcm,
		}
		this.Person.SetListen(opt)
		this.Report.AddFlow(this.Cur, voice.Text, voice.Hash, flow.Type)
	}

Listen:

	opt = &Operation{}
	voice = &Voice{}
	match = false
	next = ""

	// 机器人一直听，直到到人类说话
	for !this.isQuit && !this.Ai.IsListened() {
		time.Sleep(time.Millisecond * 10)
	}

	opt = this.Ai.Listen()
	if opt == nil {
		goto End
	}
	// 根据返回的操作类型判断是否是超时
	if opt.Opt == QUIET1 || opt.Opt == QUIET2 {
		opt1 := opt.Opt
		if _, ok := this.Keyword[opt1]; !ok {
			// 如果流程不存在的话，直接结束
			err = errors.New("该关键词组不存在：" + opt1)
			goto End
		}

		if len(this.Keyword[opt1].Voice) > 0 {
			voice, err = this.choiceVoice(this.Keyword[opt1].Voice, this.Keyword[opt1].Choice)
		}

		if opt1 == QUIET2 {
			opt1 = END
		}
		this.Person.SetListen(&Operation{
			Opt:   opt1,
			Text:  voice.Text,
			Sound: voice.LocalPcm,
		})

		if opt1 == END {
			goto End
		}
		this.Report.AddKeyword(opt1, voice.Text, voice.Hash, this.Keyword[opt1].Type)

		// 提醒用户之后，继续等待用户说话
		if this.Keyword[opt1].Next == RETURN {
			// 重新执行一遍当前流程
			goto Loop
		} else if this.Keyword[opt1].Next == WAIT {
			// 跳转到听客户说话的地方重新等待用户说话
			goto Listen
		} else if this.Keyword[opt1].Next == NEXT {
			// 跳转到下一个流程
			this.Cur = flow.Next
			goto Loop
		}
	} else {
		this.Report.AddHuman(opt.Text, opt.Sound)
	}

	// 如果当前流程关注全局关键词，则根据人类说话，去匹配全局关键词
	if flow.Hook {

		match, next, err = this.hookKeyword(opt)
		if err != nil {
			goto End
		}
		if match && next != "" {
			noMatchTimes = 0
			// 根据下一步流程的类型设置this.Cur 或返回对应操作，比如，重新执行当前流程，或下一个流程
			if next == RETURN {
				// 重新执行一遍当前流程
				goto Loop
			} else if next == WAIT {
				// 跳转到听客户说话的地方重新等待用户说话
				goto Listen
			} else if next == NEXT {
				// 跳转到下一个流程
				this.Cur = flow.Next
				goto Loop
			} else {
				this.Cur = next
				goto Loop
			}
		}
	}

	// 如果流程不是条件类型，直接进入下一步
	if flow.Section.Type != CONDITION {
		this.Cur = flow.Next
		goto Loop
	}

	// 进行到这里就一定是条件类型，进行条件处理
	for _, v := range flow.Section.KeywordMap {
		if strings.Contains(opt.Text, v.Word) {
			// 匹配到分支
			if len(flow.Section.Conds)-1 < v.Index {
				err = errors.New(fmt.Sprintf("[%v]节点条件[%v]下标越界", this.Cur, v.Index))
				goto End
			}
			this.Cur = flow.Section.Conds[v.Index].To

			noMatchTimes = 0
			goto Loop
		}
	}

	// 如果设置了下一步，就跳转到下一步
	if flow.Next != "" {
		this.Cur = flow.Next
		goto Loop
	}

	// 没有匹配任何关键词，分别处理三次没匹配到任何内容
	noMatchTimes++

	if noMatchTimes > 0 && noMatchTimes <= 3 {
		noword := []string{"noword1", "noword2", "noword3"}
		keyword := this.Keyword[noword[noMatchTimes-1]]

		// 如果有语音文件，就发送语音文件
		if len(keyword.Voice) > 0 {
			voice, err = this.choiceVoice(keyword.Voice, keyword.Choice)
			if err != nil {
				goto End
			}
			if voice != nil {
				this.Person.SetListen(&Operation{
					Opt:   noword[noMatchTimes-1],
					Text:  voice.Text,
					Sound: voice.LocalPcm,
				})
				this.Report.AddKeyword(noword[noMatchTimes-1], voice.Text, voice.Hash, keyword.Type)
			}
		}

		goto Listen
	}

	// 最多判断3次，超过3次还就直接结束
	if noMatchTimes > 3 {
		err = nil
		goto End
	}

	// 没结束就一直循环
	if !this.isQuit {
		goto Loop
	}

	err = nil
End:
	this.Person.SetListen(&Operation{
		Opt: END,
	})
	return err
}

func (this *Tpl) Listen(opt *Operation) {
	go this.Ai.SetListen(opt)
}

func (this *Tpl) Speak() *Operation {
	return this.Person.Listen()
}

func (this *Tpl) SetStartTime() {
	this.Report.TimeStart = time.Now()
}

func (this *Tpl) SetConnectTime() {
	this.Report.TimeConnect = time.Now()
}

func (this *Tpl) SetDisconnectTime() {
	this.Report.TimeDisconnect = time.Now()
	this.Report.Time = this.Report.TimeDisconnect.Sub(this.Report.TimeConnect)
}

func (this *Tpl) SpeakOver() {
	go func() {
		times := 0
		i := 0
		for {
			time.Sleep(time.Millisecond * 10)
			times++
			if times == 1000 {
				this.Ai.SetListen(&Operation{
					Opt: QUIET1,
				})
			} else if times == 2000 {
				this.Ai.SetListen(&Operation{
					Opt: QUIET2,
				})
			}

			if this.Person.IsSpeaking() {
				i++
			} else {
				i = 0
			}

			// 连续1秒 或者 机器人听到说话 才结束
			if i > 100 || this.Ai.IsListened() {
				break
			}
		}
	}()
}
