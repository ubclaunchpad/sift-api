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

	encoded, err := dm.CookieManager.Encode("session", value)

	cookie := http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}

	if err != nil {
		log.Fatal("dm.CookieManager.Encode", err)
		return cookie, err
	}

	return cookie, nil
}

// DecodeCookieHelper takes a cookie and returns the associated session ID
func (dm *DataManager) DecodeCookieHelper(cookie http.Cookie) (uint, error) {

	value := make(map[string]uint)

	err := dm.CookieManager.Decode("session", cookie.Value, &value)

	id := value["id"]

	if err != nil {
		log.Fatal("dm.CookieManager.Decode", err)
		return id, err
	}

	return value["id"], nil

}

// Helper Methods

// CreateSessionHelper validates session object then pushes to session table
func (dm *DataManager) CreateSessionHelper(sesh Session) (err error) {
	return dm.DB.Create(&sesh).Error
}

// GetSessionByIdHelper retrieves using id primary key
func (dm *DataManager) GetSessionByIdHelper(id uint) (sesh Session, err error) {
	err = dm.DB.First(&sesh, id).Error
	return sesh, err
}

// GetSessionByUserHelper retrieves using UserID
func (dm *DataManager) GetSessionByUserHelper(id uint) (sesh Session, err error) {
	err = dm.DB.Where("user_id = ?", id).First(&sesh).Error
	return sesh, err
}

// DeleteSessionByIdHelper deletes a single session by id
func (dm *DataManager) DeleteSessionByIdHelper(id uint) error {
	return dm.DB.Where("id = ?", id).Delete(Session{}).Error
}

// DeleteSessionsByUserHelper delete all sessions for a given UserID
func (dm *DataManager) DeleteSessionsByUserHelper(id uint) error {
	return dm.DB.Where("user_id = ?", id).Delete(Session{}).Error
}

// Handler Functions

//func (dm *DataManager) GetFromCookie(w http.ResponseWriter, r *http.Request) {

//}
