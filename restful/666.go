package restful

import (
	"github.com/ZYallers/rpcx-framework/service"
	v666 "github.com/ZYallers/rpcx-service/service/v666"
)

func init() {
	service.Register(&v666.HeadModel{}, &v666.HeadBanner{})
}
