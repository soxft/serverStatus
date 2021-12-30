package config

// 获取服务器基本信息
type ServerInfo struct {
	Type string         `json:"type"`
	Data ServerInfoData `json:"data"`
}

type ServerInfoData struct {
	Cpu    CpuData
	Load   LoadData
	Memory MemData
	Swap   SwapData
	Host   HostData
}

type CpuData struct {
	Percent       float64
	LogicalCores  int
	PhysicalCores int
	ModalName     string
}

type LoadData struct {
	M1  float64
	M5  float64
	M15 float64
}

type MemData struct {
	Percent float64
	Total   float64
	Free    float64
	Used    float64
}

type SwapData struct {
	Percent float64
	Total   float64
	Free    float64
	Used    float64
}

type HostData struct {
	HostName             string
	UpTime               uint64
	BootTime             uint64
	Procs                uint64
	Os                   string
	Platform             string
	PlatformFamily       string
	PlatformVersion      string
	KernelArch           string
	KernelVersion        string
	VirtualizationRole   string
	VirtualizationSystem string
}

/*
host {"hostname":"xcsoftsMBP","uptime":687914,"bootTime":1640163711,"procs":543,"os":"darwin","platform":"darwin","platformFamily":"Standalone Workstation","platformVersion":"12.1","kernelVersion":"21.2.0","kernelArch":"arm64","virtualizationSystem":"","virtualizationRole":"","hostid":"587318cb-673c-59d4-a062-31388a36e0ae"}
*/
