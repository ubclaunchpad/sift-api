// Package always goes at the top of the file
// `main` package gets compiled, any other package name is importable.
package main

import (
	"fmt" // standard formatting library, `printf` and so on
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

// `main` function is the entry point, just like in C
func main() {
	// Create a new router, routers handle sets of logically related routes
	router := mux.NewRouter()
	// Add a handler for the root route
	router.HandleFunc("/", IndexHandler)
	// Create an http server on port 9090 and start serving using our router.
	http.ListenAndServe(":9090", router)
}
