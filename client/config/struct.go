package config

// 登录到服务器
type Login struct {
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Tag      string `json:"tag"`
	Token    string `json:"token"`
}
