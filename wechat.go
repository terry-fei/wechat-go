package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Wechat struct {
	token string
}

func (w Wechat) getSignature(signature, nonce, timestamp string) bool {
	// sort by dict
	arr := sort.StringSlice{w.token, nonce, timestamp}
	arr.Sort()

	// sha1 hex
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(strings.Join(arr, "")))
	cipherStr := hex.EncodeToString(sha1Ctx.Sum(nil))

	if signature == cipherStr {
		return true
	}
	return false
}

func main() {
	wechat := Wechat{token: "feit"}
	wechat.getSignature("aaa", "bbb", "ccc")
}
