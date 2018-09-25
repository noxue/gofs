package fs

import (
	"github.com/golang/glog"
	"gofs/config"
)

var Fs *Phone

func InitFs(){
	var err error
	Fs, err = New(config.Config.Fs.Host, config.Config.Fs.Port, config.Config.Fs.Password, 10)
	if err != nil {
		glog.Error(err)
	}
	go Fs.Handle()
}
