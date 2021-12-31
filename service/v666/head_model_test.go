package v666

import (
	framework "github.com/ZYallers/rpcx-framework"
	"github.com/ZYallers/rpcx-framework/define"
	"github.com/ZYallers/rpcx-framework/util/client"
	"testing"
)

func init() {
	framework.LoadConfig("../../")
}

func TestHeadModel_LatestModel(t *testing.T) {
	args := define.M{}
	if reply, err := client.XClient(framework.ServiceName(), "head/model/latest", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}

func TestHeadModel_SwitchModel(t *testing.T) {
	args := define.M{"model": 1, "admin_user_id": 1}
	if reply, err := client.XClient(framework.ServiceName(), "head/model/switch", args); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", string(reply.([]byte)))
	}
}
