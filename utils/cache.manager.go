package utils

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheManager struct {
	//saving token in cache
	cacheToken *cache.Cache
	//saving other data in cache
	cacheJson *cache.Cache
}

func InitCache() *CacheManager {
	cToken := cache.New(5*time.Minute, 30*time.Minute)
	cJson := cache.New(5*time.Minute, 30*time.Minute)
	return &CacheManager{cacheToken: cToken, cacheJson: cJson}
}

//use this function for caching data
func (c *CacheManager) Cacheable(key string, callback func() interface{}) interface{} {
	data := c.CacheGet(key)
	if data == nil {
		result := callback()
		c.cacheJson.Set(key, result, 4380*time.Hour)
		return result
	} else {
		return data
	}
}

//use this function for get data
func (c *CacheManager) CacheGet(key string) interface{} {
	if x, found := c.cacheJson.Get(key); found {
		return x
	}

	return nil
}

func (c *CacheManager) CacheSet(key string, data string, d time.Duration) {
	c.cacheToken.Set(key, data, d)
}

//use this function for clearing data
func (c *CacheManager) CacheEvict(key string) {
	c.cacheJson.Delete(key)
	c.cacheJson.DeleteExpired()
}

//use this function for save token for key and user json for value after created
func (c *CacheManager) SaveUserToken(key string, data string, d time.Duration) {
	c.cacheToken.Set(key, data, d)
}

func (c *CacheManager) UpdateUserToken(key string, data string) {
	c.cacheToken.SetDefault(key, data)
}

//use this function for save token after created if not expired
func (c *CacheManager) SaveUserTokenNoExpired(key string, data string) {
	c.cacheToken.Set(key, data, cache.NoExpiration)
}

//use this function for get token
func (c *CacheManager) GetUserToken(key string) interface{} {
	if x, found := c.cacheToken.Get(key); found {
		return x
	}

	return nil
}

//use this function for delete token
func (c *CacheManager) DeleteToken(key string) {
	c.cacheToken.Delete(key)
	c.cacheToken.DeleteExpired()
}
