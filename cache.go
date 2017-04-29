package main

import (
	"log"
	"time"

	"github.com/streamrail/concurrent-map"
)

var cache cmap.ConcurrentMap

const cacheSize = 10000

type CacheItem struct {
	Key    string
	Value  []byte
	Expire time.Time
}

func init() {
	cache = cmap.New()
	go doEvery(time.Minute, removeStaleCacheItems)
}

// AddItemToCache add a new item to cache with expiration time
func AddItemToCache(key string, value []byte, expire time.Duration) {
	log.Println("adding key", key, "to cache")
	if cache.Count() <= cacheSize {
		cache.Set(key, CacheItem{Key: key, Value: value, Expire: time.Now().Add(expire)})
	} else {
		panic("cacheSize reached")
	}
}

// GetItemFromCache get item from cache
func GetItemFromCache(key string) *CacheItem {
	if cache.Has(key) {
		if tmp, ok := cache.Get(key); ok {
			item := tmp.(CacheItem)
			if item.Expire.After(time.Now()) {
				log.Println("getting key", item.Key, "from cache")
				return &item
			}
			removeItem(item.Key)
		}
	}
	return nil
}

func removeItem(key string) {
	log.Println("removing key", key, "from cache")
	cache.Remove(key)
}

func removeStaleCacheItems() {
	for _, item := range cache.Items() {
		cItem := item.(CacheItem)
		if time.Now().After(cItem.Expire) {
			removeItem(cItem.Key)
		}
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
