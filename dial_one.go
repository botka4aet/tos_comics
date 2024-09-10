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
	fi, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fi.Close()

	for url, _ := range Check {
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
								close(ch_close)
							}
							break
						} else if strings.HasPrefix(answer, "HTTP/1.1 200") {
							close(ch_close)
							_, _ = fi.WriteString(url + "_" + result + suffix + "\n")
							close(ch_close)
							break
						}
					}
				}
			}()
		}
		ch_scramble_o("", &letterRunes, 4, ch, ch_close)
	}
}
