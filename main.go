package main

import (
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn"
	"log"
)

func main() {
	log.Println("IM service start running ...")
	s := conn.NewServer(config.Conf)
	s.Run()
}
