wechat
======

简单高效的微信公众平台消息处理库

Installation
======
```
$ go get github.com/feit/wechat
```

Usage
======
```go
package main
import (
  "net/http"
  "wechat"
)

func wechatHandler(this *wechat.Message) {

  switch this.Msg["Content"] {
    case "text":
    this.Reply("hello text")

    case "news":
      news := []wechat.News{}
      news = append(news, wechat.News{
        Title:  "Hello Golang",
        Desc:   "Golang is Good",
        Url:    "feit.me",
        PicUrl: "feit.me/aa.png",
      })
      this.Reply(news)

    case "image":
      replyMsg := wechat.Image{
        MediaId: "mediaid",
      }
      this.Reply(replyMsg)

    case "voice":
      replyMsg := wechat.Voice{
        MediaId: "mediaid",
      }
      this.Reply(replyMsg)

    case "video":
      replyMsg := wechat.Video{
        Title:   "VideoMsg",
        Desc:    "this is a video msg",
        MediaId: "mediaid",
      }
      this.Reply(replyMsg)

    case "music":
      replyMsg := wechat.Music{
        Title:        "GO GO GO",
        Desc:         "song for Golang",
        MusicUrl:     "music.com/golang.mp3",
        HQMusicUrl:   "music.com/golang.flac",
        ThumbMediaId: "ididid",
      }
      this.Reply(replyMsg)

    case "cs":
      replyMsg := wechat.CustomService{
        KfAccount: "test1@test",
      }
      this.Reply(replyMsg)

    default:
      this.Reply("hehe")

  }
}

func main() {
  http.HandleFunc("/wechat/api", wechat.HandleMessage("token", wechatHandler))
  http.ListenAndServe(":80", nil)
}
```

License
======
The MIT licence
