package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"net/http/httptest"
	"testing"
)

/* -------------------------- HELPER METHOD TESTS --------------------------- */

func TestUserExists(t *testing.T) {
	prof := Profile{
		UserName:    "test",
		CompanyName: "test company",
		PwHash:      []byte("1234"),
		Address:     "1234 lane",
	}
	if err := db.Create(&prof).Error; err != nil {
		t.Errorf("Error creating profile not expected. err: %v", err)
	}
	defer db.Unscoped().Delete(&Profile{})
	if !dm.userExists("test", "test company") {
		t.Error("User should exist but does not.")
	}
	if dm.userExists("not test", "test company") {
		t.Error("User should not exist but does.")
	}
}

func TestParseProfileQuerySuccess(t *testing.T) {
	un, cn := "test", "test company"
	furl := url.QueryEscape(fmt.Sprintf("/profile/%s/%s", cn, un))
	rsrc := parseProfileQuery(furl)
	if len(rsrc) != 4 {
		t.Errorf("Number of query URL arguments should be 4, was %d", len(rsrc))
	}
	if rsrc[2] != cn {
		t.Errorf("Submitted company_name %s != %s.", cn, rsrc[2])
	}
	if rsrc[3] != un {
		t.Errorf("Submitted user_name %s != %s.", un, rsrc[3])
	}
}

func TestParseProfileQueryEmpty(t *testing.T) {
	rsrc := parseProfileQuery(url.QueryEscape("/"))
	if len(rsrc) != 2 {
		t.Errorf("Number of query URL arguments should be 2, was %d", len(rsrc))
	}
	if rsrc[0] != "" {
		t.Errorf("rsrc[0] expected to be empty, got %s.", rsrc[0])
	}
	if rsrc[1] != "" {
		t.Errorf("rsrc[1] expected to be empty, got %s.", rsrc[1])
	}
}

