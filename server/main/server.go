package main

import (
	"flag"
	"github.com/weizhe0422/TCPServerWithGolang/server"
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
		err error
	)

	initArgs()

	if err = server.InitConfig(confFilePath); err != nil {
		log.Fatal("failed to load configuration: ", err.Error())
		goto ERR
	}
	log.Println("Initial configuration success")

	server.InitTCPServer()
	log.Println("Initial TCP server success")

	server.InitApiServer()
	log.Println("Initial API server success")

	go server.G_ApiServer.StartToService()
	server.G_TCPServer.StartToService()

	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	log.Fatalln(err.Error())
}
