package main

import (
	"bytes"
	"encoding/json"
	"net"
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

	req, err := http.NewRequest("POST", "http://"+cconfig.ServerIp + cconfig.ServerPort + "/send_result", bytes.NewBuffer([]byte(out)))
	if err != nil {
		log.Printf("Error: Creating POST: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error: Post result: %v", err)
		return
	}
	
	defer res.Body.Close()
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

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}