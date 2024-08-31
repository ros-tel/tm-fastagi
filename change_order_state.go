package main

import (
	"log"
	"strconv"

	tm "github.com/ros-tel/taximaster/common_api"
)

// Смена состояния заказа
func changeOrderState(order_id string, need_state int) bool {
	orderId, err := strconv.Atoi(order_id)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return false
	}

	_, err = common_api_client.UpdateOrder(
		tm.UpdateOrderRequest{
			OrderID: orderId,
			StateID: need_state,
		},
	)
	if err != nil {
		log.Printf("[WARNING] %+v\n", err)
		return false
	}

	return true
}
