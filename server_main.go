package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var TaskList = make(map[string]string)

//Answer: 1 - данные приняты, 0 - ошибка
type Task struct {
    Link string
    Code string
	Secret string
	Answer uint8
}

func server_main() {
	defer wg.Done()
	r := mux.NewRouter()

	r.HandleFunc("/send_result", send_link).Methods("POST")
	r.HandleFunc("/get_task", get_task).Methods("GET")

	srv := &http.Server{
		Addr:    cconfig.ServerPort,
		Handler: r,
	}
	srv.ListenAndServe()
}

func send_link(w http.ResponseWriter, req *http.Request) {
	var result Task
	err := json.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		logmes(1,"Decoding Json from POST", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		answer := &Task{Answer: 0}
		//Проверяем - есть ли такая запись в очереди
		tasks, ok := TaskList[result.Link]
		if ok && tasks == result.Secret && check_link(result) {
			answer.Answer = 1
			ch_sql <- result
			delete(TaskList, result.Link)
			queue(result,2)
		}
		json.NewEncoder(w).Encode(answer)
	}
}

func get_task(w http.ResponseWriter, req *http.Request) {
	new_task := &Task{Secret: strconv.Itoa(rand.Int())}
	for {
		new_link := get_link_sql()
		_, ok := TaskList[new_link]
		if !ok {
			new_task.Link = new_link
			break
		}
	}
	ch_sql <- *new_task
	//Пытаемся отправить ответ
	err := json.NewEncoder(w).Encode(new_task)
	if err != nil {
		logmes(1,"New task encode", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Добавляем в список
	TaskList[new_task.Link] = new_task.Secret
	queue(*new_task,1)
	logmes(2,"Send task "+new_task.Link, "")
}


