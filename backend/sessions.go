package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Получение сессии
func getSession(r *http.Request) (*sessions.Session, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return nil, err
	}
	return session, nil
}
