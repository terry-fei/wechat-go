package wechat

import (
	"regexp"
	"strings"
	"text/template"
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
	ReplyMsgTemplate = template.Must(template.New("replyMessage").Parse(replyMsgTemplateStr))
}

func XMLToMessage(xml []byte) (msg map[string]string) {
	msg = make(map[string]string)

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

type News struct {
	Title  string
	Desc   string
	PicUrl string
	Url    string
}

type ReplyMessage struct {
	ToUserName        string
	FromUserName      string
	CreateTime        int64
	MsgType           string
	IsText            bool
	IsNews            bool
	IsMusic           bool
	IsVoice           bool
	IsImage           bool
	IsVideo           bool
	NeedCustomService bool
	Content           string
	NewsLength        int
	NewsList          []News
	Title             string
	Desc              string
	MusicUrl          string
	HQMusicUrl        string
	MediaId           string
	KfAccount         string
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
