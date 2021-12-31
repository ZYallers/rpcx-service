package env

import "github.com/ZYallers/rpcx-framework/define"

var RedisKey = &define.RedisKey{
	String: map[string]string{
		"LatestHeadModel": "rpcx-homepage@head:model:latest:string",
		"HeadNotBanner":   "rpcx-homepage@head:notbanner:string",
	},
	Hash: map[string]string{
		"HeadBanner": "rpcx-homepage@head:banner:hash",
	},
}
