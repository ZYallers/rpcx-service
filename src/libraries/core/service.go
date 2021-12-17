package core

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ZYallers/zgin/libraries/tool"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/smallnest/rpcx/server"
	"net/http"
	"src/config/define"
	"strconv"
	"time"
)

const (
	debugValue = "xxxxxxx"
	tokenKey   = "sess_token"
	signKey    = "sign"
	devSignKey = "xxxxxx"
	utimeKey   = "utime"
)

type Service struct {
	debug   bool
	ctx     context.Context
	service *define.Service
	args    map[string]interface{}
	reply   *interface{}
}

func (s *Service) Construct(service *define.Service, ctx context.Context, args map[string]interface{}, reply *interface{}) {
	s.service = service
	s.ctx = ctx
	s.args = args
	s.reply = reply
	rep := &define.Reply{}
	if debug := s.GetArgs("debug"); debug == debugValue {
		s.debug = true
		now := time.Now()
		rep.Service = &define.ReplyService{
			Name:     s.service.Name,
			Hostname: s.service.HostName,
			Ip:       s.service.SystemIP,
			Addr:     s.service.Addr,
			Start:    &now,
		}
	}
	*s.reply = rep
}

//  GetArgs 获取客户端传参值
//  @receiver s *Service
//  @author Cloud|2021-12-02 16:18:05
//  @param key string ...
//  @param defaultValue ...interface{} ...
//  @return interface{} ...
func (s *Service) GetArgs(key string, defaultValue ...interface{}) interface{} {
	if val, exist := s.args[key]; exist {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (s *Service) GetString(key string, defaultValue ...string) string {
	if val, exist := s.args[key]; exist {
		var res string
		switch v := val.(type) {
		case int:
			res = strconv.Itoa(v)
		case int8:
			res = strconv.Itoa(int(v))
		default:
			res = fmt.Sprintf("%v", v)
		}
		return res
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (s *Service) GetInt(key string, defaultValue ...int) int {
	if val, exist := s.args[key]; exist {
		var res int
		switch v := val.(type) {
		case string:
			res, _ = strconv.Atoi(v)
		case int8:
			res = int(v)
		default:
			res, _ = strconv.Atoi(fmt.Sprintf("%d", v))
		}
		return res
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (s *Service) Finish(data interface{}) error {
	resp := (*s.reply).(*define.Reply)
	resp.Data = data
	if s.debug {
		now := time.Now()
		resp.Service.End = &now
		resp.Service.Runtime = time.Since(*resp.Service.Start).String()
	}
	*s.reply = resp
	return nil
}

//  Record ...
//  @receiver s *Service
//  @author Cloud|2021-12-07 15:31:39
//  @param r define.Record ...
func (s *Service) Record(r define.Record) {
	resp := (*s.reply).(*define.Reply)
	resp.Record = &r
	*s.reply = resp
}

//  Json ...
//  @receiver s *Service
//  @author Cloud|2021-12-07 14:12:03
//  @param args ...interface{} ...
//  @return error ...
func (s *Service) Json(args ...interface{}) error {
	rep := (*s.reply).(*define.Reply)
	if len(args) == 0 {
		rep.Code = http.StatusOK
	}
	if len(args) > 0 {
		rep.Code = args[0].(int)
	}
	if len(args) > 1 {
		switch value := args[1].(type) {
		case error:
			rep.Msg = value.Error()
		case string:
			rep.Msg = value
		default:
			rep.Msg = fmt.Sprintf("%v", value)
		}
	}
	if len(args) > 2 {
		rep.Data = args[2]
	}
	if s.debug {
		now := time.Now()
		rep.Service.End = &now
		rep.Service.Runtime = time.Since(*rep.Service.Start).String()
	}
	if s.ctx.Value(server.HttpConnContextKey) == nil {
		bte, err := jsontime.ConfigWithCustomTimeFormat.Marshal(rep)
		if err != nil {
			return s.Json(http.StatusInternalServerError, err)
		}
		*s.reply = bte
	} else {
		*s.reply = rep
	}
	return nil
}

//  SignCheck APP签名验证
//  @receiver s *Service
//  @author Cloud|2021-12-02 16:57:44
//  @return bool ...
func (s *Service) SignCheck() bool {
	sign := s.GetArgs(signKey, "").(string)
	if sign == "" {
		return false
	}
	// 开发测试模式下，用固定sign判断
	if s.service.Env == define.DevelopMode && sign == devSignKey {
		return true
	}
	timestampStr := s.GetArgs(utimeKey, "").(string)
	if timestampStr == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return false
	}
	if time.Now().Unix()-timestamp > int64(define.SignTimeExpiration) {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(timestampStr + define.TokenKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	if sign == base64.StdEncoding.EncodeToString([]byte(md5str)) {
		return true
	}
	return false
}

//  GetLoggedUserData 获取APP登录用户数据
//  @receiver s *Service
//  @author Cloud|2021-12-02 17:08:20
//  @param values ...string ...
//  @return map[string]interface{} ...
func (s *Service) GetLoggedUserData(values ...string) map[string]interface{} {
	var token string
	switch len(values) {
	case 1:
		token = values[0]
	default:
		token = s.GetArgs(tokenKey, "").(string)
	}
	client := new(Redis).GetSession()
	if client == nil {
		return nil
	}
	if str, _ := client.Get("ci_session:" + token).Result(); str == "" {
		return nil
	} else {
		return tool.PhpUnserialize(str)
	}
}

//  LoginCheck APP登录检查
//  @receiver s *Service
//  @author Cloud|2021-12-02 17:07:53
//  @param values ...string ...
//  @return bool ...
func (s *Service) LoginCheck(values ...string) bool {
	if vars := s.GetLoggedUserData(values...); vars != nil {
		return true
	}
	return false
}

//  GetLoggedUserId 获取APP登陆用户的user_id
//  @receiver s *Service
//  @author Cloud|2021-12-02 16:01:32
//  @param token string
//  @return int user_id
func (s *Service) GetLoggedUserId(values ...string) int {
	vars := s.GetLoggedUserData(values...)
	if vars == nil {
		return 0
	}
	if data, ok := vars["userinfo"].(map[string]interface{}); ok {
		if str, ok := data["userid"].(string); ok && str != "" {
			userId, _ := strconv.Atoi(str)
			return userId
		}
	}
	return 0
}
