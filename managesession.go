package main

// createSession validates session object then pushes to session table
func (dm *DataManager) createSession(sesh *Session) {
	// ValidateSession() if more data is eventually attached to sesh object
	dm.DB.Create(sesh)
}

// getSessionByID retrieves using id primary key
func (dm *DataManager) getSessionByID(id uint) *Session {
	sesh := Session{}
	dm.DB.First(&sesh, id)
	return &sesh
}

// getSessionByUser retrieves using UserID
func (dm *DataManager) getSessionByUser(id uint) *Session {
	sesh := Session{}
	dm.DB.Where("user_id = ?", id).First(&sesh)
	return &sesh
}

// deleteSessionsByUser delete all sessions for a given UserID
func (dm *DataManager) deleteSessionsByUser(id uint) {
	dm.DB.Where("user_id = ?", id).Delete(Session{})
}
