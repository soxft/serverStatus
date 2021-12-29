/// @author xcsoft<contact@xcsoft.top>

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"serverStatus-client/config"
	"serverStatus-client/proc"
	"sync"

	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var connectingDown chan bool //服务器是否爆炸
var interrupt chan os.Signal
var lock sync.Mutex

// 服务器信息
var TOKEN string = "xcsoft"
var TAG string = "MBP"
var SERVER string = "10.11.11.11:8282"

func main() {

Exit:
	for {
		interrupt = make(chan os.Signal)
		connectingDown = make(chan bool)
		signal.Notify(interrupt, os.Interrupt)

		log.Println("尝试连接到服务器")
		wsconn, _, err := websocket.DefaultDialer.Dial("ws://"+SERVER, nil)
		if err != nil {
			log.Println("无法连接到服务器,将在10秒后尝试重连 > ", err)
			select {
			case <-time.After(time.Duration(1) * time.Millisecond * 5000):
				continue
			case <-interrupt:
				log.Println("Received SIGINT signal.Exiting...")
				break Exit
			}
		}

		defer wsconn.Close()

		conn := config.WsConn{
			Conn: wsconn,
			Lock: &lock,
			Down: &connectingDown,
		}

		go receiveHandler(&conn)

		log.Println("尝试进行服务器认证")
		data, _ := json.Marshal(config.Login{
			Type:     "login",
			Platform: "server",
			Tag:      TAG,
			Token:    TOKEN,
		})
		err = conn.Conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("登录到服务器失败")
		}

	mainLoop:
		for {
			select {
			case <-time.After(time.Duration(1) * time.Millisecond * 15000):
				//心跳
				conn.Lock.Lock()
				err := conn.Conn.WriteMessage(websocket.TextMessage, []byte("{\"type\": \"ping\"}"))
				conn.Lock.Unlock()
				if err != nil {
					log.Println("Error during writing to websocket:", err)
					break mainLoop
				}
			case <-connectingDown:
				log.Println("与服务器失去连接,5秒后尝试重连")
				select {
				case <-time.After(time.Duration(1) * time.Millisecond * 5000):
					break mainLoop
				case <-interrupt:
					log.Println("Received SIGINT signal. Closing all pending connections and exiting...")
					conn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

					break Exit
				}
			case <-interrupt:
				log.Println("Received SIGINT signal. Closing all pending connections and exiting...")
				conn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

				break Exit
			}
		}
	}
}

// 接受处理Handler
func receiveHandler(conn *config.WsConn) {
	for {
		msgType, msg, err := conn.Conn.ReadMessage()
		if msgType == -1 {
			close(connectingDown) //掉线 重连
			return
		}
		if err != nil {
			log.Println("Error in receive:", err)
			close(connectingDown) //重连
			return
		}

		jsonData, _ := ioutil.ReadAll(strings.NewReader(string(msg)))
		var re map[string]interface{}
		err = json.Unmarshal(jsonData, &re)
		if err != nil {
			log.Println("dataErr")
		}

		if re["type"] == "login_success" {
			log.Println("认证成功")
			// 服务器主动推送信息 >

			// 主动推送 服务器信息
			go proc.GetServerInfo(conn)
		} else if re["type"] == "ping" {
			go proc.Ping(conn)
		}
	}
}
