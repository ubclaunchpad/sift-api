package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetSessionById(t *testing.T) {
	s := Session{UserID: 1}

	ns := dm.Create(&s).Value
	defer dm.Unscoped().Delete(&s)

	id := (ns.(*Session)).ID
	if dm.First(&s).RecordNotFound() {
		t.Error("Record not found")
	}

	if sesh, err := dm.GetSessionByIdHelper(id); err != nil {
		t.Error("Failed to get session: ", err)
	} else {
		assert.Equal(t, id, sesh.ID)
	}
}

func TestCreateSession(t *testing.T) {

	var userID uint = 1337
	sesh := Session{UserID: userID}

	// Create sesh
	if err := dm.CreateSessionHelper(sesh); err != nil {
		t.Error("dm.CreateSessionHelper", err)
	}

	defer dm.Unscoped().Delete(&sesh)

	// Check sesh was added
	if seshRetrieved, err := dm.GetSessionByUserHelper(userID); err != nil {
		t.Error("dm.GetSessionByUserHelper", err)
	} else {
		assert.Equal(t, sesh.UserID, seshRetrieved.UserID)
	}

}

func TestSessionMiddlewareGood(t *testing.T) {

	// Create a profile

	prof := Profile{
		UserName:    "test_user",
		CompanyName: "test_company",
		PwHash:      []byte("1234"),
		Address:     "1234 lane",
	}

	if err := dm.Create(&prof).Error; err != nil {
		t.Error("Creation of profile failed with err: ", err)
	}

	defer dm.Unscoped().Delete(&prof)

	// Log the new user in

	formData := url.Values{
		"user_name":    {"test_user"},
		"company_name": {"test_company"},
		"pw_hash":      {"1234"},
	}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Login)
	handler.ServeHTTP(rr, req)

	// Get cookie from login response

	respWithCookie := http.Response{Header: rr.Header()}
	cookies := respWithCookie.Cookies()

	if len(cookies) != 1 {
		t.Errorf("Failed to get cookies from login")
	}

	// Get session so we can defer a delete of it when
	// test ends

	seshID, err := dm.DecodeCookieHelper(*cookies[0])

	if err != nil {
		t.Error("dm.DecodeCookieHelper", err)
	}

	defer dm.DeleteSessionByIdHelper(seshID)

	// Set up a mock handler function that will be wrapped by
	// our session MW and get a reference to the context
	// of incoming requests of this mock handler via closure

	var ctx context.Context

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		w.WriteHeader(http.StatusOK)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", mockHandler)
	rr = httptest.NewRecorder()
	routerWithMW := dm.SessionMiddleware(router)

	// Send a request through the whole chain
	// Request -> SessionMiddleware -> Router -> Mock handler

	req, err = http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.AddCookie(cookies[0])

	routerWithMW.ServeHTTP(rr, req)

	// Check that the mock handler function recieved the
	// user profile in the context of the request

	ctxProfile, ok := ctx.Value("profile").(*Profile)

	if !ok {
		t.Error("Failed to get profile from context in mock")
	}

	// Check profile in context is same as created one

	assert.Equal(t, prof.UserName, ctxProfile.UserName)
	assert.Equal(t, prof.CompanyName, ctxProfile.CompanyName)
	assert.Equal(t, prof.Address, ctxProfile.Address)
	assert.Equal(t, []byte(""), ctxProfile.PwHash)

}

func TestSessionMiddlewareNoCookie(t *testing.T) {

	// Set up a mock handler function that will be wrapped by
	// our session MW and get a reference to the context
	// of incoming requests of this mock handler via closure

	var ctx context.Context

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ctx = r.Context()
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", mockHandler)

	// Send a request through the whole chain
	// Request -> SessionMiddleware -> Router -> Mock handler

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	rr := httptest.NewRecorder()
	handler := dm.SessionMiddleware(router)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	profile := ctx.Value("profile")

	// Check that there is no profile in the context
	// of the request

	assert.Nil(t, profile)
	assert.Equal(t, http.StatusOK, rr.Code)

}
