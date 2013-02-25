package cache

import (
  "fmt"
  "github.com/kdar/cache/msgpack"
)

type DataSource struct {
  Host  string
  Port  int
  Coder Coder
}

type Coder interface {
  Marshal(v interface{}) ([]byte, error)
  Unmarshal(data []byte, v interface{}) error
}

type CacheFunc func() (interface{}, error)

type Cache interface {
  // Gets a cache value and instead of returning it, puts it into "ret".
  // If no cache is found, the CacheFunc is called and put into cache.
  // It returns the Cachefunc error, and caching error
  Into(key string, expires int, ret interface{}, cFunc CacheFunc) (funcErr error, cacheErr error)

  // Gets a cache value.
  // If no cache is found, the CacheFunc is called and put into cache.
  // It returns the Cachefunc error, and caching error
  Get(key string, expires int, cFunc CacheFunc) (ret interface{}, funcErr error, cacheErr error)

  // Delete a key from the cache.
  // Returns a caching error.
  Delete(key string) (cacheErr error)

  Open() error
  Setup(DataSource) error
}

var wrappers = make(map[string]Cache)

func Register(name string, driver Cache) {
  if name == "" {
    panic("cache: Wrapper name cannot be nil.")
  }
  if _, ok := wrappers[name]; ok != false {
    panic("cache: Wrapper was already registered.")
  }
  wrappers[name] = driver
}

func Open(name string, settings DataSource) (Cache, error) {
  if _, ok := wrappers[name]; ok == false {
    panic(fmt.Sprintf("cache: Unknown wrapper: %s.", name))
  }

  if settings.Coder == nil {
    settings.Coder = &msgpack.MsgPack{}
  }

  err := wrappers[name].Setup(settings)
  if err != nil {
    return nil, err
  }
  return wrappers[name], nil
}

func Wrapper(name string) Cache {
  return wrappers[name]
}
