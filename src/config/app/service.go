package app

import (
	"fmt"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"src/config/define"
	"src/libraries/util/zap"
	"strings"
	"time"
)

var Service *define.Service

func Env() string {
	if Service == nil {
		if modeKey := viper.GetString("global.modeKey"); modeKey != "" {
			if val := os.Getenv(modeKey); val != "" {
				return val
			}
		}
	} else {
		return Service.Env
	}
	return define.DevelopMode
}

func init() {
	viper.SetConfigName("service")
	viper.SetConfigType("json")
	_, fullFilename, _, _ := runtime.Caller(0)
	viper.AddConfigPath(path.Join(path.Dir(fullFilename), "/../../"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config file error: %s", err))
	}
}

func Bootstrap() {
	defer func() {
		if err := recover(); err != nil {
			PushSimpleMessage(fmt.Sprintf("%v", err), true, "panic")
		}
	}()
	serviceName := viper.GetString("service.name")
	if serviceName == "" {
		panic("service name is empty")
	}

	serviceLogDir := viper.GetString("service.logDir")
	if serviceLogDir == "" {
		serviceLogDir = "/apps/logs/go/" + serviceName
	}
	zap.SetLoggerDir(serviceLogDir)
	log.SetLogger(NewLogger(serviceName))

	runMode := Env()

	etcdBasePath := viper.GetString(fmt.Sprintf("service.etcd.%s.basePath", runMode))
	if etcdBasePath == "" {
		panic("etcd base path is empty")
	}

	etcdAddr := viper.GetString(fmt.Sprintf("service.etcd.%s.addr", runMode))
	if etcdAddr == "" {
		panic("etcd address is empty")
	}

	etcdUpdateInterval := viper.GetInt64(fmt.Sprintf("service.etcd.%s.updateInterval", runMode))
	if etcdUpdateInterval <= 0 {
		etcdUpdateInterval = 30
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		panic("system hostname is empty")
	} else {
		hostname = strings.ToLower(hostname)
	}

	if runMode == define.DevelopMode && hostname != viper.GetString("global.server.development.hostname") {
		etcdBasePath = strings.Replace(etcdBasePath, define.DevelopMode, "developer@"+hostname, 1)
		developmentServerIP := viper.GetString("global.server.development.ip")
		etcdAddr = strings.Replace(etcdAddr, "127.0.0.1", developmentServerIP, 1)
	}

	sip := tool.SystemIP()
	if sip == "unknown" || sip == "" {
		panic("system ip is unknown or empty")
	}

	Service = &define.Service{
		Env:                runMode,
		Name:               serviceName,
		HostName:           hostname,
		SystemIP:           sip,
		LogDir:             serviceLogDir,
		Addr:               strings.Replace(viper.GetString("service.addr"), "0.0.0.0", sip, 1),
		Version:            viper.GetString("service.version"),
		VersionKey:         viper.GetString("service.versionKey"),
		ErrorRobotToken:    viper.GetString("service.errorRobotToken"),
		GracefulRobotToken: viper.GetString("service.gracefulRobotToken"),
		Etcd: &define.Discovery{
			BasePath:       etcdBasePath,
			Addr:           strings.Split(etcdAddr, ","),
			UpdateInterval: time.Duration(etcdUpdateInterval) * time.Second,
		},
	}

	log.Infof("Service-> %+v; Etcd-> %+v", *Service, *(Service.Etcd))
}
