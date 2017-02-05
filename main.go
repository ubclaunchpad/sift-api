package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	// URL for accessing RabbitMQ
	AMQP_URL = "amqp://sift:sift@rabbitmq:5672/sift"
	// URL for accessing Redis
	REDIS_URL = "redis:6379"
	// Max file size to store in memory. 100MB
	MAX_FILE_SIZE = 6 << 24
)

// Configures the databse with user, password, host, name, and SSL encryption
// type
type DBConfig struct {
    DBUser          string
    DBPassword      string
    DBHost          string
    DBName          string
    DBSSLType       string
}

func (cfg DBConfig) createDBQueryString() string {
	return fmt.Sprintf("user=%s password=%s host=%s name=%s sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName, cfg.DBSSLType)
}

// Handles uploads of multipart forms. Files should have form name `feedback`.
// Uploaded files are stored in `./uploads`
func FeedbackFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/feedback")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(MAX_FILE_SIZE)
	file, _, err := r.FormFile("feedback")
	if err != nil {
		fmt.Println("Error parsing form: " + err.Error())
		return
	}
	defer file.Close()

	// payload, err := ProcessJSON(file)
	payload := map[string]interface{}{"something": 3}
	if err != nil {
		fmt.Println("Error parsing JSON payload: " + err.Error())
		return
	}

	api, err := NewCeleryAPI(AMQP_URL, REDIS_URL)
	if err != nil {
		fmt.Println("Error creating celery API: ", err.Error())
	}

	resultChannel := make(chan *CeleryResult)
	go api.RunJob("sift.jobrunner.jobs.sample.run", payload, resultChannel)
	result := <-resultChannel
	close(resultChannel)

	if result.Error != nil {
		fmt.Println("Error running job: " + result.Error.Error())
		return
	}

	fmt.Println("Job result: ", result.Body)
	body, err := json.Marshal(result.Body)
	if err != nil {
		fmt.Println("Error mashalling job response: " + err.Error())
	}
	w.Write(body)
}

// `main` function is the entry point, just like in C
func main() {
	// Database configuration
	cfg := DBConfig{
		DBUser:		"test",
		DBPassword:	"testpw",
		DBHost:		"localhost",
		DBName:		"sift_user_info",
		DBSSLType:	"disable",			// switch to 'require' in production
	}
	// Lazily open a connection to the database. The database will only
	// be opened when the first query/exec statement is made against it
	// while service client requests
	db, err := gorm.Open("postgres", cfg.createDBQueryString())
	// Create a DataManager with our lazy connection (see datalayer.go)
	dm := DataManager{DB: db}
	if err != nil {
		log.Fatal(err)
	}
	// Close the connection on main() exit
	defer db.Close()
	// Migration of native types, which can be added as arguments as needed
	dm.DB.AutoMigrate(&Profile{})
	// Create a new router, routers handle sets of logically related routes
	router := mux.NewRouter()
	// Handler for the feedback upload route
	router.HandleFunc("/feedback", FeedbackFormHandler).Methods("POST")
	// Example DataManager handler:
	// router.HandleFunc("/profile/{id: [0-9]+}", dm.GetExistingProfile).Methods("GET")
	// Create an http server on port 9090 and start serving using our router.
	fmt.Println("Sift API running on port 9090...")
	http.ListenAndServe(":9090", router)
}
