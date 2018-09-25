package main

import (
	"gofs/fs"
	"time"
	"gofs/service"
	"os"
	"os/signal"
	"fmt"
)

//func main() {
//	go Api.Handle()
//	time.Sleep(time.Second*3)
//	for {
//		tasks := Api.Tasks()
//		for id, _ := range tasks {
//			Api.TaskUser(id)
//		}
//		time.Sleep(time.Millisecond * 10)
//	}
//	for {
//		time.Sleep(time.Second)
//	}
//}

func main() {
	call := fs.NewCall("13101907101@192.168.4.102", &service.EndPoint{}, time.Minute*5)
	fs.Fs.MakeSimCall("19a", call)
	for {
		time.Sleep(time.Second)
	}

	waitClose()
}


func waitClose(){
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\n 收到终端信号，停止服务... \n")
			cleanup()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
func cleanup() {
	fmt.Println("清理...\n")
}