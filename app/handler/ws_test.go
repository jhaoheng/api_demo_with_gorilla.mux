package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func Test_WebsocketConnection(t *testing.T) {

	s := httptest.NewServer(http.HandlerFunc(NewWebsocketConnection()))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	//
	go send_msg()
	ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if _, p, err := ws.ReadMessage(); err != nil {
		t.Fatalf("%v", err)
	} else {
		fmt.Println(string(p))
	}
}

func send_msg() {
	wsmsg := WSErrMessage{
		Account: "max",
		Err:     fmt.Errorf("error").Error(),
		From:    "anyway",
		At:      time.Now(),
	}
	WSChannel <- wsmsg
}
