package client

import (
	"Mirror/swagger"
	"Mirror/websocket"
	"fmt"
)

func (c *Client) CurrentMargin() websocket.MarginSlice {
	InfoLogger.Println("Fetching Current margin for client ", c.ApiKey)
	c.marginLock.Lock()
	defer c.marginLock.Unlock()
	return c.currentMargin
}

func (c *Client) RestMargin() float32 {
	InfoLogger.Println("Fetching Current margin for client ", c.ApiKey)
	var currency swagger.UserGetMarginOpts

L:
	for {
		Margin, response, err := c.Rest.UserApi.UserGetMargin(&currency)
		switch c.SwaggerError(err, response) {
		case 0:
			return Margin.MarginBalance.Value
		case 1:
			continue L
		case 2:
			fmt.Println("Remove the current sub client")
			c.CloseConnection()
			return -404
			//break function
		case 3:
			fmt.Println("Restart the bot")
			return -404
		}

	}

}
