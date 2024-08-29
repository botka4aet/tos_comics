package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"math"
	"runtime"
	"time"
)

var baseurl = "https://cdn.townofsins.com/media/assets/images/"

type Semaphore struct {
	C chan struct{}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}
func (s *Semaphore) Release() {
	<-s.C
}

var sem = Semaphore{
	C: make(chan struct{}, 50),
}

func main() {
	fmt.Println(runtime.NumCPU())

	for url, _ := range check {
		fmt.Println("Bruteforcing ", url)
		sgnlCh := make(chan struct{})
		var link [5]int
		var sufix string
		var letterRunes []rune
		if strings.HasPrefix(url, "comics_adventure/th/") {
			sufix = ".webp"
			letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
		} else if strings.HasPrefix(url, "comics_events/th/") {
			sufix = ".jpg"
			letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
		} else if strings.HasPrefix(url, "comics/th/") {
			sufix = "@2x.webp"
			letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		}
		lenr := len(letterRunes)

		dtime := time.Now()
		for {
			sem.Acquire()
			go scramble_link(url, sufix, &letterRunes, lenr, link, sgnlCh)
			link[0] = lenr
			var loop bool
			for i, _ := range link {
				if link[i] >= lenr {
					if i == len(link)-1 {
						loop = true
						break
					}
					link[i+1]++
					if i+1 == 4 && link[4] < lenr {
						fmt.Printf("%v Now trying ****%v. Speed - %.2f per second\n", time.Now().Format("[15:04:05]"), string(letterRunes[link[4]]), float64(math.Pow(26, 4))/float64(runtime.NumCPU())/time.Since(dtime).Seconds())
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
}

func scramble_link(url string, sufix string, runes *[]rune, lenr int, link [5]int, sgnlCh chan struct{}) {
	defer sem.Release()
	for i := 0; i < lenr; i++ {
		link[0] = i
		var nlink string
		for _, t := range link {
			nlink += string((*runes)[t])
		}
		check_url(url+"_"+nlink+sufix, sgnlCh)
	}
}

func check_url(link string, sgnlCh chan struct{}) {
	var client = &http.Client{}
	var res *http.Response
	for res == nil {
		res, _ = client.Head(baseurl + link)
		if res != nil && res.StatusCode == 200 {
			fmt.Println("Solved: " + link)
			fi, err := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
			if err == nil {
				_, _ = fi.WriteString(link + "\n")
			}
			defer fi.Close()
			close(sgnlCh)
			return
		} else if res != nil && res.StatusCode == 500 {
			break
		}
		time.Sleep(time.Second)
	}
}
