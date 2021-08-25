package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WSErrMessage struct {
	Account string    `json:"account"`
	Err     string    `json:"err"`
	From    string    `json:"from"`
	At      time.Time `json:"at"`
}

var WSConnections int = 0
var WSChannel = make(chan WSErrMessage, 1)
var upgrader = websocket.Upgrader{}

func WebsocketConnection(w http.ResponseWriter, r *http.Request) {
	WSConnections++
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("upgrade:", err)
		return
	}
	defer func() {
		logrus.Info("ws_client connection close")
		WSConnections--
		conn.Close()
	}()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			logrus.Info("ws_read:", err)
			break
		}

		select {
		case msg := <-WSChannel:
			b, _ := json.Marshal(msg)
			err = conn.WriteMessage(messageType, b)
			if err != nil {
				logrus.Info("ws_write:", err)
				goto stop
			}
		default:
		}
	stop:
	}
}
