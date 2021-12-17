package env

import (
	zapp "github.com/ZYallers/zgin/app"
	"os"
	"src/config/app"
	"src/config/define"
	"time"
)

type redis struct {
	Cache, Session, Gym zapp.RedisClient
	CommonExpiration    time.Duration
	TTL                 ttlType
}

type ttlType struct {
	Forever, NotExist float64
}

var Redis = &redis{
	CommonExpiration: 432000 * time.Second, // 5d
	TTL:              ttlType{Forever: -1, NotExist: -2},
}

func init() {
	switch app.Env() {
	case define.DevelopMode, define.GrayMode:
		Redis.Cache = zapp.RedisClient{Host: os.Getenv("redis_host"), Port: os.Getenv("redis_port"),
			Pwd: os.Getenv("redis_password"), Db: 0}
		Redis.Session = zapp.RedisClient{Host: os.Getenv("redis_session_host"), Port: os.Getenv("redis_session_port"),
			Pwd: os.Getenv("redis_session_password"), Db: 0}
	case define.ProduceMode:
		Redis.Cache = zapp.RedisClient{Host: "xxxx", Port: "6379", Pwd: "xxxx", Db: 0}
		Redis.Session = zapp.RedisClient{Host: "xxxx", Port: "6379", Pwd: "xxxx", Db: 0}
	}
}
