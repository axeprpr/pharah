package monitor

import (
    // "fmt"
    "time"
    "strconv"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"

)

type SysMonitor struct {
    monitors []Monitor
}

type Monitor struct{
    monitor string
    threshold int
}

func(s *SysMonitor) Init(cpu_alarm_threshold,memory_alarm_threshold,disk_alarm_threshold int){
    if cpu_alarm_threshold > 0 {
        s.monitors = append(s.monitors, Monitor{monitor:"cpu", threshold: cpu_alarm_threshold})
    }
    if memory_alarm_threshold > 0 {
        s.monitors = append(s.monitors, Monitor{monitor:"memory", threshold: memory_alarm_threshold})
    }
    if disk_alarm_threshold > 0 {
        s.monitors = append(s.monitors, Monitor{monitor:"disk", threshold: disk_alarm_threshold})
    }
}

func(s *SysMonitor) SystemAlarm()(string,string,string) {
    title := ""
    body := ""
    time := time.Now().Format("2006-01-02 15:04:05")
    Info, _ := host.Info()
    title = Info.Hostname + "-" + Info.OS

    for _, item := range s.monitors {
        if item.monitor == "cpu" {
            cpu_usage, _ := cpu.Percent(0, false)
            if int(cpu_usage[0]) > item.threshold {
                body += "当前CPU使用率：" + strconv.Itoa(int(cpu_usage[0])) + "%\n"

            }
        }

        if item.monitor == "memory" {
            memory_usage, _ := mem.VirtualMemory()
            if int(memory_usage.UsedPercent) > item.threshold {
                body += "当前内存使用率：" + strconv.Itoa(int(memory_usage.UsedPercent)) + "%;"
                body += "总内存：" + strconv.Itoa(int(memory_usage.Total/1024/1024/1024)) + "G\n"
            }
        }

        if item.monitor == "disk" {
            disk_partitions, _ := disk.Partitions(false)
            for index := range disk_partitions {
                disk_usage, _ := disk.Usage(disk_partitions[index].Mountpoint)
                if int(disk_usage.UsedPercent) > item.threshold {
                    body += "分区[" + disk_usage.Path + "]使用率：" +strconv.Itoa(int(disk_usage.UsedPercent)) + "%，总大小：" + strconv.Itoa(int(disk_usage.Total/1024/1024/1024)) + "G\n"
                }
            }
        }
    }

    if body != "" {
        return title, time, body
    }
    return "","",""     
}
