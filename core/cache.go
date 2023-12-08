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
	CC CacheClient
)

//go:generate mockery --name CacheClient --filename cache_client_mock.go
type CacheClient interface {
	Get(key string) (*memcache.Item, error)
	Set(item *memcache.Item) error
}

//go:generate mockery --name Cache --filename cache_mock.go
type Cache[T any] interface {
	Get(key interface{}) (*T, error)
	Set(key interface{}, value *T) error
}

type cache[T any] struct {
	prefix string
	client CacheClient
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

	CC = mcc
}

// Create a new instance of Cache.  Requires a CacheClient.
func NewCache[T any](cc CacheClient, prefix string) Cache[T] {
	return &cache[T]{
		client: cc,
		prefix: prefix,
	}
}

// Attempt to get a value from the cache.  Returns nil on a cache miss.
// Error will only be returned if it's not a CacheMiss.
// This function will log its errors, so its optional to handle them.
func (c *cache[T]) Get(key interface{}) (*T, error) {
	ks, err := c.keyString(key)

	if err != nil {
		return nil, err
	}

	Log.WithField("key", ks).Debug("Getting key from cache")

	var value T
	item, err := c.client.Get(ks)

	if err != nil && err != memcache.ErrCacheMiss {
		Log.Error("Cache Get error: ", err)
		return nil, err
	} else if err != nil {
		Log.WithField("key", ks).Debug("Cache miss")
		// on a cache miss, just return nil for both
		return nil, nil
	}
	
	err = json.Unmarshal(item.Value, &value)

	if err != nil {
		Log.Error("Cache Unmarshal error: ", err)
		return nil, err
	}

	Log.WithField("key", ks).Debug("Cache hit")

	return &value, nil
}

// Set a value in the Cache.  Returns an error if there is one, and also logs it.
func (c *cache[T]) Set(key interface{}, value *T) error {
	ks, err := c.keyString(key)

	if err != nil {
		return err
	}

	// due this method having type safety, I don't forsee errors occurring here
	bytes, _ := json.Marshal(value)

	err = c.client.Set(&memcache.Item{
		Key: ks,
		Value: bytes,
		Expiration: expiration,
	})

	if err != nil {
		Log.Error("Cache Set error: ", err)
		return err
	}

	Log.WithField("key", ks).Debug("Cache set successfully")

	return nil
}

// Converts a key of any type to a string.
func (c *cache[T]) keyString(key interface{}) (string, error) {
	bytes, err := json.Marshal(key)

	if err != nil {
		Log.Error("Cache key Marshall error: ", err)
		return "", err
	}

	return fmt.Sprintf(keyFormat, c.prefix, string(bytes)), nil
}