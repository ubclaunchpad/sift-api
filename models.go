package main

import "github.com/jinzhu/gorm"

// UserName and CompanyName must be unique in combination
type Profile struct {
	gorm.Model
	UserName	string `gorm:"primary_key"`
	CompanyName string `gorm:"primary_key"`
	PwHash      []byte
	Address     string
}

type Session struct {
	gorm.Model
	UserID uint
}
