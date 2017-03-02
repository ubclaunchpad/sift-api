package main

import (
	"log"
	"net/http"
)

// Cookie Helpers

// CreateCookieHelper takes a session id and returns an encoded cookie
func (dm *DataManager) CreateCookieHelper(id uint) (http.Cookie, error) {

	value := map[string]uint{
		"id": id,
	}

	encoded, err := dm.Encode("session", value)
	if err != nil {
		log.Fatal("dm.Encode, err:", err)
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}

	return cookie, nil
}

// DecodeCookieHelper takes a cookie and returns the associated session ID
func (dm *DataManager) DecodeCookieHelper(cookie http.Cookie) (uint, error) {

	value := make(map[string]uint)

	if err := dm.Decode("session", cookie.Value, &value); err != nil {
		log.Fatal("dm.Decode, err:", err)
		return 0, err
	}

	return value["id"], nil
}

// Helper Methods

// CreateSessionHelper validates session object then pushes to session table
func (dm *DataManager) CreateSessionHelper(sesh Session) error {
	return dm.Create(&sesh).Error
}

// GetSessionByIdHelper retrieves using id primary key
func (dm *DataManager) GetSessionByIdHelper(id uint) (sesh Session, err error) {
	err = dm.First(&sesh, id).Error
	return
}

// GetSessionByUserHelper retrieves using UserID
func (dm *DataManager) GetSessionByUserHelper(id uint) (sesh Session, err error) {
	err = dm.Where("user_id = ?", id).First(&sesh).Error
	return
}

// DeleteSessionByIdHelper deletes a single session by id
func (dm *DataManager) DeleteSessionByIdHelper(id uint) error {
	return dm.Where("id = ?", id).Delete(Session{}).Error
}

// DeleteSessionsByUserHelper delete all sessions for a given UserID
func (dm *DataManager) DeleteSessionsByUserHelper(id uint) error {
	return dm.Where("user_id = ?", id).Delete(Session{}).Error
}

// Handler Functions

//func (dm *DataManager) GetFromCookie(w http.ResponseWriter, r *http.Request) {

//}
