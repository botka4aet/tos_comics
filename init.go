package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var сheck = make(map[string]bool)
var imode int

func init() {
	var done = make(map[string]bool)
	suffix := flag.String("suffix", "", "Отбор строк")
	mode := flag.Int("mode", 0, "Режим перебора")
	flag.Parse()
	imode = *mode

	fdone, err := os.Open("txtfiles\\done.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fdone.Close()
	scanner := bufio.NewScanner(fdone)
	for scanner.Scan() {
		text := scanner.Text()
		_, ok := done[text]
		if ok {
			continue
		}
		done[text] = true
	}
	fdone.Close()
	file, err := os.Open("txtfiles\\links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fdone, err = os.OpenFile("txtfiles\\done.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer fdone.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		ntext := strings.TrimSuffix(text, "@2x.webp")
		ntext = strings.TrimSuffix(ntext, ".webp")
		ntext = strings.TrimSuffix(ntext, ".jpg")
		if text != ntext {
			text = ntext[:len(ntext)-6]
		}
		_, ok := done[text]
		if ok {
			continue
		}
		_, err = fdone.WriteString(text + "\n")
		if err != nil {
			panic(err)
		}
		done[text] = true
	}

	var needcheck int
	file, err = os.Open("txtfiles\\check.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		_, ok := done[text]
		if ok || strings.HasPrefix(text, "comics/th/") || !strings.HasSuffix(text, *suffix) {
			continue
		}
		_, ok = сheck[text]
		if ok {
			continue
		}
		сheck[text] = true
		needcheck++
	}
	fmt.Println("Need check - ", needcheck)
}
