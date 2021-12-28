package config

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	Conn *websocket.Conn
	Lock *sync.Mutex
	Down *chan bool
}
