// Process input files into specified pre-process format
package main

import (
  "io"
  "regexp"
  "encoding/json"
)

type Feedback struct {
    FBody string  `json:"fb_body"`
    ID    uint64  `json:"fb_id"`
}

func (f *Feedback) MarshalJSON(v interface{}) ([]byte, error) {
  return json.Marshal(map[string]interface{}{
    "fb_body" : f.FBody,
    "fb_id"   : f.ID,
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
func ProcessJSON(file io.Reader) (interface{}, error) {
  dec := json.NewDecoder(file)

  var (
    temp  interface{}
    fb    []interface{}
    f     Feedback
  )

  for i := uint64(0);; i++ {
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
