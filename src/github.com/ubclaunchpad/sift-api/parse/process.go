// Process input files into specified pre-process format
package process

import (
  "io"
  "bufio"
  "regexp"
  "encoding/json"
)

// Detect whether input data is 'loose' JSON
func IsLooseJSON(file io.Reader) (bool, error) {
  // Any set of dictionaries without a ',' between brackets is loose
  re, _ := regexp.Compile(".*}[^,]*{.*")
  br := bufio.NewReader(file)
  buf := make([]byte, 1024)
  for {
    n, err := br.Read(buf)
    if err != nil && err != io.EOF {
      panic(err)
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
func ProcessJSON(file io.Reader, out *interface{}) error {
  dec := json.NewDecoder(file)
  temp := new(interface{})
  var fb []interface{}

  for i := uint64(0);; i++ {
    // Decode each JSON object
    if err := dec.Decode(&temp); err == io.EOF {
      break
    } else if err != nil {
      panic(err)
    }
    f := map[string]interface{}{
      "fb_id":  i,
      "fb_body":((*temp).(map[string]interface{})["reviewText"].(string)),
    }
    fb = append(fb, f)
  }
  *out = fb
  return nil
}
