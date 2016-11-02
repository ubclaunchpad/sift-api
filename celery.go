// Manages running NLP jobs with Celery.

package main

import (
	"errors"
	"fmt"
	"time"

	celery "github.com/shicky/gocelery"
)

const (
	// How frequently we ask Celery if it's done processing the task.
	QUERY_PERIOD = time.Millisecond * 50
	// How much time we allow for receiving a task's result, once we know it's completed.
	RETRIEVAL_TIMEOUT = time.Second
	// How much time we will spend querying Celery for completion status after dispatching a task.
	TIMEOUT = time.Second * 10
)

// CeleryAPI contains references to the Celery backend, broker, and client.
// It exposes a number of methods for running Celery jobs.
type CeleryAPI struct {
	Backend celery.CeleryBackend
	Broker  celery.CeleryBroker
	Client  *celery.CeleryClient
}

// CeleryResult is the type returned from job running functions.
// It contains the result object or an error, if one occurred.
type CeleryResult struct {
	Error  error
	Result interface{}
}

// Returns a new Celery API, connected to Celery at the given URL.
// Celery and RabbitMQ must be running for this to succeed.
func NewCeleryAPI(amqpURL string) (*CeleryAPI, error) {
	backend := celery.NewAMQPCeleryBackend(amqpURL)
	broker := celery.NewAMQPCeleryBroker(amqpURL)
	client, err := celery.NewCeleryClient(broker, backend, 0)
	if err != nil {
		return nil, err
	}

	return &CeleryAPI{backend, broker, client}, nil
}

// Runs a Celery job asynchronously and returns the result through the `result` channel.
// This function should be run as a goroutine.
func (api *CeleryAPI) RunJob(name string, payload interface{}, result chan *CeleryResult) {
	// Send the job to Celery to be run.
	job, err := api.Client.Delay(name, payload)
	if err != nil {
		result <- &CeleryResult{err, nil}
		return
	}

	beganPollingAt := time.Now()
	for {
		// Check for timeout
		if time.Now().Sub(beganPollingAt) > TIMEOUT {
			err := errors.New(fmt.Sprintf("Request timed out to retrieve job %s with timeout %s.",
				name, time.Duration(TIMEOUT).String()))
			result <- &CeleryResult{err, nil}
			return
		}

		// Check for job completion
		ready, err := job.Ready()
		if err != nil {
			result <- &CeleryResult{err, nil}
			return
		}

		// Retrieve result
		if ready {
			res, err := job.Get(RETRIEVAL_TIMEOUT)
			result <- &CeleryResult{err, res}
			return
		}
		time.Sleep(QUERY_PERIOD)
	}
}
