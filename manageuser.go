// Database API for user administration utilities, ex. creating a Profile
package main

import (
    "log"
    "fmt"
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
)

// Struct type to manage the user info database
type DataManager struct {
    *gorm.DB
}

const (
    MAX_PROF_SIZE = 1 << 20 // 1 MiB
)

// IndexNewProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain a new Profile, and creates a db 
// record of this new profile. Writes the result of that query, ID or error, 
// to w.
func (dm *DataManager) IndexNewProfile(w http.ResponseWriter, r *http.Request) {
    // Handle error here as PostFormValue ignores errors
    if err := r.ParseMultipartForm(MAX_PROF_SIZE); err != nil {
        log.Fatal(err)
        return
    }
    p := Profile{
        // new(gorm.Model),
        CompanyName:    r.PostFormValue("company_name"),
        PwHash:         []byte(r.PostFormValue("pw_hash")),
        Address:        r.PostFormValue("company_address"),
    }
    if p.CompanyName == "" || p.PwHash == nil || p.Address == "" {
        fmt.Printf("Could not index new profile, at least one field was blank:\ncompany_name: %v\npw_hash: %v\ncompany_address: %v\n",
            p.CompanyName, p.PwHash, p.Address)
        return
    }
    // Create a new profile record and check for 
    if err := dm.Create(&p).Error; err != nil {
        log.Fatal(err)
        return
    }
    // Create response via ResponseWriter
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    body := []byte(`{"Message":"Profile successfully created."}`)
    w.Write(body)
}

func (dm *DataManager) GetProfileHelper(cn, pwh string) (Profile, error) {
    var prof Profile
    d := dm.Where("company_name = ? AND pw_hash = ?", cn, pwh).Find(&prof)
    if d.Error != nil {
        log.Fatal(err)
        return nil, err
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
