package main

import (
	"context"
	"fmt"
	"net/http"

	"cache-controller/redis"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Cache struct {
	Address string
}

type CacheManager struct {
	Caches map[string]Cache
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Caches: make(map[string]Cache),
	}
}

var cacheManager = NewCacheManager()

func (c *CacheManager) AddCache(addr string) error {
	if addr == "" {
		return fmt.Errorf("addr is empty")
	}
	// c.Addresses = append(cluster.Addresses, addr)
	return nil
}

var ctx = context.TODO()

func main() {
	fmt.Println("=== controller start ===")

	rdb := redis.NewClient()
	if err := rdb.Set(ctx, "key", "value", 0).Err(); err != nil {
		log.Errorf("redisClient.Set failed: %v", err)
	} else if val, err := rdb.Get(ctx, "key").Result(); err != nil {
		log.Errorf("redisClient.Get failed: %v", err)
	} else {
		fmt.Println("key", val)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/register", func(c *gin.Context) {
		addr := c.Query("addr")
		fmt.Println("register for addr: ", addr)
		cacheManager.AddCache(addr)
		c.JSON(http.StatusOK, gin.H{
			"received-message": gin.H{
				"addr": addr,
			},
		})
	})

	r.GET("/caches", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"addresses": cacheManager.Caches,
		})
	})

	r.GET("/unregister", func(c *gin.Context) {

	})

	r.Run(":3037")
}
