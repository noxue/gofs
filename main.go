package main

import (
	"os"
	"os/signal"
	"fmt"
	"gofs/api"
	"time"
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
	time.Sleep(time.Second*4)
	go func() {
		for {
			simId := api.TaskApi.GetTaskInfo().GetFreeSim()
			if simId == 0 {
				time.Sleep(time.Millisecond * 200)
				continue
			}
			taskId:=api.TaskApi.GetTaskInfo().GetTaskIdBySimid(simId)
			if taskId == 0 {
				time.Sleep(time.Millisecond * 200)
				continue
			}

			// make a new call success, so set the sim not free
			api.TaskApi.GetTaskInfo().SimFree(simId, false)

			api.TaskApi.SimTaskUser(simId,taskId)
			time.Sleep(time.Millisecond*500)
		}
	}()

	waitClose()
}

func waitClose() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\n 程序结束 \n")
			cleanup()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
func cleanup() {
	api.TaskApi.Close()
}
