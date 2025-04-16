package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	// Хэширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}

	users[user.Username] = string(hash)

	w.WriteHeader(http.StatusOK)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	hash, ok := users[user.Username]
	if !ok || bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password)) != nil {
		http.Error(w, "Неверные данные для входа", http.StatusUnauthorized)
		return
	}

	// Создание сессии
	err = createSession(w, r, user.Username)
	if err != nil {
		http.Error(w, "Ошибка при создании сессии", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil || session.Values["userID"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Write([]byte("Добро пожаловать, " + session.Values["userID"].(string)))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := destroySession(w, r)
	if err != nil {
		http.Error(w, "Ошибка при выходе", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
