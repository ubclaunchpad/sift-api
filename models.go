package main

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	CompanyName string
	PwHash      []byte
	Address     string
}

type Session struct {
	gorm.Model
	UserID uint
}
