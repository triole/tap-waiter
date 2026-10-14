package indexer

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

// storeTapIndexCache stores a TapIndex in the cache for the given key.
func (ind Indexer) storeTapIndexCache(key string, val TapIndex) {
	ind.Cache.Set(key, val, cache.DefaultExpiration)
}

// lookupTapIndexCache retrieves a TapIndex from the cache by key.
// Returns an empty TapIndex if key not found.
func (ind Indexer) lookupTapIndexCache(key string) (val TapIndex) {
	if x, found := ind.Cache.Get(key); found {
		val = x.(TapIndex)
	}
	return
}

// lookupTapIndexCacheWithExpiration retrieves a TapIndex from cache along with its expiration time.
// Returns empty TapIndex and zero time if key not found.
func (ind Indexer) lookupTapIndexCacheWithExpiration(key string) (val TapIndex, tim time.Time) {
	if x, t, found := ind.Cache.GetWithExpiration(key); found {
		val = x.(TapIndex)
		tim = t
	}
	return
}

func (ind Indexer) flushCache() {
	ind.Cache.Flush()
}
