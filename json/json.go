package json

import (
  "bytes"
  "encoding/json"
)

type Json struct{}

func (j *Json) Marshal(v interface{}) ([]byte, error) {
  buf := &bytes.Buffer{}
  err := json.NewEncoder(buf).Encode(v)
  return buf.Bytes(), err
}

func (j *Json) Unmarshal(data []byte, v interface{}) error {
  buf := bytes.NewBuffer(data)
  return json.NewDecoder(buf).dec.Decode(v)
}
