package main 

import (
  "testing"

  "os"
  "bufio"
  "encoding/json"
)

func TestIsLooseJSONTrue(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted.json")
  if err != nil {
    panic(err)
  }
  defer uffile.Close()
  ufr := bufio.NewReader(uffile)
  // Expect unformatted json to return true
  loose, err := IsLooseJSON(ufr)
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  } else if !loose {
    t.Error("Expected loose, but was not loose")
  }
}

func TestIsLooseJSONFalse(t *testing.T) {
  ffile, err := os.Open("test_data/test_formatted.json")
  if err != nil {
    panic(err)
  }
  defer ffile.Close()
  fr := bufio.NewReader(ffile)
  // Expect formatted json to return false
  loose, err := IsLooseJSON(fr)
  if err != nil {
    t.Error("Error returned, not expected: ", err)
  } else if loose {
    t.Error("Expected not loose, but was loose")
  }
}

func TestStrictifyJSONSmall(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted_small.json")
  if err != nil {
    panic(err)
  }
  defer uffile.Close()
  ffile, err := os.Open("test_data/test_processed_small.json")
  if err != nil {
    panic(err)
  }
  defer ffile.Close()

  var desired interface{}
  _ = json.NewDecoder(ffile).Decode(&desired)

  var check interface{}

  ufr := bufio.NewReader(uffile)
  if err := ProcessJSON(ufr, &check); err != nil {
    t.Error("Error returned, not expected: ", err)
  }

  temp1,_ := json.Marshal(&check)
  temp2,_ := json.Marshal(&desired)
  for i,v := range temp1 {
    if v != temp2[i] {
      t.Error("JSON strictify did not work.")
    }
  }
}

func TestStrictifyJSONLarge(t *testing.T) {
  uffile, err := os.Open("test_data/test_unformatted.json")
  if err != nil {
    panic(err)
  }
  defer uffile.Close()
  ffile, err := os.Open("test_data/test_processed.json")
  if err != nil {
    panic(err)
  }
  defer ffile.Close()

  var desired interface{}
  _ = json.NewDecoder(ffile).Decode(&desired)

  var check interface{}

  ufr := bufio.NewReader(uffile)
  if err := ProcessJSON(ufr, &check); err != nil {
    t.Error("Error returned, not expected: ", err)
  }

  temp1,_ := json.Marshal(&check)
  temp2,_ := json.Marshal(&desired)
  for i,v := range temp1 {
    if v != temp2[i] {
      t.Error("JSON strictify did not work.")
    }
  }
}

func BenchmarkJSONFull(b *testing.B) {
  // NOTE: hk_feedback.json is a local file containing all Home and Kitchen review data from Amazon (see README), which I did not commit because of file size
  hk, err := os.Open("test_data/hk_feedback.json")
  if err != nil {
    panic(err)
  }
  defer hk.Close()

  var check interface{}

  ufr := bufio.NewReader(hk)
  b.StartTimer()
  if err := ProcessJSON(ufr, &check); err != nil {
    panic(err)
  }
  b.StopTimer()

}
