package main

import (
	framework "github.com/ZYallers/rpcx-framework"
	_ "github.com/ZYallers/rpcx-service/restful"
	"github.com/smallnest/rpcx/log"
)

func init() {
	framework.LoadConfig()
}

func main() {
	//share.Trace = true
	s := framework.NewService()
	log.Infof("Service-> %+v; Etcd-> %+v", *s, *(s.Etcd))
	s.Serve()
}
