package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]User) // Хранилище пользователей

// Регистрация нового пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка на существование пользователя
	if _, exists := users[credentials.Username]; exists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Сохраняем хеш пароля и имя пользователя в "базу данных"
	users[credentials.Username] = User{Username: credentials.Username, Password: string(hash)}
	w.WriteHeader(http.StatusOK)
}

// Вход в систему
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка пароля
	user, exists := users[credentials.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Создание сессии
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}

// Профиль пользователя (только для авторизованных)
func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	// Проверка авторизации
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Возвращаем информацию о пользователе
	sessionUser := session.Values["user"].(string)
	user, exists := users[sessionUser]
	if !exists {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Выход из системы
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "authenticated")
	delete(session.Values, "user")
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}
