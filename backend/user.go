package main

// Структура пользователя для хранения в базе данных
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
