package dummy

import (
  "errors"
  "github.com/kdar/cache"
  "reflect"
)

func init() {
  cache.Register("dummy", &Source{})
}

type Source struct{}

func (s *Source) Setup(config cache.DataSource) error {
  return nil
}

func (s *Source) Open() error {
  return nil
}

func (s *Source) OutSetFn(key string, expires int, ret interface{}, f cache.CacheFunc) (error, error) {
  val := reflect.ValueOf(ret)
  if val.Kind() != reflect.Ptr {
    return nil, errors.New("ret must be a pointer")
  }

  fret, ferr, cerr := s.GetSetFn(key, expires, f)
  if fret != nil {
    val.Elem().Set(reflect.ValueOf(fret))
  }

  return ferr, cerr
}

func (s *Source) GetSetFn(key string, expires int, f cache.CacheFunc) (interface{}, error, error) {
  fret, ferr := f()
  return fret, ferr, nil
}

func (s *Source) Out(key string, out interface{}) error {
  return nil
}

func (s *Source) Get(key string) (interface{}, error) {
  return nil, errors.New("could not find cache")
}

func (s *Source) Set(key string, value interface{}, expires int) error {
  return nil
}

func (s *Source) Delete(key string) error {
  return nil
}
