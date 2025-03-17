// Package core implements the main business logic for the Pong0 application.
// It coordinates the data retrieval, parsing, and processing workflow for IP information.
package core

import (
	"fmt"
	"log"
	"time"

	"ping0/internal/client"
	"ping0/internal/constants"
	"ping0/internal/models"
	"ping0/internal/parser"
)

// ProcessIPInfo 处理获取IP信息的完整流程
// 该函数协调整个IP信息检索和解析过程的工作流程：
// 1. 获取初始页面并提取关键参数
// 2. 生成必要的访问密钥
// 3. 获取并解析包含IP信息的最终页面
//
// 参数:
//   - queryIP: 要查询的IP地址，如果为空则查询当前IP
//
// 返回:
//   - *models.IPInfo: 包含IP详细信息的结构体
//   - error: 如果过程中出现错误则返回对应错误信息
func ProcessIPInfo(queryIP string) (*models.IPInfo, error) {
	// 设置当前查询的IP
	constants.QueryIP = queryIP

	// 记录开始时间，用于性能分析
	startTime := time.Now()
	if constants.Verbose {
		log.Printf("开始查询IP信息: %s", queryIP)
	}

	// 步骤1: 获取初始页面，提取x1值和JavaScript路径
	stepStartTime := time.Now()
	x1Value, jsPath, err := client.GetInitialPage()
	if err != nil {
		return nil, fmt.Errorf("Step 1 失败: %w", err)
	}
	if constants.Verbose {
		log.Printf("成功获取x1值: %s", x1Value)
		log.Printf("JS路径: %s", jsPath)
		log.Printf("Step 1 完成，耗时: %s", time.Since(stepStartTime))
	}

	// 步骤2: 生成访问密钥并获取包含IP信息的最终页面
	stepStartTime = time.Now()
	key, err := parser.GenerateKey(jsPath, x1Value)
	if err != nil {
		return nil, fmt.Errorf("Step 2 失败: %w", err)
	}
	if constants.Verbose {
		log.Printf("成功生成key: %s", key)
	}

	finalHtml, err := client.GetFinalPage(key)
	if err != nil {
		return nil, fmt.Errorf("Step 2 失败: %w", err)
	}
	if constants.Verbose {
		log.Printf("成功获取最终页面，长度: %d", len(finalHtml))
		log.Printf("Step 2 完成，耗时: %s", time.Since(stepStartTime))
	}

	// 步骤3: 解析HTML获取IP信息
	stepStartTime = time.Now()
	ipInfo, err := parser.ParseIPInfo(finalHtml)
	if err != nil {
		if constants.Verbose {
			log.Printf("解析IP信息失败: %v", err)
		}
		return nil, fmt.Errorf("Step 3 失败: %v", err)
	}
	if constants.Verbose {
		log.Printf("解析IP信息完成，耗时: %s", time.Since(stepStartTime))
		log.Printf("总耗时: %s", time.Since(startTime))
	}

	// 清除本次查询状态，为下一次查询准备
	constants.QueryIP = ""
	constants.ManualX1Value = ""

	return ipInfo, nil
}
