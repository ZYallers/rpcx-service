package app

import (
	"context"
	"github.com/smallnest/rpcx/server"
	"github.com/syyongx/php2go"
	"reflect"
	"src/config/define"
)

const stateActive = "state=active"

func RegisterFuncName(s *server.Server, services define.Restful) error {
	if err := registerHealthFunc(s); err != nil {
		return err
	}
	if len(services) <= 0 {
		return define.ErrRegisterServiceEmpty
	}
	if err := registerServiceMethod(s, &services); err != nil {
		return err
	}
	return nil
}

func registerHealthFunc(s *server.Server) error {
	return s.RegisterFunctionName(Service.Name, "health", func(ctx context.Context,
		args map[string]interface{}, reply *interface{}) error {
		*reply = "ok"
		return nil
	}, stateActive)
}

func registerServiceMethod(s *server.Server, services *define.Restful) error {
	for path, handlers := range *services {
		if err := s.RegisterFunctionName(Service.Name, path, dispatchHandler(handlers), stateActive); err != nil {
			return err
		}
	}
	return nil
}

func dispatchHandler(handlers []define.RestHandler) func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
	return func(ctx context.Context, args map[string]interface{}, reply *interface{}) error {
		argsVersion := Service.Version
		if ver, ok := args[Service.VersionKey].(string); ok && ver != "" {
			argsVersion = ver
		}
		if handler := versionCompare(&handlers, argsVersion); handler == nil {
			return define.ErrVersionCompare
		} else {
			v := reflect.ValueOf(handler.Service)
			ptr := reflect.New(v.Type().Elem())
			ptr.Elem().Set(v.Elem())
			s := ptr.Interface().(define.IService)
			s.Construct(Service, ctx, args, reply)
			if handler.Signed && !s.SignCheck() {
				return define.ErrSignature
			}
			if handler.Logged && !s.LoginCheck() {
				return define.ErrNeedLogin
			}
			result := ptr.MethodByName(handler.Method).Call(nil)
			if result[0].IsNil() {
				return nil
			}
			return result[0].Interface().(error)
		}
	}
}

func versionCompare(handlers *[]define.RestHandler, version string) *define.RestHandler {
	for _, handler := range *handlers {
		if handler.Version == "" || handler.Version == version {
			return &handler
		}
		if le := len(handler.Version); handler.Version[le-1:] == "+" {
			vs := handler.Version[0 : le-1]
			if version == vs {
				return &handler
			}
			if php2go.VersionCompare(version, vs, ">") {
				return &handler
			}
		}
	}
	return nil
}
