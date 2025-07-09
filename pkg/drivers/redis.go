package drivers

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(host string, port string) *redis.Client {
	address := fmt.Sprintf("%v:%v", host, port)
	return redis.NewClient(&redis.Options{
		Addr: address,
	})
}
