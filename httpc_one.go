package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"runtime"
	"time"
)

func httpc_one() {
	fd, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fd.Close()
	ff, _ := os.OpenFile("txtfiles\\failed.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer ff.Close()

	for url, _ := range сheck {
		fmt.Println("Bruteforcing ", url)
		suffix, runes := get_suffix_rune(url)

		ch := make(chan string, 50)
		ch_close := make(chan bool)

		conn := &http.Client{
			//	Timeout:   4 * time.Second,
		}

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					result, ok := <-ch
					if !ok {
						return
					}
					for {
						res, err := conn.Head("https://cdn.townofsins.com/media/assets/images/" + url + "_" + result + suffix)
						if err == nil && res != nil && res.StatusCode == 200 {
							_, _ = fd.WriteString(url + "_" + result + suffix + "\n")
							close(ch_close)
							select {
							case <-ch:
								return
							default:
							}
							break
						} else if err == nil && res != nil && res.StatusCode == 500 {
							if result == "zzzzz" {
								_, _ = ff.WriteString(url + "\n")
								close(ch_close)
							}
							break
						}
					}
				}
			}()
		}
		timer_g = time.Now()
		ch_scramble_o("", runes, 4, ch, ch_close)
	}
}

func httpc_one_fh() {
	fd, _ := os.OpenFile("txtfiles\\links.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer fd.Close()
	ff, _ := os.OpenFile("txtfiles\\failed.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer ff.Close()

	for url, _ := range сheck {
		fmt.Println("Bruteforcing ", url)
		suffix, runes := get_suffix_rune(url)

		ch := make(chan string, 50)
		ch_close := make(chan bool)

		for i := 0; i < runtime.NumCPU(); i++ {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			resp.SkipBody = true
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)
			req.Header.SetMethod("HEAD")
			go func() {
				for {
					result, ok := <-ch
					if !ok {
						return
					}
					for {
						req.SetRequestURI("https://cdn.townofsins.com/media/assets/images/" + url + "_" + result + suffix)
						err := fasthttp.Do(req, resp)
						if err == nil && resp.StatusCode() == 200 {
							_, _ = fd.WriteString(url + "_" + result + suffix + "\n")
							close(ch_close)
							break
						} else if err == nil && resp.StatusCode() == 500 {
							if result == "zzzzz" {
								_, _ = ff.WriteString(url + "\n")
								close(ch_close)
							}
							break
						}
					}
				}
			}()
		}
		timer_g = time.Now()
		ch_scramble_o("", runes, 4, ch, ch_close)
	}
}
