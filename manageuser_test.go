package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
        
    "testing"
    "net/http/httptest"
)


func TestIndexNewProfileSuccess(t *testing.T) {

    formdata := url.Values{
        "user_name":        {"super_sifter"},
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

func TestIndexNewProfileMissingField(t *testing.T) {

    formdata := url.Values{
        "user_name":        {"super_sifter"},
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


func TestGetProfileHelperSuccess(t *testing.T) {
    un := "super_sifter"
    cn := "Sift Technologies, Inc."
    adr := "4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"
    // sha1 hash of: @BlueBalls123!
    pwh := []byte("cd026ec28d7976550a52da2520660bd8e26b5b40")
    p, err := dm.GetProfileHelper(un, cn)
    if err != nil {
        t.Errorf("Error (%v) encountered when retrieving profile for %s", err, un)
    }
    if p.UserName != un {
        t.Errorf("Incorrect user_name: received %s, expected %s", p.UserName, un)
    }
    if p.CompanyName != cn {
        t.Errorf("Incorrect company_name: received %s, expected %s", p.CompanyName, cn)
    }
    if p.Address != adr {
        t.Errorf("Incorrect company_address: received %s, expected %s", p.Address, adr)
    }
    if string(p.PwHash) != string(pwh) {
        t.Errorf("Incorrect pw_hash: received %s, expected %s", p.PwHash, pwh)
    }
}

func TestGetProfileHelperIncorrectUserName(t *testing.T) {
    un := "soggy_sifter"
    cn := "Sift Technologies, Inc."
    // sha1 hash of: @BlueBalls123!
    // pwh := []byte("@BlueBalls123!")
    _, err := dm.GetProfileHelper(un, cn)
    if err == nil {
        t.Error("Error was nil, expected non-nil error and no profile to be returned")
    }
}

func TestGetProfileHelperIncorrectCompanyName(t *testing.T) {
    un := "super_sifter"
    cn := "Sift Technologies, LLC"
    // Equivalent of and empty password
    // pwh := []byte("")
    _, err := dm.GetProfileHelper(un, cn)
    if err == nil {
        t.Error("Error was nil, expected non-nil error and no profile to be returned")
    }
}

func TestUpdateExistingProfileSuccess(t *testing.T) {
    un := url.QueryEscape("super_sifter")
    cn := url.QueryEscape("Sift Technologies, Inc.")
    furl := fmt.Sprintf("/profile/%s/%s", cn, un)
    formdata := url.Values{
        "user_name":        {"solid_sifter"}, // new user_name
        "company_name":     {"Sift Technologies, LLC"}, // new company_name is ignored
        "company_address":  {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
        // sha1 hash of: @BlueBalls123!
        "pw_hash":          {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
    }
    req, err := http.NewRequest("PUT", furl, strings.NewReader(formdata.Encode()))
    if err != nil {
        t.Errorf("Request unsuccessfully created.\nerr: %v", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.UpdateExistingProfile)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusNoContent {
        t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusNoContent)
    }

}

func TestDeleteExistingProfileSuccess(t *testing.T) {
    un := url.QueryEscape("solid_sifter")
    cn := url.QueryEscape("Sift Technologies, Inc.")
    furl := fmt.Sprintf("/profile/%s/%s", cn, un)
    req, err := http.NewRequest("DELETE", furl, nil)
    if err != nil {
        t.Errorf("Request unsuccessfully created.\nerr: %v", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.DeleteExistingProfile)
    handler.ServeHTTP(rr, req)
    if dm.userExists(un, cn) {
        t.Errorf("User still exists but should have been deleted")
    }
    if rr.Code != http.StatusNoContent {
        t.Errorf("HTTP status code recieved: %d expected %d\nerror: %v", rr.Code, http.StatusNoContent, rr.Body)
    }
}

func TestDeleteExistingProfileNotExistFail(t *testing.T) {
    un := url.QueryEscape("solid_sifter")
    cn := url.QueryEscape("Sift Technologies, BARF")
    furl := fmt.Sprintf("/profile/%s/%s", cn, un)
    req, err := http.NewRequest("DELETE", furl, nil)
    if err != nil {
        t.Errorf("Request unsuccessfully created.\nerr: %v", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.DeleteExistingProfile)
    handler.ServeHTTP(rr, req)
    un, _ = url.QueryUnescape(un)
    cn, _ = url.QueryUnescape(cn)
    if dm.userExists(un, cn) {
        t.Errorf("User should not exist")
    }
    if rr.Code != http.StatusInternalServerError {
        t.Errorf("HTTP status code recieved: %d expected %d\nerror: %v", rr.Code, http.StatusInternalServerError, rr.Body)
    }
}