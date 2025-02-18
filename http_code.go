package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

func c_get_link() (data Task, err error) {
	resp, err := http.Get("http://" + cconfig.ServerIp + cconfig.ServerPort + "/get_task")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	return
}

func c_send_result(result Task) (good bool) {
	if result.Code == "" {
		return true
	}
	out, err := json.Marshal(result)
	if err != nil {
		logmes(10,"Struct convert", err.Error())
		return
	}
	req, err := http.NewRequest("POST", "http://"+cconfig.ServerIp+cconfig.ServerPort+"/send_result", bytes.NewBuffer([]byte(out)))
	if err != nil {
		logmes(10,"Creating POST", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logmes(10,"Post result", err.Error())
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		logmes(10,"Decoding json", err.Error())
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
