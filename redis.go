package gomatchmaker

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var rdb redis.Cmdable
var rdbOnce sync.Once

func RegisterRedisClient(instance redis.Cmdable) {
	rdbOnce.Do(func() {
		rdb = instance
	})
}
