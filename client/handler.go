package client

import (
	"crypto/rand"
	"fmt"
	"github.com/adi1382/bitmex-mirror-cli/bitmex"
	"github.com/adi1382/bitmex-mirror-cli/swagger"
	"github.com/adi1382/bitmex-mirror-cli/websocket"
	"log"
	"strings"
	"time"
)

func (c *Client) OrderHandler() {

	defer func() {
		InfoLogger.Println("Order Handler Closed for client ", c.ApiKey)
		//fmt.Println("Order Handler Closed for client ", c.ApiKey)
	}()

	//fmt.Println("Order Handler Started for client ", c.ApiKey)

	c.calibrate()

	calibrateBool := true
	calibrateBoolReset := time.Now().Add(time.Second * time.Duration(c.calibrationTime))
	calibrateBoolClose := false

	go func() {

		defer func() {
			InfoLogger.Println("Calibrator timer closed for client ", c.ApiKey)
			//fmt.Println("Calibrator timer closed for client ", c.ApiKey)
		}()

		//fmt.Println("Calibrator timer Started for client ", c.ApiKey)

		for {
			time.Sleep(time.Nanosecond)
			if calibrateBoolClose {
				break
			}

			if !c.RunningStatus() {
				break
			}

			if calibrateBool == false && calibrateBoolReset.Unix() < time.Now().Unix() {
				calibrateBoolReset = time.Now().Add(time.Second * time.Duration(c.calibrationTime))
				calibrateBool = true
			}

			//if calibrateBool == false {
			//	time.Sleep(time.Second*time.Duration(c.calibrationTime))
			//	calibrateBool = true
			//}
		}
	}()

	for {

		if !c.RunningStatus() {
			calibrateBoolClose = true
			break
		}

		time.Sleep(time.Nanosecond)
		select {

		case message := <-c.hostUpdatesFetcher:
			c.mirroring(&message, &calibrateBoolReset, &calibrateBool)
			continue
		default:
			if calibrateBool {
				c.calibrate()
				calibrateBool = false
				continue
			}
			continue
		}
	}
}

