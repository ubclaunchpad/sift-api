// Manages running NLP jobs.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func RunJob(name string, payload interface{}) (interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:5000/"+name, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
