package client

import (
	"Mirror/bitmex"
	"Mirror/swagger"
	"Mirror/websocket"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	_, err := os.Stat("logs")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("logs", 0755)
		if errDir != nil {
			ErrorLogger.Fatal(err)
		}
	}

	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

var AllClients []*Client
var AllClientsLock sync.Mutex

type Client struct {
	active             bool
	marginUpdateTime   int
	BalanceProportion  bool
	FixedRatio         float32
	Test               bool
	marginUpdated      bool
	partials           int
	marginBalance      float32
	LimitFilledTimeout int
	activeOrders       websocket.OrderSlice
	activePositions    websocket.PositionSlice
	currentMargin      websocket.MarginSlice
	ordersLock         sync.Mutex
	positionsLock      sync.Mutex
	marginLock         sync.Mutex
	marginBalanceLock  sync.Mutex
	runningStatusLock  sync.Mutex
	ApiKey             string
	apiSecret          string
	WebsocketTopic     string
	chWriteToWSClient  chan<- interface{}
	chReadFromWSClient chan []byte
	Rest               *swagger.APIClient
	hostClient         *Client
	calibrationTime    int
	hostUpdatesFetcher chan []byte
	restartCounter     *int32
}

//////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////
func (c *Client) CloseConnection() {
	AllClientsLock.Lock()
	defer AllClientsLock.Unlock()

	if len(AllClients) == 0 {
		return
	}

	if c.RunningStatus() {
		InfoLogger.Println("Close connection request initiated for client ", c.ApiKey)
		c.runningStatusLock.Lock()
		c.active = false
		c.runningStatusLock.Unlock()
		var message []interface{}
		message = append(message, 2, c.ApiKey, c.WebsocketTopic)
		c.chWriteToWSClient <- message
		c.chReadFromWSClient <- []byte("quit")
		c.removeCurrentClient()
	}
}

func (c *Client) DropConnection() {
	if c.RunningStatus() {
		InfoLogger.Println("Drop connection request initiated for client ", c.ApiKey)
		c.runningStatusLock.Lock()
		c.active = false
		c.runningStatusLock.Unlock()
		c.chReadFromWSClient <- []byte("quit")
	}
}

//////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////

func (c *Client) WaitForPartial() {
	InfoLogger.Println("Waiting For partials")
	for {
		if c.partials >= 3 {
			break
		}
	}
	InfoLogger.Println("Partials received")
}

func (c *Client) Initialize(apiKey, apiSecret, websocketTopic string,
	ch chan<- interface{}, balanceProportion bool, fixedRatio float32, test bool, marginUpdateTime int,
	hostClient *Client, calibrationTime int, LimitFilledTimeout int, RestartCounter *int32) {

	InfoLogger.Println("New Client Initializing | ", websocketTopic, " | ", apiKey)

	c.restartCounter = RestartCounter
	c.hostUpdatesFetcher = make(chan []byte, 100)
	c.hostClient = hostClient
	c.calibrationTime = calibrationTime
	c.active = true
	c.marginUpdateTime = marginUpdateTime
	c.BalanceProportion = balanceProportion
	c.FixedRatio = fixedRatio
	c.Test = test
	c.ApiKey = apiKey
	c.apiSecret = apiSecret
	c.WebsocketTopic = websocketTopic
	c.chWriteToWSClient = ch
	c.LimitFilledTimeout = LimitFilledTimeout
	c.initiateRest()
	c.StartConnection()
	c.Authorize()
	c.chReadFromWSClient = make(chan []byte, 100)
	AllClientsLock.Lock()
	AllClients = append(AllClients, c)
	AllClientsLock.Unlock()
	go c.marginUpdate()
	go c.dataHandler()

	InfoLogger.Println("New Client Initialized | ", websocketTopic, " | ", apiKey)
}

