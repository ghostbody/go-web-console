package sshclient

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type SSHWriter struct {
	Ws *websocket.Conn
}

func (writer *SSHWriter) Write(p []byte) (int, error) {
	err := writer.Ws.WriteMessage(websocket.BinaryMessage, p)
	fmt.Println(string(p[:]))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type SSHReader struct {
	Ws *websocket.Conn
}

func (reader *SSHReader) Read(p []byte) (int, error) {
	// TODO: ignore mt
	_, message, err := reader.Ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	fmt.Println(string(message[:]))
	copy(p, message)
	return len(message), nil
}
