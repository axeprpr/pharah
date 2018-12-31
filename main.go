package main

import (
    "fmt"
    "time"
    "os/exec"
    "encoding/json"
    "github.com/jmoiron/jsonq"
    "strings"
    "github.com/pharah/utils"
    "github.com/pharah/report"
    "github.com/pharah/monitor"  
)

const (
    CONF = "conf.ini"
)

type WechatData struct {
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

// GET INTERVAL
func get_interval_value(interval_key string) int {
    ini_parser := utils.IniParser{}
    if err := ini_parser.Load(CONF); err != nil {
        fmt.Printf("try load config file[%s] error[%s]\n", CONF, err.Error())
        return 0
    }
    interval_value := int(ini_parser.GetInt64("INTERVALS", interval_key))
    return interval_value
}

// GET MONITOR CONFIGURATION
func get_plugins_config() []string {
    ini_parser := utils.IniParser{}
    if err := ini_parser.Load(CONF); err != nil {
        fmt.Printf("try load config file[%s] error[%s]\n", CONF, err.Error())
        return nil
    }
    plugins := ini_parser.GetSectionKeys("PLUGINS")
    return plugins
}

func get_plugin_value(plugin_key string) string {
    ini_parser := utils.IniParser{}
    if err := ini_parser.Load(CONF); err != nil {
        fmt.Printf("try load config file[%s] error[%s]\n", CONF, err.Error())
        return ""
    }
    plugin_value := ini_parser.GetString("PLUGINS", plugin_key)
    return plugin_value
}

func get_system_config() (int,int,int) {
    ini_parser := utils.IniParser{}
    if err := ini_parser.Load(CONF); err != nil {
        fmt.Printf("try load config file[%s] error[%s]\n", CONF, err.Error())
        return 0,0,0
    }
    c := int(ini_parser.GetInt64("SYSTEM", "cpu_alarm_threshold"))
    m := int(ini_parser.GetInt64("SYSTEM", "memory_alarm_threshold"))
    d := int(ini_parser.GetInt64("SYSTEM", "disk_alarm_threshold"))
    
    return c,m,d
}

// GET REPORT CONFIGURATION
func get_wechat_config() *WechatData {
    data := new(WechatData)
    ini_parser := utils.IniParser{}
    if err := ini_parser.Load(CONF); err != nil {
        fmt.Printf("try load config file[%s] error[%s]\n", CONF, err.Error())
        return data
    }
    data.report_to = ini_parser.GetString("WECHAT", "report_to")
    if data.report_to == "" {
        return data
    }
    data.token_get_url = ini_parser.GetString("WECHAT", "token_get_url")
    data.message_send_url = ini_parser.GetString("WECHAT", "message_send_url")
    data.corpid = ini_parser.GetString("WECHAT", "corpid")
    data.corpsecret = ini_parser.GetString("WECHAT", "corpsecret")
    data.agentid = int(ini_parser.GetInt64("WECHAT", "agentid"))
    return data    
}

func report_to_wechat(title, time, body string) {
    data := get_wechat_config()
    if data.report_to == "" { 
        return 
    }
    w := report.Wechat{}
    w.Init(data.report_to, 
           data.token_get_url,
           data.message_send_url,
           data.corpid,
           data.corpsecret,
           data.agentid,
           title,
           time,
           body)
    w.ReportToWechat()
}

func system_monitor(system_interval int){
    if system_interval == 0 {
        return
    }
    for {
        fmt.Println("Start system monitor...")
        c,m,d := get_system_config()
        s:= monitor.SysMonitor{}
        s.Init(c,m,d)
        title, t, body := s.SystemAlarm()
        if title != "" && t != "" && body != "" {
            report_to_wechat(title, t, body)
        }
        time.Sleep(time.Duration(system_interval)*time.Second)
    }
}

func plugin_monitor(plugin string, interval int){
    if interval == 0 {
        return
    }
    for {
        fmt.Println("Start plugin monitor...")
        cmd := exec.Command(plugin)
        out, err := cmd.Output()  
        if err != nil {  
            fmt.Println(err)  
        }
        data := map[string]interface{}{}
        dec := json.NewDecoder(strings.NewReader(string(out)))
        dec.Decode(&data)
        jq := jsonq.NewQuery(data)
        title,_ := jq.String("title")
        t, _ := jq.String("time")
        body, _ := jq.String("body")
        if title != "" && t != "" && body != "" {
            report_to_wechat(title, t, body)
        }
        time.Sleep(time.Duration(interval)*time.Second)
    }
}

func start_monitor(){
    go system_monitor(get_interval_value("system"))
    plugins := get_plugins_config()
    for plugin_index := range plugins{
        go plugin_monitor(get_plugin_value(plugins[plugin_index]), get_interval_value(plugins[plugin_index]))
    }
}

func main() {
    start_monitor()
    var str string
    fmt.Scan(&str)
}