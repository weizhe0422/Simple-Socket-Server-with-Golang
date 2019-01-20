# sSSG - simple Socket Server with Golang

sSSG is a socket server written in Go (Golang). It features a configurable light-weight with dashboard to check the server status. If you want to build a socket server with no more limits, then you can try this.

## Requirements
Golang 1.9 and above

## Installation
`go get -u -v github.com/weizhe0422/Simple-Socket-Server-with-Golang`

## Features
### Server
- [X] Server can serve multiple connections at the same time.
- [X] Configurable setting 
- [X] Support rate limit machanism to limit connection request count per second
- [X] Dashboard to display the server status include 
 - Current total connection count
 - Session status with every session id
    - Request Count: Count the total request count
    - Request Rate: Requests per second
    - Time per request: Time per request
    - Detail history
- [X] Can also get the status raw dat with HTTP `GET` method.
![RespWithReqCnt](https://github.com/weizhe0422/TCPServerWithGolang/blob/develop/img/RespWithReqCnt.png)
- [X] Gracfully disconnect with client
- [ ] Use etcd to save the configuration setting, then you can change it with remote and works immediately.
- [ ] Support log to save into Mongo.

### Client
- [X] Client can input multi-line text and send to server with `quit` command.
- [ ] Web form interface for user to input message

