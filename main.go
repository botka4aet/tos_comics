package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

var speed_counter = 100000
var timer_g time.Time
var counter_g int
var mutex = &sync.Mutex{}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var letternRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

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

	switch gmodef {
	case 1:
		dial_one()
	case 2:
		dial_all()
	case 3:
		httpc_one()
	default:
		httpc_one_fh()
	}
	fmt.Scanln()
}

func ch_scramble_o(suffix string, runes *[]rune, step int, ch chan string, ch_close chan bool) {
	var i int
	for i < len(*runes) {
		select {
		case <-ch_close:
			if step == 0 {
				close(ch)
				counter_g = 0
			}
			return
		default:
		}
		if step > 0 {
			ch_scramble_o(string((*runes)[i])+suffix, runes, step-1, ch, ch_close)
		} else {
			ch <- string((*runes)[i]) + suffix
			counter_g++
			if counter_g >= speed_counter {
				counter_g = 0
				fmt.Printf("%vSpeed - %.2f per second\n", time.Now().Format("[15:04:05]"), float64(speed_counter)/time.Since(timer_g).Seconds())
				timer_g = time.Now()
			}
		}
		i++
	}
}

func get_suffix_rune(url string) (suffix string, runes *[]rune) {
	if strings.HasPrefix(url, "comics_adventure/th/") {
		suffix = ".webp"
		runes = &letterRunes
	} else if strings.HasPrefix(url, "comics_events/th/") || strings.HasPrefix(url, "comics_mythic/th/") {
		suffix = ".jpg"
		runes = &letterRunes
	} else if strings.HasPrefix(url, "comics/th/") {
		suffix = "@2x.webp"
		runes = &letternRunes
	}
	return
}

func clear_url(url string) string {
	return url[strings.Index(url, ":")+1:]
}
