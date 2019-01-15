package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/weizhe0422/TCPServerWithGolang/client"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	confFilePath string
)

func initArgs() {
	flag.StringVar(&confFilePath, "config", "./client.json", "Config file path")
	flag.Parse()
}

func main() {
	var (
		err  error
		conn net.Conn
		inputReader *bufio.Reader
		clientName, clientNameTrim string
		inputMsg, inputMsgTrim string
		respString *bytes.Buffer
	)

	initArgs()

	if err = client.InitConfig(confFilePath); err != nil {
		log.Fatal("failed to load configuration: ", err.Error())
		goto ERR
	}
	log.Println("Initial configuration success")


	if err = client.InitTCPServer(client.G_Config.ConnectMethod, client.G_Config.ConnectionPort); err != nil {
		log.Fatal("failed to initial TCP server: ")
		goto ERR
	}
	log.Println("Initial TCP server success")

	if conn, err = client.G_TCPServer.CreateDialer(client.G_Config.ClientAddress); err != nil {
		goto ERR
	}
	log.Println("Success to dial to " + client.G_Config.ClientAddress + ":" + strconv.Itoa(client.G_Config.ConnectionPort))


	inputReader = bufio.NewReader(os.Stdin)
	fmt.Println("Input your name: ")
	if clientName, err = inputReader.ReadString('\n'); err!=nil{
		log.Fatal("failed to get client's name: ")
		return
	}
	clientNameTrim = strings.Trim(clientName,"\r\n")

	for{
		fmt.Println("Start to send message until you type quit to quit:")
		respString = bytes.NewBufferString("")

		for {
			if inputMsg, err = inputReader.ReadString('\n'); err!=nil{
				log.Println("failed to read string from user:", err.Error())
				continue
			}
			inputMsgTrim = strings.Trim(inputMsg,"\r\n")

			if (inputMsgTrim == "quit"){
				break
			}

			respString.Write([]byte(inputMsgTrim))
		}
		conn.Write([]byte(clientNameTrim + ":" + respString.String()))
	}



ERR:
	log.Fatalln(err.Error())
}