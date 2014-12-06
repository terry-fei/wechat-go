package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func wechatHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nonce = r.Form["nonce"]
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello feit!")
}

func main() {
	http.HandleFunc("/", wechatHandle)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
