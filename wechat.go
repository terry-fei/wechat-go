package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type WechatHandleFunc func(info map[string]string)

type Wechat struct {
	Token   string
	handler WechatHandleFunc
}

func New(token string, handler WechatHandleFunc) (wechat *Wechat) {
	wechat = &Wechat{}
	wechat.Token = token
	wechat.handler = handler
	return
}

func (w *Wechat) CheckSignature(query map[string]string) bool {
	// sort by dict
	arr := sort.StringSlice{w.Token, query["nonce"], query["timestamp"]}
	arr.Sort()

	// sha1 hex
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(strings.Join(arr, "")))
	cipherStr := hex.EncodeToString(sha1Ctx.Sum(nil))

	if query["signature"] == cipherStr {
		return true
	}
	return false
}

func HandleMessage(token string, handler WechatHandleFunc) http.HandlerFunc {
	wechat := New(token, handler)
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Read Request Content Fail")
		}
		fmt.Println(string(body))
		info := make(map[string]string)
		info["FromUserName"] = "feit"
		wechat.handler(info)
		fmt.Fprintf(w, "Hello Wechat")
	}
}
