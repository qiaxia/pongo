// Package constants defines global configuration variables and constants
// used throughout the Pong0 application. This includes runtime settings,
// command-line options, and HTTP-related constants.
package constants

// 全局配置变量，存储应用程序的运行时状态和配置
var (
	// 命令行参数和运行时配置
	Verbose       bool   // 是否显示详细日志信息
	ManualX1Value string // 手动指定的x1值，用于调试或绕过自动获取
	QueryIP       string // 要查询的IP地址，为空时查询当前IP
	ServerMode    bool   // 是否启动HTTP服务器模式
	APIPort       string // HTTP服务器监听的端口号
	APIKey        string // API验证密钥，用于限制API访问
	Version       string // 应用程序版本号
	UpdateDate    string // 最近更新日期

	// HTTP服务相关常量
	BaseURL   = "https://ping0.cc"               // Ping0服务的基础URL
	UserAgent = "Mozilla/5.0 Pong0/1.0.0 Golang" // HTTP请求的User-Agent头
)
