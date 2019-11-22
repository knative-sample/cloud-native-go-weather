package tracing

import (
	"testing"
	"time"
)

func TestExampleNewTracer(t *testing.T) {
	t.Log("TestExampleNewTracer test")

	// 手动埋点
	ExampleNewTracer()

	// 用http框架埋点
	HttpExample()

	// 等待一下，trace 数据上报是异步的，如果不等待就无法得到 trace 数据
	time.Sleep(time.Second * 5)
}
