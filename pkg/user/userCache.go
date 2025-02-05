package user

import (
	"sync"
	"time"

	"github.com/Monkhai/shwipe-server.git/pkg/db"
)

type CacheEntry struct {
	createdAt time.Time
	user      *db.DBUser
}

type UserCache struct {
	cacheMap map[string]CacheEntry
	Mutex    *sync.RWMutex
}

func (c *UserCache) Add(key string, user *db.DBUser) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.cacheMap[key] = CacheEntry{
		createdAt: time.Now().UTC(),
		user:      user,
	}
}

func (c *UserCache) Get(key string) (*db.DBUser, bool) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	entry, exists := c.cacheMap[key]
	if !exists {
		return nil, false
	}
	return entry.user, true
}

func (c *UserCache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		c.Mutex.Lock()
		for key, entry := range c.cacheMap {
			freshTime := time.Now().Add(-duration)
			if entry.createdAt.Before(freshTime) {
				delete(c.cacheMap, key)
			}
		}
		c.Mutex.Unlock()
	}
}

func NewUserCache(duration time.Duration) *UserCache {
	newCache := UserCache{
		cacheMap: map[string]CacheEntry{},
		Mutex:    &sync.RWMutex{},
	}
	go newCache.reapLoop(duration)
	return &newCache
}
