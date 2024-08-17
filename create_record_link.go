package main

import (
	"log"
	"net/url"
	"strconv"

	tmt "github.com/ros-tel/taximaster/tm_tapi"
)

func createRecordLink(query url.Values) {
	phone := query.Get("phone")
	call_id := query.Get("call_id")

	callType, err := strconv.Atoi(query.Get("type"))
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	recordLength, err := strconv.Atoi(query.Get("length"))
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	_, err = tm_tapi_client.CreateRecordLink(
		tmt.CreateRecordLinkRequest{
			Phone:        phone,
			RecordDate:   query.Get("date"),
			FilePath:     query.Get("path"),
			CallType:     callType,
			RecordLength: recordLength,
			CallID:       call_id,
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	if query.Get("order_id") == "" {
		v := url.Values{}
		v.Add("PHONE", phone)
		v.Add("RECORD_DATE", query.Get("date"))
		v.Add("FILE_PATH", query.Get("path"))
		v.Add("CALL_TYPE", query.Get("type"))
		v.Add("RECORD_LENGTH", query.Get("length"))
		// меняем параметр CALL_ID
		v.Set("CALL_ID", call_id+"0")
		if err := RedisByteSETEX("record:"+phone, 512, []byte(v.Encode())); err != nil {
			log.Fatal("[ERROR] record", err)
		}
	}
}
