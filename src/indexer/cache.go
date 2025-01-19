package indexer

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

func (ind Indexer) setTapIndexCache(key string, val TapIndex) {
	ind.Cache.Set(key, val, cache.DefaultExpiration)
}

func (ind Indexer) getTapIndexCache(key string) (val TapIndex) {
	if x, found := ind.Cache.Get(key); found {
		val = x.(TapIndex)
	}
	return
}

func (ind Indexer) getTapIndexCacheWithExpiration(key string) (val TapIndex, tim time.Time) {
	if x, t, found := ind.Cache.GetWithExpiration(key); found {
		val = x.(TapIndex)
		tim = t
	}
	return
}

func (ind Indexer) flushCache() {
	ind.Cache.Flush()
}
