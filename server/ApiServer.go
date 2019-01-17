package server

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpSvr  *http.Server
	Method   string
	Address  string
	Port     int
	Listener net.Listener
	StopCh   chan error
}

var (
	G_ApiServer *ApiServer
)

func InitApiServer() {
	var (
		apiServer *ApiServer
	)

	apiServer = &ApiServer{
		httpSvr: &http.Server{
			ReadTimeout:  time.Duration(G_Config.ApiSvrReadTimeOut) * time.Millisecond,
			WriteTimeout: time.Duration(G_Config.ApiSvrWriteTimeOut) * time.Millisecond,
		},
		Method:  G_Config.ConnectMethod,
		Address: G_Config.ServerAddress,
		Port:    G_Config.SocketPort,
		StopCh:  make(chan error),
	}

	G_ApiServer = apiServer
}

func (a *ApiServer) StartToService() (err error) {
	var (
		listener net.Listener
		mux      *http.ServeMux
	)
	if listener, err = a.CreateListener(); err != nil {
		return
	}
	a.Listener = listener
	log.Println("create API server listener success")

	mux = createHandleFunc(G_Config.ServerStatusPath, chkServerStatus)
	a.httpSvr.Handler = mux
	log.Println("create API server HandleFunc success")

	log.Println("start to API server service")
	a.httpSvr.Serve(a.Listener)

	return
}

func (a *ApiServer) CreateListener() (listener net.Listener, err error) {

	if listener, err = net.Listen(G_Config.ConnectMethod, G_Config.ServerAddress+":"+strconv.Itoa(G_Config.HttpPort)); err != nil {
		log.Fatal("failed to create a listener:", err.Error())
	}
	log.Println("[API SERVER] start " + G_Config.ConnectMethod + " at " + G_Config.ServerAddress + ":" + strconv.Itoa(G_Config.HttpPort))
	return
}

func (a *ApiServer) Stop(reason string) {
	a.StopCh <- errors.New(reason)
}

func createHandleFunc(routerPath string, handlFunc func(http.ResponseWriter, *http.Request)) (mux *http.ServeMux) {

	mux = http.NewServeMux()
	mux.HandleFunc(routerPath, handlFunc)

	return
}

func chkServerStatus(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte(strconv.Itoa(G_TCPServer.GetConnsCount())))
}
