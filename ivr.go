package main

import (
	"log"
	"strconv"
	"strings"

	tm "github.com/ros-tel/taximaster/common_api"
	tmt "github.com/ros-tel/taximaster/tm_tapi"
	"github.com/zaf/agi"
)

func ivr(myAgi *agi.Session) {
	var rep agi.Reply

	phone := myAgi.Env["callerid"]

	tresp, err := tm_tapi_client.GetInfoByPhone(
		tmt.GetInfoByPhoneRequest{
			Phone:  phone,
			Fields: "PHONE_TYPE-PHONE_SYSTEM_CATEGORY-CATEGORYID-CLIENT_GROUP_ID-DRIVER_PHONE-ORDER_ID-ORDER_STATE-DRIVER_TIMECOUNT-SOUND_COLOR-SOUND_MARK-GOSNUMBER-CAR_COLOR-CAR_MARK",
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
						Fields: "PHONE_TYPE-PHONE_SYSTEM_CATEGORY-CATEGORYID-CLIENT_GROUP_ID-DRIVER_PHONE-ORDER_ID-ORDER_STATE-DRIVER_TIMECOUNT-SOUND_COLOR-SOUND_MARK-GOSNUMBER-CAR_COLOR-CAR_MARK",
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

	// по ID категории
	if tresp.CategoryID == config.BlackPhoneCategoryId {
		rep, err = myAgi.Exec("Goto", "busy-call,s,1")
		if err != nil || rep.Res == -1 {
			log.Printf("Goto: %+v\n", err)
			return
		}
	}
	if tresp.CategoryID == config.WhitePhoneCategoryId {
		rep, err = myAgi.SetVariable("QUEUE_PRIO", "10")
		if err != nil || rep.Res == -1 {
			log.Printf("SetVariable: %+v\n", err)
			return
		}
	}

	// по ID группы клиентов
	if tresp.ClientGroupID == config.BlackClientGroupId {
		rep, err = myAgi.Exec("Goto", "busy-call,s,1")
		if err != nil || rep.Res == -1 {
			log.Printf("Goto: %+v\n", err)
			return
		}
	}
	if tresp.ClientGroupID == config.WhiteClientGroupId {
		rep, err = myAgi.SetVariable("QUEUE_PRIO", "10")
		if err != nil || rep.Res == -1 {
			log.Printf("SetVariable: %+v\n", err)
			return
		}
	}

	if tresp.PhoneType == 1 || tresp.PhoneType == 5 || tresp.PhoneType == 6 {
		rep, err = myAgi.SetVariable("IS_DRIVER", 1)
		if err != nil || rep.Res == -1 {
			log.Printf("SetVariable: %+v\n", err)
			return
		}
	} else if (tresp.PhoneType == 2 || tresp.PhoneType == 3 || tresp.PhoneType == 4) && tresp.OrderID != 0 {
		rep, err = myAgi.SetVariable("ID_ZAKAZ", tresp.OrderID)
		if err != nil || rep.Res == -1 {
			log.Printf("SetVariable: %+v\n", err)
			return
		}
		confirm := false
		in_place := false
		in_car := false
		sound := "custom/obrabotka_perv"
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

		// Уже в машине выходим
		if in_car {
			return
		}

		if confirm {
			if tresp.DriverPhone != "" {
				sound = "custom/obrabotka_full"
			} else if tresp.GosNumber != "" {
				sound = "custom/obrabotka_info"
			}
		}

		rep, err = myAgi.Answer()
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] answer: %+v\n", err)
			return
		}

		if *debug {
			log.Println("[DEBUG] StreamFile silence/1", rep)
		}
		rep, err = myAgi.StreamFile("silence/1", "*")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}

		if confirm && config.DriverInfoBeforeIvr {
			ivrOrderInfo(myAgi, tresp, in_place)
		}

		for count := 0; count < 2; count++ {
			if *debug {
				log.Printf("[DEBUG] GetData %s", sound)
			}
			rep, err = myAgi.GetData(sound, 7000, 1)
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] GetData: %+v\n", err)
				break
			}
			if *debug {
				log.Printf("[DEBUG] %+v\n", rep)
			}

			orderID_str := strconv.Itoa(tresp.OrderID)

			menu := rep.Res
			switch menu {
			case 0: // соединение с оператором
				if *debug {
					log.Println("[DEBUG] StreamFile queue-callswaiting", rep)
				}
				rep, err = myAgi.StreamFile("queue-callswaiting", "*")
				if err != nil || rep.Res == -1 {
					log.Printf("[ERROR] StreamFile: %+v\n", err)
					return
				}
				return
			case 1: // информация о заявке
				if confirm {
					ivrOrderInfo(myAgi, tresp, in_place)

					rep, err = myAgi.Hangup()
					if err != nil || rep.Res == -1 {
						log.Printf("[ERROR] Hangup: %+v\n", err)
						return
					}
				}
				return
			case 2: // соединение с водителем, если заказ подтвержден
				if confirm {
					if tresp.DriverPhone != "" {
						ivrCallToDriver(myAgi, orderID_str, tresp.DriverPhone)
						return
					}
				}

				rep, err = myAgi.Hangup()
				if err != nil || rep.Res == -1 {
					log.Printf("[ERROR] Hangup: %+v\n", err)
					return
				}
				return
			case 3: // отказ от заявки
				// отказ когда водитель уже подтвердил (в конфиге надо или нет) или на месте только через оператора
				if (config.NotCancelInStateConfirm && confirm) || in_place {
					if *debug {
						log.Println("[DEBUG] in_place not cancel")
					}
					return
				}
				ivrCancelOrder(myAgi, orderID_str)
				return
			case 5: // повтор
				if *debug {
					log.Println("[DEBUG] continue")
				}
				continue
			default:
				if *debug {
					log.Println("[DEBUG] break")
				}
				return
			}
		}
	}
}

