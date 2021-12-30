/// @author xcsoft<contact@xcsoft.top>

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"serverStatus/config"
	"serverStatus/proc"
	"sync"

	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var connectingDown chan bool //服务器是否掉线
var connected chan bool      //标识 是否连接到服务器
var interrupt chan os.Signal
var lock sync.Mutex

// 定义命令行参数 > 服务器信息
var cliToken string
var cliTag string
var cliHost string
var cliDuration int64

func main() {

	// 解析命令行
	tag, err := os.Hostname() //获取主机名 用于default tag
	if err != nil {
		tag = "unknown"
	}
	flag.StringVar(&cliToken, "token", "", "Server token")
	flag.StringVar(&cliHost, "host", "", "Server ip:port, ex: 127.0.0.1:8282")
	flag.StringVar(&cliTag, "tag", tag, "Server tag")
	flag.Int64Var(&cliDuration, "duration", 5000, "Data send Duration, Unit:ms")

	flag.Parse()
	if flag.NArg() > 0 || cliToken == "" || cliHost == "" { //未知信息
		flag.Usage()
		os.Exit(0)
	}

	// 解析命令行END

	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

Exit:
	for {
		connected = make(chan bool)
		connectingDown = make(chan bool)

		var wsconn *websocket.Conn
		var err error

		log.Println("尝试连接到服务器")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		go func(_ context.Context) {
			defer cancel()
			wsconn, _, err = websocket.DefaultDialer.Dial("ws://"+cliHost, nil)
			if err == nil {
				connected <- true
			}
		}(ctx)

	ConnectedDone: // 防止连接超时
		for {
			select {
			case <-connected:
				break ConnectedDone
			case <-ctx.Done():
				err = errors.New("连接超时")
				break ConnectedDone
			case <-interrupt:
				log.Println("Received SIGINT signal. Closing all pending connections and exiting...")
				break Exit
			}
		}

		if err != nil { //检测连接信息
			log.Println("无法连接到服务器,将在5秒后尝试重连 > ", err)
			select {
			case <-time.After(time.Duration(5) * time.Second):
				continue
			case <-interrupt:
				log.Println("Received SIGINT signal.Exiting...")
				break Exit
			}
		}

		defer wsconn.Close()

		//连接成功
		conn := config.WsConn{
			Conn: wsconn,
			Lock: &lock,
			Down: &connectingDown,
		}

		go receiveHandler(&conn)

		data, _ := json.Marshal(config.Login{
			Type:     "login",
			Platform: "server",
			Tag:      cliTag,
			Token:    cliToken,
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
				log.Println("与服务器失去连接,2秒后重连")
				select {
				case <-time.After(time.Duration(2) * time.Second):
					break mainLoop
				case <-interrupt:
					log.Println("Received SIGINT signal. Closing all pending connections and exiting...")
					_ = conn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

					break Exit
				}
			case <-interrupt:
				log.Println("Received SIGINT signal. Closing all pending connections and exiting...")
				_ = conn.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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
			log.Println("连接已建立")
			// 主动推送信息 >
			go proc.GetServerInfo(conn, proc.GetBaseInfo(), cliDuration)
		} else if re["type"] == "ping" {
			go proc.Ping(conn)
		} else if re["type"] == "invalid_token" {
			log.Fatal("错误的TOKEN,请检查token是否正确")
		}
	}
}
