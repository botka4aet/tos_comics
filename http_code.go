package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"log"
)

func c_get_link() (data Task, err error){
	resp, err := http.Get("http://"+cconfig.ServerIp + cconfig.ServerPort + "/get_task")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	return
}

func c_send_result(result Task) (good bool){
	out, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error: Struct convert: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "http://"+cconfig.ServerIp + cconfig.ServerPort + "/send_link", bytes.NewBuffer([]byte(out)))
	if err != nil {
		log.Printf("Error: Post result: %v", err)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		log.Printf("Error: decoding json: %v", err)
		return
	}
	if result.Answer == 1 {
		good = true
	}
	return
}