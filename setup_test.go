package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
    "testing"
)

var cfg = DBConfig{
    DBUser:		"postgres",
    DBPassword:	"testpw",
    DBHost:		"localhost",
    DBName:		"sift_user_data",
    DBSSLType:	"disable",
}

var db, _ = gorm.Open("postgres", cfg.createDBQueryString())
var dm = DataManager{db}

// Setup and teardown for all tests
func TestMain(m *testing.M) {
    // Add new models here
	dm.AutoMigrate(&Session{})
	dm.AutoMigrate(&Profile{})
	defer dm.Close()
	m.Run()
}

func TestConnectDB(t *testing.T) {
	var _, err = gorm.Open("postgres", cfg.createDBQueryString())
	if err != nil {
		t.Fail()
	}
}

