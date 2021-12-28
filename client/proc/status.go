package proc

import (
	"encoding/json"
	"fmt"
	"log"
	"serverStatus-client/config"
	"serverStatus-client/tool"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
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
		case <-time.After(time.Duration(2) * time.Second):
			memInfo, _ := mem.VirtualMemory() //内存信息
			memUsedPer := tool.Decimal(memInfo.UsedPercent)
			//memAll := memInfo.Total
			//memUsed := memInfo.Used
			//memFree := memInfo.Free

			//log.Println(memUsedPer, memAll, memFree, memUsed)

			percent, _ := cpu.Percent(time.Second, false)
			cpuPercent := tool.Decimal(percent[0])
			//log.Println("cpu", cpuPercent)

			serverBaseInfo, _ := json.Marshal(config.ServerInfo{
				Type: "server_info",
				Data: config.ServerInfoData{
					CpuPercent: cpuPercent,
					MemPercent: memUsedPer,
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
func GetSaerverInfo(conn *config.WsConn) {
	percent, _ := cpu.Percent(time.Second, false)
	log.Println("cpu", percent)
	memInfo, _ := mem.VirtualMemory()
	log.Println("mem", memInfo)
	info2, _ := mem.SwapMemory() //SWAP
	log.Println("swap", info2)
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	log.Println("disk", diskInfo)
	log.Println("ddisk", parts)
	info, _ := host.Info()
	log.Println("host", info)
	cp, _ := cpu.Info() //总体信息
	log.Println("cpuInfo", cp)
	c, _ := cpu.Counts(true) //cpu逻辑数量
	fmt.Println(c)           //4
	c, _ = cpu.Counts(false) //cpu物理核心
	fmt.Println(c)           //如果是2说明是双核超线程, 如果是4则是4核非超线程
	netInfo, _ := net.IOCounters(false)
	fmt.Println("netInfo", netInfo)
}
