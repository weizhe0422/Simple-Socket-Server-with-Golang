package main

import (
	"flag"
	"fmt"
	"github.com/weizhe0422/TCPServerWithGolang/server"
	"log"
	"net"
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
		err      error
		listener net.Listener
	)

	initArgs()

	if err = server.InitConfig(confFilePath); err != nil {
		log.Fatal("failed to load configuration: ",err.Error())
		goto ERR
	}
	log.Println("Initial configuration success")

	if err = server.InitTCPServer(); err != nil {
		log.Fatal("failed to initial TCP server: ")
		goto ERR
	}
	log.Println("Initial TCP server success")

	if listener, err = server.G_TCPServer.CreateListener(); err != nil {
		goto ERR
	}
	log.Println("Create TCP listener success")

	log.Println("Start to accept request and do action...")
	for {
		server.G_TCPServer.ListenAndAction(listener,doReceiveMessage)
	}

ERR:
	log.Fatalln(err.Error())
}

func doReceiveMessage(conn net.Conn) {
	var (
		msgBuf    []byte
		msgLength int
		err       error
		sess *server.Session

	)

	sess = server.NewSession(&conn)
	server.G_TCPServer.Sessions.Store(sess.GetSessionID(),sess)

	defer func(){
		conn.Close()
		server.G_TCPServer.Sessions.Delete(sess.GetSessionID())
	}()

	log.Println(server.G_Config.ReceiveBuffer)

	for {
		msgBuf = make([]byte, server.G_Config.ReceiveBuffer)
		if msgLength, err = conn.Read(msgBuf); err != nil {
			log.Fatalln("failed to read message: ", err.Error())
			continue
		}
		fmt.Println("Received message: ", string(msgBuf[:msgLength]))
		fmt.Println("GetConnsCount: ",server.G_TCPServer.GetConnsCount())
	}
}