package cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"static-mirrors/pkg/config"
)

// Cache 缓存接口
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items map[string]*cacheItem
	mutex sync.RWMutex
	size  int
}

// cacheItem 内存缓存项
type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewCache 创建新的缓存实例
func NewCache(cfg config.Config) (Cache, error) {
	if !cfg.Cache.Enabled {
		return nil, nil
	}

	switch cfg.Cache.Type {
	case "memory":
		return NewMemoryCache(cfg.Cache.Memory.Size), nil
	default:
		return nil, fmt.Errorf("不支持的缓存类型: %s", cfg.Cache.Type)
	}
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache(size int) *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*cacheItem),
		size:  size,
	}

	// 启动清理过期项的协程
	go cache.cleanupExpired()

	log.Println("内存缓存初始化成功")
	return cache
}

// Get 从内存缓存获取数据
func (c *MemoryCache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, nil // 缓存未命中
	}

	// 检查是否过期
	if time.Now().After(item.expiration) {
		return nil, nil // 缓存已过期
	}

	return item.value, nil
}

// Set 向内存缓存设置数据
func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查缓存大小
	if len(c.items) >= c.size {
		// 简单的LRU策略：删除最早的项
		c.evictOldest()
	}

	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	return nil
}

// Delete 从内存缓存删除数据
func (c *MemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	return nil
}

// Exists 检查内存缓存中是否存在键
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if time.Now().After(item.expiration) {
		return false, nil // 缓存已过期
	}

	return true, nil
}

// cleanupExpired 清理过期的缓存项
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// evictOldest 删除最早的缓存项
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestKey == "" || item.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiration
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// GenerateCacheKey 生成缓存键
func GenerateCacheKey(url string, method string) string {
	return fmt.Sprintf("%s:%s", method, url)
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits   int64 `json:"hits"`
	Misses int64 `json:"misses"`
	Size   int64 `json:"size"`
}

// GetStats 获取缓存统计信息
func GetStats(cache Cache) (CacheStats, error) {
	// 这里可以根据实际需要实现统计功能
	// 目前返回默认值
	return CacheStats{
		Hits:   0,
		Misses: 0,
		Size:   0,
	}, nil
}
