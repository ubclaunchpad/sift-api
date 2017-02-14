package main

import (
	"log"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// Database configuration
var cfg = DBConfig{
	DBUser:     "dev",
	DBPassword: "123",
	DBHost:     "localhost",
	DBName:     "siftapi",
	DBSSLType:  "disable", // switch to 'require' in production
}

var db, err = gorm.Open("postgres", cfg.createDBQueryString())
var dm = DataManager{DB: db}

// setup/teardown
func TestMain(m *testing.M) {
	dm.DB.AutoMigrate(&Session{})
	defer db.Close()
	m.Run()
}

func TestConnectDB(t *testing.T) {

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
}

func TestGetSession(t *testing.T) {

	// assuming session with ID 1 exists and
	// hasn't been deleted
	var id uint = 1
	sesh := dm.GetSessionByID(id)
	t.Log("Got session:")
	t.Log(sesh.ID)
	t.Log(sesh.UserID)

	assert.Equal(t, id, sesh.ID)

}

func TestCreateSession(t *testing.T) {

	// create sesh
	var userID uint = 1337
	sesh := Session{UserID: userID}
	dm.CreateSession(&sesh)

	// check sesh was added
	seshRetrieved := dm.GetSessionByUser(userID)
	assert.Equal(t, sesh.UserID, seshRetrieved.UserID)

	// delete sesh and check it was deleted
	dm.DeleteSessionsByUser(userID)
	seshRetrieved = dm.GetSessionByUser(userID)
	assert.Equal(t, seshRetrieved.ID, uint(0))

}
