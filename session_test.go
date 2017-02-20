package main

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestGetSession(t *testing.T) {

	// assuming session with ID 1 exists and
	// hasn't been deleted
	var id uint = 1
	sesh := dm.getSessionByID(id)
	t.Log("Got session:")
	t.Log(sesh.ID)
	t.Log(sesh.UserID)

	assert.Equal(t, id, sesh.ID)

}

func TestCreateSession(t *testing.T) {

	// create sesh
	var userID uint = 1337
	sesh := Session{UserID: userID}
	dm.createSession(&sesh)

	// check sesh was added
	seshRetrieved := dm.getSessionByUser(userID)
	assert.Equal(t, sesh.UserID, seshRetrieved.UserID)

	// delete sesh and check it was deleted
	dm.deleteSessionsByUser(userID)
	seshRetrieved = dm.getSessionByUser(userID)
	assert.Equal(t, seshRetrieved.ID, uint(0))

}
