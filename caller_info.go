package main

import (
	"log"
	"net/url"
	"strings"

	tm "github.com/ros-tel/taximaster/common_api"
	tmt "github.com/ros-tel/taximaster/tm_tapi"
	"github.com/zaf/agi"
)

func caller_info(myAgi *agi.Session) {
	get_info(myAgi, myAgi.Env["callerid"])
}

func phone_info(myAgi *agi.Session, query url.Values) {
	get_info(myAgi, query.Get("phone"))
}

func get_info(myAgi *agi.Session, phone string) {
	var rep agi.Reply

	tresp, err := tm_tapi_client.GetInfoByPhone(
		tmt.GetInfoByPhoneRequest{
			Phone:  phone,
			Fields: "PHONE_TYPE-PHONE_SYSTEM_CATEGORY-CATEGORYID-CLIENT_GROUP_ID-IS_PRIOR-CREW_ID-CREW_GROUP_ID-DRIVER_PHONE-ORDER_ID-ORDER_STATE-DRIVER_TIMECOUNT-SOUND_COLOR-SOUND_MARK-GOSNUMBER-CAR_COLOR-CAR_MARK-CREATION_WAY",
		},
	)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)

		orders, err := common_api_client.GetCurrentOrders(tm.GetCurrentOrdersRequest{})
		if err != nil {
			log.Printf("[ERROR] %+v\n", err)
			return
		}

		for _, order := range orders.Orders {
			if order.PhoneToDial == phone {
				tresp, err = tm_tapi_client.GetInfoByPhone(
					tmt.GetInfoByPhoneRequest{
						Phone:  order.Phone,
						Fields: "PHONE_TYPE-PHONE_SYSTEM_CATEGORY-CATEGORYID-CLIENT_GROUP_ID-IS_PRIOR-CREW_ID-CREW_GROUP_ID-DRIVER_PHONE-ORDER_ID-ORDER_STATE-DRIVER_TIMECOUNT-SOUND_COLOR-SOUND_MARK-GOSNUMBER-CAR_COLOR-CAR_MARK-CREATION_WAY",
					},
				)
				if err != nil {
					log.Printf("[ERROR] %+v\n", err)
					return
				}
			}
		}
	}

	// неизвестный номер
	if tresp.PhoneType == 0 {
		log.Printf("[INFO] Number unknown: %+v\n", tresp)
		return
	}

	// в черном списке
	if tresp.PhoneSystemCategory == 1 {
		rep, err = myAgi.Exec("Goto", "busy-call,s,1")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}
	}

	rep, err = myAgi.SetVariable("CATEGORYID", tresp.CategoryID)
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] SetVariable: %+v\n", err)
		return
	}

	rep, err = myAgi.SetVariable("CLIENT_GROUP_ID", tresp.ClientGroupID)
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] SetVariable: %+v\n", err)
		return
	}

	if tresp.PhoneType == 1 || tresp.PhoneType == 5 || tresp.PhoneType == 6 {
		rep, err = myAgi.SetVariable("IS_DRIVER", 1)
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] SetVariable: %+v\n", err)
			return
		}

		rep, err = myAgi.SetVariable("CREW_ID", tresp.CrewID)
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] SetVariable: %+v\n", err)
			return
		}

		rep, err = myAgi.SetVariable("CREW_GROUP_ID", tresp.CrewGroupID)
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] SetVariable: %+v\n", err)
			return
		}
	} else if tresp.PhoneType == 2 || tresp.PhoneType == 3 || tresp.PhoneType == 4 {
		if tresp.OrderID != 0 {
			rep, err = myAgi.SetVariable("ID_ZAKAZ", tresp.OrderID)
			if err != nil || rep.Res == -1 {
				log.Printf("SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("ORDER_ID", tresp.OrderID)
			if err != nil || rep.Res == -1 {
				log.Printf("SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("IS_PRIOR", tresp.IsPrior)
			if err != nil || rep.Res == -1 {
				log.Printf("SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("ORDER_STATE", tresp.OrderState)
			if err != nil || rep.Res == -1 {
				log.Printf("SetVariable: %+v\n", err)
				return
			}

			confirm := false
			in_place := false
			in_car := false
			for _, v := range config.StateConfirm {
				if v == tresp.OrderState {
					confirm = true
					break
				}
			}
			for _, v := range config.StateInPlace {
				if v == tresp.OrderState {
					confirm = true
					in_place = true
					break
				}
			}
			for _, v := range config.StateInCar {
				if v == tresp.OrderState {
					confirm = true
					in_car = true
					break
				}
			}

			if confirm {
				rep, err = myAgi.SetVariable("ORDER_CONFIRM", 1)
				if err != nil || rep.Res == -1 {
					log.Printf("SetVariable: %+v\n", err)
					return
				}

			}
			if in_place {
				rep, err = myAgi.SetVariable("IN_PLACE", 1)
				if err != nil || rep.Res == -1 {
					log.Printf("SetVariable: %+v\n", err)
					return
				}
			}
			if in_car {
				rep, err = myAgi.SetVariable("IN_CAR", 1)
				if err != nil || rep.Res == -1 {
					log.Printf("SetVariable: %+v\n", err)
					return
				}
			}

			rep, err = myAgi.SetVariable("CREW_ID", tresp.CrewID)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("CREW_GROUP_ID", tresp.CrewGroupID)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("DRIVER_PHONE", tresp.DriverPhone)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("CREATION_WAY", tresp.CreationWay)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] SetVariable: %+v\n", err)
				return
			}

			rep, err = myAgi.SetVariable("DRIVER_TIMECOUNT", tresp.DriverTimeCount)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] SetVariable: %+v\n", err)
				return
			}

			if tresp.CarColor != "" && tresp.CarMark != "" {
				if color, ok := config.Colors[strings.ToLower(tresp.CarColor)]; ok {
					rep, err = myAgi.SetVariable("SOUND_COLOR", color)
					if err != nil || rep.Res == -1 {
						log.Printf("[ERROR] SetVariable: %+v\n", err)
						return
					}
				}

				if mark, ok := config.CarMarks[strings.ToLower(tresp.CarMark)]; ok {
					rep, err = myAgi.SetVariable("SOUND_MARK", mark)
					if err != nil || rep.Res == -1 {
						log.Printf("[ERROR] SetVariable: %+v\n", err)
						return
					}
				}

				if tresp.GosNumber != "" {
					gosNumber := carNumberToSay(tresp.GosNumber)
					if *debug {
						log.Printf("[DEBUG] gosNumber %s", gosNumber)
					}
					rep, err = myAgi.SetVariable("GOSNUMBER", gosNumber)
					if err != nil || rep.Res == -1 {
						log.Printf("[ERROR] SetVariable: %+v\n", err)
						return
					}
				}
			}
		}
	}
}
