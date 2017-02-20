package main

import (
    "net/http"
    "net/url"
    "strings"
        
    "testing"
    "net/http/httptest"
)


func TestIndexNewProfileOK(t *testing.T) {

    formdata := url.Values{
        "company_name":     {"Sift Technologies, Inc."},
        "company_address":  {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
        // sha1 hash of: @BlueBalls123!
        "pw_hash":          {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
    }
    req, err := http.NewRequest("POST", "/profile", strings.NewReader(formdata.Encode()))
    if err != nil {
        t.Errorf("Request unsuccessfully created.\nerr: %v", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.IndexNewProfile)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusOK)
    }
}

func TestIndexNewProfileFail(t *testing.T) {

    formdata := url.Values{
        // Missing company_name
        "company_address":  {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
        // sha1 hash of: @BlueBalls123!
        "pw_hash":          {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
    }
    req, err := http.NewRequest("POST", "/profile", strings.NewReader(formdata.Encode()))
    if err != nil {
        t.Errorf("Request unsuccessfully created.\nerr: %v", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.IndexNewProfile)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusBadRequest)
    }
}


func TestGetProfileHelper(t *testing.T) {
    
}