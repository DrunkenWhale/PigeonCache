package pigeoncache

import (
	"sync"
)

type cache struct {
	mutex      sync.Mutex
	lru        *PigeonCache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = c.lru.New(c.cacheBytes, nil)
	}
	c.lru.Put(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}



