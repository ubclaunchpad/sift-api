package main

// CreateSession validates session object then pushes to session table
func (dm *DataManager) CreateSession(sesh *Session) {
	// ValidateSession() if more data is eventually attached to sesh object
	dm.DB.Create(sesh)
}

// GetSessionByID retrieves using id primary key
func (dm *DataManager) GetSessionByID(id uint) *Session {
	sesh := Session{}
	dm.DB.First(&sesh, id)
	return &sesh
}

// GetSessionByUser retrieves using UserID
func (dm *DataManager) GetSessionByUser(id uint) *Session {
	sesh := Session{}
	dm.DB.Where("user_id = ?", id).First(&sesh)
	return &sesh
}

// DeleteSessionsByUser deletes all sessions for a given UserID
func (dm *DataManager) DeleteSessionsByUser(id uint) {
	dm.DB.Where("user_id = ?", id).Delete(Session{})
}
