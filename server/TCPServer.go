package server

import (
	"log"
	"net"
	"strconv"
)

type TCPServer struct {
	Method string
	Port   int
}

var (
	G_TCPServer *TCPServer
)

func InitTCPServer(ConnMethod string, ConnPort int) (err error) {
	var (
		tcpSvr *TCPServer
	)


	tcpSvr = &TCPServer{
		Method: ConnMethod,
		Port:   ConnPort,
	}

	G_TCPServer = tcpSvr

	return
}

func (t *TCPServer) CreateListener(ListenAddress string) (listener net.Listener, err error) {

	if listener, err = net.Listen(G_TCPServer.Method, ListenAddress+":"+strconv.Itoa(G_TCPServer.Port)); err != nil {
		log.Fatal("failed to create a listener:", err.Error())
	}
	log.Println("start " + G_TCPServer.Method + " at " + ListenAddress + ":" + strconv.Itoa(G_TCPServer.Port))
	return
}

func (t *TCPServer) CreateDialer(ListenAddress string) (listener net.Listener, err error) {

	if listener, err = net.Listen(G_TCPServer.Method, ListenAddress); err != nil {
		log.Fatal("failed to create a listener:", err.Error())
	}
	return
}

func (t *TCPServer) ListenAndAction(listener net.Listener, Action func(conn net.Conn)) (err error) {
	var (
		conn net.Conn
	)

	if conn, err = listener.Accept(); err != nil {
		log.Fatal("failed to accept request: ")
		return
	}

	go Action(conn)

	return
}
