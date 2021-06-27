package cache

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

//var RedisConnectionPool = sync.Pool{
//	// function to execute when no instance of a buffer is not found
//	New: func() interface{} {
//		fmt.Println("allocating new redis connection")
//		return redis.NewClient(&redis.Options{
//			Addr:     "localhost:6379",
//			Password: "", // no password set
//			DB:       0,  // use default DB
//		})
//	},
//}

var RedisCachePool = sync.Pool{
	// function to execute when no instance of a buffer is not found
	New: func() interface{} {
		return cache.New(&cache.Options{
			Redis: redis.NewRing(&redis.RingOptions{
				Addrs: map[string]string{
					"server1": ":6379",
				},
			}),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	},
}
