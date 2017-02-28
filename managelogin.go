package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Login takes a request with a username, company name, and password hash,
// it then authenticates the credentials, creates a session for the user
// if successful, then redirects the user to the landing page with a cookie
// attached containing the session ID
func (dm *DataManager) Login(w http.ResponseWriter, r *http.Request) {

	// Get credentials from request form values
	name := r.FormValue("user_name")
	company := r.FormValue("company_name")
	pass := r.FormValue("pw_hash")

	// Make sure credentials not empty
	if name == "" || company == "" || pass == "" {
		http.Error(w, "One or more credentials were blank", http.StatusBadRequest)
		return
	}

	passHash := []byte(pass)

	// Failed authentication
	if !dm.UserPwAuthSuccess(name, company, passHash) {
		http.Error(w, "Failed authentication", http.StatusUnauthorized)
		return
	}

	// Get user profile from DB
	profile, err := dm.GetProfileHelper(name, company)

	if err != nil {
		log.Fatal("dm.GetProfileHelper", err)
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Clear out any existing sessions for user
	err = dm.DeleteSessionsByUserHelper(profile.ID)

	if err != nil {
		log.Fatal("dm.DeleteSessionsByUserHelper", err)
		http.Error(w, "Database error on clearing sessions for user login", http.StatusInternalServerError)
		return
	}

	// Create new session for user

	sesh := Session{UserID: profile.ID}

	err = dm.CreateSessionHelper(sesh)

	if err != nil {
		log.Fatal("dm.CreateSessionHelperHelper", err)
		http.Error(w, "Database error on creating new session for login", http.StatusInternalServerError)
		return
	}

	// Retrieve created session to get session ID

	createdSesh, err := dm.GetSessionByUserHelper(profile.ID)

	if err != nil {
		log.Fatal("dm.GetSessionByUserHelper", err)
		http.Error(w, "Database error on creating new session for login", http.StatusInternalServerError)
		return
	}

	// Create cookie with encoded session ID

	cookieValue := map[string]uint{
		"id": createdSesh.ID,
	}

	encoded, err := dm.CookieManager.Encode("session", cookieValue)

	if err != nil {
		log.Fatal("dm.CookieManager.Encode", err)
		http.Error(w, "Error creating cookie for user", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}

	// Attach cookie to response and redirect to landing

	http.SetCookie(w, cookie)

	redirect := "/dashboard"
	http.Redirect(w, r, redirect, http.StatusFound)

}

// Logout takes a request with a cookie in the header and deletes the
// session found in the decoded cookie
func (dm *DataManager) Logout(w http.ResponseWriter, r *http.Request) {

	cookies := r.Cookies()
	var seshCookie http.Cookie

	switch len(cookies) {
	case 0:
		// no cookie present, user must be already logged out
		http.Redirect(w, r, "/", http.StatusFound)
		return
	case 1:
		// has cookie
		seshCookie = *cookies[0]
	default:
		log.Fatal("Length of request cookie array is not 0 or 1")
		http.Error(w, "Error logging user out", http.StatusInternalServerError)
		return
	}

	// decode cookie

	seshID, err := dm.DecodeCookieHelper(seshCookie)

	if err != nil {
		log.Fatal("dm.DecodeCookieHelper", err)
		http.Error(w, "Could not get session ID from cookie", http.StatusBadRequest)
		return
	}

	// delete session ID

	err = dm.DeleteSessionByIdHelper(seshID)

	if err != nil {
		log.Fatal("dm.DeleteSessionByIdHelper", err)
		http.Error(w, "Error clearing sessions for logout", http.StatusInternalServerError)
		return
	}

	redirect := "/"

	http.Redirect(w, r, redirect, http.StatusFound)
}

// GetProfileFromCookie takes a request with a cookie in the header
// checks if the user is logged in and returns the profile if
// they are, else it returns a 401 Not Authorized
func (dm *DataManager) GetProfileFromCookie(w http.ResponseWriter, r *http.Request) {

	cookies := r.Cookies()
	var seshCookie http.Cookie

	switch len(cookies) {
	case 0:
		// no cookie present
		http.Error(w, "No cookie in request header", http.StatusBadRequest)
		return
	case 1:
		// has cookie
		seshCookie = *cookies[0]
	default:
		log.Fatal("Length of request cookie array is not 0 or 1")
		http.Error(w, "Error logging user out", http.StatusInternalServerError)
		return
	}

	// decode cookie

	seshID, err := dm.DecodeCookieHelper(seshCookie)

	if err != nil {
		log.Fatal("dm.DecodeCookieHelper", err)
		http.Error(w, "Could not get session ID from cookie", http.StatusBadRequest)
		return
	}

	sesh, err := dm.GetSessionByIdHelper(seshID)

	// check if user is logged in
	if err != nil {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}

	// get profile object
	Profile, err := dm.GetProfileByIdHelper(sesh.UserID)

	if err != nil {
		log.Fatal("dm.GetProfileByIdHelper", err)
		http.Error(w, "Could not find user from that session", http.StatusInternalServerError)
		return
	}

	// strip password
	Profile.PwHash = []byte("")

	body, err := json.Marshal(Profile)
	if err != nil {
		log.Fatal("json.Marshal: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