func TestGetProfileHelperSuccess(t *testing.T) {
	un := "super_sifter"
	cn := "Sift Technologies, Inc."
	adr := "4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"
	pwh := []byte("cd026ec28d7976550a52da2520660bd8e26b5b40")
	prof := Profile{
		UserName:    un,
		CompanyName: cn,
		Address:     adr,
		PwHash:      pwh,
	}
	if err := db.Create(&prof).Error; err != nil {
		t.Errorf("Error creating profile not expected. err: %v", err)
	}
	defer db.Unscoped().Delete(&Profile{})
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

func TestGetProfileHelperIncorrectInputs(t *testing.T) {
	prof := Profile{
		UserName:    "super_sifter",
		CompanyName: "Sift Technologies, Inc.",
	}
	if err := db.Create(&prof).Error; err != nil {
		t.Errorf("Error creating profile not expected. err: %v", err)
	}
	defer db.Unscoped().Delete(&Profile{})
	_, err := dm.GetProfileHelper("soggy_sifter", "Sift Technologies, Inc.")
	if err == nil {
		t.Error("Error was nil, expected non-nil error and no profile to be returned")
	}
	_, err = dm.GetProfileHelper("super_sifter", "Sift Technologies, LLC")
	if err == nil {
		t.Error("Error was nil, expected non-nil error and no profile to be returned")
	}
}

func TestUserPwAuthSuccess(t *testing.T) {
	un := "super_sifter"
	cn := "Sift Technologies, Inc."
	pwh := []byte("abc123")
	prof := Profile{
		UserName:    un,
		CompanyName: cn,
		PwHash:      pwh,
	}
	if err := db.Create(&prof).Error; err != nil {
		t.Errorf("Error creating profile not expected. err: %v", err)
	}
	defer db.Unscoped().Delete(&Profile{})
	if !dm.UserPwAuthSuccess(un, cn, []byte("abc123")) {
		t.Error("Expected user to log in successfully")
	}
	if dm.UserPwAuthSuccess(un, cn, []byte("")) {
		t.Error("Expected user login to be rejected with empty password hash")
	}
	if dm.UserPwAuthSuccess(un, cn, []byte("xyz123")) {
		t.Error("Expected user login to be rejected with bad password hash")
	}
}

/* -------------------------------------------------------------------------- */

func TestIndexNewProfileSuccess(t *testing.T) {
	formdata := url.Values{
		"user_name":       {"super_sifter"},
		"company_name":    {"Sift Technologies, Inc."},
		"company_address": {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
		"pw_hash":         {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
	}
	req, err := http.NewRequest("POST", "/profile", strings.NewReader(formdata.Encode()))
	if err != nil {
		t.Errorf("Request unsuccessfully created.\nerr: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.IndexNewProfile)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusCreated)
	}
	var p Profile
	err = json.Unmarshal(rr.Body.Bytes(), &p)
	if err != nil {
		t.Errorf("Error (%v) encountered when unmarshalling profile", err)
	}
	if p.UserName != formdata["user_name"][0] {
		t.Errorf("Incorrect user_name: received %s, expected %s", p.UserName, formdata["user_name"][0])
	}
	if p.CompanyName != formdata["company_name"][0] {
		t.Errorf("Incorrect company_name: received %s, expected %s", p.CompanyName, formdata["company_name"][0])
	}
	if p.Address != formdata["company_address"][0] {
		t.Errorf("Incorrect company_address: received %s, expected %s", p.Address, formdata["company_address"][0])
	}
	if string(p.PwHash) != formdata["pw_hash"][0] {
		t.Errorf("Incorrect pw_hash: received %s, expected %s", p.PwHash, formdata["pw_hash"][0])
	}
}

func TestIndexNewProfileMissingField(t *testing.T) {

	formdata := url.Values{
		"user_name": {"super_sifter"},
		// Missing company_name
		"company_address": {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
		"pw_hash":         {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
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

func TestUpdateExistingProfileSuccess(t *testing.T) {
	un := "super_sifter"
	cn := "Sift Technologies, Inc."
	furl := url.QueryEscape(fmt.Sprintf("/profile/%s/%s", cn, un))
	formdata := url.Values{
		"user_name":       {"solid_sifter"},           // new user_name
		"company_name":    {"Sift Technologies, LLC"}, // new company_name is ignored
		"company_address": {"4321 Pleasantown Rd, Pleasantville, PV, UPV, V1A 1X1"},
		"pw_hash":         {"cd026ec28d7976550a52da2520660bd8e26b5b40"},
	}
	req, err := http.NewRequest("PUT", furl, strings.NewReader(formdata.Encode()))
	if err != nil {
		t.Errorf("Request unsuccessfully created.\nerr: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dm.UpdateExistingProfile)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusOK)
	}
	var p Profile
	err = json.Unmarshal(rr.Body.Bytes(), &p)
	if err != nil {
		t.Errorf("Error (%v) encountered when unmarshalling profile", err)
	}
	if p.UserName != formdata["user_name"][0] {
		t.Errorf("Incorrect user_name: received %s, expected %s", p.UserName, formdata["user_name"][0])
	}
	// Should not have changed as UpdateExistingProfile ignores updates to company_name
	if p.CompanyName != cn {
		t.Errorf("Incorrect company_name: received %s, expected %s", p.CompanyName, cn)
	}
	if p.Address != formdata["company_address"][0] {
		t.Errorf("Incorrect company_address: received %s, expected %s", p.Address, formdata["company_address"][0])
	}
	if string(p.PwHash) != formdata["pw_hash"][0] {
		t.Errorf("Incorrect pw_hash: received %s, expected %s", p.PwHash, formdata["pw_hash"][0])
	}
}

func TestDeleteExistingProfileSuccess(t *testing.T) {
	un := "solid_sifter"
	cn := "Sift Technologies, Inc."
	furl := url.QueryEscape(fmt.Sprintf("/profile/%s/%s", cn, un))
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
	un := "solid_sifter"
	cn := "Sift Technologies, BARF"
	furl := url.QueryEscape(fmt.Sprintf("/profile/%s/%s", cn, un))
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
	if rr.Code != http.StatusBadRequest {
		t.Errorf("HTTP status code recieved: %d expected %d\nerror: %v", rr.Code, http.StatusBadRequest, rr.Body)
	}
}
