package cache

import "sync"

type CacheAddress struct {
	mtx   sync.Mutex
	addrs map[string]bool
}

func NewCacheAddress() *CacheAddress {

	return &CacheAddress{
		addrs: make(map[string]bool),
	}
}

func (c *CacheAddress) IsExist(addr string) bool {
	_, ok := c.addrs[addr]
	return ok
}

func (c *CacheAddress) WriteCache(addr string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.addrs[addr] = true
}
