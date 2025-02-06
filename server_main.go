package main

import (
	"encoding/json"
	"log"
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
	r := mux.NewRouter()

	r.HandleFunc("/get_task", get_task).Methods("GET")
	r.HandleFunc("/send_result", send_link).Methods("POST")

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
		log.Printf("Error: decoding json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		//Проверяем - есть ли такая запись в очереди
		result.Answer = 0
		tasks, ok := TaskList[result.Link]
		if ok && tasks == result.Secret && check_link(result) {
			result.Answer = 1
			delete(TaskList, result.Link)
			queue(result,2)
		}
		json.NewEncoder(w).Encode(result)
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

	//Пытаемся отправить ответ
	err := json.NewEncoder(w).Encode(new_task)
	if err != nil {
		log.Printf("Error: encode error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Добавляем в список
	TaskList[new_task.Link] = new_task.Secret
	queue(*new_task,1)
	log.Printf("Info: send task: %v", new_task.Link)
}


