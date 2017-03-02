package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSessionById(t *testing.T) {

	// assuming session with ID 1 exists and
	// hasn't been deleted
	var id uint = 1
	if sesh, err := dm.GetSessionByIdHelper(id); err != nil {
		t.Log("Failed to get session: ", err)
		t.Fail()
	} else {
		assert.Equal(t, id, sesh.ID)
	}
}

func TestCreateSession(t *testing.T) {

	var userID uint = 1337
	sesh := Session{UserID: userID}

	// Create sesh
	if err := dm.CreateSessionHelper(sesh); err != nil {
		t.Log("dm.CreateSessionHelper", err)
		t.Fail()
	}

	// Check sesh was added
	if seshRetrieved, err := dm.GetSessionByUserHelper(userID); err != nil {
		t.Log("dm.GetSessionByUserHelper", err)
		t.Fail()
	} else {
		assert.Equal(t, sesh.UserID, seshRetrieved.UserID)
	}

	// Delete sesh and check it was deleted
	if err := dm.DeleteSessionsByUserHelper(userID); err != nil {
		t.Log("dm.DeleteSessionsByUserHelper", err)
		t.Fail()
	}

	if seshRetrieved, err := dm.GetSessionByUserHelper(userID); err == nil {
		t.Log("Session was not deleted")
		t.Fail()
	} else {
		assert.Equal(t, seshRetrieved.ID, uint(0))
	}

}
