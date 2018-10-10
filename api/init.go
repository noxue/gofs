package api

import (
	"time"
	"gofs/config"
)

var TaskApi *Api

func InitApi() {
	ac := config.Config.Api
	TaskApi = New(ac.Url, ac.Origin, ac.AppId, ac.Key)
	go TaskApi.Handle()

	time.Sleep(time.Second)
}
