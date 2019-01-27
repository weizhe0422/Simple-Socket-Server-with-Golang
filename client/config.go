package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	ConnectMethod  string        `json:"connectMethod"`
	ClientAddress  string        `json:"clientAddress"`
	ConnectionPort int           `json:"connectionPort"`
	ConnectTimeOut time.Duration `json:"connectTimeOut"`
	ConnectAlive   time.Duration `json:"connectAlive"`
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
