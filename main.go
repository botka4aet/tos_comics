package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	a := GetOutboundIP()
	_ = a

	if cconfig.Server & 1 > 0 {
		logmes(8,"Starting server with port "+cconfig.ServerPort, "")
		server_init()
		wg.Add(1)
		go server_main()
	}
	if cconfig.Server & 1 > 0 {
		logmes(8,"Starting client", "")
		wg.Add(1)
		go client_main()
	}
	wg.Wait()
}