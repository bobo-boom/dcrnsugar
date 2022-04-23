package cache

import (
	"github.com/bobo-boom/dcrnsugar/db/dbtypes"
	"sync"
)

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

func (c *CacheAddress) WriteCache(addr string, isFinished bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.addrs[addr] = isFinished
}
func (c *CacheAddress) GetAddressStatus(addr string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	status := c.addrs[addr]
	return status
}
func (c *CacheAddress) HasAddress() bool {
	if len(c.addrs) > 0 {
		return true
	}
	return false
}

type AddressQueue struct {
	mtx   sync.Mutex
	addrs []string
}

func (a *AddressQueue) AddAddressToQueue(addr []string) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.addrs = append(a.addrs, addr...)
}

type BalanceInfoCache struct {
	mtx         sync.Mutex
	balanceInfo []dbtypes.BalanceInfo
}

func (b *BalanceInfoCache) WriteCache(balanceInfo dbtypes.BalanceInfo) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.balanceInfo = append(b.balanceInfo, balanceInfo)
}
func NewBalanceInfoCache() *BalanceInfoCache {
	return &BalanceInfoCache{
		balanceInfo: make([]dbtypes.BalanceInfo, 10000),
	}
}
