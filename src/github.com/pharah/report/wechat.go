package report

import (
    "fmt"
    "bytes"
    "io/ioutil"  
    "net/http"
    "net/url"
    "encoding/json"
    "strings"
    "github.com/jmoiron/jsonq"

)

const (
    MESSAGE_TITLE = "【警告】%s"
    MESSAGE_DESCRIPTION = "<div class=\"gray\">%s</div> <div class=\"normal\"></div><div class=\"highlight\">%s</div>"
)

type Wechat struct {
    WechatConf WechatConf
    Message Message
}

type WechatConf struct {
    report_to string
    token_get_url string
    message_send_url string
    corpid string
    corpsecret string
    agentid int
}

type Message struct {
    title string
    time string
    body string
}


func (w *Wechat) Init(report_to string,
                      token_get_url string,
                      message_send_url string,
                      corpid string,
                      corpsecret string,
                      agentid int,
                      title string,
                      time string,
                      body string){
    w.WechatConf.report_to = report_to
    w.WechatConf.token_get_url = token_get_url 
    w.WechatConf.message_send_url = message_send_url
    w.WechatConf.corpid = corpid
    w.WechatConf.corpsecret = corpsecret
    w.WechatConf.agentid = agentid  
    w.Message.title = title
    w.Message.time = time
    w.Message.body = body
}

func (w *Wechat) GetToken() string {
    conf := w.WechatConf
    u, _ := url.Parse(conf.token_get_url)
    q := u.Query()
    q.Set("corpid", conf.corpid)
    q.Set("corpsecret", conf.corpsecret)
    u.RawQuery = q.Encode()
    res, err := http.Get(u.String())
    if err != nil { 
    }
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
    }
    data := map[string]interface{}{}
    dec := json.NewDecoder(strings.NewReader(string(body)))
    dec.Decode(&data)
    jq := jsonq.NewQuery(data)
    token, err := jq.String("access_token")
    return token
}

func (w *Wechat) ReportToWechat() {
    conf := w.WechatConf
    message := w.Message

    token := w.GetToken()
    u, _ := url.Parse(conf.message_send_url)
    q := u.Query()
    q.Set("access_token", token)
    u.RawQuery = q.Encode()

    // 嵌套struct需要通过map转json
    data, textcard := make(map[string]interface{}),make(map[string]interface{})
    textcard["title"] = fmt.Sprintf(MESSAGE_TITLE, message.title)
    textcard["description"] = fmt.Sprintf(MESSAGE_DESCRIPTION, message.time, message.body)
    textcard["url"] = "www.astute-tec.com"
    data["touser"] = conf.report_to
    data["msgtype"] = "textcard"
    data["safe"] = 0
    data["agentid"] = conf.agentid
    data["textcard"] = textcard
    buf, _ := json.Marshal(data)
    body := bytes.NewBuffer([]byte(buf))
    _, err := http.Post(u.String(), "application/json;charset=utf-8", body)
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Println(res)
    return
}