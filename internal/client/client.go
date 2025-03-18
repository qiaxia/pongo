// Package client implements HTTP client functionality for the Pong0 application.
// It handles all communication with the Ping0.cc service, including session management,
// request construction, and response handling.
package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"ping0/internal/constants"
	"ping0/internal/parser"

	"github.com/PuerkitoBio/goquery"
)

// 全局HTTP客户端实例，用于在整个应用程序中复用
var httpClient *http.Client

// init 初始化HTTP客户端，配置cookie存储和超时设置
func init() {
	// 创建cookie jar以管理会话cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化全局HTTP客户端，设置cookie jar和超时时间
	httpClient = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
}

// GetInitialPage 获取初始页面并提取关键参数
// 该函数向Ping0.cc发送初始请求，并从响应中提取x1参数、difficulty参数和JavaScript路径，
// 这些参数对于后续请求是必需的。
//
// 返回:
//   - string: 提取的x1值，用于生成访问密钥
//   - string: 提取的difficulty值，用于生成访问密钥
//   - string: JavaScript文件路径，用于解析生成密钥的算法
//   - error: 如果请求失败或解析失败则返回相应错误
func GetInitialPage() (string, string, string, error) {
	// 如果在服务器模式下，每次创建新的HTTP客户端，避免会话状态问题
	if constants.ServerMode {
		resetHTTPClient()
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建初始请求
	req, err := http.NewRequestWithContext(ctx, "GET", constants.BaseURL, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", constants.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	if constants.Verbose {
		log.Printf("请求初始页面: %s", constants.BaseURL)
		log.Printf("请求头:")
		for k, v := range req.Header {
			log.Printf("- %s: %s", k, v)
		}
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if constants.Verbose {
		log.Printf("响应状态码: %d", resp.StatusCode)
		log.Printf("响应头:")
		for k, v := range resp.Header {
			log.Printf("- %s: %s", k, v)
		}
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	if constants.Verbose {
		log.Printf("响应内容长度: %d", len(body))
	}

	// 如果提供了手动x1值，直接返回
	if constants.ManualX1Value != "" {
		if constants.Verbose {
			log.Printf("使用手动指定的x1值: %s\n", constants.ManualX1Value)
		}

		// 获取手动指定的difficulty值或使用默认值（x1的前3位）
		difficultyValue := constants.ManualDiffValue
		if difficultyValue == "" {
			difficultyValue = constants.ManualX1Value[:3]
			if constants.Verbose {
				log.Printf("未指定difficulty值，使用x1值的前3个字符作为默认值: %s\n", difficultyValue)
			}
		} else {
			if constants.Verbose {
				log.Printf("使用手动指定的difficulty值: %s\n", difficultyValue)
			}
		}

		return constants.ManualX1Value, difficultyValue, "/js/main.js", nil
	}

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", "", "", fmt.Errorf("解析HTML失败: %w", err)
	}

	// 提取x1值和difficulty值
	var x1Value, difficultyValue string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if strings.Contains(content, "window.x1") {
			// 使用字符串函数提取x1值
			x1Start := strings.Index(content, "window.x1 = '") + 13
			if x1Start == 12 { // not found
				x1Start = strings.Index(content, `window.x1 = "`) + 13
			}
			if x1Start > 12 { // found
				x1End := strings.Index(content[x1Start:], "'")
				if x1End == -1 {
					x1End = strings.Index(content[x1Start:], `"`)
				}
				if x1End > 0 {
					x1Value = content[x1Start : x1Start+x1End]
					if constants.Verbose {
						log.Printf("找到x1值: %s", x1Value)
					}
				}
			}
		}
		if strings.Contains(content, "window.difficulty") {
			// 使用字符串函数提取difficulty值
			diffStart := strings.Index(content, "window.difficulty = '") + 21
			if diffStart == 20 { // not found
				diffStart = strings.Index(content, `window.difficulty = "`) + 21
			}
			if diffStart > 20 { // found
				diffEnd := strings.Index(content[diffStart:], "'")
				if diffEnd == -1 {
					diffEnd = strings.Index(content[diffStart:], `"`)
				}
				if diffEnd > 0 {
					difficultyValue = content[diffStart : diffStart+diffEnd]
					if constants.Verbose {
						log.Printf("找到difficulty值: %s", difficultyValue)
					}
				}
			}
		}
	})

	if x1Value == "" {
		if constants.Verbose {
			// 打印响应内容的前200个字符作为预览
			preview := string(body)
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			log.Printf("无法找到x1值，响应内容预览: %s", preview)
		}
		return "", "", "", fmt.Errorf("未找到x1值")
	}

	if difficultyValue == "" {
		if constants.Verbose {
			log.Printf("未找到difficulty值，使用x1值的前3个字符作为默认值")
		}
		// 使用x1值的前3个字符作为默认difficulty值
		if len(x1Value) >= 3 {
			difficultyValue = x1Value[:3]
		} else {
			return "", "", "", fmt.Errorf("无法设置默认difficulty值")
		}
	}

	// 查找js路径
	jsPath := "/js/main.js"
	doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && strings.Contains(src, "main.js") {
			jsPath = src
			if constants.Verbose {
				log.Printf("找到JS路径: %s", jsPath)
			}
		}
	})

	if jsPath == "" {
		if constants.Verbose {
			log.Printf("使用默认的JS路径: /js/main.js")
		}
		jsPath = "/js/main.js"
	}

	return x1Value, difficultyValue, jsPath, nil
}

