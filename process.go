// Process input files into specified pre-process format
package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const timeFormat = "20060102030405"

type Feedback struct {
	ID    uint64 `json:"fb_id"`
	FBody string `json:"fb_body"`
}

func (f *Feedback) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"fb_id":   f.ID,
		"fb_body": f.FBody,
	})
}

// Detect whether input data is 'loose' JSON
func IsMalformedJSON(file io.Reader) bool {
	// Any set of dictionaries without a ',' between brackets is malformed
	re, err := regexp.Compile(".*}[^,]*{.*")
	if err != nil {
		panic("IsMalformedJSON: regex failed to compile")
	}
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		// Return true if a non-EOF error is returned or match to invalid JSON syntax found
		if (err != nil && err != io.EOF) || re.Match(buf) {
			return true
		}
		if n == 0 || err == io.EOF {
			break
		}
	}
	return false
}

// Convert JSON into a predetermined JSON format with assigned ID's
// TODO: parameterize to take a text tag argument for mapping feedback bodies
func ProcessJSON(file io.Reader) ([]Feedback, error) {
	dec := json.NewDecoder(file)

	var (
		temp interface{}
		fb   []Feedback
		f    Feedback
	)

	// Decode first to get type of data
	if err := dec.Decode(&temp); err != nil {
		return nil, err
	}
	// Data should only be one of: []interface{}, representing valid JSON, or
	// map[string]interface{}, representing invalid JSON
	switch temp.(type) {
	case []interface{}:
		// Data has valid syntax, just loop through JSON array and convert field names
		for i, val := range temp.([]interface{}) {
			// Decode each JSON object into a Feedback struct
			f.ID = uint64(i)
			f.FBody = ((val).(map[string]interface{})["reviewText"].(string))
			fb = append(fb, f)
		}
	case map[string]interface{}:
		// As first value in set of values was decoded already, append temp to the
		// slice first
		bod := (temp).(map[string]interface{})["reviewText"].(string)
		fb = append(fb, Feedback{ID: uint64(0), FBody: bod})
		for i := uint64(1); ; i++ {
			// Decode each JSON object into a Feedback struct
			if err := dec.Decode(&temp); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			f.ID = i
			f.FBody = (temp).(map[string]interface{})["reviewText"].(string)
			fb = append(fb, f)
		}
	default:
		return nil, errors.New("ProcessJSON: incorrect data type, was " + reflect.TypeOf(temp).Name())
	}

	return fb, nil
}

func WriteToJSON(fblist []Feedback) error {
	f, err := os.Create("hk_feedback_processed_" + time.Now().Format(timeFormat) + ".json")
	if err != nil {
		return err
	}

	jw := json.NewEncoder(f)
	if err := jw.Encode(fblist); err != nil {
		return err
	}

	return nil
}

func WriteToCSV(fblist []Feedback) error {
	f, err := os.Create("hk_feedback_processed_" + time.Now().Format(timeFormat) + ".csv")
	if err != nil {
		return err
	}
	cw := csv.NewWriter(f)
	if err := cw.Write([]string{"fb_id", "fb_body"}); err != nil {
		return err
	}

	for _, record := range fblist {
		temp := []string{strconv.FormatUint(record.ID, 10), string(record.FBody)}
		if err := cw.Write(temp); err != nil {
			return err
		}
	}

	cw.Flush()

	if err := cw.Error(); err != nil {
		return err
	}
	return nil
}

// Handles uploads of multipart forms. Files should have form name `feedback`.
// Uploaded files are stored in `./uploads`
func FeedbackFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/feedback")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(MAX_FILE_SIZE); err != nil {
		fmt.Println("Error parsing form: " + err.Error())
		http.Error(w, "Could not parse file upload", http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("feedback")
	if err != nil {
		fmt.Println("Error creating form file: " + err.Error())
		return
	}
	defer file.Close()

	var payload interface{}
	// Perform malformed JSON check first so file position can be reset
	malformed := IsMalformedJSON(file)
	// Reset file position
	file.Seek(0, 0)
	if malformed {
		payload, err = ProcessJSON(file)
		if err != nil {
			fmt.Println("Error processing invalid JSON payload: " + err.Error())
			http.Error(w, "Could not process file upload", http.StatusInternalServerError)
			return
		}
	} else {
		if err = json.NewDecoder(file).Decode(payload); err != nil {
			fmt.Println("Error decoding valid JSON payload: " + err.Error())
			http.Error(w, "Could not process file upload", http.StatusInternalServerError)
			return
		}
	}
	// payload, err := ProcessJSON(file)
	if err != nil {
		fmt.Println("Error processing invalid JSON payload: " + err.Error())
		http.Error(w, "Could not process file upload", http.StatusInternalServerError)
		return
	}

	if err != nil {
		fmt.Println("Error parsing JSON payload: " + err.Error())
		return
	}

	api, err := NewCeleryAPI(AMQP_URL, REDIS_URL)
	if err != nil {
		fmt.Println("Error creating celery API: ", err.Error())
		return
	}

	resultChannel := make(chan *CeleryResult)
	go api.RunJob("sift.jobrunner.jobs.lda_nlp.run", payload, resultChannel)
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
		return
	}
	w.Write(body)
}
