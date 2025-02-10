package main

import (
	"bufio"
	"log"
	"os"
)

func logmes(lvl uint8, mes string, err string) {
	if cconfig.LogLevel&lvl == 0 {
		return
	}
	mestype := ""
	switch lvl {
	case 1:
		mestype = "Error"
	case 2:
		mestype = "Http"
	case 4:
		mestype = "Result"
	case 8:
		mestype = "Info"
	default:
		mestype = "Unknown"
	}
	if err == "" {
		log.Printf("%v: %v", mestype, mes)
	} else {
		log.Printf("%v: %v: %v", mestype, mes, err)
	}
}

func logfatal(mes string, err error) {
	log.Printf("Fatal Error: %v: %v", mes, err)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}