func (c *Client) marginUpdate() {

	defer func() {
		//fmt.Println("Margin Update Stopped for client ", c.ApiKey)
	}()
	c.marginUpdated = false
	c.WaitForPartial()
	//fmt.Println("Margin Update Started for client ", c.ApiKey)
	for {
		if !c.RunningStatus() {
			break
		}

		marginBalance := c.RestMargin()
		InfoLogger.Println("Margin Balance is: ", marginBalance)

		InfoLogger.Println("Updating margin on ", c.ApiKey)
		c.marginBalanceLock.Lock()
		c.marginLock.Lock()
		//c.marginBalance = c.currentMargin[0].MarginBalance.Value
		c.marginBalance = marginBalance
		c.marginUpdated = true
		c.marginLock.Unlock()
		c.marginBalanceLock.Unlock()
		InfoLogger.Println("Margin updated on ", c.ApiKey)
		//fmt.Println("Margin updated on ", c.ApiKey)
		// Calibration time
		//time.Sleep(time.Second * time.Duration(c.marginUpdateTime))

		time.Sleep(time.Second * 5)

		resetTime := time.Now().Add(time.Second * time.Duration(c.marginUpdateTime))
		for {
			time.Sleep(time.Nanosecond)
			if time.Now().Unix() > resetTime.Unix() {
				break
			} else {
				if !c.RunningStatus() {
					break
				}
			}
		}
	}
}

func (c *Client) GetMarginBalance() float32 {
	InfoLogger.Println("Fetching balance of ", c.ApiKey)
	defer c.marginBalanceLock.Unlock()
	for {
		if c.marginUpdated {
			c.marginBalanceLock.Lock()
			return c.marginBalance
		}
	}
}

func (c *Client) initiateRest() {
	InfoLogger.Println("Initiating Rest object on client ", c.ApiKey)
	if c.Test {
		c.Rest = swagger.NewAPIClient(swagger.NewTestnetConfiguration())
	} else {
		c.Rest = swagger.NewAPIClient(swagger.NewConfiguration())
	}
	c.Rest.InitializeAuth(c.ApiKey, c.apiSecret)
	InfoLogger.Println("Initiated Rest object on client ", c.ApiKey)
}

func (c *Client) StartConnection() {
	InfoLogger.Println("Initiating connection of client ", c.ApiKey, " to sockets.")
	var message []interface{}
	message = append(message, 1, c.ApiKey, c.WebsocketTopic)
	c.chWriteToWSClient <- message
	InfoLogger.Println("Initiated connection of client ", c.ApiKey, " to sockets.")
}

func (c *Client) Authorize() {
	var message []interface{}
	InfoLogger.Println("Authenticating websocket connection of client ", c.ApiKey)
	message = append(message, 0, c.ApiKey, c.WebsocketTopic)
	message = append(message, websocket.GetAuthMessage(c.ApiKey, c.apiSecret))
	c.chWriteToWSClient <- message
}

func (c *Client) Subscribe(tables ...string) {
	var message []interface{}
	command := websocket.Message{Op: "subscribe"}

	for _, v := range tables {
		command.AddArgument(v)
	}

	InfoLogger.Println("Subscribing tables ", tables, " on client ", c.ApiKey)

	message = append(message, 0, c.ApiKey, c.WebsocketTopic)
	message = append(message, command)
	c.chWriteToWSClient <- message
}

func (c *Client) Unsubscribe(tables ...string) {

	InfoLogger.Println("Unsubscribing tables ", tables, " on client ", c.ApiKey)

	var message []interface{}
	command := websocket.Message{Op: "unsubscribe"}

	for _, v := range tables {
		command.AddArgument(v)
	}

	message = append(message, 0, c.ApiKey, c.WebsocketTopic)
	message = append(message, command)
	c.chWriteToWSClient <- message
}

func (c *Client) Push(message *[]byte) {
	c.chReadFromWSClient <- *message
}

func (c *Client) HostUpdatePush(message *[]byte) {
	c.hostUpdatesFetcher <- *message
}

