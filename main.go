// Package always goes at the top of the file
// `main` package gets compiled, any other package name is importable.
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	// URL for accessing RabbitMQ
	AMQP_URL = "amqp://sift:sift@localhost:5672/sift"
	// Max file size to store in memory. 100MB
	MAX_FILE_SIZE = 6 << 24
)

// Functions with lowercase names are private to the package.
// Uppercase names are public
func findARandomNumber(c chan int) {
	// Seed a random number generator with the current Unix time
	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	// Send a random int across the channel
	c <- r.Int()
}

// Handler functions for HTTP routes take two inputs
// You write strings to `w` to send responses back to the client
// `r` contains info like headers, body, method, url, etc.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Make a new channel of `ints`. Channels are used to communicate between goroutines.
	c := make(chan int)
	// Create a new goroutine and run `findARandomNumber` on it.
	go findARandomNumber(c)
	// Channels block by default, so we wait for our goroutine to send a number back.
	n := <-c

	// Print to the responsewriter and terminate the connection.
	fmt.Fprintf(w, "Hello Sift. Your random number is %d", n)
}

// Handles uploads of multipart forms. Files should have form name `feedback`.
// Uploaded files are stored in `./uploads`
func FeedbackFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(MAX_FILE_SIZE)
	file, _, err := r.FormFile("feedback")
	if err != nil {
		fmt.Println("Error parsing form: " + err.Error())
		return
	}
	defer file.Close()

	var payload interface{}
	err = json.NewDecoder(file).Decode(&payload)
	if err != nil {
		fmt.Println("Error parsing JSON payload: " + err.Error())
		return
	}

	api, err := NewCeleryAPI(AMQP_URL)
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

	body, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("Error mashalling job response: " + err.Error())
	}
	w.Write(body)
}

// `main` function is the entry point, just like in C
func main() {
	// Create a new router, routers handle sets of logically related routes
	router := mux.NewRouter()
	// Add a handler for the root route
	router.HandleFunc("/", IndexHandler)
	// Handler for the feedback upload route
	router.HandleFunc("/feedback", FeedbackFormHandler).Methods("POST")
	// Create an http server on port 9090 and start serving using our router.
	fmt.Println("Sift API running on port 9090...")
	http.ListenAndServe(":9090", router)
}
