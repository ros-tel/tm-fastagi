package main

import (
	"log"
	"net/url"

	tmt "github.com/ros-tel/taximaster/tm_tapi"
	"github.com/zaf/agi"
)

func getDriverPhoneByCrew(myAgi *agi.Session, query url.Values) {
	tresp, err := tm_tapi_client.GetDriverPhonesByCrewCode(
		tmt.GetDriverPhonesByCrewCodeRequest{
			CrewCode: query.Get("crew_code"),
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	rep, err := myAgi.SetVariable("DRIVER_PHONE", tresp.MobilePhone)
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] SetVariable: %+v\n", err)
		return
	}
}

func getDriverPhoneByCaller(myAgi *agi.Session, query url.Values) {
	tresp, err := tm_tapi_client.GetInfoByPhone(
		tmt.GetInfoByPhoneRequest{
			Phone:  query.Get("phone"),
			Fields: "DRIVER_PHONE",
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	if tresp.DriverPhone != "" {
		rep, err := myAgi.SetVariable("DRIVER_PHONE", tresp.DriverPhone)
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] SetVariable: %+v\n", err)
			return
		}
	}
}
