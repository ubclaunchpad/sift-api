package main

// Profile struct used to hold user data
type Profile struct {
	id          uint64
	companyName string
	pwHash      []byte
	address     string
}
