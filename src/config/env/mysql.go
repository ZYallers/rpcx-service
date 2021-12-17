package env

import (
	zapp "github.com/ZYallers/zgin/app"
	"os"
	"src/config/app"
	"src/config/define"
)

var Mysql = &struct {
	EnjoyThin zapp.MysqlDialect
}{}

func init() {
	switch app.Env() {
	case define.DevelopMode, define.GrayMode:
		Mysql.EnjoyThin = zapp.MysqlDialect{Host: os.Getenv("mysql_host"), Port: os.Getenv("mysql_port"),
			User: os.Getenv("mysql_username"), Pwd: os.Getenv("mysql_password"), Db: os.Getenv("mysql_database")}
	case define.ProduceMode:
		Mysql.EnjoyThin = zapp.MysqlDialect{Host: "xxx", Port: "3306", User: "enjoythin", Pwd: "xxx", Db: "enjoythin"}
	}
}
