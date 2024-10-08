package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"slices"
)

var сheck = make(map[string]bool)
var checka []string
var gmodef int
var gmodes bool

func init() {
	var done = make(map[string]bool)
	suffixo := flag.String("suffix", "", "Отбор строк")
	modeo := flag.Int("mode", 0, "Режим перебора")
	//Нужно ли сортировать? 0 - нет, 1 - да, 2 - сортировать по суффиксу, 3 - реверс по суффиксу
	sorto := flag.Int("sort", 0, "Режим сортировки")
	flag.Parse()
	gmodef = *modeo
	if *sorto > 1 {
		gmodes = true
	}

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
		if ok || strings.HasPrefix(text, "comics/th/") || strings.HasPrefix(text, "comics_mythic/th/") || !strings.HasSuffix(text, *suffixo) {
			continue
		}
		_, ok = сheck[text]
		if ok {
			continue
		}
		сheck[text] = true
		needcheck++
	}

	ffail, err := os.Open("txtfiles\\failed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer ffail.Close()
	scanner = bufio.NewScanner(ffail)
	for scanner.Scan() {
		text := scanner.Text()
		_, ok := сheck[text]
		if !ok {
			continue
		}
		delete(сheck, text)
		needcheck--
	}
	for url, _ := range сheck {
		if *sorto > 1 {
			tpos := strings.LastIndex(url, "_")
			if tpos != -1 {
				numb := url[tpos+1:]
				if len(numb) == 1 {
					numb = "0"+numb
				}
				url = numb+":"+url		
			}
		}
		checka = append(checka, url)
	}
	switch *sorto {
	case 1,2:
		slices.Sort(checka)
	case 3:
		slices.Sort(checka)
		slices.Reverse(checka)
	}
	fmt.Println("Need check - ", needcheck)
}
