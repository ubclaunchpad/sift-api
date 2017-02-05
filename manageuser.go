// Database API for user administration utilities, ex. creating a Profile
package main

import (
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
)

// Struct type to manage the user info database
type DataManager struct {
    DB  *gorm.DB
}

// IndexNewProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain a new Profile, and creates a db 
// record of this new profile. Writes the result of that query, ID or error, 
// to rw.
func (dm *DataManager) IndexNewProfile(rw http.ResponseWriter, req *http.Request) {
    return
}

// GetExistingProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain the ID of an existing user profile
// record, and writes the result of that query, user Profile or error, to rw.
func (dm *DataManager) GetExistingProfile(rw http.ResponseWriter, req *http.Request) {
    return
}

// UpdateExistingProfile operates on a DataManager struct and takes a
// ResponseWriter and a Request, whose body should contain the user's unique ID,
// and updates the existing record for that user in the db. Either an error or
// success message on record update will be written to rw.
func (dm *DataManager) UpdateExistingProfile(rw http.ResponseWriter, req *http.Request) {
    return
}

// DeleteExistingProfile operates on a DataManager struct and takes a user's 
// unique ID and deletes the corresponding db record. Either an error or 
// success message onm successful deletion of the record will be written to rw.
func (dm *DataManager) DeleteExistingProfile(rw http.ResponseWriter, req *http.Request) {
    return
}
