package main

import (
	"log"
	"net/url"

	"github.com/zaf/agi"
)

func cancelByOrderID(myAgi *agi.Session, query url.Values) {
	ivrCancelOrder(myAgi, query.Get("order_id"))
}

func ivrCancelOrder(myAgi *agi.Session, order_id string) {
	if changeOrderState(order_id, config.StateCancel) {
		if *debug {
			log.Println("[DEBUG] StreamFile custom/zayavka_otmenena")
		}
		rep, err := myAgi.StreamFile("custom/zayavka_otmenena", "*")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}
	}
	rep, err := myAgi.Hangup()
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] Hangup: %+v\n", err)
		return
	}
}
