package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/zaf/agi"
)

func auth(myAgi *agi.Session, query url.Values) {
	var rep agi.Reply

	phone := query.Get("phone")

	val, err := nats_js.KvGetString(phone)
	switch {
	case err == nats.ErrKeyNotFound:
		return
	case err != nil:
		log.Printf("[ERROR] StreamFile: %+v\n", err)
		return
	}

	if val == "" {
		log.Println("[WARNING] val empty")
		return
	}

	var otp_slice []string
	for _, v := range val {
		otp_slice = append(otp_slice, "digits/"+string(v))
	}

	otp := strings.Join(otp_slice, "&")
	if otp == "" {
		log.Println("[WARNING] otp empty")
		return
	}

	rep, err = myAgi.SetVariable("OTP", otp)
	if err != nil || rep.Res == -1 {
		log.Printf("SetVariable: %+v\n", err)
		return
	}

	rep, err = myAgi.Exec("Goto", "taxo-phone-auth-code,s,1")
	if err != nil || rep.Res == -1 {
		log.Printf("Goto: %+v\n", err)
		return
	}
}
