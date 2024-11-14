package cache

import (
	"sync"
	"time"
)

type ReferralCache struct {
	mu    sync.Mutex
	cache map[string]CacheItem
}

type CacheItem struct {
	Code      string
	CreatedAt time.Time
}

func NewReferralCache() *ReferralCache {
	return &ReferralCache{
		cache: make(map[string]CacheItem),
	}
}

func (rc *ReferralCache) Set(key string, code string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cache[key] = CacheItem{Code: code, CreatedAt: time.Now()}
}

func (rc *ReferralCache) Get(key string) (string, bool) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	item, exists := rc.cache[key]
	return item.Code, exists
}
