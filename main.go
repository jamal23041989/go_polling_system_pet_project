package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
			handleGetPolls(w, r) // GET получить список всех опросов.
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

func handleGetPolls(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// 2. Подготавливаем список опросов для ответа
	pollsList := make([]Poll, 0, len(polls))
	for _, poll := range polls {
		pollsList = append(pollsList, poll)
	}

	// 3. Возвращаем данные
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pollsList)
}

func handleCreatePoll(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// 2. Проверяем Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ожидается application/json", http.StatusBadRequest)
		return
	}

	// 3. Читаем и декодируем тело запроса
	var newPoll struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
	}

	err := json.NewDecoder(r.Body).Decode(&newPoll)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// 4. Валидация данных
	if newPoll.Question == "" {
		http.Error(w, "Поле question обязательно", http.StatusBadRequest)
		return
	}

	if len(newPoll.Options) < 2 {
		http.Error(w, "Должно быть минимум 2 варианта ответа", http.StatusBadRequest)
		return
	}

	// 5. Создаём новый опрос
	poll := Poll{
		ID:       len(polls) + 1, // Генерируем ID
		Question: newPoll.Question,
		Options:  newPoll.Options,
		Votes:    make([]int, len(newPoll.Options)), // Инициализируем счётчики нулями
	}

	// 6. Сохраняем опрос
	polls[poll.ID] = poll

	// 7. Возвращаем созданный опрос
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(poll)
}

func handleGetPoll(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// 2. Извлекаем ID из URL
	idStr := strings.TrimPrefix(r.URL.Path, "/polls/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// 3. Ищем опрос в хранилище
	poll, exists := polls[id]
	if !exists {
		http.Error(w, "Опрос не найден", http.StatusNotFound)
		return
	}

	// 4. Возвращаем данные опроса
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// 2. Извлекаем ID опроса из URL
	idStr := strings.TrimPrefix(r.URL.Path, "/polls/")
	idStr = strings.TrimSuffix(idStr, "/vote")
	pollID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID опроса", http.StatusBadRequest)
		return
	}

	// 3. Проверяем Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ожидается application/json", http.StatusBadRequest)
		return
	}

	// 4. Читаем и декодируем тело запроса
	var voteRequest struct {
		OptionIndex int `json:"option_index"`
	}

	err = json.NewDecoder(r.Body).Decode(&voteRequest)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// 5. Проверяем существование опроса
	poll, exists := polls[pollID]
	if !exists {
		http.Error(w, "Опрос не найден", http.StatusNotFound)
		return
	}

	// 6. Валидация option_index
	if voteRequest.OptionIndex < 0 || voteRequest.OptionIndex >= len(poll.Options) {
		http.Error(w, "Неверный индекс варианта ответа", http.StatusBadRequest)
		return
	}

	// 7. Обновляем счетчик голосов
	poll.Votes[voteRequest.OptionIndex]++
	polls[pollID] = poll // Сохраняем изменения

	// 8. Возвращаем обновленный опрос
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
}
