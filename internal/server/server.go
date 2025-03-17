// Package server implements the HTTP API server functionality for the Pong0 application.
// It provides endpoints for querying IP information and handles authentication,
// request validation, and response formatting.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"ping0/internal/constants"
	"ping0/internal/core"
)

// StartServer 启动HTTP API服务器
// 该函数配置并启动一个HTTP服务器，提供IP信息查询API。
// 它设置路由处理器、超时配置，并监听指定端口。
//
// 返回:
//   - error: 如果服务器启动失败则返回相应错误
func StartServer() error {
	// 设置服务器地址
	serverAddr := fmt.Sprintf(":%s", constants.APIPort)

	// 端口检测
	if !isPortAvailable(constants.APIPort) {
		return fmt.Errorf("端口 %s 已被占用，请使用 -p 参数指定其他端口", constants.APIPort)
	}

	// 设置路由
	http.HandleFunc("/query", handleIPQuery)

	// 打印启动信息
	fmt.Printf("Pong0 v%s 服务器模式已启动，监听端口 %s\n", constants.Version, constants.APIPort)

	if constants.APIKey != "" && constants.Verbose {
		fmt.Println("已启用API密钥验证")
	}

	fmt.Println("服务器已准备就绪，按Ctrl+C停止服务...")

	// 添加超时设置
	server := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("服务器启动失败: %v", err)
	}

	return nil
}

// handleIPQuery 处理IP查询请求
func handleIPQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 设置CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 处理OPTIONS请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 仅接受POST或GET请求
	if r.Method != "POST" && r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error":    "仅支持POST和GET请求",
			"princess": "https://linux.do/u/amna",
		})
		return
	}

	// 检查API密钥（如果配置了的话）
	if constants.APIKey != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") || authHeader[7:] != constants.APIKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":    "未授权：无效或缺失的API密钥",
				"princess": "https://linux.do/u/amna",
			})
			return
		}
	}

	var ipToQuery string

	// 处理POST请求
	if r.Method == "POST" {
		// 检查内容类型
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// 处理JSON格式请求
			var requestBody map[string]string
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error":    "无法解析请求体：" + err.Error(),
					"princess": "https://linux.do/u/amna",
				})
				return
			}
			ipToQuery = requestBody["ip"]
		} else {
			// 处理表单格式请求
			r.ParseForm()
			ipToQuery = r.FormValue("ip")
		}
	} else if r.Method == "GET" {
		// 处理GET请求
		ipToQuery = r.URL.Query().Get("ip")
	}

	// 记录处理请求
	if constants.Verbose {
		if ipToQuery == "" {
			log.Printf("处理查询：当前IP")
		} else {
			log.Printf("处理IP查询: %s", ipToQuery)
		}
	}

	// 执行IP查询，确保传递IP参数
	ipInfo, err := core.ProcessIPInfo(ipToQuery)
	if err != nil {
		if constants.Verbose {
			log.Printf("查询失败: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":    err.Error(),
			"princess": "https://linux.do/u/amna",
		})
		return
	}

	// 返回结果
	w.WriteHeader(http.StatusOK)
	// 确保IPInfo结构体有Princess字段
	if ipInfo.Princess == "" {
		ipInfo.Princess = "https://linux.do/u/amna"
	}
	json.NewEncoder(w).Encode(ipInfo)
}

// isPortAvailable 检查端口是否可用
func isPortAvailable(port string) bool {
	// 尝试监听指定端口，与服务器相同的地址
	addr := fmt.Sprintf(":%s", port)
	server, err := net.Listen("tcp", addr)

	// 如果有错误，说明端口不可用
	if err != nil {
		// 输出详细的错误信息
		if constants.Verbose {
			log.Printf("端口检测失败: %v", err)
		}
		return false
	}

	// 成功监听，关闭并返回true
	server.Close()
	return true
}

// getClientIP 获取客户端IP地址
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For获取
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// 尝试从X-Real-IP获取
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 如果上述都没有，使用RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// 如果无法解析，直接返回RemoteAddr
		return r.RemoteAddr
	}

	return ip
}
