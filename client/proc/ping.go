package proc

import (
	"encoding/json"
	"serverStatus-client/config"

	"github.com/gorilla/websocket"
)

// 查询服务器是否在线
func Ping(conn *config.WsConn) {
	serverBaseInfo, _ := json.Marshal(config.Ping{
		Type: "ping",
		Msg:  "pong",
	})

	conn.Lock.Lock() //锁
	conn.Conn.WriteMessage(websocket.TextMessage, serverBaseInfo)
	conn.Lock.Unlock()
}
