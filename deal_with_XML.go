package wechat

import (
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	ToUserNameRe     *regexp.Regexp
	FromUserNameRe   *regexp.Regexp
	CreateTimeRe     *regexp.Regexp
	MsgTypeRe        *regexp.Regexp
	ContentRe        *regexp.Regexp
	MsgIdRe          *regexp.Regexp
	ReplyMsgTemplate *template.Template
)

func init() {
	ToUserNameRe = regexp.MustCompile(`<ToUserName><!\[CDATA\[(\S+)\]\]></ToUserName>`)
	FromUserNameRe = regexp.MustCompile(`<FromUserName><!\[CDATA\[(\S+)\]\]></FromUserName>`)
	CreateTimeRe = regexp.MustCompile(`<CreateTime>(\d+)</CreateTime>`)
	MsgTypeRe = regexp.MustCompile(`<MsgType><!\[CDATA\[(\S+)\]\]></MsgType>`)
	MsgIdRe = regexp.MustCompile(`<MsgId>(\d+)</MsgId>`)
	ContentRe = regexp.MustCompile(`<Content><!\[CDATA\[(.*)\]\]></Content>`)
	replyMsgTemplateStr := strings.Join(replyMsgTemplateStrArr, "")
	ReplyMsgTemplate = template.Must(template.New("replyMessageTpl").Parse(replyMsgTemplateStr))
}

func XMLToMessage(xml []byte) (msg map[string]interface{}) {
	msg = make(map[string]interface{})

	// get commom attrbute
	msg["ToUserName"] = string(ToUserNameRe.FindSubmatch(xml)[1])
	msg["FromUserName"] = string(FromUserNameRe.FindSubmatch(xml)[1])
	msg["CreateTime"] = string(CreateTimeRe.FindSubmatch(xml)[1])
	msg["MsgType"] = string(MsgTypeRe.FindSubmatch(xml)[1])

	// get other attrbute by MsgType
	switch msg["MsgType"] {
	case "text":
		msg["Content"] = string(ContentRe.FindSubmatch(xml)[1])
		msg["MsgId"] = string(MsgIdRe.FindSubmatch(xml)[1])
	case "event":
	case "link":
	case "image":
	case "voice":
	case "video":
	case "location":
	default:
		// TODO handle err
	}

	return
}

type Message struct {
	Msg      map[string]interface{}
	ReplyMsg map[string]interface{}
}

func (msg *Message) Reply(replyMsg interface{}) {
	msg.ReplyMsg["ToUserName"] = msg.Msg["FromUserName"]
	msg.ReplyMsg["FromUserName"] = msg.Msg["ToUserName"]
	msg.ReplyMsg["CreateTime"] = time.Now().Unix()

	if value, ok := replyMsg.(string); ok {
		msg.ReplyMsg["MsgType"] = "text"
		msg.ReplyMsg["IsText"] = true
		msg.ReplyMsg["Content"] = value

	} else if value, ok := replyMsg.([]News); ok {
		msg.ReplyMsg["MsgType"] = "news"
		msg.ReplyMsg["IsNews"] = true
		msg.ReplyMsg["NewsLength"] = len(value)
		msg.ReplyMsg["News"] = value

	} else if value, ok := replyMsg.(CustomService); ok {
		msg.ReplyMsg["MsgType"] = "transfer_customer_service"
		msg.ReplyMsg["NeedCustomService"] = true
		msg.ReplyMsg["KfAccount"] = value.KfAccount

	} else if value, ok := replyMsg.(Image); ok {
		msg.ReplyMsg["MsgType"] = "image"
		msg.ReplyMsg["IsImage"] = true
		msg.ReplyMsg["MediaId"] = value.MediaId

	} else if value, ok := replyMsg.(Voice); ok {
		msg.ReplyMsg["MsgType"] = "voice"
		msg.ReplyMsg["IsVoice"] = true
		msg.ReplyMsg["MediaId"] = value.MediaId

	} else if value, ok := replyMsg.(Video); ok {
		msg.ReplyMsg["MsgType"] = "video"
		msg.ReplyMsg["IsVideo"] = true
		msg.ReplyMsg["Title"] = value.Title
		msg.ReplyMsg["Desc"] = value.Desc
		msg.ReplyMsg["MediaId"] = value.MediaId

	} else if value, ok := replyMsg.(Music); ok {
		msg.ReplyMsg["MsgType"] = "music"
		msg.ReplyMsg["IsMusic"] = true
		msg.ReplyMsg["Title"] = value.Title
		msg.ReplyMsg["Desc"] = value.Desc
		msg.ReplyMsg["MusicUrl"] = value.MusicUrl
		msg.ReplyMsg["HQMusicUrl"] = value.HQMusicUrl
		msg.ReplyMsg["ThumbMediaId"] = value.ThumbMediaId

	}
}

type News struct {
	Title  string
	Desc   string
	PicUrl string
	Url    string
}

type Image struct {
	MediaId string
}

type Voice struct {
	MediaId string
}

type Video struct {
	Title   string
	Desc    string
	MediaId string
}

type Music struct {
	Title        string
	Desc         string
	MusicUrl     string
	HQMusicUrl   string
	ThumbMediaId string
}

type CustomService struct {
	KfAccount string
}

var replyMsgTemplateStrArr = []string{
	"<xml>",
	"<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>",
	"<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>",
	"<CreateTime>{{.CreateTime}}</CreateTime>",
	"<MsgType><![CDATA[{{.MsgType}}]]></MsgType>",
	"{{if .IsText}}",
	"<Content><![CDATA[{{.Content}}]]></Content>",
	"{{else if .IsNews}}",
	"<ArticleCount>{{.NewsLength}}</ArticleCount>",
	"<Articles>",
	"{{range .News}}",
	"<item>",
	"<Title><![CDATA[{{.Title}}]]></Title>",
	"<Description><![CDATA[{{.Desc}}]]></Description>",
	"<PicUrl><![CDATA[{{.PicUrl}}]]></PicUrl>",
	"<Url><![CDATA[{{.Url}}]]></Url>",
	"</item>",
	"{{end}}",
	"</Articles>",
	"{{else if .IsMusic}}",
	"<Music>",
	"<Title><![CDATA[{{.Title}}]]></Title>",
	"<Description><![CDATA[{{.Desc}}]]></Description>",
	"<MusicUrl><![CDATA[{{.MusicUrl}}]]></MusicUrl>",
	"<HQMusicUrl><![CDATA[{{.HQMusicUrl}}]]></HQMusicUrl>",
	"<ThumbMediaId><![CDATA[{{.ThumbMediaId}}]]></ThumbMediaId>",
	"</Music>",
	"{{else if .IsVoice}}",
	"<Voice>",
	"<MediaId><![CDATA[{{.MediaId}}]]></MediaId>",
	"</Voice>",
	"{{else if .IsImage}}",
	"<Image>",
	"<MediaId><![CDATA[{{.MediaId}}]]></MediaId>",
	"</Image>",
	"{{else if .IsVideo}}",
	"<Video>",
	"<MediaId><![CDATA[{{.MediaId}}]]></MediaId>",
	"<Title><![CDATA[{{.Title}}]]></Title>",
	"<Description><![CDATA[{{.Desc}}]]></Description>",
	"</Video>",
	"{{else if .NeedCustomService}}",
	"{{if .KfAccount}}",
	"<TransInfo>",
	"<KfAccount><![CDATA[{{.KfAccount}}]]></KfAccount>",
	"</TransInfo>",
	"{{end}}",
	"{{end}}",
	"</xml>",
}
