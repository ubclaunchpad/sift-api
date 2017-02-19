package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/securecookie"
)

// init cookie encryption keys
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pass := r.FormValue("pass")

	if name == "" || pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
		redirect := "/"

		profile, err := dm.GetProfileHelper(name, pass)

		if err != nil {
			http.Redirect(w, r, redirect, 302)
		}

	*/

	// passed
	// get user's id and put into a new session

	user := Profile{CompanyName: "Planet Express",
		PwHash: []byte("Super secret"), Address: "123 Fake Street"}

	user.ID = 10

	// userID := dm.GetUser(name) requires user profile table setup
	dm.deleteSessionsByUser(user.ID)
	newSesh := Session{UserID: user.ID}
	dm.createSession(&newSesh)
	createdSesh := dm.getSessionByUser(user.ID)

	// put seshID in a cookie and attach to response
	cookieValue := map[string]uint{
		"id": createdSesh.ID,
	}

	encoded, err := cookieHandler.Encode("session", cookieValue)

	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}

	// redirect to wherever login takes us
	//redirect = "/dashboard"
	//http.Redirect(w, r, redirect, 302)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	fmt.Println(name)
	fmt.Println(pass)

	body := []byte(name + pass + "\n")

	w.Write(body)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	cookie := r.Cookies()
	// make sure cookie exists

	// decode cookie

	// delete session ID

	// redirect to /

	fmt.Println(cookie)
}
