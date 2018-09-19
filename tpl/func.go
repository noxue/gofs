package tpl

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 是否在说话
func (this *BaseStatus) IsSpeaking() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.isSpeaking
}

// 设置说话状态
func (this *BaseStatus) SetSpeakStatus(isSpeak bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.isSpeaking = isSpeak
}

// 听
func (this *BaseStatus) SetListen(opt *Operation) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if opt == nil {
		return errors.New("参数不能为空指针")
	}
	this.listen = opt
	if opt.Opt != QUIET1 || opt.Opt != QUIET2 {
		this.isListened = true
	}

	return nil
}

// 是否听到了内容
func (this *BaseStatus) IsListened() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.isListened
}

// 获取听到的内容
func (this *BaseStatus) Listen() *Operation {
	this.lock.Lock()
	defer this.lock.Unlock()

	if !this.isListened {
		return nil
	}
	this.isListened = false
	return this.listen
}

// 是否打断
func (this *Tpl) IsBreak() bool {
	return this.Break != 0
}

// 是否允许声音打断
func (this *Tpl) IsVoiceBreak() bool {
	return this.Break == -1
}

// 是否允许关键词打断
func (this *Tpl) IsKeywordBreak() bool {
	return this.Break == -2
}

// 初始化全局关键词，关键词排序
func (this *Tpl) initKeyword() error {
	// 分析全局关键词
	for k, v := range this.Keyword {
		for _, v1 := range v.Keyword {
			// 解析关键词，如果关键词格式错误，返回错误信息
			wordRank := strings.Split(v1, "#")
			if len(wordRank) != 2 {
				return errors.New("关键词格式不对:" + v1)
			}

			word := wordRank[0]
			rank, err := strconv.Atoi(wordRank[1])
			if err != nil {
				return errors.New("关键词格式不对:" + v1)
			}

			keywordMap := &KeywordMap{
				Word: word,
				Id:   k,
				Rank: rank,
			}
			this.KeywordMap = append(this.KeywordMap, *keywordMap)
		}
	}

	// 排序关键词
	l := len(this.KeywordMap)
	for i := 0; i < l-1; i++ {
		max := this.KeywordMap[i].Rank
		maxIndex := i
		// 从未排序的元素中找到最大的
		for j := i + 1; j < l; j++ {
			if this.KeywordMap[j].Rank > max {
				max = this.KeywordMap[j].Rank
				maxIndex = j
			}
		}

		// 如果找到更大的，交换
		if maxIndex != i {
			var t KeywordMap
			t.Id = this.KeywordMap[i].Id
			t.Rank = this.KeywordMap[i].Rank
			t.Word = this.KeywordMap[i].Word

			this.KeywordMap[i].Id = this.KeywordMap[maxIndex].Id
			this.KeywordMap[i].Rank = this.KeywordMap[maxIndex].Rank
			this.KeywordMap[i].Word = this.KeywordMap[maxIndex].Word

			this.KeywordMap[maxIndex].Id = t.Id
			this.KeywordMap[maxIndex].Rank = t.Rank
			this.KeywordMap[maxIndex].Word = t.Word
		}

	}
	return nil
}

// 对每组全局关键词内部的条件关键词做排序操作
func (this *Tpl) initKeywordCondKeyword() error {

	// 处理全局关键词里面的关键词排序
	for k, v := range this.Keyword {
		for index, cond := range v.Conds {
			for _, keyword := range cond.Keyword {

				// 解析关键词，如果关键词格式错误，返回错误信息
				wordRank := strings.Split(keyword, "#")
				if len(wordRank) != 2 {
					return errors.New("关键词格式不对:" + keyword)
				}
				word := wordRank[0]
				rank, err := strconv.Atoi(wordRank[1])
				if err != nil {
					return errors.New("关键词格式不对:" + keyword)
				}

				this.Keyword[k].KeywordMap = append(this.Keyword[k].KeywordMap, *&CondKeywordMap{
					Word:  word,
					Index: index,
					Rank:  rank,
				})
			}
		}

		// 排序条件关键词
		l := len(this.Keyword[k].KeywordMap)
		for i := 0; i < l-1; i++ {
			max := this.Keyword[k].KeywordMap[i].Rank
			maxIndex := i
			// 从未排序的元素中找到最大的
			for j := i + 1; j < l; j++ {
				if this.Keyword[k].KeywordMap[j].Rank > max {
					max = this.Keyword[k].KeywordMap[j].Rank
					maxIndex = j
				}
			}

			// 如果找到更大的，交换
			if maxIndex != i {
				var t CondKeywordMap
				t.Index = this.Keyword[k].KeywordMap[i].Index
				t.Rank = this.KeywordMap[i].Rank
				t.Word = this.KeywordMap[i].Word

				this.Keyword[k].KeywordMap[i].Index = this.Keyword[k].KeywordMap[maxIndex].Index
				this.Keyword[k].KeywordMap[i].Rank = this.Keyword[k].KeywordMap[maxIndex].Rank
				this.Keyword[k].KeywordMap[i].Word = this.Keyword[k].KeywordMap[maxIndex].Word

				this.Keyword[k].KeywordMap[maxIndex].Index = t.Index
				this.Keyword[k].KeywordMap[maxIndex].Rank = t.Rank
				this.Keyword[k].KeywordMap[maxIndex].Word = t.Word
			}

		}
	}
	return nil
}

