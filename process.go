// Process input files into specified pre-process format
package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
)

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
func IsLooseJSON(file io.Reader) (bool, error) {
	// Any set of dictionaries without a ',' between brackets is loose
	re, _ := regexp.Compile(".*}[^,]*{.*")
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return false, err
		}
		if n == 0 {
			break
		}
		// If a match is found, return true
		if matched := re.Match(buf); matched {
			return true, nil
		}
	}
	return false, nil
}

// Convert JSON into a predetermined JSON format with assigned ID's
func ProcessJSON(file io.Reader) ([]Feedback, error) {
	dec := json.NewDecoder(file)

	var (
		temp interface{}
		fb   []Feedback
		f    Feedback
	)

	for i := uint64(0); ; i++ {
		// Decode each JSON object
		if err := dec.Decode(&temp); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		f.ID = i
		tbody := ((temp).(map[string]interface{})["reviewText"].(string))
		f.FBody = tbody
		fb = append(fb, f)
	}
	return fb, nil
}

func WriteCSVToStdOut(fblist []Feedback) error {
	new := csv.NewWriter(os.Stdout)
	headers := []string{"fb_id", "fb_body"}

	if err := new.Write(headers); err != nil {
		return err
	}

	for _, record := range fblist {
		temp := []string{strconv.FormatUint(record.ID, 10), string(record.FBody)}
		if err := new.Write(temp); err != nil {
			return err
		}
	}

	new.Flush()

	if err := new.Error(); err != nil {
		return err
	}
	return nil
}
