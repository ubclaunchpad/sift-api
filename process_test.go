package main

import (
	"testing"

	"encoding/json"
	"os"
)

func TestIsMalformedJSONTrue(t *testing.T) {
	mf, err := os.Open("test_data/test_malformed.json")
	if err != nil {
		return
	}
	defer mf.Close()
	// Expect unformatted json to return true
	if !IsMalformedJSON(mf) {
		t.Error("Expected malformed, but was not malformed")
	}
}

func TestIsMalformedJSONFalse(t *testing.T) {
	ff, err := os.Open("test_data/test_formatted.json")
	if err != nil {
		return
	}
	defer ff.Close()
	// Expect formatted json to return false
	if IsMalformedJSON(ff) {
		t.Error("Expected not malformed, but was malformed")
	}
}

func TestProcessJSOMalformed(t *testing.T) {
	mf, err := os.Open("test_data/test_malformed.json")
	if err != nil {
		return
	}
	defer mf.Close()
	ff, err := os.Open("test_data/test_processed.json")
	if err != nil {
		return
	}
	defer ff.Close()

	var desired interface{}
	if err = json.NewDecoder(ff).Decode(&desired); err != nil {
		t.Error("Error returned decoding json, not expected: ", err)
	}

	check, err := ProcessJSON(mf)
	if err != nil {
		t.Error("Error returned, not expected: ", err)
	}
	temp1, _ := json.Marshal(&check)
	temp2, _ := json.Marshal(&desired)

	if string(temp1) != string(temp2) {
		t.Errorf("StrictifyJSON: results not equal: %s != %s", temp1, temp2)
	}
}

func TestProcessJSONFormatted(t *testing.T) {
	ff, err := os.Open("test_data/test_formatted.json")
	if err != nil {
		return
	}
	defer ff.Close()

	_, err = ProcessJSON(ff)
	if err != nil {
		t.Error("Error returned, not expected: ", err)
	}
}

func BenchmarkJSONFull(b *testing.B) {
	// NOTE: hk_feedback.json is a local file containing all Home and Kitchen review
	// data from Amazon (see README), which I did not commit because of file size.
	b.StopTimer()
	hk, err := os.Open("test_data/hk_feedback.json")
	if err != nil {
		return
	}
	defer hk.Close()

	b.StartTimer()
	// Uncomment following two commented sections to write data to stdout as csv
	_, err = ProcessJSON(hk)
	b.StopTimer()
	if err != nil {
		b.Error("ProcessJSON did not work on a yuge file")
	}

}

// func TestFeedbackFormHandlerValidJSON(t *testing.T) {
// 	fp := "test_data/hk_feedback_valid.json"
// 	f, err := os.Open(fp)
// 	if err != nil {
// 		t.Error("Couldn't open file")
// 	}
//
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("feedback", filepath.Base(fp))
// 	if err != nil {
// 		t.Error("Couldn't create form file")
// 	}
// 	_, err = io.Copy(part, f)
//
// 	err = writer.Close()
// 	if err != nil {
// 		t.Error("Couldn't close writer")
// 	}
//
// 	req, _ := http.NewRequest("POST", "/feedback", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(FeedbackFormHandler)
// 	handler.ServeHTTP(rr, req)
// 	if rr.Code != http.StatusOK {
// 		t.Errorf("HTTP status code recieved: %d, expected %d", rr.Code, http.StatusOK)
// 	}
// 	// err = json.Unmarshal(rr.Body.Bytes(), &p)
// 	// if err != nil {
// 	// 	t.Errorf("Error (%v) encountered when unmarshalling profile", err)
// 	// }
//
// }
