package server

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type channelData struct {
	session     *Session
	sessionID   string
	sessionInfo SessionInfo
	content     []byte
}

type ServerStatus struct {
	ConnCount    int
	SessInfoSumm map[string]SessionReqInfo
	ConnHist     map[string][]SessionInfo
}
type SessionReqInfo struct {
	RequestCount int
	RequestRate  float64
	TimePerReq   float64
}
type TCPServer struct {
	Method       string
	Address      string
	Port         int
	Sessions     *sync.Map
	Listener     net.Listener
	Connects     map[string]int
	Limiter      *rate.Limiter
	SvrStatus    *ServerStatus
	SessInfoSumm map[string]SessionReqInfo
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
		Limiter:  rate.NewLimiter(rate.Every(time.Duration(G_Config.RateLimitPerSecond)), G_Config.RateLimitBuffer),
		SvrStatus: &ServerStatus{
			ConnCount: 0,
			ConnHist:  make(map[string][]SessionInfo, 0),
		},
		SessInfoSumm: make(map[string]SessionReqInfo, 0),
	}

	G_TCPServer = tcpSvr

}

func (t *TCPServer) StartToService() {
	var (
		err      error
		listener net.Listener
		ctx      context.Context
	)
	if listener, err = t.CreateListener(); err != nil {
		log.Println("failed to create listener: ", err.Error())
	}
	t.Listener = listener
	log.Println("Create TCP listener success")

	log.Println("Start to accept request and do action...")

	ctx, _ = context.WithCancel(context.TODO())

	for {
		t.Limiter.Wait(ctx)
		t.ListenAndAction(listener, doReceiveMessage)
	}
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

func (t *TCPServer) GetConnHistBySessId(sessionId string) (connHist []SessionInfo) {
	var (
		ok bool
	)

	if connHist, ok = t.SvrStatus.ConnHist[sessionId]; ok {
		return connHist
	}

	return nil
}

func (t *TCPServer) GetConnHistALL() (connHist map[string][]SessionInfo) {
	return t.SvrStatus.ConnHist
}

func (t *TCPServer) GetProcTimeSum(sessionId string) float64 {
	var (
		procTimeSum float64
		infoItem    SessionInfo
	)

	for _, infoItem = range t.GetConnHistBySessId(sessionId) {
		procTimeSum += infoItem.Duration
	}

	return procTimeSum
}

func (t *TCPServer) SetConnHist(sessionId string, data SessionInfo) {
	t.SvrStatus.ConnHist[sessionId] = append(t.SvrStatus.ConnHist[sessionId], data)
}

func (t *TCPServer) UpdateServerSummry(sessionId string, reqCnt int) {
	t.SessInfoSumm[sessionId] = SessionReqInfo{
		RequestCount: reqCnt,
		RequestRate:  float64(reqCnt) / t.GetProcTimeSum(sessionId),
		TimePerReq:   t.GetProcTimeSum(sessionId) / float64(reqCnt),
	}
}

func (t *TCPServer) GetServerSummry() map[string]SessionReqInfo {
	return t.SessInfoSumm
}

func doReceiveMessage(conn net.Conn) {
	var (
		sess         *Session
		sessInfo     SessionInfo
		sessionID    string
		readChannel  chan []byte
		writeChannel chan *channelData
		readData     []byte
		passData     *channelData
	)

	G_TCPServer.Connects[conn.RemoteAddr().String()] = 0

	sess = NewSession(&conn)
	sessionID = sess.GetSessionID()

	G_TCPServer.Sessions.Store(sessionID, sess)
	log.Println("Address(" + conn.RemoteAddr().String() + "):  Dial in: Session ID:" + sess.GetSessionID())

	defer func() {
		conn.Close()
		G_TCPServer.Sessions.Delete(sessionID)
	}()

	readChannel = make(chan []byte, 1024)
	writeChannel = make(chan *channelData, 1024)

	go readCoroutine(conn, readChannel, sessionID, &sessInfo)
	go writeCoroutine(conn, writeChannel)

	for {
		select {
		case readData = <-readChannel:
			if string(readData) == "bye" {
				return
			}

			passData = &channelData{
				session:     sess,
				sessionID:   sessionID,
				sessionInfo: sessInfo,
				content:     readData,
			}

			writeChannel <- passData
		}
	}
}

func readCoroutine(conn net.Conn, readChannel chan []byte, sessionID string, sessInfo *SessionInfo) {
	var (
		msgBuf    []byte
		msgLength int
		err       error
	)

	for {
		msgBuf = make([]byte, G_Config.ReceiveBuffer)
		if msgLength, err = conn.Read(msgBuf); err != nil {
			if err == io.EOF {
				log.Println("Address(" + conn.RemoteAddr().String() + "): Close this connection! Session ID:" + sessionID)
				return
			}
			log.Fatalln("failed to read message: ", err.Error())
			continue
		}
		G_TCPServer.Connects[conn.RemoteAddr().String()]++
		sessInfo.RemoteAddress = conn.RemoteAddr().String()
		sessInfo.ReqTime = time.Now()

		fmt.Println("Received message: ", string(msgBuf[:msgLength]))

		sessInfo.Data = string(msgBuf[:msgLength])
		readChannel <- msgBuf[:msgLength]
	}

}

func writeCoroutine(conn net.Conn, writeChannel chan *channelData) {
	var (
		readData *channelData
	)

	for {
		select {
		case readData = <-writeChannel:
			log.Println("Current Connection Count: ", G_TCPServer.GetConnsCount())
			fmt.Println("Address(" + conn.RemoteAddr().String() + "): " + strconv.Itoa(G_TCPServer.Connects[conn.RemoteAddr().String()]))

			readData.sessionInfo.RespTime = time.Now()
			readData.sessionInfo.Duration = readData.sessionInfo.RespTime.Sub(readData.sessionInfo.ReqTime).Seconds()
			readData.session.SetSetting(readData.sessionID, readData.sessionInfo)
			G_TCPServer.SetConnHist(readData.sessionID, readData.sessionInfo)
			G_TCPServer.UpdateServerSummry(readData.sessionID, G_TCPServer.Connects[conn.RemoteAddr().String()])
			log.Println(readData.session.GetSetting(readData.sessionID))

			//simulate send the request message to another external API
			mockRedirect(string(readData.content))
		}
	}
}

func mockRedirect(message string) {
	var (
		redirectURL  string
		resp         *http.Response
		redirectResp []byte
		err          error
	)

	message = strings.Replace(message, " ", "_", -1)

	redirectURL = "http://" + G_Config.ServerAddress + ":" + strconv.Itoa(G_Config.HttpPort) + "/mock?ReceiveMSG=" + message
	log.Println(redirectURL)
	if resp, err = http.Get(redirectURL); err != nil {
		log.Println("failed to link to ", redirectURL, ":", err.Error())
	}
	if redirectResp, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("failed to get redirect URL response", err.Error())
	}
	log.Println("wait to get the response......")
	time.Sleep(3 * time.Second)
	fmt.Println("Redirect Response Content", string(redirectResp))
}
