package config

// 获取服务器基本信息
type ServerInfo struct {
	Type string         `json:"type"`
	Data ServerInfoData `json:"data"`
}

type ServerInfoData struct {
	CpuPercent float64 `json:"cpu_percent"` // cpu占用百分比
	MemPercent float64 `json:"mem_percent"` // 内存占用百分比
}
