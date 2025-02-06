package main

import (
	"crypto/tls"
	"fmt"
	"net"
//	"runtime"
	"strings"
	"time"
)

func dial_one(url *string) (code string) {
	fmt.Println("Bruteforcing ", *url)
	suffix, runes := get_suffix_rune(url)

	ch := make(chan string, 50)
	ch_close := make(chan bool)

	for i := 0; i < 500; i++ {
		//		for i := 0; i < runtime.NumCPU(); i++ {
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
					data := []byte("HEAD /media/assets/images/" + *url + "_" + result + suffix + " HTTP/1.1\r\nHost: cdn.townofsins.com\r\n\r\n")
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
						code = result
						close(ch_close)
						break
					}
				}
			}
		}()
	}
	timer_g = time.Now()
	ch_scramble_o("", runes, 4, ch, ch_close)
	return
}
