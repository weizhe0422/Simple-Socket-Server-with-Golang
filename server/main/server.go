package main

import (
	"flag"
	"github.com/weizhe0422/Simple-Socket-Server-with-Golang/server"
	"log"
	"time"
)

var (
	confFilePath string
)

func initArgs() {
	flag.StringVar(&confFilePath, "config", "./server.json", "Config file path")
	flag.Parse()
}

func main() {
	var (
		err    error
		apiSvr *server.ApiServer
	)

	initArgs()

	if err = server.InitConfig(confFilePath); err != nil {
		log.Fatal("failed to load configuration: ", err.Error())
		goto ERR
	}
	log.Println("Initial configuration success")

	server.InitTCPServer()
	log.Println("Initial TCP server success")

	apiSvr = server.InitApiServer()
	log.Println("Initial API server success")

	go apiSvr.StartToService()
	server.G_TCPServer.StartToService()

	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	log.Fatalln(err.Error())
}
