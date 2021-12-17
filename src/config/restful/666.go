package restful

import (
	v666 "src/service/v666"
)

func init() {
	register(&v666.HeadModel{}, &v666.HeadBanner{})
}
