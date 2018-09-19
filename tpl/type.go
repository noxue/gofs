package tpl

import "sync"

// section类型
const (
	START     = "start"
	CONDITION = "condition"
	KEYWORD   = "keyword"
	END       = "end"
	RETURN    = "return"
	WAIT      = "wait"
	NEXT      = "next"
	OPERATION = "operation"
	QUIET1    = "quiet1"
	QUIET2    = "quiet2"
	NOWORD1   = "noword1"
	NOWORD2   = "noword2"
	NOWORD3   = "noword3"
)

type Operation struct {
	Opt      string // 操作类型
	Name     string // 流程名称或，匹配的关键词
	Text     string // 保存文字
	Sound    string // 对应的语音
	UserType int    // 用户类型
}

// 声音信息
type Voice struct {
	Hash     string `json:"hash"` // 文件hash
	Text     string `json:"text"` // 声音对应的文本
	Path     string `json:"path"` // wav文件的远程路径
	Pcm      string `json:"pcm"`  // pcm声音文件远程路径
	LocalPcm string // 下载远程文件到本地的保存路径
}

// 全局关键词
type Keyword struct {
	Keyword    []string `json:"keyword"`
	Voice      []string `json:"voice"`
	Choice     string   `json:"choice"`
	Conds      []*Cond  `json:"conds"`
	KeywordMap []CondKeywordMap
	Type       int    `json:"type"`
	Next       string `json:"next"`
}

// 用于记录每个关键字的优先级，关键字和键名对应的关系
type KeywordMap struct {
	Word string // 关键字
	Id   string // 关键字在map中的键名
	Rank int    // 优先级权重，就是关键词#后面的数值，数值越大越优先匹配
}

type CondKeywordMap struct {
	Word  string // 关键字
	Index int    // 第几个分支的关键词
	Rank  int    // 优先级，数字越大越靠前，越优先匹配
}

// 条件分支
type Cond struct {
	Keyword []string `json:"keyword"`
	To      string   `json:"to"`
}

type Section struct {
	Type       string   `json:"type"`
	Voice      []string `json:"voice,omitempty"` // tag里面加上omitempy，可以在序列化的时候忽略0值或者空值，因为pass是没有其他信息的
	Choice     string   `json:"choice,omitempty"`
	Conds      []*Cond  `json:"conds,omitempty"`
	KeywordMap []CondKeywordMap
}

type Flow struct {
	Section *Section `json:"section"`
	Hook    bool     `json:"hook"`
	Max     int      `json:"max"`
	Type    int      `json:"type"` // 记录执行到当前节点是什么类型的客户
	Next    string   `json:"next"` // 如果都不命中选择分支，则进入这个节点
}

// 客户分类
type Type struct {
	Name string `json:"name"`
	Rule string `json:"rule"`
}

type Tpl struct {
	Main       string              `json:"main"`
	Flow       map[string]*Flow    `json:"flow"`
	Keyword    map[string]*Keyword `json:"keyword"` // keyword数组元素连起来形成字符串之后，sha1加密后的字符串
	Voice      map[string]*Voice   `json:"voice"`   // 声音文件的hash值作为key
	Type       []Type              `json:"type"`
	KeywordMap []KeywordMap        // 关键词优先级排序以及对应关系

	// 记录是否允许被打断
	// 后期可能支持最多被打断次数，所以用int类型
	// 由于每个任务都不一定要支持打断，所以放在任务发放的时候设定，不写死在模板中，每次执行任务都重新赋值
	// 0:不打断; -1:任意声音打断; -2:全局关键词打断
	Break     int
	BreakChan chan bool
	Cur       string // 当前执行到哪个位置

	VoiceRemoteBaseUrl string // 远程声音文件下载的基础url地址
	LocalBasePath      string // 本地声音文件保存的基础地址

	isQuit    bool // 是否已结束通话
	IsRunning bool // 已经启动，防止多次被调用

	Ai     *Ai
	Person *Person

	Report Report // 记录拨打报告
}

// 机器人和人类基本的状态信息
type BaseStatus struct {
	isSpeaking bool       // 说话状态，true表示正在说话
	listen     *Operation // 保存听到的内容
	isListened bool       // 已经听到内容,表示上面的listen里面是否有内容，true表示有内容
	lock       sync.Mutex
}

// 机器人状态
type Ai struct {
	BaseStatus
}

// 人类状态
type Person struct {
	BaseStatus
}
