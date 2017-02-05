package main

import (
	"time"
)

// Profile struct used to hold user data
type Profile struct {
	ID          uint64
	CompanyName string
	PwHash      []byte
	Address     string
	Created		time.Time
}
