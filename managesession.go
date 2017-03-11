package main

import (
	"context"
	"fmt"
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
		fmt.Println("dm.Encode, err:", err)
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
		fmt.Println("dm.Decode, err:", err)
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
	return dm.Unscoped().Where("id = ?", id).Delete(Session{}).Error
}

// DeleteSessionsByUserHelper delete all sessions for a given UserID
func (dm *DataManager) DeleteSessionsByUserHelper(id uint) error {
	return dm.Unscoped().Where("user_id = ?", id).Delete(Session{}).Error
}

// SessionMiddleware takes a request, checks if it has a cookie,
// decodes it if it does, then finds the profile associated with the
// decoded sessionID and attaches it to the context struct, it then
// calls the next middleware in the chain
func (dm *DataManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookies := r.Cookies()
		var seshCookie http.Cookie

		switch len(cookies) {
		case 0:
			// no cookie present, user must be already logged out
			next.ServeHTTP(w, r)
			return
		case 1:
			// has cookie
			seshCookie = *cookies[0]
		default:
			fmt.Println("Length of request cookie array is not 0 or 1")
			http.Error(w, "Error logging user out", http.StatusInternalServerError)
			return
		}

		// Decode cookie

		seshID, err := dm.DecodeCookieHelper(seshCookie)

		if err != nil {
			fmt.Println("dm.DecodeCookieHelper", err)
			http.Error(w, "Could not get session ID from cookie", http.StatusBadRequest)
			return
		}

		// Get session record

		sesh, err := dm.GetSessionByIdHelper(seshID)

		if err != nil {
			fmt.Println("dm.GetSessionByIdHelper", err)
			http.Error(w, "Could not find session with that ID", http.StatusBadRequest)
			return
		}

		// Get user profile for session

		profile, err := dm.GetProfileByIdHelper(sesh.UserID)

		if err != nil {
			fmt.Println("dm.GetProfileByIdHelper")
			http.Error(w, "Could not find profile with that ID", http.StatusBadRequest)
			return
		}

		// Remove password for security
		profile.PwHash = []byte("")

		// Attach user profile to request context

		// TODO: custom 'key' type for profile to avoid possible
		// collisions with other packages (recommended by docs)
		ctx := context.WithValue(r.Context(), "profile", &profile)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		return
	})
}
