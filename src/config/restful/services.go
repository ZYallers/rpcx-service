package restful

import (
	"fmt"
	"reflect"
	"sort"
	"src/config/define"
	"strconv"
	"sync"
)

var (
	lock     sync.Mutex
	services []define.IService
)

func GetServices() define.Restful {
	res := define.Restful{}
	for _, service := range services {
		serviceValueOf := reflect.ValueOf(service)
		serviceName := serviceValueOf.Elem().Type().Name()
		if _, exist := serviceValueOf.Elem().Type().FieldByName("tag"); !exist {
			continue
		}
		tagVal := serviceValueOf.Elem().FieldByName("tag")
		if tagVal.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < tagVal.NumField(); i++ {
			if tagVal.Field(i).Kind() != reflect.Func {
				continue
			}
			methodName := tagVal.Type().Field(i).Name
			fieldTagVal := tagVal.Type().Field(i).Tag
			path := fieldTagVal.Get("path")
			if path == "" {
				panic(fmt.Errorf("restHandler.Path is empty: %s.%s\n", serviceName, methodName))
			}
			if _, exist := serviceValueOf.Type().MethodByName(methodName); !exist {
				panic(fmt.Errorf("restHandler.Method does not exist: %s.%s\n", serviceName, methodName))
			}

			resHandler := define.RestHandler{
				Path:    path,
				Service: service,
				Method:  methodName,
				Version: fieldTagVal.Get("ver"),
				Signed:  fieldTagVal.Get("sign") == "on",
				Logged:  fieldTagVal.Get("login") == "on",
			}
			if sortStr := fieldTagVal.Get("sort"); sortStr != "" {
				if sortInt, err := strconv.Atoi(sortStr); err != nil {
					panic(fmt.Errorf("restHandler sort is invalid: %s", sortStr))
				} else {
					resHandler.Sort = sortInt
				}
			}
			res[path] = append(res[path], resHandler)
			if len(res[path]) > 1 {
				resHandlers := res[path]
				sort.SliceStable(resHandlers, func(i, j int) bool {
					return resHandlers[i].Sort > resHandlers[j].Sort
				})
				res[path] = resHandlers
			}
		}
	}
	return res
}

func register(s ...define.IService) {
	lock.Lock()
	defer lock.Unlock()
	services = append(services, s...)
}
