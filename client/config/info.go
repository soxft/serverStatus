package config

// 获取服务器基本信息
type ServerInfo struct {
	Type string         `json:"type"`
	Data ServerInfoData `json:"data"`
}

type ServerInfoData struct {
	CpuPercent float64  `json:"cpu_percent"`
	Load       LoadData `json:"load"`
	Memory     MemData  `json:"memory"`
}

type MemData struct {
	Percent float64 `json:"percent"`
	Total   float64 `json:"total"`
	Free    float64 `json:"free"`
	Used    float64 `json:"used"`
}

type LoadData struct {
	M1  float64 `json:"1"`
	M5  float64 `json:"5"`
	M15 float64 `json:"15"`
}
