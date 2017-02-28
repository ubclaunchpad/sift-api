// Database API for user administration utilities, ex. creating a Profile
package main

import (
    "log"
    "net/http"
    "net/url"
    "encoding/json"
    "strings"
    "errors"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    
)

// Struct type to manage the user info database
type DataManager struct {
    *gorm.DB
}

/* ----------------------------- HELPER METHODS ----------------------------- */

// Tests if user exists already. Returns true if the combination of user_name 
// and company_name is already taken, false otherwise. Assumes un and cn are not
// empty strings.
func (dm *DataManager) userExists(un, cn string) bool {
    qstring := "user_name = ? AND company_name = ?"
    return !dm.Where(qstring, un, cn).First(&Profile{}).RecordNotFound()
}

// Helper to extract company_name and user_name from requests to resources with
// the form: /profile/{company_name}/{user_name}
func parseProfileQuery(path string) []string {
    upath, err := url.QueryUnescape(path)
    if err != nil {
        log.Fatal("url.QueryUnescape: ", err)
        return nil
    }
    params := strings.Split(upath, "/")
    resources := make([]string, 0)
    for _, val := range params {
        resources = append(resources, val)
    }
    return resources
}

// Queries db for user profile matching user_name and company_name. 
// Assumes un, cn, and pwh are not empty strings or nil.
func (dm *DataManager) GetProfileHelper(un, cn string) (Profile, error) {
    var prof Profile
    qstring := "user_name = ? AND company_name = ?"
    if err := dm.Where(qstring, un, cn).First(&prof).Error; err != nil {
        return Profile{}, err
    }
    return prof, nil
}

// Queries db for user profile matching user's id. Useful once user is logged in.
// Assumes id is not an empty string or nil.
func (dm *DataManager) GetProfileByIdHelper(id uint) (Profile, error) {
    var prof Profile
    qstring := "id = ?"
    if err := dm.Where(qstring, id).First(&prof).Error; err != nil {
        return Profile{}, err
    }
    return prof, nil
}

// Authenticates users by company_name, user_name, and pw_hash. Returns true if
// record is found and fields match, false otherwise. If false, user did not enter
// the correct credentials.
func (dm *DataManager) UserPwAuthSuccess(un, cn string, pwh []byte) bool {
    qstring := "user_name = ? AND company_name = ? AND pw_hash = ?"
    return !dm.Where(qstring, un, cn, pwh).First(&Profile{}).RecordNotFound()
}

// Updates given use field based on field name, which must be one of {user_name,
// company_name, company_address, pw_hash}. Can be used for changing passwords.
// Assumes user is already logged in.
// NOTE: Validation of password strength when changing passwords should occur in frontend
func (dm *DataManager) UpdateProfileHelper(un, cn, key string, val interface{}) error  {
    switch key {
    case "company_name":
        // company_name is final and cannot be updated
        return errors.New("company_name cannot be changed.")        
    case "user_name":
        // A user cannot change their user_name to conflict with an existing
        // user_name/company_name combination
        if dm.userExists(val.(string), cn) {
            return errors.New("user_name and company_name combination already exists")
        }
    case "company_address":
    case "pw_hash":
    default:
        return errors.New("Provided field is invalid.")
    }
    qstring := "user_name = ? AND company_name = ?"
    if err := dm.Model(&Profile{}).Where(qstring, un, cn).Update(key, val).Error; err != nil {
        return err
    }
    return nil
}

/* -------------------------------------------------------------------------- */

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
        UserName:       r.PostFormValue("user_name"),
        CompanyName:    r.PostFormValue("company_name"),
        PwHash:         []byte(r.PostFormValue("pw_hash")),
        Address:        r.PostFormValue("company_address"),
    }
    if p.UserName == "" || p.CompanyName == "" || p.PwHash == nil || p.Address == "" {
        http.Error(w, "One or more profile data fields were blank.", http.StatusBadRequest)
        return
    }
    // TODO: replace with 'unique' key check, such that user_name/company_name is
    // unique by db schema
    if dm.userExists(p.UserName, p.CompanyName) {
        http.Error(w, "User already exists", http.StatusBadRequest)
        return        
    }
    // Create a new profile record
    if err := dm.Create(&p).Error; err != nil {
        http.Error(w, "Database error on profile indexing.", http.StatusInternalServerError)
        log.Fatal("dm.Create: ", err)
        return
    }
    // Profile record created successfully, send OK
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    body, err := json.Marshal(p)
    if err != nil {
        log.Fatal("json.Marshal: ", err)
        return
    }
    w.Write(body)
}

