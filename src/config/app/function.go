package app

import (
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/smallnest/rpcx/log"
	"os"
	"src/config/define"
	"strings"
	"time"
)

// PushSimpleMessage
func PushSimpleMessage(msg string, isAtAll bool, logType ...interface{}) {
	hostname, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"Env: " + Service.Env,
		"App: " + Service.Name,
		"Addr: " + Service.Addr,
		"HostName: " + hostname,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + tool.SystemIP(),
		"PublicIP: " + tool.PublicIP(),
	}
	if Service.Env != define.ProduceMode {
		isAtAll = false // 开发环境下，不需要@所有人，减少干扰!
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + Service.GracefulRobotToken
	if resp, err := tool.NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json"}).
		SetPostData(postData).SetTimeOut(time.Second).Post(); err != nil {
		log.Errorf("push simple message error: %v, resp: %v", err, resp)
	}

	if len(logType) == 1 {
		switch logType[0] {
		case "debug":
			log.Debug(msg)
		case "info":
			log.Info(msg)
		case "warn":
			log.Warn(msg)
		case "error":
			log.Error(msg)
		case "fatal":
			log.Fatal(msg)
		case "panic":
			log.Panic(msg)
		}
	}
}

// PushContextMessage
func PushContextMessage(msg string, stack string, isAtAll bool, logType ...interface{}) {
	hostname, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"Env: " + Service.Env,
		"App: " + Service.Name,
		"Addr: " + Service.Addr,
		"HostName: " + hostname,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + tool.SystemIP(),
		"PublicIP: " + tool.PublicIP(),
	}
	if stack != "" {
		text = append(text, "\nStack:\n"+stack)
	}
	if Service.Env != define.ProduceMode {
		isAtAll = false // 开发环境下，不需要@所有人，减少干扰!
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + Service.ErrorRobotToken

	if resp, err := tool.NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json"}).
		SetPostData(postData).SetTimeOut(time.Second).Post(); err != nil {
		log.Errorf("push context message error: %v, resp: %v", err, resp)
	}

	if len(logType) == 1 {
		switch logType[0] {
		case "debug":
			log.Debug(msg)
		case "info":
			log.Info(msg)
		case "warn":
			log.Warn(msg)
		case "error":
			log.Error(msg)
		case "fatal":
			log.Fatal(msg)
		case "panic":
			log.Panic(msg)
		}
	}
}
