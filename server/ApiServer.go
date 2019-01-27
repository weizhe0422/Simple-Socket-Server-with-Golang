package server

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type errHandler func(http.ResponseWriter, *http.Request) error

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

type ApiServer struct {
	httpSvr  *http.Server
	Method   string
	Address  string
	Port     int
	Listener net.Listener
	StopCh   chan error
}

func InitApiServer() (apiServer *ApiServer) {

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

	return apiServer
}

func (a *ApiServer) StartToService() (err error) {
	var (
		listener      net.Listener
		mux           *http.ServeMux
		staticDir     http.Dir
		staticHandler http.Handler
	)
	if listener, err = a.CreateListener(); err != nil {
		return
	}
	a.Listener = listener
	log.Println("create API server listener success")

	mux = createHandleFunc(G_Config.ServerStatusPath, chkServerStatus)
	mux.HandleFunc("/mock", mockExternAPI)

	staticDir = http.Dir(G_Config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))
	log.Println("staticDir", staticDir)

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

func createHandleFunc(routerPath string, handlFunc errHandler) (mux *http.ServeMux) {

	mux = http.NewServeMux()
	mux.HandleFunc(routerPath, errWrapper(handlFunc))

	return
}

func chkServerStatus(resp http.ResponseWriter, req *http.Request) (err error) {
	var (
		respSvrStatus *ServerStatus
		respJson      []byte
	)
	respSvrStatus = &ServerStatus{
		ConnCount:    G_TCPServer.GetConnsCount(),
		SessInfoSumm: G_TCPServer.GetServerSummry(),
		ConnHist:     G_TCPServer.GetConnHistALL(),
	}

	resp.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if respSvrStatus == nil {
		return userError("There is no connection history")
	}

	if respJson, err = json.Marshal(respSvrStatus); err != nil {
		return userError("Failed to convert connection history as JSON format")
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(respJson)
	return nil
}

func mockExternAPI(resp http.ResponseWriter, req *http.Request) {
	var (
		respMsg string
		reqKeys []string
		ok      bool
	)
	if reqKeys, ok = req.URL.Query()["ReceiveMSG"]; !ok || len(reqKeys) < 1 {
		resp.Write([]byte("can not get valud of ReceiveMSG"))
		return
	}
	respMsg = req.RemoteAddr + ":" + reqKeys[0]

	resp.Write([]byte(respMsg))
}

func errWrapper(handler errHandler) func(http.ResponseWriter, *http.Request) {
	var (
		code int
	)
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				log.Print(r)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handler(w, r)
		if userErr, ok := err.(userError); ok {
			http.Error(w, userErr.Message(), http.StatusBadRequest)
			log.Println("user error")
		}

		//code := http.StatusOK
		if err != nil {
			code = http.StatusInternalServerError
		}

		http.Error(w, http.StatusText(code), code)

	}
}