// GetExistingProfile operates on a DataManager struct and takes a ResponseWriter
// and a Request, whose body should contain the ID of an existing user profile
// record, and writes the result of that query, user Profile or error, to w.
func (dm *DataManager) GetExistingProfile(w http.ResponseWriter, r *http.Request) {
    rsrc := parseProfileQuery(r.URL.Path)
    cn, un := rsrc[2], rsrc[3]
    if un == "" || cn == "" {
        http.Error(w, "One or more credentials were blank", http.StatusBadRequest)
        return
    }
    // TODO: ensure authentication via cookie
    // if !IsUserAuthenticated(un, cn) {
    //     http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //     return        
    // }
    p, err := dm.GetProfileHelper(un, cn)
    if err != nil {
        http.Error(w, "Database error on profile retrieval", http.StatusInternalServerError)
        log.Fatal("dm.GetProfileHelper: ", err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    body, err := json.Marshal(p)
    if err != nil {
        log.Fatal("json.Marshal: ", err)
        return
    }
    w.Write(body)
}

// UpdateExistingProfile operates on a DataManager struct and takes a
// ResponseWriter and a Request, whose body should contain the user's unique ID,
// and updates the existing record for that user in the db. Either an error or
// success message on record update will be written to w.
func (dm *DataManager) UpdateExistingProfile(w http.ResponseWriter, r *http.Request) {
    // TODO: replace un/cn matching with cookie lookup
    rsrc := parseProfileQuery(r.URL.Path)
    cn, un := rsrc[2], rsrc[3]
    if un == "" || cn == "" {
        http.Error(w, "One or more credentials were blank", http.StatusBadRequest)
        return
    }
    if err := r.ParseForm(); err != nil {
        log.Fatal("r.ParseForm: ", err)
        return
    }
    // TODO: ensure authentication via cookie
    // if !IsUserAuthenticated(un, cn) {
    //     http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //     return        
    // }
    // user_name, company_address, and pw_hash are the only updatable profile fields
    // NOTE: ignore company_name (attempt to change should be handled by frontend)
    temp := map[string]interface{}{
        "user_name":        r.PostFormValue("user_name"),
        "company_address":  r.PostFormValue("company_address"),
        "pw_hash":          []byte(r.PostFormValue("pw_hash")),
    }
    // Potentially multiple fields updated, we want to ensure all or none are updated
    tx := dm.Begin()
    for k, v := range temp {
        if err := (&DataManager{tx}).UpdateProfileHelper(un, cn, k, v); err != nil {
            tx.Rollback()
            http.Error(w, "Database error on profile update", http.StatusInternalServerError)
            log.Fatal("dm.UpdateProfileHelper: ", err)
            return
        }
    }
    tx.Commit()
    w.Header().Set("Content-Type", "application/json")
    p := Profile{
        CompanyName:    cn,
        UserName:       temp["user_name"].(string),
        Address:        temp["company_address"].(string),
        PwHash:         temp["pw_hash"].([]byte),
    }
    body, err := json.Marshal(p)
    if err != nil {
        log.Fatal("json.Marshal: ", err)
        return
    }
    w.Write(body)
}

// DeleteExistingProfile operates on a DataManager struct and takes a user's 
// unique ID and deletes the corresponding db record. Either an error or 
// success message onm successful deletion of the record will be written to w.
func (dm *DataManager) DeleteExistingProfile(w http.ResponseWriter, r *http.Request) {
    rsrc := parseProfileQuery(r.URL.Path)
    cn, un := rsrc[2], rsrc[3]
    if un == "" || cn == "" {
        http.Error(w, "One or more credentials were blank", http.StatusBadRequest)
        return
    }
    // TODO: replace with 'unique' key check, such that user_name/company_name is
    // unique by db schema
    if !dm.userExists(un, cn) {
        http.Error(w, "User does not exist", http.StatusBadRequest)
        return        
    }
    // TODO: ensure authentication via cookie
    // if !IsUserAuthenticated(un, cn) {
    //     http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //     return        
    // }
    // BUG: Delete() does not throw an error when un or cn are not in the db;
    // seen by running with above userExists statement commented via 
    // TestDeleteExistingProfileNotExistFail()
    qstring := "user_name = ? AND company_name = ?"
    if err := dm.Unscoped().Where(qstring, un, cn).Delete(&Profile{}).Error; err != nil {
        http.Error(w, "Database error on profile delete", http.StatusInternalServerError)
        log.Fatal("dm.DeleteExistingProfile: ", err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
