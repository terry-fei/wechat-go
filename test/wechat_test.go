package wechat_test

import (
	"testing"
	"wechat"
)

func Test_getSignature(t *testing.T) {
	wechat := Wechat{token: "feit"}
	sign := wechat.getSignature("haha", "djdjdj", "djdjds")
	t.Log(sign)
}
