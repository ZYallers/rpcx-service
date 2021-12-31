package v666

import (
	framework "github.com/ZYallers/rpcx-framework"
	"github.com/ZYallers/rpcx-framework/define"
	"github.com/ZYallers/rpcx-framework/util/client"
	"github.com/ZYallers/zgin/libraries/json"
	"testing"
)

func init() {
	framework.LoadConfig("../../")
}

func TestHeadBanner_Rotate(t *testing.T) {
	args := define.M{}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/rotate", args); err != nil {
		t.Fatal(err)
	} else {
		var res map[string]interface{}
		_ = json.Unmarshal(reply.([]byte), &res)
		t.Logf("%+v\n", res)
		for _, v := range res["data"].([]interface{}) {
			t.Logf("%+v\n", v)
		}
	}
}

func TestHeadBanner_Save(t *testing.T) {
	//args := define.M{
	//	//"id":            2,
	//	"title":         "开发测试标题24",
	//	"image":         "2016-07-01/4cdcd196c68448baaa05fb5f3ba8019e.jpg",
	//	"url":           "https://www.hxsapp.com",
	//	"sort":          1,
	//	"start_time":    "2021-12-02 19:50:00",
	//	"end_time":      "2021-12-08 18:00:00",
	//	"model":         1,
	//	"admin_user_id": 1,
	//}
	args := define.M{
		"title":         "开发测试标题25",
		"image":         "2016-07-01/4cdcd196c68448baaa05fb5f3ba8019e.jpg",
		"admin_user_id": 1,
		"model":         2,
	}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/save", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadBanner_Edit(t *testing.T) {
	args := define.M{"id": 1}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/edit", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadBanner_Delete(t *testing.T) {
	args := define.M{"id": 2}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/delete", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadBanner_OnLine(t *testing.T) {
	args := define.M{"id": 2}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/online", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadBanner_OffLine(t *testing.T) {
	args := define.M{"id": 10}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/offline", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadBanner_List(t *testing.T) {
	args := define.M{"model": 1, "state": 1}
	if reply, err := client.XClient(framework.ServiceName(), "head/banner/list", args); err != nil {
		t.Fatal(err)
	} else {
		var res map[string]interface{}
		_ = json.Unmarshal(reply.([]byte), &res)
		t.Logf("%+v\n", res)
		for _, v := range res["data"].([]interface{}) {
			t.Logf("%+v\n", v)
		}
	}
}
