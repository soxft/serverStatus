package proc

import (
	"encoding/json"
	"log"
	"serverStatus/config"
	"serverStatus/tool"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// @description	 获取服务器基本信息
// @param     conn  			*config.WsConn   		"ws连接信息"
// @param     baseInfo  		*config.ServerBaseInfo 	"服务器基本信息"
// @param     conn  			 int64   				"发送间隔"
// @return    void
func GetServerInfo(conn *config.WsConn, baseInfo *config.ServerBaseInfo, duration int64) {
Exit:
	for {
		select {
		case <-*conn.Down:
			break Exit
		case <-time.After(time.Duration(duration) * time.Millisecond):
			memInfo, _ := mem.VirtualMemory()                //内存信息
			cpuPercent, _ := cpu.Percent(time.Second, false) //cpu占用率
			loadInfo, _ := load.Avg()                        //负载信息
			swapInfo, _ := mem.SwapMemory()                  //SWAP
			hostInfo, _ := host.Info()                       //系统信息

			serverBaseInfo, _ := json.Marshal(config.ServerInfo{
				Type: "server_info",
				Data: config.ServerInfoData{
					Cpu: config.CpuData{
						Percent:       tool.Decimal(cpuPercent[0], 2),
						PhysicalCores: baseInfo.PhysicalCores,
						LogicalCores:  baseInfo.LogicalCores,
						ModalName:     baseInfo.ModalName,
					},
					Memory: config.MemData{
						Percent: tool.Decimal(memInfo.UsedPercent, 2),
						Total:   memInfo.Total,
						Free:    memInfo.Free,
						Used:    memInfo.Used,
					},
					Swap: config.SwapData{
						Percent: tool.Decimal(swapInfo.UsedPercent, 2),
						Total:   swapInfo.Total,
						Free:    swapInfo.Free,
						Used:    swapInfo.Used,
					},
					Load: config.LoadData{
						M1:  tool.Decimal(loadInfo.Load1, 2),
						M5:  tool.Decimal(loadInfo.Load5, 2),
						M15: tool.Decimal(loadInfo.Load15, 2),
					},
					Host: config.HostData{
						HostName:             hostInfo.Hostname,
						UpTime:               uint64(time.Now().Unix()) - hostInfo.BootTime,
						BootTime:             hostInfo.BootTime,
						Procs:                hostInfo.Procs,
						Os:                   hostInfo.OS,
						Platform:             hostInfo.Platform,
						PlatformFamily:       hostInfo.PlatformFamily,
						PlatformVersion:      hostInfo.PlatformVersion,
						KernelArch:           hostInfo.KernelArch,
						KernelVersion:        hostInfo.KernelVersion,
						VirtualizationRole:   hostInfo.VirtualizationRole,
						VirtualizationSystem: hostInfo.VirtualizationSystem,
					},
					Time: time.Now().Unix(),
				},
			})

			conn.Lock.Lock() //锁
			err := conn.Conn.WriteMessage(websocket.TextMessage, serverBaseInfo)
			conn.Lock.Unlock()

			if err != nil {
				log.Println("Error during send server base info", err)
			}
		}
	}
}

// 获取一些基本信息
func GetBaseInfo() *config.ServerBaseInfo {
	cpuInfo, cpuErr := cpu.Info()            //cpu详细信息
	cpuPhysicalCores, _ := cpu.Counts(false) //cpu物理核心
	cpuLogicalCores, _ := cpu.Counts(true)   //cpu逻辑数量

	var cpuModalName string
	if cpuErr == nil {
		cpuModalName = cpuInfo[0].ModelName
	} else {
		cpuModalName = ""
	}
	return &config.ServerBaseInfo{
		LogicalCores:  cpuLogicalCores,
		PhysicalCores: cpuPhysicalCores,
		ModalName:     cpuModalName,
	}
}
