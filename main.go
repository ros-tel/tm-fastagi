package main

import (
	"bufio"
	"flag"
	"log"
	"net"

	"tm-fastagi/pkg/queue"

	tm "github.com/ros-tel/taximaster/common_api"
	tmt "github.com/ros-tel/taximaster/tm_tapi"
	"github.com/zaf/agi"
)

var (
	// config the settings variable
	config = &conf{}

	common_api_client *tm.Client
	tm_tapi_client    *tmt.Client

	nats_js *queue.Nats

	config_file = flag.String("config", "", "Usage: -config=<config_file>")
	debug       = flag.Bool("debug", false, "Print debug information on stderr")
)

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()

	getConfig()

	redisConnect(config.Redis)

	common_api_client = tm.NewClient(config.Api.Host+":"+config.Api.Port, config.Api.ApiKey)
	tm_tapi_client = tmt.NewClient(config.Api.Host+":"+config.Api.Port, config.Api.TApiKey)

	if *debug {
		log.Printf("[DEBUG] CONFIG: %+v", config)
	}

	if config.Nats.Uri != "" {
		nats_js = queue.Connect(config.Nats)
		nats_js.KeyValue("tm-contract_auth")
	}

	fagiserv, err := net.Listen("tcp", config.ListenAddr)
	if fagiserv == nil {
		log.Fatalf("Cannot listen: %v", err)
	}

	log.Println("Server started")

	for {
		conn, err := fagiserv.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}
		go handleFastAgiConnection(conn)
	}
}

func handleFastAgiConnection(client net.Conn) {
	defer client.Close()

	myAgi := agi.New()
	rw := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	err := myAgi.Init(rw)
	if err != nil {
		log.Printf("Error Init: %+v\n", err)
		return
	}

	if *debug {
		// Print AGI environment
		log.Println("[DEBUG] AGI environment vars:")
		for key, value := range myAgi.Env {
			log.Printf("[DEBUG] %-15s: %s\n", key, value)
		}
	}

	// Parse AGI request
	path, query, err := parseAgiReq(myAgi.Env["request"])
	if err != nil {
		log.Println(err)
		return
	}
	if *debug {
		log.Println("[DEBUG]", path, query)
	}

	if path == "" {
		log.Println("[ERROR] path not set")
		return
	}

	switch path {
	case "/record":
		createRecordLink(query)
	case "/cancel_by_order_id":
		cancelByOrderID(myAgi, query)
	case "/driver_phone_by_crew":
		getDriverPhoneByCrew(myAgi, query)
	case "/driver_phone_by_caller":
		getDriverPhoneByCaller(myAgi, query)
	case "/ivr":
		ivr(myAgi)
	case "/caller_info":
		caller_info(myAgi)
	case "/phone_info":
		phone_info(myAgi, query)
	case "/auth":
		auth(myAgi, query)
	case "/set_client_group":
		setClientGroup(query)
	case "/show_tm_message":
		showTmMessage(query)
	}
}
