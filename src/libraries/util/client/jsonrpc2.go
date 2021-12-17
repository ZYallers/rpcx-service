package client

import (
	"bytes"
	"fmt"
	"github.com/ZYallers/zgin/libraries/tool"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"math/rand"
	"src/config/app"
	"src/libraries/util/helper"
	"time"
)

// {\"jsonrpc\":\"2.0\",\"error\":{\"code\":-32601,\"message\":\"must contains servicepath and method\",\"data\":null},\"id\":8674665223082153551}
type jsonRpc struct {
	Id      int                    `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Result  interface{}            `json:"result"`
	Error   struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
}

func JsonRpc2(serviceMethod string, args map[string]interface{}, other ...interface{}) (interface{}, error) {
	req := tool.NewRequest(fmt.Sprintf("http://%s", app.Service.Addr))
	req.SetHeaders(map[string]string{"X-JSONRPC-2.0": "true"})
	data := jsonRpc{
		Id:      rand.Int(),
		Jsonrpc: "2.0",
		Method:  fmt.Sprintf("%s.%s", app.Service.Name, serviceMethod),
		Params:  args,
	}
	b, _ := jsontime.ConfigWithCustomTimeFormat.Marshal(data)
	req.SetBody(bytes.NewReader(b))
	if len(other) > 0 {
		req.SetTimeOut(other[0].(time.Duration))
	}
	resp, err := req.Post()
	if err != nil {
		return nil, err
	}
	if resp.Body == "" {
		return nil, nil
	}
	var res jsonRpc
	if err := jsontime.ConfigWithCustomTimeFormat.Unmarshal(helper.String2Bytes(resp.Body), &res); err != nil {
		return nil, err
	}
	if res.Error.Message != "" {
		return nil, fmt.Errorf("jsonRpc2 error: code:%d, message:%s, data:%v", res.Error.Code, res.Error.Message, res.Error.Data)
	}
	return res.Result, nil
}
