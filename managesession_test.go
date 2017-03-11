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
	id := (ns.(*Session)).ID
	if dm.First(&s).RecordNotFound() {
		t.Error("Record not found")
	}

	if sesh, err := dm.GetSessionByIdHelper(id); err != nil {
		t.Log("Failed to get session: ", err)
		t.Fail()
	} else {
		assert.Equal(t, id, sesh.ID)
	}
}

func TestCreateSession(t *testing.T) {

	var userID uint = 1337
	sesh := Session{UserID: userID}

	// Create sesh
	if err := dm.CreateSessionHelper(sesh); err != nil {
		t.Log("dm.CreateSessionHelper", err)
		t.Fail()
	}

	// Check sesh was added
	if seshRetrieved, err := dm.GetSessionByUserHelper(userID); err != nil {
		t.Log("dm.GetSessionByUserHelper", err)
		t.Fail()
	} else {
		assert.Equal(t, sesh.UserID, seshRetrieved.UserID)
	}

	// Delete sesh and check it was deleted
	if err := dm.DeleteSessionsByUserHelper(userID); err != nil {
		t.Log("dm.DeleteSessionsByUserHelper", err)
		t.Fail()
	}

	if seshRetrieved, err := dm.GetSessionByUserHelper(userID); err == nil {
		t.Log("Session was not deleted")
		t.Fail()
	} else {
		assert.Equal(t, seshRetrieved.ID, uint(0))
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

	// Log the user in
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

	// Get cookie from login
	respWithCookie := http.Response{Header: rr.Header()}
	cookies := respWithCookie.Cookies()

	if len(cookies) == 0 {
		t.Errorf("Failed to get cookies from login")
	}

	// Set up a inner handler function and get a reference to the context
	// of incoming requests via closure
	var ctx context.Context

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		w.WriteHeader(http.StatusOK)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", fakeHandler)
	rr = httptest.NewRecorder()
	routerWithMW := dm.SessionMiddleware(router)

	// Send a request through the whole chain
	req, err = http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Log("http.NewRequest", err)
		t.Fail()
	}

	req.AddCookie(cookies[0])

	routerWithMW.ServeHTTP(rr, req)

	// Check that user profile was attached to the context
	// of the inner handler function

	ctxProfile, ok := ctx.Value("profile").(*Profile)

	if !ok {
		t.Error("Failed to get profile from context in inner handler")
	}

	// Check profile in context is same as created one
	assert.Equal(t, prof.UserName, ctxProfile.UserName)
	assert.Equal(t, prof.CompanyName, ctxProfile.CompanyName)
	assert.Equal(t, prof.Address, ctxProfile.Address)
	assert.Equal(t, []byte(""), ctxProfile.PwHash)

}

func TestSessionMiddlewareNoCookie(t *testing.T) {

	var ctx context.Context

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ctx = r.Context()
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", fakeHandler)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Log("http.NewRequest", err)
		t.Fail()
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := dm.SessionMiddleware(router)
	handler.ServeHTTP(rr, req)

	profile := ctx.Value("profile")

	assert.Nil(t, profile)
	assert.Equal(t, http.StatusOK, rr.Code)

}
