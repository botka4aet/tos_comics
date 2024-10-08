package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

func dial_one() {
	fd, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fd.Close()
	ff, _ := os.OpenFile("txtfiles\\failed.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer ff.Close()

	for url, _ := range сheck {
		if gmodes {
			url = clear_url(url)
		}
		fmt.Println("Bruteforcing ", url)
		suffix, runes := get_suffix_rune(url)

		ch := make(chan string, 50)
		ch_close := make(chan bool)

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				var answer string
				conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", "cdn.townofsins.com:443", &tls.Config{})
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
							_, _ = fd.WriteString(url + "_" + result + suffix + "\n")
							close(ch_close)
							break
						}
					}
				}
			}()
		}
		ch_scramble_o("", runes, 4, ch, ch_close)
	}
}
