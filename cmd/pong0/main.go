// Package main provides the entry point for the Pong0 application.
// Pong0 is a tool for retrieving and analyzing IP information from the Ping0.cc service.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"ping0/internal/constants"
	"ping0/internal/core"
	"ping0/internal/server"
)

// 命令行选项定义
var (
	ip              string // 要查询的IP地址
	port            string // API服务器端口
	apiKey          string // API访问密钥
	serverMode      bool   // 是否启动API服务器模式
	verbose         bool   // 详细输出模式
	manualX1Value   string // 手动指定x1值
	manualDiffValue string // 手动指定difficulty值
	showVersion     bool   // 显示版本信息
)

// 构建信息，在编译时通过-ldflags注入
var (
	Version   = "dev"     // 程序版本号
	buildDate = "unknown" // 构建日期
)

// init 函数初始化应用程序设置并解析命令行参数
func init() {
	// 设置程序版本
	constants.Version = Version

	// 注册命令行选项
	flag.StringVar(&ip, "ip", "", "要查询的IP地址，不提供则查询本机IP")
	flag.StringVar(&port, "p", "8080", "API服务器监听端口")
	flag.StringVar(&apiKey, "k", "", "API访问密钥")
	flag.StringVar(&manualX1Value, "x1", "", "手动指定x1值")
	flag.StringVar(&manualDiffValue, "diff", "", "手动指定difficulty值")
	flag.BoolVar(&serverMode, "c", false, "启动API服务器模式")
	flag.BoolVar(&verbose, "all", false, "输出详细日志")
	flag.BoolVar(&showVersion, "v", false, "显示版本信息")

	// 解析命令行参数
	flag.Parse()
}

// main 函数是程序的入口点，处理命令行参数并执行相应功能
func main() {
	// 检查是否显示版本信息
	if showVersion {
		fmt.Printf("Pong0 %s (构建日期: %s)\n", constants.Version, buildDate)
		return
	}

	// 验证参数组合是否合法
	validateCommandLineOptions()

	// 将命令行参数应用到全局配置
	applyCommandLineOptions()

	// 根据运行模式执行不同功能
	if constants.ServerMode {
		runServerMode()
	} else {
		runQueryMode()
	}
}

// validateCommandLineOptions 验证命令行参数组合的有效性
func validateCommandLineOptions() {
	// 检查 -c 和 -all 参数是否同时使用
	if serverMode && verbose {
		fmt.Println("错误: -c 和 -all 参数不能同时使用")
		fmt.Println("用法示例:")
		fmt.Println("  服务器模式: pong0 -c -p 8080 -k your_api_key")
		fmt.Println("  查询模式: pong0 -ip 1.1.1.1")
		os.Exit(1)
	}

	// 检查 -p 和 -k 参数是否在没有 -c 参数的情况下使用
	if !serverMode && (port != "8080" || apiKey != "") {
		fmt.Println("错误: -p 和 -k 参数只能在服务器模式(-c)下使用")
		fmt.Println("用法示例:")
		fmt.Println("  服务器模式: pong0 -c -p 8080 -k your_api_key")
		fmt.Println("  查询模式: pong0 -ip 1.1.1.1")
		os.Exit(1)
	}
}

// applyCommandLineOptions 将命令行参数应用到全局配置
func applyCommandLineOptions() {
	if verbose {
		constants.Verbose = true
	}

	if serverMode {
		constants.ServerMode = true
	}

	if manualX1Value != "" {
		constants.ManualX1Value = manualX1Value
	}

	if manualDiffValue != "" {
		constants.ManualDiffValue = manualDiffValue
	}

	if apiKey != "" {
		constants.APIKey = apiKey
	}

	if ip != "" {
		constants.QueryIP = ip
	}
}

// runServerMode 在服务器模式下运行程序
func runServerMode() {
	// 设置服务器模式
	constants.ServerMode = true

	// 设置API服务器端口
	if port != "" {
		constants.APIPort = port
	}

	if constants.Verbose {
		fmt.Printf("启动API服务器，监听端口 %s...\n", constants.APIPort)
	}

	// 启动服务器并处理错误
	if err := server.StartServer(); err != nil {
		fmt.Printf("启动服务器失败: %v\n", err)
		os.Exit(1)
	}
}

// runQueryMode 在查询模式下运行程序
func runQueryMode() {
	// 输出详细信息头
	if constants.Verbose {
		fmt.Println("-------------------------------------")
		fmt.Println("Pong0 Pong0 Pong0")
		fmt.Println("-------------------------------------")
		if constants.QueryIP != "" {
			fmt.Printf("查询IP: %s\n", constants.QueryIP)
		} else {
			fmt.Println("查询当前IP")
		}
	}

	// 执行查询，获取IP信息
	ipInfo, err := core.ProcessIPInfo(constants.QueryIP)
	if err != nil {
		if constants.Verbose {
			fmt.Printf("获取IP信息失败: %v\n", err)
		} else {
			// 输出带Princess字段的错误信息JSON
			errorJSON := map[string]string{
				"error":    err.Error(),
				"princess": "https://linux.do/u/amna",
			}
			jsonData, _ := json.MarshalIndent(errorJSON, "", "  ")
			fmt.Println(string(jsonData))
		}
		os.Exit(1)
	}

	// 输出结果
	if constants.Verbose {
		fmt.Println("-------------------------------------")
	}

	// 确保IPInfo中有Princess字段
	if ipInfo.Princess == "" {
		ipInfo.Princess = "https://linux.do/u/amna"
	}

	// 输出JSON结果
	jsonData, _ := json.MarshalIndent(ipInfo, "", "  ")
	fmt.Println(string(jsonData))
}
