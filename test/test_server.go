//go:build ignore
// +build ignore

// Package main implements a simple test server for the Pong0 application.
// This server is used for development and testing purposes, allowing developers
// to simulate the behavior of the Ping0.cc service locally.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// main 函数启动测试服务器并配置请求处理
// 该函数创建一个简单的HTTP服务器，监听8081端口，并记录所有接收到的请求。
// 它将请求信息写入日志文件，便于调试和测试。
func main() {
	// 创建日志文件
	logFile, err := os.Create("test_server.log")
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}
	defer logFile.Close()

	// 设置日志输出
	logger := log.New(logFile, "", log.LstdFlags)

	// 配置测试端点处理函数
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("收到请求:", r.Method)
		logger.Println("Content-Type:", r.Header.Get("Content-Type"))

		if r.Method == "POST" {
			// 读取请求体
			body, _ := io.ReadAll(r.Body)
			logger.Println("请求体:", string(body))

			// 解析表单
			r.ParseForm()
			logger.Println("Form:", r.Form)
			logger.Println("PostForm:", r.PostForm)

			// 获取IP参数
			ip := r.FormValue("ip")
			logger.Println("FormValue(ip):", ip)
		}

		// 返回响应
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	// 启动服务器
	fmt.Println("测试服务器启动在 http://localhost:8081/test")
	fmt.Println("日志输出到 test_server.log")
	http.ListenAndServe(":8081", nil)
}
