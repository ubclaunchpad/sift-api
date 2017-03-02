package main

import (
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
)

// DataManager contains the DB connection and functions
// for managing and handling request to the database
type DataManager struct {
	*gorm.DB
	*securecookie.SecureCookie
}

// NewDataManager constructs the DataManager
// struct and initializes the secure cookie keys
func NewDataManager(db *gorm.DB) DataManager {
	key := securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))

	dm := DataManager{db, key}
	return dm
}

// UserName and CompanyName must be unique in combination
type Profile struct {
	gorm.Model
	UserName    string `gorm:"primary_key"`
	CompanyName string `gorm:"primary_key"`
	PwHash      []byte
	Address     string
}

type Session struct {
	gorm.Model
	UserID uint
}
