package main

import (
	"fmt"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/soheilhy/cmux"
	"os"
	"src/config/app"
	"src/config/define"
	"src/config/restful"
	"time"
)

func init() {
	timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
	jsontime.SetDefaultTimeFormat(define.TimeFormat, timeZoneShanghai)
}

func main() {
	/*if app.Service.Env == define.DevelopMode {
		share.Trace = true
	}*/
	app.Bootstrap()
	app.Service.Server = server.NewServer(
		//server.WithReadTimeout(10*time.Second),
		//server.WithWriteTimeout(15*time.Second),
		server.WithTCPKeepAlivePeriod(time.Minute),
	)
	app.Service.Server.DisableJSONRPC = false
	app.Service.Server.DisableHTTPGateway = true

	if err := addRegistryPlugin(app.Service.Server); err != nil {
		msg := fmt.Sprintf("register etcdv3 plugin error: %s", err)
		app.PushSimpleMessage(msg, true, "panic")
	}

	if err := app.RegisterFuncName(app.Service.Server, restful.GetServices()); err != nil {
		msg := fmt.Sprintf("register function name error: %s", err)
		app.PushSimpleMessage(msg, true, "panic")
	}

	app.Service.Server.RegisterOnRestart(func(s *server.Server) {
		msg := fmt.Sprintf("%s service(%d) is restarting...", app.Service.Name, os.Getpid())
		app.PushSimpleMessage(msg, true, "info")
	})

	app.Service.Server.RegisterOnShutdown(func(s *server.Server) {
		msg := fmt.Sprintf("%s service(%d) is shutting down...", app.Service.Name, os.Getpid())
		app.PushSimpleMessage(msg, true, "info")
	})

	msg := fmt.Sprintf("%s service(%d) is ready to serve", app.Service.Name, os.Getpid())
	app.PushSimpleMessage(msg, true, "info")
	if err := app.Service.Server.Serve("tcp", app.Service.Addr);
		err != nil && err != server.ErrServerClosed && err != cmux.ErrServerClosed {
		msg := fmt.Sprintf("%s service serve error: %v", app.Service.Name, err)
		app.PushSimpleMessage(msg, true, "panic")
	}
}

func addRegistryPlugin(s *server.Server) error {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + app.Service.Addr,
		EtcdServers:    app.Service.Etcd.Addr,
		BasePath:       app.Service.Etcd.BasePath,
		UpdateInterval: app.Service.Etcd.UpdateInterval,
	}
	if err := r.Start(); err != nil {
		return err
	}
	s.Plugins.Add(r)
	return nil
}