func ivrOrderInfo(myAgi *agi.Session, tresp tmt.GetInfoByPhoneResponse, in_place bool) {
	if *debug {
		log.Printf("[DEBUG] SayNumber driverTimecount %d", tresp.DriverTimeCount)
	}
	if tresp.DriverTimeCount > 0 {
		if *debug {
			log.Println("[DEBUG] StreamFile taxi/in-t")
		}
		rep, err := myAgi.StreamFile("taxi/in-t", "*")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}

		driverTimeCount_str := strconv.Itoa(tresp.DriverTimeCount)
		if *debug {
			log.Println("[DEBUG] StreamFile taxi/min-" + driverTimeCount_str)
		}

		rep, err = myAgi.StreamFile("taxi/min-"+driverTimeCount_str, "*")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}
	}
	if tresp.CarColor != "" && tresp.CarMark != "" {
		if in_place {
			if *debug {
				log.Println("[DEBUG] StreamFile taxi/car-wait")
			}
			rep, err := myAgi.StreamFile("taxi/car-wait", "*")
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] StreamFile: %+v\n", err)
				return
			}
		} else {
			if *debug {
				log.Println("[DEBUG] StreamFile taxi/car-will-come")
			}
			rep, err := myAgi.StreamFile("taxi/car-will-come", "*")
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] StreamFile: %+v\n", err)
				return
			}
		}
		if color, ok := config.Colors[strings.ToLower(tresp.CarColor)]; ok {
			if *debug {
				log.Println("[DEBUG] StreamFile " + color)
			}
			rep, err := myAgi.StreamFile(color, "*")
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] StreamFile: %+v\n", err)
				return
			}
		}
		if *debug {
			log.Println("[DEBUG] StreamFile silence/03")
		}
		rep, err := myAgi.StreamFile("silence/03", "*")
		if err != nil || rep.Res == -1 {
			log.Printf("[ERROR] StreamFile: %+v\n", err)
			return
		}
		if mark, ok := config.CarMarks[strings.ToLower(tresp.CarMark)]; ok {
			if *debug {
				log.Println("[DEBUG] StreamFile " + mark)
			}
			rep, err = myAgi.StreamFile(mark, "*")
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] StreamFile: %+v\n", err)
				return
			}
		}
		if tresp.GosNumber != "" {
			if *debug {
				log.Println("[DEBUG] StreamFile taxi/number")
			}
			rep, err = myAgi.StreamFile("taxi/number", "*")
			if err != nil || rep.Res == -1 {
				log.Printf("[ERROR] StreamFile: %+v\n", err)
				return
			}
			gosNumber := strings.Split(carNumberToSay(tresp.GosNumber), "&")
			if *debug {
				log.Printf("[DEBUG] SayNumber gosNumber %s", gosNumber)
			}
			for _, v := range gosNumber {
				rep, err = myAgi.StreamFile(v, "*")
				if err != nil || rep.Res == -1 {
					log.Printf("[ERROR] StreamFile: %+v\n", err)
					return
				}
			}
		}
	}
}

func ivrCallToDriver(myAgi *agi.Session, order_id, driver_phone string) {
	if *debug {
		log.Println("[DEBUG] StreamFile priv-introsaved")
	}
	rep, err := myAgi.StreamFile("priv-introsaved", "*")
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] StreamFile: %+v\n", err)
		return
	}
	if *debug {
		log.Println("[DEBUG] Dial to " + driver_phone)
	}
	if config.StateClientDialToDriver != 0 {
		changeOrderState(order_id, config.StateClientDialToDriver)
	}
	rep, err = myAgi.Exec("Dial", "Local/"+driver_phone+"@local")
	if err != nil || rep.Res == -1 {
		log.Printf("[ERROR] StreamFile: %+v\n", err)
		return
	}
}
