package client

import (
	"Mirror/swagger"
	"Mirror/tools"
	"Mirror/websocket"
	"fmt"
	"net/http"
)

func (c *Client) UpdateLeverage(symbol string, leverage float64) {

	InfoLogger.Printf("Updating leverage of %s to %f for client %s\n", symbol, leverage, c.ApiKey)


	//c.SwaggerError(err, response)

	var res swagger.Position
	var response *http.Response
	var err error

	L :for {
		res, response, err = c.Rest.PositionApi.PositionUpdateLeverage(symbol, leverage)
		switch c.SwaggerError(err, response) {
		case 0:
			break L
		case 1:
			continue L
		case 2:
			c.CloseConnection()
			break L
		case 3:
			fmt.Println("Restart the bot")
			break L
		}

	}

	c.positionsLock.Lock()
	defer c.positionsLock.Unlock()

	for idx := range c.activePositions {
		if res.Symbol.Value == c.activePositions[idx].Symbol.Value {
			c.activePositions[idx].Leverage.Value = res.Leverage.Value
		}
	}

	//tools.CheckErr(err)
}

func (c *Client) ActivePositions() websocket.PositionSlice {

	InfoLogger.Printf("Fetching active positons for client %s\n", c.ApiKey)

	c.positionsLock.Lock()
	defer c.positionsLock.Unlock()
	return c.activePositions
}

func (c *Client) TransferMargin(symbol string, margin int) {

	InfoLogger.Printf("Transfering margin on %s by %d for client %s\n", symbol, margin, c.ApiKey)

	res, _, err := c.Rest.PositionApi.PositionTransferIsolatedMargin(symbol, margin)
	tools.CheckErr(err)

	c.positionsLock.Lock()
	defer c.positionsLock.Unlock()

	for idx := range c.activePositions {
		if c.activePositions[idx].Symbol.Value == res.Symbol.Value {
			c.activePositions[idx].PosMargin.Value = res.PosMargin.Value
		}
	}

}