// 对每个条件流程的关键词做排序操作
func (this *Tpl) initFlowKeyword() error {

	// 处理全局关键词里面的关键词排序
	for k, v := range this.Flow {
		// 不是条件流程就没有关键词，不需要处理
		if v.Section.Type != CONDITION {
			continue
		}
		for index, cond := range v.Section.Conds {
			for _, keyword := range cond.Keyword {

				// 解析关键词，如果关键词格式错误，返回错误信息
				wordRank := strings.Split(keyword, "#")
				if len(wordRank) != 2 {
					return errors.New("关键词格式不对:" + keyword)
				}
				word := wordRank[0]
				rank, err := strconv.Atoi(wordRank[1])
				if err != nil {
					return errors.New("关键词格式不对:" + keyword)
				}

				this.Flow[k].Section.KeywordMap = append(this.Flow[k].Section.KeywordMap, *&CondKeywordMap{
					Word:  word,
					Index: index,
					Rank:  rank,
				})
			}
		}

		// 排序条件关键词
		l := len(this.Flow[k].Section.KeywordMap)
		for i := 0; i < l-1; i++ {
			max := this.Flow[k].Section.KeywordMap[i].Rank
			maxIndex := i
			// 从未排序的元素中找到最大的
			for j := i + 1; j < l; j++ {
				if this.Flow[k].Section.KeywordMap[j].Rank > max {
					max = this.Flow[k].Section.KeywordMap[j].Rank
					maxIndex = j
				}
			}

			// 如果找到更大的，交换
			if maxIndex != i {
				var t CondKeywordMap
				t.Index = this.Flow[k].Section.KeywordMap[i].Index
				t.Rank = this.Flow[k].Section.KeywordMap[i].Rank
				t.Word = this.Flow[k].Section.KeywordMap[i].Word

				this.Flow[k].Section.KeywordMap[i].Index = this.Flow[k].Section.KeywordMap[maxIndex].Index
				this.Flow[k].Section.KeywordMap[i].Rank = this.Flow[k].Section.KeywordMap[maxIndex].Rank
				this.Flow[k].Section.KeywordMap[i].Word = this.Flow[k].Section.KeywordMap[maxIndex].Word

				this.Flow[k].Section.KeywordMap[maxIndex].Index = t.Index
				this.Flow[k].Section.KeywordMap[maxIndex].Rank = t.Rank
				this.Flow[k].Section.KeywordMap[maxIndex].Word = t.Word
			}
		}
	}
	return nil
}

// 检查流程是否合法，如果不合法，就设置为退出状态，并返回对应错误
func (this *Tpl) validFlow(name string) error {

	if name == "" {
		return errors.New("模板出现空流")
	}

	// 如果模板设计流程有问题，遇到了非结束流程并且有空的下一个流程，或者流程名称错误，导致获取失败。
	// 就退出模板，电话根据这个挂机

	if _, ok := this.Flow[name]; !ok {
		this.isQuit = true
		return errors.New("模板出现空流程或流程名称错误，请检查模板是否合法：" + name)
	}

	return nil
}

func (this *Tpl) getFlow(name string) (*Flow, error) {
	err := this.validFlow(name)
	if err != nil {
		return nil, err
	}
	return this.Flow[name], nil
}

func (this Tpl) choiceVoice(voices []string, choice string) (voice *Voice, err error) {
	voice = &Voice{}
	var ok bool
	if len(voices) > 0 {
		if len(voices) == 1 { // 如果只有一条语音，直接选用第一条
			voice, ok = this.Voice[voices[0]]
			if !ok {
				err = errors.New("找不到语音:" + voices[0])
			}
		} else if choice == "random" {
			// 随机选择一条语音
			hash := voices[rand.Intn(len(voices))]
			voice, ok = this.Voice[hash]
			if !ok {
				err = errors.New("找不到语音:" + hash)
			}
		} else {
			// 如果该key对应的声音存在，就是用这个声音；否则使用第一条
			if v, ok := this.Voice[choice]; ok {
				voice = v
			} else {
				err = errors.New("找不到语音:" + choice)
			}
		}
		return
	}
	err = errors.New("没有语音")
	return
}

