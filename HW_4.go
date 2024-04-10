package main

import (
	"encoding/json"
	"fmt"
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

var users = make(map[string]User)

func main() {
	r := chi.NewRouter()
	r.Post("/create", CreateUser)
	r.Get("/user/{userID}", GetUser)
	r.Post("/make_friends", MakeFriends)

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
	users[id] = newUser

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

	user, ok := users[userID]
	if !ok {
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

	sourceUser, ok1 := users[request.SourceID]
	targetUser, ok2 := users[request.TargetID]
	if !ok1 || !ok2 {
		http.Error(w, "One or both users not found", http.StatusBadRequest)
		return
	}

	responseMessage := fmt.Sprintf("%s и %s теперь друзья", sourceUser.UserName, targetUser.UserName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}
