package ws

/*
	c,_, err := ws.NewWs(u)
*/

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WS struct {
	u        url.URL
	WS       *websocket.Conn
	IsClosed bool
	mu       *sync.Mutex
}

func NewWs(u url.URL) (ws *WS, res *http.Response, err error) {
	ws = &WS{}
	ws.u = u
	ws.mu = &sync.Mutex{}
	ws.WS, res, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	return
}

func (ws *WS) Close() {
	ws.Close()
}

func (ws *WS) ReDial() {
	if ws.IsClosed {
		return
	}
	ws.IsClosed = true
	no := 0
	var err error
	for {
		no++
		time.Sleep(1 * 1e9)
		log.Printf("redial no(%d)(%s)\n", no, time.Now())
		ws.WS, _, err = websocket.DefaultDialer.Dial(ws.u.String(), nil)
		if err != nil {
			continue
		}
		break
	}
	ws.IsClosed = false
	return
}

func (ws *WS) ReadMessage() (messageType int, message []byte, err error) {
	if ws.IsClosed {
		for {
			time.Sleep(2 * 1e9)
			if ws.IsClosed {
				continue
			}
			break
		}
	}
	messageType, message, err = ws.WS.ReadMessage()
	if err != nil {
		log.Printf("ws read err(%v)\n", err)
		ws.ReDial()
		return ws.ReadMessage()
	}
	return
}

func (ws *WS) WriteMessage(messageType int, data []byte) (err error) {
	if ws.IsClosed {
		for {
			time.Sleep(2 * 1e9)
			if ws.IsClosed {
				continue
			}
			break
		}
	}
	err = ws.WS.WriteMessage(messageType, data)
	if err != nil {
		log.Printf("ws write err(%v)\n", err)
		ws.ReDial()
		return ws.WriteMessage(messageType, data)
	}
	return
}
