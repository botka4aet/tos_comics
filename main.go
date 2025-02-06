package main

import (
	"fmt"
)

func main() {
	if cconfig.Server {
		fmt.Println("Starting server with port", cconfig.ServerPort)
		server_init()
		server_main()
	} else {
		fmt.Println("Starting client")
		client_main()
	}
}