// GetFinalPage 获取最终页面
// 该函数使用生成的js1key和pow值作为cookie发送请求，
// 获取包含IP信息的最终页面。
//
// 参数:
//   - keys: 包含js1key和pow值的结构体
//
// 返回:
//   - string: 获取的HTML内容
//   - error: 如果请求失败则返回相应错误
func GetFinalPage(keys *parser.Keys) (string, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建请求URL
	reqURL := constants.BaseURL
	if constants.QueryIP != "" {
		// 如果指定了IP，使用/ip/路径
		reqURL = fmt.Sprintf("%s/ip/%s", constants.BaseURL, constants.QueryIP)
		if constants.Verbose {
			log.Printf("使用特定IP查询URL: %s", reqURL)
		}
	} else {
		// 未指定IP，直接使用基础URL
		if constants.Verbose {
			log.Printf("使用当前IP查询URL: %s", reqURL)
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", constants.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", constants.BaseURL)

	if constants.Verbose {
		log.Printf("请求头:")
		for k, v := range req.Header {
			log.Printf("- %s: %s", k, v)
		}
	}

	// 设置cookie：同时设置js1key和pow
	u, _ := url.Parse(constants.BaseURL)
	httpClient.Jar.SetCookies(u, []*http.Cookie{
		{
			Name:  "js1key",
			Value: keys.Js1key,
		},
		{
			Name:  "pow",
			Value: keys.Pow,
		},
	})

	if constants.Verbose {
		log.Printf("设置Cookie: js1key=%s, pow=%s", keys.Js1key, keys.Pow)
		cookies := httpClient.Jar.Cookies(u)
		log.Printf("当前所有Cookie:")
		for _, cookie := range cookies {
			log.Printf("- %s=%s", cookie.Name, cookie.Value)
		}
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if constants.Verbose {
		log.Printf("响应状态码: %d", resp.StatusCode)
		log.Printf("响应头:")
		for k, v := range resp.Header {
			log.Printf("- %s: %s", k, v)
		}
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if constants.Verbose {
		log.Printf("响应内容长度: %d", len(body))
		if len(body) > 0 {
			// 打印前100个字符作为预览
			preview := string(body)
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			log.Printf("响应内容预览: %s", preview)
		}
	}

	return string(body), nil
}

// extractX1Value 从HTML中提取x1值
func extractX1Value(html string) string {
	// 查找x1值
	start := strings.Index(html, `var x1 = "`)
	if start == -1 {
		return ""
	}
	start += 10 // len(`var x1 = "`)
	end := strings.Index(html[start:], `"`)
	if end == -1 {
		return ""
	}
	return html[start : start+end]
}

// extractJSPath 从HTML中提取js路径
func extractJSPath(html string) string {
	// 查找js路径
	start := strings.Index(html, `src="/js/main.js?v=`)
	if start == -1 {
		return ""
	}
	start += 5 // len(`src="`)
	end := strings.Index(html[start:], `"`)
	if end == -1 {
		return ""
	}
	return html[start : start+end]
}

// 重置HTTP客户端，用于在API模式下每次请求前调用
func resetHTTPClient() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Printf("创建新的cookie jar失败: %v", err)
		return
	}

	httpClient = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	if constants.Verbose {
		log.Printf("已重置HTTP客户端")
	}
}