func (c *Client) mirroring(message *[]byte, calibrateBoolReset *time.Time, calibrateBool *bool) {

	InfoLogger.Println("Starting Mirror for client ", c.ApiKey)

	strResponse := string(*message)
	prefix := fmt.Sprintf(`[0,"%s","%s",`, c.hostClient.ApiKey, "hostAccount")
	suffix := fmt.Sprintf("]")
	strResponse = strings.TrimPrefix(strResponse, prefix)
	strResponse = strings.TrimSuffix(strResponse, suffix)

	var ratio float32
	InfoLogger.Println("Calculating ratio on client ", c.ApiKey)
	if c.BalanceProportion {
		ratio = c.GetMarginBalance() / c.hostClient.GetMarginBalance()
	} else {
		ratio = c.FixedRatio
	}

	response, table := bitmex.DecodeMessage([]byte(strResponse))

	InfoLogger.Println("Manipulating ", table, " from mirror on client ", c.ApiKey)

	// Potential Race Condition when fetching
	if table == "order" {
		orderResponse, ok := response.(bitmex.OrderResponse)

		if !ok {
			ErrorLogger.Println("Invalid Interface Conversion of ", orderResponse)
			panic("Invalid Conversion")
		}

		if orderResponse.Action == "insert" {

			InfoLogger.Println("New Order Inserted for Client ", c.ApiKey)

			//fmt.Println("Host Margin Balance: ", hostClient.GetMarginBalance())
			//fmt.Println("Sub Margin Balance: ", subClient.GetMarginBalance())

			orders := make([]map[string]interface{}, 0, 5)

			for _, h := range orderResponse.Data {

				ord := make(map[string]interface{})
				ord["clOrdID"] = random() + h.OrderID.Value[8:]
				if h.Symbol.Valid {
					ord["symbol"] = h.Symbol.Value
				}
				if h.Side.Valid {
					ord["side"] = h.Side.Value
				}
				if h.LeavesQty.Valid {
					ord["orderQty"] = int(h.LeavesQty.Value * ratio)
				}
				if h.Price.Valid {
					ord["price"] = h.Price.Value
				}
				if h.DisplayQty.Valid {
					ord["displayQty"] = int(h.DisplayQty.Value * ratio)
				}
				if h.StopPx.Valid {
					ord["stopPx"] = h.StopPx.Value
				}
				if h.PegOffsetValue.Valid {
					ord["pegOffsetValue"] = h.PegOffsetValue.Value
				}
				if h.PegPriceType.Valid {
					ord["pegPriceType"] = h.PegPriceType.Value
				}
				if h.OrdType.Valid {
					ord["ordType"] = h.OrdType.Value
				}
				if h.TimeInForce.Valid {
					ord["timeInForce"] = h.TimeInForce.Value
				}
				if h.ExecInst.Valid {
					ord["execInst"] = h.ExecInst.Value
				}
				if h.Text.Valid {
					ord["text"] = h.Text.Value
				}
				orders = append(orders, ord)
			}

			var symbols []string

			for _, order := range orders {
				if !isIn(order["symbol"].(string), symbols) {
					symbols = append(symbols, order["symbol"].(string))
				}
			}

			for _, symbol := range symbols {
				placeNewOrders := make([]interface{}, 0, 5)
				for _, order := range orders {
					if order["symbol"].(string) == symbol {
						placeNewOrders = append(placeNewOrders, order)
					}
				}

				InfoLogger.Println("New Order placed for symbol ", symbol, " on client ", c.ApiKey)

				c.OrderNewBulk(&placeNewOrders)

			}

		} else if orderResponse.Action == "update" {
			InfoLogger.Println("New Order update received for Client ", c.ApiKey)

			amendOrders := make([]map[string]interface{}, 0, 5)

			activeOrders := c.ActiveOrders()

			var toCancel []string
			for _, h := range orderResponse.Data {

				if h.OrdStatus.Valid {
					if h.OrdStatus.Value == "Filled" || h.OrdStatus.Value == "PartiallyFilled" {
						InfoLogger.Println("Order " + h.OrdStatus.Value)
						*calibrateBool = false
						*calibrateBoolReset = time.Now().Add(time.Second * time.Duration(c.LimitFilledTimeout))
					}
				}

				if !h.OrdStatus.Valid {
					if h.Price.Valid || h.StopPx.Valid || h.LeavesQty.Valid || h.PegOffsetValue.Valid {

						InfoLogger.Println("Amended order detected for Client ", c.ApiKey)

						subOrders := getSubOrder(h.OrderID.Value, activeOrders)

						if len(subOrders) == 0 {
							continue
						}

						subOrder := subOrders[0]

						if len(subOrders) > 1 {
							for _, v := range subOrders[1:] {
								toCancel = append(toCancel, v.OrderID.Value)
							}
						}

						if (h.OrdType.Value == "StopLimit" || h.OrdType.Value == "LimitIfTouched") &&
							h.Triggered.Value != "" {

							if h.Price.Value != subOrder.Price.Value ||
								int(h.LeavesQty.Value*ratio) != int(subOrder.LeavesQty.Value) {

								amend := make(map[string]interface{})

								amend["symbol"] = subOrder.Symbol.Value
								amend["orderID"] = subOrder.OrderID.Value
								if h.Price.Valid {
									amend["price"] = h.Price.Value
								}
								if h.LeavesQty.Valid {
									amend["orderQty"] = int(h.LeavesQty.Value * ratio)
								}

								amendOrders = append(amendOrders, amend)
							}

						} else {
							if h.Price.Value != subOrder.Price.Value ||
								int(h.LeavesQty.Value*ratio) != int(subOrder.LeavesQty.Value) ||
								h.StopPx.Value != subOrder.StopPx.Value ||
								h.PegOffsetValue.Value != subOrder.PegOffsetValue.Value {

								amend := make(map[string]interface{})

								amend["symbol"] = subOrder.Symbol.Value
								amend["orderID"] = subOrder.OrderID.Value
								if h.Price.Valid {
									amend["price"] = h.Price.Value
								}
								if h.LeavesQty.Valid {
									amend["orderQty"] = int(h.LeavesQty.Value * ratio)
								}
								if h.StopPx.Valid {
									amend["stopPx"] = h.StopPx.Value
								}
								if h.PegOffsetValue.Valid {
									amend["pegOffsetValue"] = h.PegOffsetValue.Value
								}
								amendOrders = append(amendOrders, amend)
							}
						}

					}
				} else if h.OrdStatus.Valid {

					subOrders := getSubOrder(h.OrderID.Value, activeOrders)

					if len(subOrders) == 0 {
						continue
					}
					subOrder := subOrders[0]

					if h.OrdStatus.Value == "Canceled" {
						toCancel = append(toCancel, subOrder.OrderID.Value)
					}

					if len(subOrders) > 1 {
						for _, v := range subOrders[1:] {
							toCancel = append(toCancel, v.OrderID.Value)
						}
					}
				}
			}

			symbols := make([]string, 0, 5)
			for _, order := range amendOrders {
				if !isIn(order["symbol"].(string), symbols) {
					symbols = append(symbols, order["symbol"].(string))
				}
			}
			for _, symbol := range symbols {
				amendOldOrders := make([]interface{}, 0, 5)
				for _, order := range amendOrders {
					if order["symbol"].(string) == symbol {
						_, ok := order["symbol"]
						if ok {
							delete(order, "symbol")
							amendOldOrders = append(amendOldOrders, order)
						}
						order["symbol"] = symbol
					}
				}

				InfoLogger.Println("Order amended for symbol ", symbol, " on client ", c.ApiKey)

				c.OrderAmendBulk(&amendOldOrders)
			}

			InfoLogger.Println(len(toCancel), " Orders canceled for client ", c.ApiKey)
			c.OrderCancelBulk(&toCancel)
		}

	} else if table == "position" {

		InfoLogger.Println("Position Update from mirror for client ", c.ApiKey)

		positionResponse := response.(bitmex.PositionResponse)

		if positionResponse.Action == "update" {

			for _, v := range positionResponse.Data {

				if v.CrossMargin.Valid {
					if v.CrossMargin.Value {
						c.UpdateLeverage(v.Symbol.Value, 0)

					} else if v.Leverage.Valid {
						c.UpdateLeverage(v.Symbol.Value, v.Leverage.Value)

					} else {
						for _, activePosition := range c.hostClient.ActivePositions() {
							if activePosition.Symbol.Value == v.Symbol.Value {
								c.UpdateLeverage(v.Symbol.Value, activePosition.Leverage.Value)
							}
						}
					}
				} else if v.Leverage.Valid {
					c.UpdateLeverage(v.Symbol.Value, v.Leverage.Value)

				}
			}
		}

	}
}

func getSubOrder(id string, orders websocket.OrderSlice) []swagger.Order {
	returnValue := make([]swagger.Order, 0, 5)
	for _, order := range orders {
		if order.ClOrdID.Value[8:] == id[8:] {
			returnValue = append(returnValue, order)
		}
	}

	return returnValue
}

func random() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid[0:8]
}

func isIn(key string, xi []string) bool {
	for _, v := range xi {
		if key == v {
			return true
		}
	}
	return false
}
