package client

import (
	"log"
	"math/rand"
	"src/config/app"
	"testing"
)

func Benchmark_RandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Log(rand.Int())
	}
}

func init() {
	app.Bootstrap()
}

func Test_HttpInvoke(t *testing.T) {
	args := map[string]interface{}{"name": "HaiMei"}
	reply, err := HttpInvoke("test/check", args)
	log.Printf("reply: %#v, error: %v\n", reply, err)
}
