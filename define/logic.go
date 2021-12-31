package define

import (
	"github.com/ZYallers/rpcx-framework/util/mtsc"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/go-redis/redis"
)

type Logic struct {
	mtsc.Redis
}

var (
	cache, session             mvcs.RdsCollector
	cacheClient, sessionClient *app.RedisClient
)

func init() {
	cacheClient = &app.RedisClient{
		Host: "",
		Port: "",
		Pwd:  "",
		Db:   0,
	}
	sessionClient = &app.RedisClient{
		Host: "",
		Port: "",
		Pwd:  "",
		Db:   0,
	}
}
func (c *Logic) Cache() *redis.Client {
	return c.NewClient(&cache, cacheClient)
}

func (c *Logic) Session() *redis.Client {
	return c.NewClient(&session, sessionClient)
}
