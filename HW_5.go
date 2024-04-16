package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       string   `json:"id"`
	UserName string   `json:"username"`
	Age      int      `json:"age"`
	Friends  []string `json:"friends"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу пользователей
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT,
		age INTEGER,
		friends TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Post("/create", CreateUser)
	r.Get("/user/{userID}", GetUser)
	r.Post("/make_friends", MakeFriends)
	r.Get("/users", GetAllUsers)

	http.ListenAndServe(":8080", r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := generateUserID()
	newUser.ID = id

	// Преобразуем слайс друзей в строку
	friendsString := strings.Join(newUser.Friends, ",")

	// Вставляем пользователя в базу данных
	_, err = db.Exec("INSERT INTO users (id, username, age, friends) VALUES (?, ?, ?, ?)",
		newUser.ID, newUser.UserName, newUser.Age, friendsString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	userID := chi.URLParam(r, "userID")

	var user User
	err := db.QueryRow("SELECT id, username, age, friends FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.UserName, &user.Age, &user.Friends)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func MakeFriends(w http.ResponseWriter, r *http.Request) {
	type FriendRequest struct {
		SourceID string `json:"source_id"`
		TargetID string `json:"target_id"`
	}

	var request FriendRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем пользователей из базы данных
	sourceUser, err := getUserByID(request.SourceID)
	if err != nil {
		http.Error(w, "Source user not found", http.StatusBadRequest)
		return
	}

	targetUser, err := getUserByID(request.TargetID)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusBadRequest)
		return
	}

	// Добавляем друзей друг другу
	sourceUser.Friends = append(sourceUser.Friends, targetUser.ID)
	targetUser.Friends = append(targetUser.Friends, sourceUser.ID)

	// Обновляем информацию о пользователях в базе данных
	err = updateUserFriends(sourceUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = updateUserFriends(targetUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseMessage := fmt.Sprintf("%s и %s теперь друзья", sourceUser.UserName, targetUser.UserName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Выполняем запрос к базе данных
	rows, err := db.Query("SELECT id, username, age, friends FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Создаем слайс для хранения пользователей
	var users []User

	// Итерируемся по результатам запроса и добавляем пользователей в слайс
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.UserName, &user.Age, &user.Friends)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Проверяем наличие ошибок после завершения итерации по результатам запроса
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Кодируем слайс пользователей в формат JSON и отправляем его в ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByID(userID string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, age, friends FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.UserName, &user.Age, &user.Friends)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func updateUserFriends(user User) error {
	// Избегаем ошибки при пустом слайсе друзей
	friendsString := strings.Join(user.Friends, ",")

	// Обновляем список друзей пользователя в базе данных
	_, err := db.Exec("UPDATE users SET friends = ? WHERE id = ?", friendsString, user.ID)
	if err != nil {
		return err
	}
	return nil
}
