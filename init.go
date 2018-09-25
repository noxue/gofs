package main

import (
	"runtime"
	"gofs/config"
	"flag"
	"gofs/api"
	"gofs/fs"
)

func init(){
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err:=config.InitConfig("./config.json");err!=nil{
		panic(err)
	}
	api.InitApi()
	fs.InitFs()
}
