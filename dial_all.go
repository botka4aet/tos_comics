package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func dial_all() {
	fd, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fd.Close()
	ff, _ := os.OpenFile("txtfiles\\failed.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer ff.Close()

	for url, _ := range сheck {
		sem.Acquire()
		go func() {
			defer sem.Release()
			fmt.Println("Bruteforcing ", url)
			suffix, runes := get_suffix_rune(url)

			dtime := time.Now()
			var counter int
			conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", "cdn.townofsins.com:443", &tls.Config{})
			ch := make(chan string, 10)
			ch_close := make(chan bool, 1)
			go ch_scramble_da("", runes, 4, ch, ch_close)

			var answer string
			buf := make([]byte, 1000)

			for {
				result, ok := <-ch
				if !ok {
					return
				}
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
						if result == "zzzzz" {
							_, _ = ff.WriteString(url + "\n")
							close(ch_close)
						}
						break
					} else if strings.HasPrefix(answer, "HTTP/1.1 200") {
						mutex.Lock()
						_, _ = fd.WriteString(url + "_" + result + suffix + "\n")
						mutex.Unlock()
						close(ch_close)
						break
					}
				}
				counter++
				if counter >= speed_counter {
					counter = 0
					fmt.Printf("%v[%v]Speed - %.2f per second.\n", time.Now().Format("[15:04:05]"), url, float64(speed_counter)/time.Since(dtime).Seconds())
					dtime = time.Now()
				}
			}
		}()
	}
}

func ch_scramble_da(suffix string, runes *[]rune, step int, ch chan string, ch_close chan bool) {
	var i int
	for i < len(*runes) {
		select {
		case <-ch_close:
			if step == 0 {
				close(ch)
			}
			return
		default:
		}
		if step > 0 {
			ch_scramble_da(string((*runes)[i])+suffix, runes, step-1, ch, ch_close)
		} else {
			ch <- string((*runes)[i]) + suffix
		}
		i++
	}
}
