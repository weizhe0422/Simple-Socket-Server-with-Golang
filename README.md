# sSSG - simple Socket Server with Golang

sSSG is a socket server written in Go (Golang). It features a configurable light-weight with dashboard to check the server status. If you want to build a socket server with no more limits, then you can try this.

## Requirements
Golang 1.9 and above

## Installation
`go get -u -v github.com/weizhe0422/Simple-Socket-Server-with-Golang`

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