func (this *Tpl) HookKeyword(word string) bool {
	for _, k := range this.KeywordMap {
		if !strings.Contains(word, k.Word) {
			continue
		}
		_, ok := this.Keyword[k.Id]
		// 如果是一个不存在的关键词，说明模板有错，返回错误信息
		if !ok {
			return false
		}
		return true
	}

	return false
}

// match 返回false则表示没有匹配到关键词，不对流程产生任何影响
// next 返回下一步流程
// err 报错信息
func (this *Tpl) hookKeyword(opt *Operation) (match bool, next string, err error) {

	for _, k := range this.KeywordMap {
		// 不匹配就下一个
		if !strings.Contains(opt.Text, k.Word) {
			continue
		}

		keyword, ok := this.Keyword[k.Id]
		// 如果是一个不存在的关键词，说明模板有错，返回错误信息
		if !ok {
			err = errors.New("key为该名字的关键词不存在:" + k.Id)
			return
		}

		match = true

		next = keyword.Next

		// 如果有语音文件，就发送语音文件
		if len(keyword.Voice) > 0 {
			var voice *Voice
			voice, err = this.choiceVoice(keyword.Voice, keyword.Choice)
			if err != nil {
				return
			}
			this.Person.SetListen(&Operation{
				Opt:   KEYWORD,
				Text:  voice.Text,
				Sound: voice.LocalPcm,
			})
			this.Report.AddKeyword(k.Id, voice.Text, voice.Hash, keyword.Type)
		}

		// 如果有条件，就等待客户说话
		if len(keyword.Conds) > 0 {

			// 机器人一直听，直到到人类说话
			for !this.isQuit && !this.Ai.IsListened() {
				time.Sleep(time.Millisecond * 10)
			}
			opt2 := this.Ai.Listen()

			if opt2 == nil {
				err = errors.New("获取机器人说话失败，可能是通话已结束")
				return
			}

			this.Report.AddHuman(opt2.Text, opt2.Sound)
			// 如果匹配了全局，或者出错，直接返回
			var next1 string
			match, next1, err = this.hookKeyword(opt2)
			if match || err != nil {
				next = next1
				return
			}

			for _, v := range keyword.KeywordMap {
				// 如果有匹配到关键词中的条件之一，则跳转到对应流程执行
				if strings.Contains(opt2.Text, v.Word) {
					next = keyword.Conds[v.Index].To
					return
				}
			}
		}

		if next == RETURN || next == NEXT {
			// 机器人一直听，直到到人类说话
			for !this.isQuit && !this.Ai.IsListened() {
				time.Sleep(time.Millisecond * 10)
			}
			opt2 := this.Ai.Listen()

			if opt2 == nil {
				err = errors.New("获取机器人说话失败，可能是通话已结束")
				return
			}

			this.Report.AddHuman(opt2.Text, opt2.Sound)
			// 如果匹配了全局，或者出错，直接返回
			var next1 string
			match, next1, err = this.hookKeyword(opt2)
			if match || err != nil {
				next = next1
				return
			}
		}

		break
	}

	return
}

// 根据拨打结果分析出客户分类
func (this *Tpl) DoType() {

	lastType := 0
	// 统计人类说话次数
	n := 0
	for _, node := range this.Report.Nodes {
		// 类型是人类说话
		if node.Type == 0 {
			n++
		}
		lastType = node.UserType
	}

	this.Report.Time = this.Report.TimeDisconnect.Sub(this.Report.TimeConnect)
	if this.Report.Time.Seconds() > 60*100 {
		this.Report.Time = 0
	}

	// status 电话拨打状态，0，空号。1，未打通/未接等。2，拒接。3，客户挂机，4，机器人挂机
	if this.Report.Status == 1 {
		this.Report.Type = 3 // 拒接
	} else if this.Report.Status == 2 {
		this.Report.Type = 1 // 未接
	} else if this.Report.Status == 4 {
		this.Report.Type = lastType // 如果是机器人挂机，以结束流程的类型为最终分类
		if lastType == 0 {          // 如果结束流程没指定用户类型，就设置为a类
			this.Report.Type = 5
		}
	} else if this.Report.Time == 0 { // 没接通
		this.Report.Type = 1 // 未接通
	} else if n == 0 {
		this.Report.Type = 10 // 如果客户没说一句话就挂了，直接f类
	} else if n > 3 {
		this.Report.Type = 6 // 客户说话超过3句  b类客户
	} else if this.Report.Time > 30 {
		this.Report.Type = 7 // 没挂机交流超过30秒的，c类客户
	} else if n > 0 && n <= 3 {
		this.Report.Type = 8 // 1-3句之间的 d类客户
	}

}
