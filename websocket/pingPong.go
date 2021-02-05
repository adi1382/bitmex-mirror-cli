package websocket

import (
	"github.com/adi1382/bitmex-mirror-cli/tools"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

func PingPong(conn *websocket.Conn, RestartCounter *int32) {

	InfoLogger.Println("Ping Pong protocol Initiated")

	pongWait := time.Second * 10
	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	tools.WebsocketErr(err, RestartCounter)
	conn.SetPongHandler(func(string) error { err = conn.SetReadDeadline(time.Now().Add(pongWait)); return err })

	go func() {
		for {
			if atomic.LoadInt32(RestartCounter) > 0 {
				//fmt.Println("Ping Pong Closed")
				atomic.AddInt32(RestartCounter, 1)
				break
			}
			err := conn.WriteMessage(9, []byte{})
			if err != nil {
				atomic.AddInt32(RestartCounter, 1)
				break
			}
			resetTime := time.Now().Add(time.Second * 5)
			for {
				time.Sleep(time.Nanosecond)
				if time.Now().Unix() > resetTime.Unix() {
					break
				} else {
					if atomic.LoadInt32(RestartCounter) > 0 {
						break
					}
				}
			}
			//time.Sleep(time.Second*5)
		}
	}()
}
