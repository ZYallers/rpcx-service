package client

import (
	"context"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"src/config/app"
	"testing"
	"time"
)

var (
	etcdAddr []string
	basePath string
)

func init() {
	app.Bootstrap()
	etcdAddr = app.Service.Etcd.Addr
	basePath = app.Service.Etcd.BasePath
}

func Test_XClient(t *testing.T) {
	share.Trace = true
	discovery, _ := etcdClient.NewEtcdV3Discovery(basePath, app.Service.Name, etcdAddr, false, nil)
	xClient := client.NewXClient(app.Service.Name, client.Failover, client.RoundRobin, discovery, client.DefaultOption)
	defer xClient.Close()

	args := map[string]interface{}{"name": "Peter", "key": "panic"}
	var reply interface{}
	ctx, cancel := context.WithTimeout(
		context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{"aaa": "from client"}),
		10*time.Second,
	)
	defer cancel()

	if err := xClient.Call(ctx, "test/check", &args, &reply); err != nil {
		t.Fatalf("call error: %v", err)
	}
	t.Logf("reply: %+v", reply)
}

func Test_XClientPool(t *testing.T) {
	share.Trace = true
	discovery, _ := etcdClient.NewEtcdV3Discovery(basePath, app.Service.Name, etcdAddr, false, nil)
	opt := client.DefaultOption
	pool := client.NewXClientPool(10, app.Service.Name, client.Failover, client.RoundRobin, discovery, opt)
	defer pool.Close()

	args := map[string]interface{}{"name": "Jack"}
	var reply interface{}
	ctx, cancel := context.WithTimeout(
		context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{"aaa": "from client"}),
		10*time.Second,
	)
	defer cancel()

	if err := pool.Get().Call(ctx, "test/check", &args, &reply); err != nil {
		t.Fatalf("call error: %v", err)
	}
	t.Logf("reply: %+v", reply)
}

func Benchmark_XClient(b *testing.B) {
	discovery, _ := etcdClient.NewEtcdV3Discovery(basePath, app.Service.Name, etcdAddr, false, nil)
	xClient := client.NewXClient(app.Service.Name, client.Failover, client.RoundRobin, discovery, client.DefaultOption)
	defer xClient.Close()

	for i := 0; i < b.N; i++ {
		args := map[string]interface{}{"name": "Peter"}
		var reply interface{}
		if err := xClient.Call(context.Background(), "test/check", &args, &reply); err != nil {
			b.Fatalf("call error: %v\n", err)
		}
		//b.Logf("reply: %+v\n", reply)
	}
}

func Benchmark_XClientPool(b *testing.B) {
	discovery, _ := etcdClient.NewEtcdV3Discovery(basePath, app.Service.Name, etcdAddr, false, nil)
	opt := client.DefaultOption
	pool := client.NewXClientPool(10, app.Service.Name, client.Failover, client.RoundRobin, discovery, opt)
	defer pool.Close()

	for i := 0; i < b.N; i++ {
		args := map[string]interface{}{"name": "Peter"}
		var reply interface{}
		if err := pool.Get().Call(context.Background(), "test/check", &args, &reply); err != nil {
			b.Fatalf("call error: %v\n", err)
		}
		//b.Logf("reply: %+v\n", reply)
	}
}
