// Package always goes at the top of the file
// `main` package gets compiled, any other package name is importable.
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/ubclaunchpad/sift-api/jobs"
	"github.com/ubclaunchpad/sift-api/parse"

	"github.com/gorilla/mux"
)

// Max file size to store in memory. 1GB
const MAX_FILE_SIZE = 1 << 30

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
	r.ParseMultipartForm(MAX_FILE_SIZE)
	file, _, err := r.FormFile("feedback")
	if err != nil {
		fmt.Println("Error parsing form: " + err.Error())
		return
	}
	defer file.Close()

	var payload interface{}

	// Pre-process into specified structure
	if err := ProcessJson(file.Body, &payload); err != nil {
		fmt.Println("Error preprocessing JSON payload: " + err.Error())
		return
	}

	res, err := RunJob("sample", &payload)
	if err != nil {
		fmt.Println("Error running job: " + err.Error())
		return
	}

	body, err := json.Marshal(res)
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
	http.ListenAndServe(":9090", router)
}
