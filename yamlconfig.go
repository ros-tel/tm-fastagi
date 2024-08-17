package main

import (
	"log"
	"os"

	"tm-fastagi/pkg/queue"

	"gopkg.in/yaml.v2"
)

// Load the YAML config file
func configLoad(configFile string, config *conf) {
	input, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer input.Close()

	decoder := yaml.NewDecoder(input)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal("error:", err)
	}
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// configuration contains the application settings
type (
	conf struct {
		ListenAddr string            `yaml:"listen_addr"`
		Redis      redisInfo         `yaml:"redis"`
		Nats       queue.Nats        `yaml:"nats"`
		Colors     map[string]string `yaml:"colors"`
		CarMarks   map[string]string `yaml:"car_marks"`

		DriverInfoBeforeIvr     bool `yaml:"driver_info_before_ivr"`
		NotCancelInStateConfirm bool `yaml:"not_cancel_in_state_confirm"`

		StateCancel             int `yaml:"state_cancel"`
		StateClientDialToDriver int `yaml:"state_client_dial_to_driver"`
		BlackPhoneCategoryId    int `yaml:"black_phone_category_id"`
		WhitePhoneCategoryId    int `yaml:"white_phone_category_id"`
		BlackClientGroupId      int `yaml:"black_client_group_id"`
		WhiteClientGroupId      int `yaml:"white_client_group_id"`

		StateConfirm []string `yaml:"state_confirm"`
		StateInPlace []string `yaml:"state_in_place"`
		StateInCar   []string `yaml:"state_in_car"`

		Api struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			ApiKey  string `yaml:"apikey"`
			TApiKey string `yaml:"tapikey"`
		} `yaml:"api"`
	}

	redisInfo struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"auth"`
	}
)

// Make config
func getConfig() {
	// Load the configuration file
	if *config_file == "" {
		*config_file = "config" + string(os.PathSeparator) + "config.yml"
	}
	configLoad(*config_file, config)
}
