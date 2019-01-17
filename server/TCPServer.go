package server

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type TCPServer struct {
	Method   string
	Address  string
	Port     int
	Sessions *sync.Map
	Listener net.Listener
	Connects map[string]int
	Limiter *rate.Limiter
}

var (
	G_TCPServer *TCPServer
)

func InitTCPServer() {
	var (
		tcpSvr *TCPServer
	)


	tcpSvr = &TCPServer{
		Method:   G_Config.ConnectMethod,
		Address:  G_Config.ServerAddress,
		Port:     G_Config.SocketPort,
		Sessions: &sync.Map{},
		Connects: make(map[string]int, 0),
		Limiter: rate.NewLimiter(rate.Every(time.Duration(G_Config.RateLimitPerSecond)),G_Config.RateLimitBuffer),
	}

	G_TCPServer = tcpSvr

}

func (t *TCPServer) StartToService() {
	var (
		err      error
		listener net.Listener
		ctx context.Context
	)
	if listener, err = t.CreateListener(); err != nil {
		goto ERR
	}
	t.Listener = listener
	log.Println("Create TCP listener success")

	log.Println("Start to accept request and do action...")

	ctx, _ = context.WithCancel(context.TODO())
	for {
		t.Limiter.Wait(ctx)
		t.ListenAndAction(listener, doReceiveMessage)
	}

ERR:
	log.Fatalln(err.Error())
}

func (t *TCPServer) CreateListener() (listener net.Listener, err error) {

	if listener, err = net.Listen(G_TCPServer.Method, G_TCPServer.Address+":"+strconv.Itoa(G_TCPServer.Port)); err != nil {
		log.Fatal("failed to create a listener:", err.Error())
	}
	log.Println("[TCP Server] start " + G_TCPServer.Method + " at " + G_TCPServer.Address + ":" + strconv.Itoa(G_TCPServer.Port))
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

func (t *TCPServer) GetConnsCount() int {
	var count int
	t.Sessions.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

func (t *TCPServer) GetConnHist() int {
	var count int
	t.Sessions.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

func doReceiveMessage(conn net.Conn) {
	var (
		msgBuf    []byte
		msgLength int
		err       error
		sess      *Session
	)

	G_TCPServer.Connects[conn.RemoteAddr().String()] = 0

	sess = NewSession(&conn)
	G_TCPServer.Sessions.Store(sess.GetSessionID(), sess)
	log.Println("Address("+conn.RemoteAddr().String()+"):  Dial in: Session ID:"+sess.GetSessionID() )

	defer func() {
		conn.Close()
		G_TCPServer.Sessions.Delete(sess.GetSessionID())
	}()

	for {
		G_TCPServer.Connects[conn.RemoteAddr().String()] ++
		log.Println("111")
		msgBuf = make([]byte, G_Config.ReceiveBuffer)
		if msgLength, err = conn.Read(msgBuf); err != nil {
			if err == io.EOF{
				log.Println("Address("+conn.RemoteAddr().String()+"): Close this connection! Session ID:"+sess.GetSessionID() )
				return
			}
			log.Fatalln("failed to read message: ", err.Error())
			continue
		}

		fmt.Println("Received message: ", string(msgBuf[:msgLength]))
		log.Println("Current Connection Count: ", G_TCPServer.GetConnsCount())
		log.Println("Address("+conn.RemoteAddr().String()+"): "+strconv.Itoa(G_TCPServer.Connects[conn.RemoteAddr().String()]))

	}
}
