package main

import (
	"log"
	"net/url"
	"strconv"

	tm "github.com/ros-tel/taximaster/common_api"
)

func showTmMessage(query url.Values) {
	var err error
	header := query.Get("header")
	text := query.Get("text")
	qtimeout := query.Get("timeout")
	timeout := 0
	if qtimeout != "" {
		timeout, err = strconv.Atoi(qtimeout)
		if err != nil {
			log.Printf("[ERROR] %+v\n", err)
			return
		}
	}

	_, err = common_api_client.ShowTmMessage(
		tm.ShowTmMessageRequest{
			Header:  header,
			Text:    text,
			Timeout: timeout,
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}
}
