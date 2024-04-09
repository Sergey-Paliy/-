package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

type User struct {
	ID       string   `json:"id"`
	UserName string   `json:"username"`
	Age      int      `json:"age"`
	Friends  []string `json:"friends"`
}

var users = make(map[string]User) // Переменная для хранения пользователей

func main() {
	r := chi.NewRouter()
	r.Post("/create", CreateUser)
	r.Get("/user/{userID}", GetUser)

	http.ListenAndServe(":8080", r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := generateUserID() // Генерация уникального ID для пользователя
	newUser.ID = id
	users[id] = newUser // Сохранение пользователя в мапе

	jsonResponse := map[string]string{"id": id}
	response, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func generateUserID() string {
	id := uuid.New().String()
	return id
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Получаем уникальный идентификатор пользователя из URL
	userID := chi.URLParam(r, "userID")

	// Получаем данные пользователя из вашего хранилища
	user, ok := users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Преобразуем данные пользователя в формат JSON и отправляем их клиенту
	json.NewEncoder(w).Encode(user)
}
