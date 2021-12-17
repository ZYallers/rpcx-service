package env

type redisKey struct {
	String stringKey
	Hash   hashKey
	Set    setKey
	ZSet   zSetKey
	List   listKey
}

type stringKey struct {
	LatestHeadModel string
	HeadNotBanner   string
}

type hashKey struct {
	HeadBanner string
}

type setKey struct {
}

type zSetKey struct {
}

type listKey struct {
}

var RedisKey = &redisKey{
	String: stringKey{
		LatestHeadModel: "rpcx-homepage@head:model:latest:string",
		HeadNotBanner:   "rpcx-homepage@head:notbanner:string",
	},
	Hash: hashKey{
		HeadBanner: "rpcx-homepage@head:banner:hash",
	},
	Set:  setKey{},
	ZSet: zSetKey{},
	List: listKey{},
}
