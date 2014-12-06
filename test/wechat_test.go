package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	"wechat"
)

func Test_checkSignature(t *testing.T) {
	wechat := new(wechat.Wechat)
	wechat.Token = "some token"
	timestamp := time.Now().Unix()
	rand.Seed(timestamp)
	nonce := rand.Intn(1e10)

	query := make(map[string]string)
	query["nonce"] = strconv.Itoa(nonce)
	query["timestamp"] = strconv.FormatInt(timestamp, 10)

	arr := sort.StringSlice{wechat.Token, query["nonce"], query["timestamp"]}
	arr.Sort()

	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(strings.Join(arr, "")))
	query["signature"] = hex.EncodeToString(sha1Ctx.Sum(nil))
	t.Log(query)

	ret := wechat.CheckSignature(query)
	if ret {
		t.Log("check signature ok")
	} else {
		t.Error("check signature fail")
	}
}
