package bitmex

import (
	"encoding/json"
	"github.com/adi1382/bitmex-mirror-cli/tools"
	"github.com/adi1382/bitmex-mirror-cli/websocket"
	"log"
	"os"
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

type Response interface {
	getTable() string
}

func (v Res) getTable() string {
	return v.Table
}

func (v OrderResponse) getTable() string {
	return v.Table
}

func (v PositionResponse) getTable() string {
	return v.Table
}

func (v MarginResponse) getTable() string {
	return v.Table
}

type Res struct {
	Success   bool        `json:"success,omitempty"`
	Subscribe string      `json:"subscribe,omitempty"`
	Request   interface{} `json:"request,omitempty"`
	Table     string      `json:"table,omitempty"`
	//Action    string      `json:"action,omitempty"`
	//Data      interface{} `json:"data,omitempty"`
}

type PositionResponse struct {
	Table  string                           `json:"table,omitempty"`
	Action string                           `json:"action,omitempty"`
	Keys   []string                         `json:"keys,omitempty"`
	Data   []websocket.PositionResponseData `json:"data,omitempty"`
}

type OrderResponse struct {
	Table  string                        `json:"table,omitempty"`
	Action string                        `json:"action,omitempty"`
	Keys   []string                      `json:"keys,omitempty"`
	Data   []websocket.OrderResponseData `json:"data,omitempty"`
}

type MarginResponse struct {
	Table  string                         `json:"table,omitempty"`
	Action string                         `json:"action,omitempty"`
	Keys   []string                       `json:"keys,omitempty"`
	Data   []websocket.MarginResponseData `json:"data,omitempty"`
}

func DecodeMessage(message []byte) (Response, string) {

	InfoLogger.Println("Decoding Socket Message")

	var response Response
	var table string

	var res Res
	var positionResponse PositionResponse
	var orderResponse OrderResponse
	var marginResponse MarginResponse
	err := json.Unmarshal(message, &res)
	tools.CheckErr(err)

	table = res.Table

	if res.Table == "position" {
		err = json.Unmarshal(message, &positionResponse)
		tools.CheckErr(err)
		response = positionResponse
	} else if res.Table == "order" {
		err = json.Unmarshal(message, &orderResponse)
		tools.CheckErr(err)
		response = orderResponse
	} else if res.Table == "margin" {
		err = json.Unmarshal(message, &marginResponse)
		tools.CheckErr(err)
		response = marginResponse
	} else {
		err = json.Unmarshal(message, &res)
		tools.CheckErr(err)
		response = res
		table = ""
	}

	InfoLogger.Println("Socket Message Decoded")

	return response, table
}
