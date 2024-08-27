package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"runtime"
	"math"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var letternumbersRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

var mu sync.Mutex
var tries_counter int

type Semaphore struct {
	C chan struct{}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}
func (s *Semaphore) Release() {
	<-s.C
	// mu.Lock()
	// tries_counter++
	// mu.Unlock()
}

var sem = Semaphore{
	C: make(chan struct{}, 50),
}

func main() {
	fmt.Println(runtime.NumCPU())

	c := make(chan string)
	go func() {
		for url, _ := range check {
			fmt.Println("Bruteforcing ", url)
			sgnlCh := make(chan struct{})
			var link [5]int
			dtime := time.Now()
			for {
				sem.Acquire()
				go scramble_link(url, link, c, sgnlCh)
				link[0] = 26
				var loop bool
				for i, _ := range link {
					if link[i] >= 26 {
						if i == len(link)-1 {
							loop = true
							break
						}
						link[i+1]++
						if i+1 == 4 && link[4] < 26 {
							fmt.Printf("%v Now trying ****%v. Speed - %.2f per second\n",time.Now().Format("[15:04:05]"),string(letterRunes[link[4]]),float64(math.Pow(26, 5))/float64(runtime.NumCPU())/time.Since(dtime).Seconds())
							dtime = time.Now()
						}
						for j := 0; j <= i; j++ {
							link[j] = 0
						}
					}
				}
				select {
				case <-sgnlCh:
					loop = true
				default:
				}
				if loop {
					break
				}
			}
		}
	}()
	// go func() {
	// 	for {
	// 		time.Sleep(time.Minute)
	// 		c <- ""
	// 	}

	// }()
	fi, err := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	for {
		text := <-c
		switch text {
		case "":
			fmt.Println("Passed 1 min, tries - ", tries_counter*26)
			mu.Lock()
			tries_counter = 0
			mu.Unlock()
		default:
			_, err = fi.WriteString(text + "\n")
			if err != nil {
				panic(err)
			}
		}

	}
}

func scramble_link(url string, link [5]int, c chan string, sgnlCh chan struct{}) {
	defer sem.Release()
	for i := 0; i < 26; i++ {
		select {
		case <-sgnlCh:
			return
		default:
		}
		link[0] = i

		var nlink string
		for _, t := range link {
			nlink += string(letterRunes[t])
		}
		check_url(url+"_"+nlink, c, sgnlCh)
	}
}

func check_url(link string, c chan string, sgnlCh chan struct{}) {
	baseurl := "https://cdn.townofsins.com/media/assets/images/"
	postfix := ""
	if strings.HasPrefix(link, "comics_adventure/th/") {
		postfix = ".webp"
	} else if strings.HasPrefix(link, "comics_events/th/") {
		postfix = ".jpg"
		//	} else if strings.HasPrefix(link, "comics/th/") {
		//		postfix = "@2x.webp"
	}
	if postfix == "" {
		return
	}
	var client = &http.Client{}
	var res *http.Response
	for res == nil {
		res, _ = client.Head(baseurl + link + postfix)
		if res != nil && res.StatusCode == 200 {
			fmt.Println("Solved: " + link)
			c <- (link + postfix)
			close(sgnlCh)
			return
		} else if res != nil && res.StatusCode == 500 {
			break
		}
		time.Sleep(time.Second)
	}
}
