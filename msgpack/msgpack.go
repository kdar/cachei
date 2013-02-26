package msgpack

import (
  "bytes"
  "github.com/vmihailenco/msgpack"
)

type Coder struct{}

func (c *Coder) Marshal(v interface{}) ([]byte, error) {
  buf := &bytes.Buffer{}
  err := msgpack.NewEncoder(buf).Encode(v)
  return buf.Bytes(), err
}

func (c *Coder) Unmarshal(data []byte, v interface{}) error {
  buf := bytes.NewBuffer(data)
  dec := msgpack.NewDecoder(buf)

  // this func allows us to decode maps as map[string]interface{},
  // which is inline with how JSON does it.
  dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
    n, err := d.DecodeMapLen()
    if err != nil {
      return nil, err
    }

    m := make(map[string]interface{}, n)
    for i := 0; i < n; i++ {
      mk, err := d.DecodeString()
      if err != nil {
        return nil, err
      }

      mv, err := d.DecodeInterface()
      if err != nil {
        return nil, err
      }

      m[mk] = mv
    }
    return m, nil
  }

  return dec.Decode(v)
}
