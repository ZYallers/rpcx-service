package define

import (
	"context"
	"github.com/smallnest/rpcx/server"
	"time"
)

type Discovery struct {
	BasePath       string
	UpdateInterval time.Duration
	Addr           []string
}

type Service struct {
	Env                string
	Version            string
	VersionKey         string
	Name               string
	HostName           string
	SystemIP           string
	Addr               string
	LogDir             string
	ErrorRobotToken    string
	GracefulRobotToken string
	Etcd               *Discovery
	Server             *server.Server
}

type IService interface {
	Construct(service *Service, ctx context.Context, args map[string]interface{}, reply *interface{})
	SignCheck() bool
	LoginCheck(values ...string) bool
}

type ReplyService struct {
	Name     string     `json:"name,omitempty"`
	Hostname string     `json:"hostname,omitempty"`
	Ip       string     `json:"ip,omitempty"`
	Addr     string     `json:"addr,omitempty"`
	Runtime  string     `json:"runtime,omitempty"`
	Start    *time.Time `json:"start,omitempty"`
	End      *time.Time `json:"end,omitempty"`
}

type Record struct {
	Type      string      `json:"type,omitempty"`
	TableName string      `json:"table_name,omitempty"`
	DataId    interface{} `json:"data_id,omitempty"`
	Intro     string      `json:"intro,omitempty"`
}

type Reply struct {
	Code    int           `json:"code"`
	Msg     string        `json:"msg"`
	Data    interface{}   `json:"data"`
	Record  *Record       `json:"record,omitempty"`
	Service *ReplyService `json:"service,omitempty"`
}
