package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandlerGoodLogin(t *testing.T) {

	formData := url.Values{"name": {"test_user"}, "pass": {"123123"}}

	req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	respWithCookie := http.Response{Header: rr.Header()}
	cookie := respWithCookie.Cookies()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)
	fmt.Println(cookie)

}

func TestLoginHandlerBadNamePass(t *testing.T) {

	req, err := http.NewRequest("POST", "/login", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

}
