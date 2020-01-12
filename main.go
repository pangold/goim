package main

import (
	"gitlab.com/pangold/goim/api"
	"gitlab.com/pangold/goim/config"
	"log"
)

func main() {
	log.Println("IM service start running ...")
	conf := config.NewYaml("config/config.yml").ReadConfig()
	s := api.NewApiServer(*conf)
	//
	// dispatcher := system.NewDispatchServer(conf.Back.Dispatch)
	// s.ResetDispatcher(dispatcher)
	// s.ResetSyncSession(dispatcher)
	s.Run()
}
