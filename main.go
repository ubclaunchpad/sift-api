package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// URL for accessing RabbitMQ
	AMQP_URL = "amqp://sift:sift@localhost:5672/sift"
	// URL for accessing Redis
	REDIS_URL = "localhost:6379"
	// Max file size to store in memory. 100MB
	MAX_FILE_SIZE = 6 << 24
)

// Handles uploads of multipart forms. Files should have form name `feedback`.
// Uploaded files are stored in `./uploads`
func FeedbackFormHandler(w http.ResponseWriter, r *http.Request) {
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
	// Create a new router, routers handle sets of logically related routes
	router := mux.NewRouter()
	// Handler for the feedback upload route
	router.HandleFunc("/feedback", FeedbackFormHandler).Methods("POST")
	// Create an http server on port 9090 and start serving using our router.
	fmt.Println("Sift API running on port 9090...")
	http.ListenAndServe(":9090", router)
}
