package main

import (
	"fmt"
	"runtime"
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

func client_main() {
	defer wg.Done()
	fmt.Println(runtime.NumCPU())

	for {
		task, err := c_get_link()
		if err != nil {
			logmes(10,"Can't get task", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		logmes(18,"Got task - "+task.Link,"")
		switch cconfig.Brute_mode {
		case 1:
			task.Code = dial_one(&task.Link)
		case 2:
			task.Code = httpc_one(&task.Link)
		default:
			task.Code = httpc_one_fh(&task.Link)
		}
		for {
			if c_send_result(task) { break }
			//Если данные не получилось отправить - ждем
			time.Sleep(10 * time.Second)
		}
	}
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
				message := 	fmt.Sprintf("Speed - %.2f per second", float64(speed_counter)/time.Since(timer_g).Seconds())
				logmes(8,message,"")
				timer_g = time.Now()
			}
		}
		i++
	}
}
