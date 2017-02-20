// Database API for user administration utilities, ex. creating a Profile
package main

import (
    "log"
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
)

// Struct type to manage the user info database
type DataManager struct {
    *gorm.DB
}

// IndexNewProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain a new Profile, and creates a db 
// record of this new profile. Writes the result of that query, ID or error, 
// to w.
func (dm *DataManager) IndexNewProfile(w http.ResponseWriter, r *http.Request) {
    // Handle error here as PostFormValue ignores errors
    if err := r.ParseForm(); err != nil {
        log.Fatal("r.ParseForm: ", err)
        return
    }
    p := Profile{
        // new(gorm.Model),
        CompanyName:    r.PostFormValue("company_name"),
        PwHash:         []byte(r.PostFormValue("pw_hash")),
        Address:        r.PostFormValue("company_address"),
    }
    if p.CompanyName == "" || p.PwHash == nil || p.Address == "" {
        http.Error(w, "One or more profile data fields were blank.", http.StatusBadRequest)
        return
    }
    // Create a new profile record
    if err := dm.Create(&p).Error; err != nil {
        http.Error(w, "Database error on profile indexing.", http.StatusInternalServerError)
        log.Fatal("dm.Create: ", err)
        return
    }
    // Profile record created successfully, send OK
    w.WriteHeader(http.StatusOK)
}

func (dm *DataManager) GetProfileHelper(cn, pwh string) (Profile, error) {
    var prof Profile
    d := dm.Where("company_name = ? AND pw_hash = ?", cn, pwh).Find(&prof)
    if d.Error != nil {
        log.Fatal("dm.Where: ", d.Error)
        return Profile{}, d.Error
    }
    return prof, nil
}

// GetExistingProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain the ID of an existing user profile
// record, and writes the result of that query, user Profile or error, to w.
func (dm *DataManager) GetExistingProfile(w http.ResponseWriter, r *http.Request) {
    return
}

// UpdateExistingProfile operates on a DataManager struct and takes a
// ResponseWriter and a Request, whose body should contain the user's unique ID,
// and updates the existing record for that user in the db. Either an error or
// success message on record update will be written to w.
func (dm *DataManager) UpdateExistingProfile(w http.ResponseWriter, r *http.Request) {
    return
}

// DeleteExistingProfile operates on a DataManager struct and takes a user's 
// unique ID and deletes the corresponding db record. Either an error or 
// success message onm successful deletion of the record will be written to w.
func (dm *DataManager) DeleteExistingProfile(w http.ResponseWriter, r *http.Request) {
    return
}
