package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

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
	C: make(chan struct{}, runtime.NumCPU()),
}

func main() {
	fmt.Println(runtime.NumCPU())
	fi, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fi.Close()

	for url, _ := range check {
		sem.Acquire()
		go func() {
			defer sem.Release()
			fmt.Println("Bruteforcing ", url)
			var suffix string
			var letterRunes []rune
			if strings.HasPrefix(url, "comics_adventure/th/") {
				suffix = ".webp"
				letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
			} else if strings.HasPrefix(url, "comics_events/th/") {
				suffix = ".jpg"
				letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
			} else if strings.HasPrefix(url, "comics/th/") {
				suffix = "@2x.webp"
				letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
			}

			dtime := time.Now()
			var counter int
			conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", "cdn.townofsins.com:443", &tls.Config{})
			ch := make(chan string, 10)
			go ch_scramble("", &letterRunes, 4, ch)

			var answer string
			buf := make([]byte, 1000)

			for {
				result := <-ch
				for {
					data := []byte("HEAD /media/assets/images/" + url + "_" + result + suffix + " HTTP/1.1\r\nHost: cdn.townofsins.com\r\n\r\n")
					if strings.HasPrefix(answer, "\x00") || err != nil {
						conn, err = tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", "cdn.townofsins.com:443", &tls.Config{})
					}
					if err != nil {
						continue
					}
					_, err = conn.Write(data)
					if err != nil {
						continue
					}
					_, err = conn.Read(buf)
					if err != nil {
						continue
					}
					answer = string(buf[:])
					if strings.HasPrefix(answer, "HTTP/1.1 500") {
						break
					} else if strings.HasPrefix(answer, "HTTP/1.1 200") {
						close(ch)
						mutex.Lock()
						_, _ = fi.WriteString(url + "_" + result + suffix + "\n")
						mutex.Unlock()
						result = "zzzzz"
						break
					}
				}
				if result == "zzzzz" {
					break
				} else if counter == 10000 {
					counter = 0
					fmt.Printf("%v[%v]Speed - %.2f per second.\n", time.Now().Format("[15:04:05]"), url, 10000/time.Since(dtime).Seconds())
					dtime = time.Now()
				} else {
					counter++
				}
			}
		}()
	}
}

func ch_scramble(suffix string, runes *[]rune, step int, Ch chan string) {
	var i int
	for i < len(*runes) {
		if step > 0 {
			ch_scramble(string((*runes)[i])+suffix, runes, step-1, Ch)
		} else {
			Ch <- string((*runes)[i]) + suffix
		}
		i++
	}
}
