package env

import (
	"github.com/ZYallers/rpcx-framework/define"
	"time"
)

var Redis = &define.Redis{
	CommonExpiration: 432000 * time.Second, // 5d
	TTL:              define.TTLType{Forever: -1, NotExist: -2},
}

