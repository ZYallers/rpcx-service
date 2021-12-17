package client

import (
	"log"
	"net/rpc/jsonrpc"
	"testing"
)

func Test_JsonRpcCall(t *testing.T) {
	client, err := jsonrpc.Dial("tcp", "172.18.28.123:8978")
	if err != nil {
		log.Fatal("dial error:", err)
	}

	args := map[string]interface{}{}
	var reply interface{}
	err = client.Call("health", args, &reply)
	if err != nil {
		t.Fatal("client.Call error:", err)
	}
	t.Logf("reply: %#v\n", reply)
}

func Test_JsonRpc2(t *testing.T) {
	resp, err := JsonRpc2("head/model/latest", map[string]interface{}{"name": "ddd"})
	log.Printf("resp: %#v, error: %v\n", resp, err)
}
