package main

import (
  "testing"

  "os"
  "encoding/json"
)

func TestIsLooseJSONTrue(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted.json")
  if err != nil {
    return
  }
  defer uffile.Close()
  // Expect unformatted json to return true
  loose, err := IsLooseJSON(uffile)
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  } else if !loose {
    t.Error("Expected loose, but was not loose")
  }
}

func TestIsLooseJSONFalse(t *testing.T) {
  ffile, err := os.Open("test_data/test_formatted.json")
  if err != nil {
    return
  }
  defer ffile.Close()
  // Expect formatted json to return false
  loose, err := IsLooseJSON(ffile)
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  } else if loose {
    t.Error("Expected not loose, but was loose")
  }
}

func TestStrictifyJSONSmall(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted_small.json")
  if err != nil {
    return
  }
  defer uffile.Close()
  ffile, err := os.Open("test_data/test_processed_small.json")
  if err != nil {
    return
  }
  defer ffile.Close()

  var desired interface{}
  _ = json.NewDecoder(ffile).Decode(&desired)

  check, err := ProcessJSON(uffile);
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  }
  temp1,_ := json.Marshal(&check)
  temp2,_ := json.Marshal(&desired)

  if string(temp1) != string(temp2) {
    t.Error("StrictifyJSON did not work on a small file")
  }
}

func TestStrictifyJSONLarge(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted.json")
  if err != nil {
    return
  }
  defer uffile.Close()
  ffile, err := os.Open("test_data/test_processed.json")
  if err != nil {
    return
  }
  defer ffile.Close()

  var desired interface{}
  _ = json.NewDecoder(ffile).Decode(&desired)

  check, err := ProcessJSON(uffile);
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  }
  temp1,_ := json.Marshal(&check)
  temp2,_ := json.Marshal(&desired)

  if string(temp1) != string(temp2) {
    t.Error("StrictifyJSON did not work on a large file")
  }
}

func BenchmarkJSONFull(b *testing.B) {
  // NOTE: hk_feedback.json is a local file containing all Home and Kitchen review data from Amazon (see README), which I did not commit because of file size
  hk, err := os.Open("test_data/hk_feedback.json")
  if err != nil {
    return
  }
  defer hk.Close()

  b.StartTimer()
  _, _ = ProcessJSON(hk)
  b.StopTimer()

}
