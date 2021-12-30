package proc

import (
	"encoding/json"
	"fmt"
	"log"
	"serverStatus/config"
	"serverStatus/tool"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// @description	 获取服务器基本信息
// @param     conn  			*config.WsConn   "ws连接信息"
// @return    void
func GetServerInfo(conn *config.WsConn) {
Exit:
	for {
		select {
		case <-*conn.Down:
			break Exit
		case <-time.After(time.Duration(1) * time.Second):
			memInfo, _ := mem.VirtualMemory()                //内存信息
			cpuPercent, _ := cpu.Percent(time.Second, false) //cpu占用率
			cpuPhysicalCores, _ := cpu.Counts(false)         //cpu物理核心
			cpuLogicalCores, _ := cpu.Counts(true)           //cpu逻辑数量
			cpuInfo, cpuErr := cpu.Info()                    //cpu详细信息
			loadInfo, _ := load.Avg()                        //负载信息
			swapInfo, _ := mem.SwapMemory()                  //SWAP
			hostInfo, _ := host.Info()                       //host信息

			var cpuModalName string
			log.Println(cpuErr)
			if cpuErr == nil {
				cpuModalName = cpuInfo[0].ModelName
			} else {
				cpuModalName = ""
			}
			log.Println(cpuModalName)

			serverBaseInfo, _ := json.Marshal(config.ServerInfo{
				Type: "server_info",
				Data: config.ServerInfoData{
					Cpu: config.CpuData{
						Percent:       tool.Decimal(cpuPercent[0], 2),
						PhysicalCores: cpuPhysicalCores,
						LogicalCores:  cpuLogicalCores,
						ModalName:     cpuModalName,
					},
					Memory: config.MemData{ //单位 兆字节
						Percent: tool.Decimal(memInfo.UsedPercent, 2),
						Total:   tool.MemTrans(memInfo.Total, 6),
						Free:    tool.MemTrans(memInfo.Free, 5),
						Used:    tool.MemTrans(memInfo.Used, 6),
					},
					Load: config.LoadData{
						M1:  tool.Decimal(loadInfo.Load1, 2),
						M5:  tool.Decimal(loadInfo.Load5, 2),
						M15: tool.Decimal(loadInfo.Load15, 2),
					},
					Swap: config.SwapData{
						Percent: tool.Decimal(swapInfo.UsedPercent, 2),
						Total:   tool.MemTrans(swapInfo.Total, 6),
						Free:    tool.MemTrans(swapInfo.Free, 5),
						Used:    tool.MemTrans(swapInfo.Used, 6),
					},
					Host: config.HostData{
						HostName:             hostInfo.Hostname,
						UpTime:               hostInfo.Uptime,
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

// 获取服务器状态信息
func GetSaerverInfo() {

	netInfo, _ := net.Interfaces()
	fmt.Println("netInfo", netInfo)
}
