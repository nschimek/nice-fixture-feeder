package core

import (
	"encoding/json"
	"fmt"

	"github.com/rainycape/memcache"
	"github.com/sirupsen/logrus"
)

const (
	locationFormat = "%s:%d"
	expiration = 60 * 15 // 15 minutes
	keyFormat = "%s-%s"
)

var (	
	MCC *memcache.Client
)

//go:generate mockery --name MC --filename mc_mock.go
type MC interface {
	Get(key string) (*memcache.Item, error)
	Set(item *memcache.Item) error
}

//go:generate mockery --name Cache --filename cache_mock.go
type Cache[T any] interface {
	Get(key string) (*T, error)
	Set(key string, value *T) error
}

type cache[T any] struct {
	prefix string
	mc MC
}

// Sets up the Memcached Client global variable for injection into NewCache().
func SetupCache(config *Config) {
	Log.WithFields(logrus.Fields{
		"host": config.Cache.Host,
		"port": config.Cache.Port,
	}).Info("Connecting to Memcached...")
	mcc, err := memcache.New(fmt.Sprintf(locationFormat, config.Cache.Host, config.Cache.Port))

	if err != nil {
		Log.Fatal(err)
	}

	MCC = mcc
}

// Create a new instance of Cache with the given type.
func NewCache[T any](mc MC, prefix string) Cache[T] {
	return &cache[T]{
		mc: mc,
		prefix: prefix,
	}
}

// Attempt to get a value from the cache.  Returns nil on a cache miss.
// Error will only be returned if it's not a CacheMiss.
// This function will log its errors, so its optional to handle them.
func (c *cache[T]) Get(key string) (*T, error) {
	var value T
	item, err := c.mc.Get(fmt.Sprintf(keyFormat, c.prefix, key))

	if err != nil && err != memcache.ErrCacheMiss {
		Log.Error("Cache Get error: ", err)
		return nil, err
	} else if err != nil {
		// on a cache miss, just return nil for both
		return nil, nil
	}
	
	err = json.Unmarshal(item.Value, &value)

	if err != nil {
		Log.Error("Cache Unmarshal error: ", err)
		return nil, err
	}

	return &value, nil
}

// Set a value in the Cache.  Returns an error if there is one, and also logs it.
func (c *cache[T]) Set(key string, value *T) error {
	bytes, err := json.Marshal(value)

	if err != nil {
		Log.Error("Cache Marshall error: ", err)
		return err
	}

	err = c.mc.Set(&memcache.Item{
		Key: fmt.Sprintf(keyFormat, c.prefix, key),
		Value: bytes,
		Expiration: expiration,
	})

	if err != nil {
		Log.Error("Cache Set error: ", err)
		return err
	}

	return nil
}