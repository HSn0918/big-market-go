package test

import (
	"testing"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

var rdb *redis.Redis

func init() {

	rdb, err := redis.NewRedis(redis.RedisConf{
		Host:        "127.0.0.1",
		Type:        "node",
		Pass:        "",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: 0,
	})
	rdb.Ping()
	if err != nil {
		panic(err)
	}
}
func TestGetRandomAwardId(t *testing.T) {

}
