package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func getSession(r *http.Request) (*sessions.Session, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return nil, err
	}
	return session, nil
}

func createSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}
	session.Values["userID"] = userID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return session.Save(r, w)
}

func destroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}
	session.Values = nil
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
