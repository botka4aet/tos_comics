package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type config struct {
	ServerPort       string `json:"ServerPort"`
	ServerIp     string `json:"ServerIp"`
	Brute_mode uint8  `json:"brute_mode"`
	Server uint8  `json:"Server"`
}

var cconfig config

func init() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cconfig); err != nil {
		log.Fatal(err)
	}
	mode := flag.Int("mode", 0, "Режим работы: 0 - из файла, 1 - сервер, 2 - клиент, 3 - 1&2(wip), 4 - solo(wip)")
	brute_mode := flag.Int("brute_mode", 0, "Режим перебора: 0 - из файла, 1 - dial, 2 - http, 3 - fast http")

	flag.Parse()

	if *mode != 0 {
		cconfig.Server = uint8(*mode)
	}
	cconfig.Server = 1

	if *brute_mode != 0 {
		cconfig.Brute_mode = uint8(*brute_mode)
	}
	

}
