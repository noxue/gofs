package config

import (
	"encoding/json"
	"io/ioutil"
	"errors"
)

var Config *AppConfig

func InitConfig(configFile string) error {
	var err error
	Config, err = NewFromFile(configFile)
	if err != nil {
		return err
	}
	Config.Asr.AppId = "5adf1c1e"
	Config.Asr.Key = "8a413009f6cfa9346692736688361bfa"
	return nil
}

type _Api struct {
	AppId  string `json:"appid"`
	Key    string `json:"key"`
	Url    string `json:"url"`
	Origin string `json:"origin"`
}

type _Local struct {
	RobotVoicePath string `json:"robot_voice_path"`
	UserVoicePath  string `json:"user_voice_path"`
}

type _Asr struct {
	AppId string `json:"appid"`
	Key   string `json:"key"`
}

type _Fs struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

type AppConfig struct {
	Api   _Api   `json:"api"`
	Local _Local `json:"local"`
	Asr   _Asr   `json:"asr"`
	Fs    _Fs    `json:"fs"`
}

func New(text string) (*AppConfig, error) {
	var config AppConfig
	if err := json.Unmarshal([]byte(text), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func NewFromFile(filepath string) (*AppConfig, error) {
	if filepath == "" {
		return nil, errors.New("config filename is empty")
	}

	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return New(string(b))
}
