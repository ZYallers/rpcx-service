package client

import (
	"context"
	client2 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"src/config/app"
	"src/config/define"
	"src/libraries/util/helper"
	"sync"
	"time"
)

var (
	xClientMap   sync.Map
	failMode     = client.Failover
	selectMode   = client.RoundRobin
	clientOption = client.Option{
		Retries:            3,
		RPCPath:            share.DefaultRPCPath,
		ConnectTimeout:     time.Second,
		SerializeType:      protocol.MsgPack,
		CompressType:       protocol.None,
		BackupLatency:      10 * time.Millisecond,
		TCPKeepAlivePeriod: time.Minute, // if it is zero we don't set keepalive
		IdleTimeout:        time.Minute, // ReadTimeout sets max idle time for underlying net.Conns
		GenBreaker: func() client.Breaker {
			// if failed 10 times, return error immediately, and will try to connect after 60 seconds
			return client.NewConsecCircuitBreaker(10, 60*time.Second)
		},
	}
)

func init() {
	switch app.Service.Env {
	case define.DevelopMode:
		failMode = client.Failfast
		selectMode = client.RandomSelect
	}
}

// XClient ...
// timeout = options[0], trace = options[1]
func XClient(service, method string, args map[string]interface{}) (interface{}, error) {
	share.Trace = false
	if val, ok := args["trace"]; ok && val.(string) == "on" {
		share.Trace = true
		log.Infof("env: %s, etcd: %s->%+v", app.Service.Env, app.Service.Etcd.BasePath, app.Service.Etcd.Addr)
	}

	var xClient client.XClient
	if val, ok := xClientMap.Load(service); ok {
		xClient = val.(client.XClient)
	} else {
		d, _ := client2.NewEtcdV3Discovery(app.Service.Etcd.BasePath, service, app.Service.Etcd.Addr, false, nil)
		xClient = client.NewXClient(service, failMode, selectMode, d, clientOption)
		xClientMap.Store(service, xClient)
	}

	var reply interface{}
	ctx, cancel := context.WithTimeout(context.Background(), helper.DefaultHttpClientTimeout)
	defer cancel()
	return reply, xClient.Call(ctx, method, args, &reply)
}
