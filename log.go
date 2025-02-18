package main

import (
	"bufio"
	"log"
	"os"
)

var log_text = []string{
	"Server",  //1
	"Client",  //2
	"undef",   //4
	"Error",   //8
	"Http",    //16
	"Result",  //32
	"Info",    //64
	"Unknown"} //128

func logmes_text(lvl int) (result string) {
	cur_exp := 1
	for _, j := range log_text {
		if lvl&cur_exp > 0 {
			result = result + j
			lvl -= cur_exp
		}
		cur_exp *= 2
		if cur_exp > lvl {
			break
		}
	}
	return
}

func logmes(lvl int, mes string, err string) {
	if cconfig.LogLevel&lvl == 0 {
		return
	}
	mestype := logmes_text(lvl)

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
