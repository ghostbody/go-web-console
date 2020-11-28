package routes

import (
	"go-web-console/exceptions"
	"go-web-console/sshclient"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

func getUpgrader() websocket.Upgrader {
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return upgrader

}

//webSocket请求ping 返回pong
func v1WsPing(c *gin.Context) {
	upgrader := getUpgrader()
	//升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = []byte("pong from server")
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func v1wsLink(c *gin.Context) {

	sshClientID := c.Query("client_id")
	sshClientSecret := c.Query("client_secret")
	if sshClientID == "" {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusBadRequest,
				"ssh client id is required",
			),
		)
		return
	}
	if sshClientSecret == "" {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusBadRequest,
				"ssh client secret is required",
			),
		)
		return
	}

	sshClientManager := sshclient.GetSSHClientsManager()
	client := sshClientManager.GetSSHClientByClientID(sshClientID)
	if client == nil {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusBadRequest,
				"ssh client not found",
			),
		)
		return
	}
	if client.ClientSecret != sshClientSecret {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusForbidden,
				"ssh client secret is not correct",
			),
		)
	}

	// create session
	var session *ssh.Session
	var err error
	if session, err = client.Client.NewSession(); err != nil {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusForbidden,
				"can not create session",
			),
		)
		return
	}

	// TODO: hard code
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// TODO: hard code
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		return
	}

	upGrader := getUpgrader()
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	session.Stdin = &sshclient.SSHReader{
		Ws: ws,
	}
	session.Stdout = &sshclient.SSHWriter{
		Ws: ws,
	}

	if err = session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
		return
	}

	session.Wait()
}
