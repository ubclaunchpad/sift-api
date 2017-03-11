package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginGood(t *testing.T) {

	prof := Profile{
		UserName:    "super_sifter",
		CompanyName: "Sift Technologies, Inc.",
		Address:     "4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1",
		PwHash:      []byte("cd026ec28d7976550a52da2520660bd8e26b5b40"),
	}

	if err := dm.Create(&prof).Error; err != nil {
		t.Error("Creation of profile failed with err: ", err)
	}

	formData := url.Values{
		"user_name":       {"super_sifter"},
		"company_name":    {"Sift Technologies, Inc."},
		"company_address": {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
		// sha1 hash of: @BlueBalls123!
		"pw_hash": {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
	}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Login)
	handler.ServeHTTP(rr, req)

	respWithCookie := http.Response{Header: rr.Header()}
	cookie := respWithCookie.Cookies()

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.True(t, len(cookie) > 0)

}

func TestLoginBadNamePass(t *testing.T) {

	req, err := http.NewRequest("POST", "/login", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Login)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)

}

func TestLogoutNoCookie(t *testing.T) {

	req, err := http.NewRequest("GET", "/logout", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Logout)
	handler.ServeHTTP(rr, req)

	// should get a redirect
	assert.Equal(t, rr.Code, http.StatusFound)

}

func TestLogoutGood(t *testing.T) {

	// first we login
	formData := url.Values{
		"user_name":    {"super_sifter"},
		"company_name": {"Sift Technologies, Inc."},
		// sha1 hash of: @BlueBalls123!
		"pw_hash": {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
	}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Login)
	handler.ServeHTTP(rr, req)

	respWithCookie := http.Response{Header: rr.Header()}
	cookie := respWithCookie.Cookies()

	// check we got a cookie back
	assert.Equal(t, rr.Code, http.StatusFound)
	assert.True(t, len(cookie) > 0)

	// check we got a session
	profile, err := dm.GetProfileHelper("super_sifter", "Sift Technologies, Inc.")

	if err != nil {
		t.Error("dm.GetProfileHelper", err)
	}

	sesh, err := dm.GetSessionByUserHelper(profile.ID)

	if err != nil {
		t.Error("dm.GetSessionByUserHelper", err)
	}

	assert.True(t, sesh.ID > 0)
	assert.Equal(t, sesh.UserID, profile.ID)

	// now we logout

	req, err = http.NewRequest("GET", "/logout", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.AddCookie(cookie[0])
	handler = http.HandlerFunc(dm.Logout)
	handler.ServeHTTP(rr, req)

	// check we got a redirect
	assert.Equal(t, rr.Code, http.StatusFound)

	// check that our session was deleted
	deletedSesh, err := dm.GetSessionByUserHelper(profile.ID)
	assert.Equal(t, deletedSesh.ID, uint(0))
	assert.NotNil(t, err)

}

func TestGetProfileFromCookie(t *testing.T) {

	// first we login
	formData := url.Values{
		"user_name":    {"super_sifter"},
		"company_name": {"Sift Technologies, Inc."},
		// sha1 hash of: @BlueBalls123!
		"pw_hash": {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
	}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.Login)
	handler.ServeHTTP(rr, req)

	respWithCookie := http.Response{Header: rr.Header()}
	cookie := respWithCookie.Cookies()

	// check we got a cookie back from login
	assert.Equal(t, rr.Code, http.StatusFound)
	assert.True(t, len(cookie) > 0)

	// check we got a session
	profile, err := dm.GetProfileHelper("super_sifter", "Sift Technologies, Inc.")

	if err != nil {
		t.Error("dm.GetProfileHelper", err)
	}

	sesh, err := dm.GetSessionByUserHelper(profile.ID)

	if err != nil {
		t.Error("dm.GetSessionByUserHelper", err)
	}

	assert.True(t, sesh.ID > 0)
	assert.Equal(t, sesh.UserID, profile.ID)

	// now we actually test getting profile from our
	// newly acquired cookie

	req, err = http.NewRequest("GET", "/auth", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	rr = httptest.NewRecorder()
	req.AddCookie(cookie[0])
	handler = http.HandlerFunc(dm.GetProfileFromCookie)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// get profile from body of response

	resProfile := Profile{}
	err = json.Unmarshal(rr.Body.Bytes(), &resProfile)
	if err != nil {
		t.Error("json.Unmarshal", err)
	}

	assert.Equal(t, profile.ID, resProfile.ID)
	assert.Equal(t, profile.UserName, resProfile.UserName)
	assert.Equal(t, profile.CompanyName, resProfile.CompanyName)
	assert.Equal(t, 0, bytes.Compare(resProfile.PwHash, []byte("")))

}

func TestGetProfileFromCookieUnauthorized(t *testing.T) {

	// make up a bad session ID and encode in cookie

	cookie, err := dm.CreateCookieHelper(99999)

	if err != nil {
		t.Error("dm.CreateCookieHelper", err)
	}

	// try to get profile using bad session
	req, err := http.NewRequest("GET", "/auth", nil)

	if err != nil {
		t.Error("http.NewRequest", err)
	}

	rr := httptest.NewRecorder()
	req.AddCookie(&cookie)
	handler := http.HandlerFunc(dm.GetProfileFromCookie)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

}