func (c *Client) dataHandler() {
	//fmt.Println("Data Handler started for client ", c.ApiKey)
	defer func() {
		InfoLogger.Println("Data Handler Closed for client ", c.ApiKey)
		//fmt.Println("Data Handler Closed for client ", c.ApiKey)
	}()
	for {

		if !c.RunningStatus() {
			break
		}

		message := <-c.chReadFromWSClient
		InfoLogger.Println("Received new message in Data Handler for client ", c.ApiKey)
		strResponse := string(message)
		if strResponse == "quit" {
			break
		}

		if strings.Contains(string(message), "Access Token expired for subscription") {
			atomic.AddInt32(c.restartCounter, 1)
			ErrorLogger.Println(string(message))
			//fmt.Println(string(message))
		}

		if strings.Contains(string(message), "Invalid API Key") {
			fmt.Println("API key ", c.ApiKey, " is invalid.")
			ErrorLogger.Println(string(message))
			ErrorLogger.Println("API key ", c.ApiKey, " is invalid.")
			if c.WebsocketTopic == "hostAccount" {
				fmt.Println("Host Account API key is Invalid. Closing the bot in 10 seconds.")
				ErrorLogger.Println("Host Account API key is Invalid. Closing the bot in 10 seconds.")
				time.Sleep(time.Second * 10)
				os.Exit(-1)
			} else {
				c.CloseConnection()
			}
		}

		if strings.Contains(string(message), "This key is disabled") {
			ErrorLogger.Println(string(message))
			fmt.Println("API key ", c.ApiKey, " is disabled.")
			ErrorLogger.Println("API key ", c.ApiKey, " is disabled.")
			if c.WebsocketTopic == "hostAccount" {
				fmt.Println("Host Account API key is disabled. Closing the bot in 10 seconds.")
				ErrorLogger.Println("Host Account API key is disabled. Closing the bot in 10 seconds.")
				time.Sleep(time.Second * 10)
				os.Exit(-1)
			} else {
				c.CloseConnection()
			}
		}

		prefix := fmt.Sprintf(`[0,"%s","%s",`, c.ApiKey, c.WebsocketTopic)
		suffix := fmt.Sprintf("]")
		strResponse = strings.TrimPrefix(strResponse, prefix)
		strResponse = strings.TrimSuffix(strResponse, suffix)
		if !strings.Contains(string(message), "table") {
			continue
		}

		response, table := bitmex.DecodeMessage([]byte(strResponse))

		InfoLogger.Println("Manipulating ", table, " table on client ", c.ApiKey)

		// Potential Race Condition when fetching
		if table == "order" {

			orderResponse := response.(bitmex.OrderResponse)

			c.ordersLock.Lock()
			if orderResponse.Action == "partial" {
				c.partials++
				c.activeOrders.OrderPartial(&orderResponse.Data)
			} else if orderResponse.Action == "insert" {
				c.activeOrders.OrderInsert(&orderResponse.Data)
			} else if orderResponse.Action == "update" {
				c.activeOrders.OrderUpdate(&orderResponse.Data)
			} else if orderResponse.Action == "delete" {
				c.activeOrders.OrderDelete(&orderResponse.Data)
			}
			c.ordersLock.Unlock()

		} else if table == "position" {
			positionResponse := response.(bitmex.PositionResponse)

			c.positionsLock.Lock()
			if positionResponse.Action == "partial" {
				c.partials++
				c.activePositions.PositionPartial(&positionResponse.Data)
			} else if positionResponse.Action == "insert" {
				c.activePositions.PositionInsert(&positionResponse.Data)
			} else if positionResponse.Action == "update" {
				c.activePositions.PositionUpdate(&positionResponse.Data)
			} else if positionResponse.Action == "delete" {
				c.activePositions.PositionDelete(&positionResponse.Data)
			}
			c.positionsLock.Unlock()

		} else if table == "margin" {
			marginResponse := response.(bitmex.MarginResponse)

			c.marginLock.Lock()
			if marginResponse.Action == "partial" {
				c.partials++
				c.currentMargin.MarginPartial(&marginResponse.Data)
			} else if marginResponse.Action == "insert" {
				c.currentMargin.MarginInsert(&marginResponse.Data)
			} else if marginResponse.Action == "update" {
				c.currentMargin.MarginUpdate(&marginResponse.Data)
			} else if marginResponse.Action == "delete" {
				c.currentMargin.MarginDelete(&marginResponse.Data)
			}
			c.marginLock.Unlock()
		}
	}
}

func (c *Client) RunningStatus() bool {
	c.runningStatusLock.Lock()
	defer c.runningStatusLock.Unlock()
	return c.active
}

func (c *Client) removeCurrentClient() {

	InfoLogger.Println("Removing client ", c.ApiKey, " from all clients")

	for i := range AllClients {
		if AllClients[i].ApiKey == c.ApiKey {
			AllClients = append(AllClients[:i], AllClients[i+1:]...)
			break
		}
	}

}

//func GetAllClients() []*Client {
//	AllClientsLock.Lock()
//	defer AllClientsLock.Unlock()
//	return AllClients
//}

//func ResetAllClients() {
//	AllClientsLock.Lock()
//	defer AllClientsLock.Unlock()
//
//	fmt.Println("Removing all clients")
//	//for i := range AllClients {
//	//	AllClients[i].CloseConnection()
//	//}
//	fmt.Println("All clients removed")
//
//	AllClients = AllClients[:0]
//}
