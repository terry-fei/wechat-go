package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func init() {
	log.Println("wechat init")
}

type WechatHandleFunc func(msg *Message)

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

func CheckSignature(query url.Values, token string) bool {
	// sort by dict
	arr := sort.StringSlice{token, query.Get("nonce"), query.Get("timestamp")}
	arr.Sort()

	// sha1 hex
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(strings.Join(arr, "")))
	cipherStr := hex.EncodeToString(sha1Ctx.Sum(nil))

	if query.Get("signature") == cipherStr {
		return true
	}
	return false
}

func HandleMessage(token string, handler WechatHandleFunc) http.HandlerFunc {
	wechat := New(token, handler)
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if !CheckSignature(r.Form, wechat.Token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Signature"))
			return
		}

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.Form.Get("echostr")))

		} else if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				// TODO handle err
				log.Println(err)
			}
			msg := Message{
				Msg:      XMLToMessage(body),
				ReplyMsg: make(map[string]interface{}),
			}

			wechat.handler(&msg)
			ReplyMsgTemplate.Execute(os.Stdout, msg.ReplyMsg)
			ReplyMsgTemplate.Execute(w, msg.ReplyMsg)

		} else {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Not Implemented"))
		}
	}
}
