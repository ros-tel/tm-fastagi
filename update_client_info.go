package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	tm "github.com/ros-tel/taximaster/common_api"
	tmt "github.com/ros-tel/taximaster/tm_tapi"
)

func setClientGroup(query url.Values) {
	phone := query.Get("phone")
	group := query.Get("group")

	groupId, err := strconv.Atoi(group)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	tresp, err := tm_tapi_client.GetInfoByPhone(
		tmt.GetInfoByPhoneRequest{
			Phone:  phone,
			Fields: "CLIENT_ID-CLIENT_GROUP_ID",
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	registerClientRequest := tm.RegisterClientRequest{
		ClientGroup: groupId,
	}

	// нет клиента
	if tresp.ClientID == 0 {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		registerClientRequest.Name = "Без имени"
		registerClientRequest.Login = phone + fmt.Sprintf("%06d", r1.Int63())[:8]
		registerClientRequest.Password = phone + fmt.Sprintf("%06d", r1.Int63())[:8]
		registerClientRequest.Phones = phone
		_, err = common_api_client.RegisterClient(registerClientRequest)
		if err != nil {
			log.Printf("[ERROR] %+v\n", err)
			return
		}
		return
	}

	// уже есть
	if groupId == tresp.ClientGroupID {
		return
	}

	_, err = common_api_client.UpdateClientInfo(
		tm.UpdateClientInfoRequest{
			ClientGroupID: groupId,
			ClientID:      tresp.ClientID,
		},
	)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
}
