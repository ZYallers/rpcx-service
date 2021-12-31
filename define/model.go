package define

import (
	"github.com/ZYallers/rpcx-framework/util/mtsc"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"gorm.io/gorm"
)

type Model struct {
	mtsc.Model
}

var (
	enjoyThin        mvcs.DbCollector
	enjoyThinDialect *app.MysqlDialect
)

func init() {
	enjoyThinDialect = &app.MysqlDialect{
		User: "",
		Pwd:  "",
		Host: "",
		Port: "",
		Db:   "",
	}
}

func (m *Model) EnjoyThin() *gorm.DB {
	return m.NewClient(&enjoyThin, enjoyThinDialect)
}
