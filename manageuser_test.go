package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
    "testing"
    "net/http/httptest"
)

// const (
//     ROOT_URL = "http://127.0.0.1:9090"
// )

var dbconf = DBConfig{
    DBUser:		"test",
    DBPassword:	"testpw",
    DBHost:		"localhost",
    DBName:		"sift_user_info",
    DBSSLType:	"disable",
}

var db, _ = gorm.Open("postgres", cfg.createDBQueryString())
var dm = DataManager{db}


func TestIndexNewProfile(t *testing.T) {

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
    req.Header.Set("Content-Type", "multipart/form-data")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(dm.IndexNewProfile)
    handler.ServeHTTP(rr, req)
    if rr.Code != 200 {
        t.Errorf("HTTP status code not 200, was %d", rr.Code)
    }
    fmt.Printf("Response was correct: %s", rr.Body)
}

func TestGetProfileHelper(t *testing.T) {
    
}