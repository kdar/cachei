package redis

import (
  "errors"
  "fmt"
  "github.com/garyburd/redigo/redis"
  "github.com/kdar/cache"
  //"github.com/vmihailenco/msgpack"
  "reflect"
)

func init() {
  cache.Register("redis", &Source{})
}

type Source struct {
  conn   redis.Conn
  config cache.DataSource
}

func (s *Source) Setup(config cache.DataSource) error {
  s.config = config
  return s.Open()
}

func (s *Source) Open() (err error) {
  if s.config.Port == 0 {
    s.config.Port = 6379
  }

  s.conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
  return
}

func (s *Source) Get(key string, expires int, f cache.CacheFunc) (interface{}, error, error) {
  if v, err := redis.Bytes(s.conn.Do("GET", key)); err != redis.ErrNil {
    var ret interface{}
    err = s.config.Coder.Unmarshal(v, &ret)
    if err != nil {
      return nil, nil, err
    }

    return ret, nil, nil
  } else {
    if f != nil {
      fret, rerr := f()

      if fret != nil {
        b, err := s.config.Coder.Marshal(fret)
        if err != nil {
          return nil, rerr, err
        }

        s.conn.Send("MULTI")
        s.conn.Send("SET", key, string(b))
        s.conn.Send("EXPIRE", key, expires)
        _, err = s.conn.Do("EXEC")
        if err != nil {
          return nil, rerr, err
        }
      }

      return fret, rerr, rerr
    }
  }

  return nil, nil, nil
}

func (s *Source) Into(key string, expires int, ret interface{}, f cache.CacheFunc) (error, error) {
  val := reflect.ValueOf(ret)
  if val.Type().Kind() != reflect.Ptr {
    return nil, errors.New("ret must be a pointer")
  }

  fret, ferr, cerr := s.Get(key, expires, f)
  if fret != nil {
    val.Elem().Set(reflect.ValueOf(fret))
  }

  return ferr, cerr
}

func (s *Source) Delete(key string) error {
  _, err := s.conn.Do("DEL", key)
  return err
}
