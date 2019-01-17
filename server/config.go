package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ConnectMethod      string `json:"connectMethod"`
	ServerAddress      string `json:"serverAddress"`
	SocketPort         int    `json:"socketPort"`
	HttpPort           int    `json:"httpPort"`
	ServerStatusPath   string `json:"serverStatusPath"`
	ReceiveBuffer      int    `json:"receiveBuffer"`
	ApiSvrReadTimeOut  int    `json:"apiSvrReadTimeOut"`
	ApiSvrWriteTimeOut int    `json:"apiSvrWriteTimeOut"`
}

var (
	G_Config *Config
)

func InitConfig(fileName string) (err error) {
	var (
		content []byte
		conf    Config
	)

	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_Config = &conf

	log.Println(G_Config)
	return
}
