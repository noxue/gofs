package main

import (
	"runtime"
	config "gofs/config"
	"gofs/api"
	"flag"
)

func init(){
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err:=config.InitConfig("./config.json");err!=nil{
		panic(err)
	}
	Api  = api.New(config.Config.Api.Url,config.Config.Api.Origin,config.Config.Api.AppId,config.Config.Api.Key)
}
