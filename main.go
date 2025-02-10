package main

import (
	"sync"
	"fmt"
)

var wg sync.WaitGroup

func main() {
	a := GetOutboundIP()
	_ = a

	if cconfig.Server & 1 == 1 {
		fmt.Println("Starting server with port", cconfig.ServerPort)
		server_init()
		wg.Add(1)
		go server_main()
	}
	if cconfig.Server & 2 == 2 {
		fmt.Println("Starting client")
		wg.Add(1)
		go client_main()
	}
	wg.Wait()
}