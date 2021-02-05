package config

import (
	"encoding/json"
	"github.com/adi1382/bitmex-mirror-cli/client"
	"github.com/adi1382/bitmex-mirror-cli/tools"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func SocketResponseDistributor(c <-chan []byte, RestartCounter *int32, wg *sync.WaitGroup) {
L:
	for {
		time.Sleep(time.Nanosecond)

		select {
		case message := <-c:
			if atomic.LoadInt32(RestartCounter) > 0 {
				client.AllClientsLock.Lock()
				for _, v := range client.AllClients {
					v.DropConnection()
				}
				client.AllClients = make([]*client.Client, 0, 5)
				client.AllClientsLock.Unlock()
				atomic.AddInt32(RestartCounter, 1)
				client.InfoLogger.Println("Distributor closed")
				break L
			}

			var u []interface{}

			err := json.Unmarshal(message, &u)

			tools.CheckErr(err)

			key := u[1].(string)

			socketTopic := u[2].(string)

			go func() {
				client.AllClientsLock.Lock()
				for _, v := range client.AllClients {
					if v.WebsocketTopic == "hostAccount" {
						continue
					}
					if !(strings.Contains(string(message), `"table":"instrument"`) || !strings.Contains(string(message), "table")) {
						if socketTopic == "hostAccount" {
							go v.HostUpdatePush(&message)
						}
					}
				}
				client.AllClientsLock.Unlock()
			}()

			go func() {
				client.AllClientsLock.Lock()
				for _, v := range client.AllClients {
					if v.ApiKey == key {
						v.Push(&message)
						break
					}
				}
				client.AllClientsLock.Unlock()
			}()
		default:
			if atomic.LoadInt32(RestartCounter) > 0 {
				client.AllClientsLock.Lock()
				for _, v := range client.AllClients {
					v.DropConnection()
				}
				client.AllClients = make([]*client.Client, 0, 5)
				client.AllClientsLock.Unlock()
				atomic.AddInt32(RestartCounter, 1)
				client.InfoLogger.Println("Distributor closed")
				break L
			}

		}
	}
	wg.Done()
}
