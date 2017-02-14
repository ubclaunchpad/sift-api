package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	name := "person"
	pass := "123123"

	redirect := "/"

	// authentication() here

	// failed
	if false {
		http.Redirect(w, r, redirect, 302)
	}

	// passed
	// get user's id and put into a new session

	// userID := dm.GetUser(name) requires user profile table setup
	dm.DeleteSessionsByUser(userID)
	dm.CreateSession(Session{UserID: userID})
	seshID := dm.GetSessionByUser(userID).ID

	// put seshID in a cookie and attach to response
	cookieValue := map[string]string{
		"id": seshID,
	}

	encoded, err := cookieHandler.Encode("session", value)

	if err != nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}

	// redirect to wherever login takes us
	redirect = "/dashboard"
	http.Redirect(w, r, redirect, 302)
}
