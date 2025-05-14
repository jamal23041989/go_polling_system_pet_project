package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Poll struct {
	ID       int      `json:"id"`       // Уникальный идентификатор
	Question string   `json:"question"` // Текст вопроса (например, "Нравится ли вам Go?")
	Options  []string `json:"options"`  // Варианты ответов (например, ["Да", "Нет"])
	Votes    []int    `json:"votes"`    // Количество голосов за каждый вариант (например, [5, 2])
}

type Vote struct {
	PollID      int `json:"poll_id"`      // ID опроса, в котором голосуют
	OptionIndex int `json:"option_index"` // Индекс выбранного варианта (0, 1, ...)
}

var polls = make(map[int]Poll)   // Хранит все опросы
var votes = make(map[int][]Vote) // Хранит голоса по poll_id

func main() {
	http.HandleFunc("/polls", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAllPolls(w, r) // GET получить список всех опросов.
		case http.MethodPost:
			handleCreatePoll(w, r) // POST создать новый опрос (принимает Question и Options).
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/polls/{id}", handleGetPoll)   // GET получить опрос по ID (с результатами голосования).
	http.HandleFunc("/polls/{id}/vote", handleVote) // POST проголосовать (принимает OptionIndex).

	fmt.Println("Listening on port 8081")
	http.ListenAndServe(":8081", nil)
}

func handleGetAllPolls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(polls)
}

func handleCreatePoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func handleGetPoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}
