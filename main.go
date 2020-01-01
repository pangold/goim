package main

import (
	"gitlab.com/pangold/goim/api"
	"gitlab.com/pangold/goim/config"
	"log"
)

func main() {
	log.Println("IM service start running ...")
	s := api.NewApiServer(config.Conf)
	s.Run()
}
