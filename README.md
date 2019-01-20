# sSSG - simple Socket Server with Golang

sSSG is a socket server written in Go (Golang). It features a configurable light-weight with dashboard to check the server status. If you want to build a socket server with no more limits, then you can try this.

## Requirements
Golang 1.9 and above

## Installation & start up
### 1. Installation
1. `go get -u -v github.com/weizhe0422/Simple-Socket-Server-with-Golang`

### 2. Server
1. check the ./server/main/server.json, and modify the setting value to meet your situation
2. type `cd ./main` to swtich to main function file : `server.go`
3. type `go run server.go` to start up the server to liesten

### 3. Client
1. check the ./client/main/client.json, and modify the setting value to meet your situation
2. type `cd ./main` to swtich to main function file : `client.go`
3. type `go run client.go` to start up the client
4. type the message even multi-line, and use `quit` to send to server

## Confifuration
### Server
1. connectMethod: Method that server supply, ex: `tcp` or `udp`.
2. serverAddress: IP Address that server host, ex: `localhost`.
3. socketPort: Port that server supply, ex: `5000`.
4. receiveBuffer: Receive buffer to reveive message, the unit is byte, ex: `512`.
5. httpPort: HTTP server port to supply external API to `GET` server status, ex: `4000`.
6. serverStatusPath: Route path to supply external API to `GET` server status, ex: `/server/status`.
7. apiSvrReadTimeOut: API server time out to read, the uite is millisecond, ex: `5000`.
8. apiSvrWriteTimeOut: API server time out to write, the uite is millisecond, ex: `5000`.
9. rateLimitPerSecond: connection count per second, ex: `30`.
10. rateLimitBuffer: buffer for rate limit, ex: `1`.
11. webRoot: folder path that static HTML page saved, ex: `./webPage`.
  
### client
1. connectMethod: Method to connect to server, ex. `tcp` or `udp`.
2. clientAddress: Server address that want to connect, ex: `localhost`.
3. connectionPort: Server address port.

## Features
![SystemStructure](https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/master/img/SystemStructure.png)
### Server
- [X] Server can serve multiple connections at the same time.
- [X] Configurable setting.

<img src="https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/ConfigurationSetting.png" width="650" height="450" alt="ConfigurationSetting"/>

- [X] Support rate limit machanism to limit connection request count per second
- [X] Dashboard to display the server status include 
 - Current total connection count
 - Session status with every session id
    - Request Count: Count the total request count
    - Request Rate: Requests per second
    - Time per request: Time per request
    - Detail history
![ServerStatusDashboard](https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/ServerStatusDashboard.png)
- [X] Can also get the status raw dat with HTTP `GET` method.
<img src="https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/RespWithReqCnt.png" width="450" height="650" alt="RespWithReqCnt"/>

- [X] Gracfully disconnect with client
<img src="https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/ClientBye.png" width="400" height="250" alt="ClientBye"/>

![ServerCloseConnection](https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/ServerCloseConnection.png)

- [ ] Use etcd to save the configuration setting, then you can change it with remote and works immediately.
- [ ] Support log to save into Mongo.

### Client
- [X] Client can input multi-line text and send to server with `quit` command.
<img src="https://github.com/weizhe0422/Simple-Socket-Server-with-Golang/blob/develop/img/ClintMultilineInput.png" width="400" height="300" alt="ClintMultilineInput"/>

- [ ] Web form interface for user to input message

