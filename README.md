# Pong0 - Ping0.cc网站数据获取工具

## 功能特点

1. 自动化获取和处理IP信息
2. 支持查询当前IP或指定IP地址
3. 内置JavaScript加密算法，无需下载外部资源
4. 高效的性能，通常在2-3秒内完成操作
5. 简洁的JSON输出格式
6. 详细的执行过程和性能统计（可选）
7. 无需浏览器或其他外部依赖
8. 支持提取多值字段（如IP类型、ASN类型、组织类型等）
9. 支持API服务器模式，提供HTTP接口

## 安装方法

### 环境

- 安装Go语言环境 (1.18或更高版本)

### 构建程序

```bash
# 克隆仓库
git clone https://github.com/yourusername/pongo.git
cd pongo

# 安装依赖
go mod tidy

# 构建程序
go build -o pongo.exe ./cmd/pongo
```

### 使用打包脚本

项目提供了跨平台打包脚本，可以一键构建多个平台的可执行文件：

```bash
# 在Windows上使用PowerShell脚本
.\scripts\build.ps1

# 在Linux/macOS上使用Shell脚本
./scripts/build.sh
```

打包脚本会在`dist`目录下生成各平台的可执行文件和ZIP压缩包。

## 使用方法

### 基本使用

```bash
# 获取当前IP信息，以JSON格式输出
.\pongo.exe

# 查询指定IP地址
.\pongo.exe -ip 1.1.1.1
```

### 参数选项

```bash
# 显示所有执行详情和性能统计
.\pong0.exe -all

# 查询指定IP并显示详细信息
.\pong0.exe -ip 8.8.8.8 -all

# 显示帮助信息
.\pong0.exe -h

# 显示版本信息
.\pong0.exe -v

# 使用自定义x1值(用于调试)
.\pong0.exe -x1 YOUR_X1_VALUE
```

### API服务器模式

```bash
# 启动API服务器模式（默认端口8080）
.\pongo.exe -c

# 指定自定义端口
./pongo.exe -c -p 3000

# 启用API密钥验证
.\pong0.exe -c -k YOUR_SECRET_KEY
```

在API服务器模式下：

- **GET请求：**
  - 查询当前IP：`GET http://localhost:8080/query`
  - 查询指定IP：`GET http://localhost:8080/query?ip=1.1.1.1`

- **POST请求：**
  - 支持JSON格式：`POST http://localhost:8080/query` 请求体: `{"ip": "1.1.1.1"}`
  - 支持表单格式：`POST http://localhost:8080/query` 表单参数: `ip=1.1.1.1`
  - 查询当前IP时，可以发送空请求体或省略IP参数

- 如果启用了API密钥验证，需要添加请求头：`Authorization: Bearer YOUR_SECRET_KEY`
- 如果指定的端口已被占用，程序会显示错误信息并退出，你可以使用 `-p` 参数指定其他可用端口

示例（使用curl）：

```bash
# GET请求 - 查询当前IP
curl http://localhost:8080/query

# GET请求 - 查询指定IP
curl http://localhost:8080/query?ip=1.1.1.1

# POST请求 - JSON格式查询指定IP
curl -X POST -H "Content-Type: application/json" -d '{"ip":"1.1.1.1"}' http://localhost:8080/query

# POST请求 - 表单格式查询指定IP
curl -X POST -d "ip=1.1.1.1" http://localhost:8080/query

# 带API密钥验证
curl -H "Authorization: Bearer YOUR_SECRET_KEY" http://localhost:8080/query
```

## 输出示例

### 标准JSON输出

```json
{
  "ip": "1.1.1.1",
  "ip_location": "美国 加州 洛杉矶",
  "asn": "AS13335",
  "asn_owner": "Cloudflare, Inc.",
  "asn_type": "IDC",
  "organization": "APNIC Research and Development",
  "org_type": "GOV",
  "longitude": "-118.24356842041",
  "latitude": "34.05286026001",
  "ip_type": "IDC机房IP; CloudFlare DNS IP",
  "risk_value": "26% 中性",
  "native_ip": "广播 IP",
  "country_flag": "us"
}
```

### 详细模式输出

详细模式(-all)会显示程序执行的每个步骤及其耗时:

```
-------------------------------------
Pong0 Pong0 Pong0
-------------------------------------
查询IP: 1.1.1.1
Received 276 bytes from initial page
Found x1Value: 3ef12496741412ab807c60c346ded5e7
JS Path: /static/js/a2296d5c180a52cf01f4b428fb97d804.js?t=1742134555
Step 1 完成，耗时: 3.5684864s
Calculated key with x1Value=3ef12496741412ab807c60c346ded5e7, URL=https://ping0.cc
Generated key: 310424
Step 2 完成，耗时: 523.7µs
查询URL: https://ping0.cc/ip/1.1.1.1
Received 49888 bytes from final page request
Step 3 完成，耗时: 2.1685182s
解析IP信息完成，耗时: 520.4µs
总耗时: 5.7385612s
-------------------------------------
{
  "ip": "1.1.1.1",
  "ip_location": "美国 加州 洛杉矶",
  "asn": "AS13335",
  "asn_owner": "Cloudflare, Inc.",
  "asn_type": "IDC",
  "organization": "APNIC Research and Development",
  "org_type": "GOV",
  "longitude": "-118.24356842041",
  "latitude": "34.05286026001",
  "ip_type": "IDC机房IP; CloudFlare DNS IP",
  "risk_value": "26% 中性",
  "native_ip": "广播 IP",
  "country_flag": "us"
}
```

## 数据字段说明

程序提取的数据字段包括：

| 字段名         | 描述                                 | 示例值                                 |
|---------------|--------------------------------------|--------------------------------------|
| ip            | IP地址                                | 1.1.1.1                              |
| ip_location   | IP地址地理位置                        | 美国 加州 洛杉矶                        |
| asn           | 自治系统编号                           | AS13335                             |
| asn_owner     | 自治系统拥有者                         | Cloudflare, Inc.                    |
| asn_type      | 自治系统类型（多值用分号分隔）            | IDC                                 |
| organization  | 组织名称                              | APNIC Research and Development      |
| org_type      | 组织类型（多值用分号分隔）                | GOV                                 |
| longitude     | 经度                                  | -118.24356842041                    |
| latitude      | 纬度                                  | 34.05286026001                       |
| ip_type       | IP类型（多值用分号分隔）                 | IDC机房IP; CloudFlare DNS IP          |
| risk_value    | 风险值                                | 26% 中性                              |
| native_ip     | 原生IP信息                            | 广播 IP                               |
| country_flag  | 国家/地区标志代码                      | us                                    |

## 技术实现

该工具实现了以下关键功能：

1. 通过HTTP请求获取初始页面并提取x1值
2. 使用内置算法生成JS1密钥
3. 使用密钥获取最终页面(当前IP或指定IP)
4. 提取并解析IP信息
5. 输出JSON格式数据
6. 提供HTTP API服务器模式

## 项目结构

项目采用标准Go模块结构：

```
pong0/
├── cmd/
│   └── pong0/           # 主程序入口
│       └── main.go      # 程序入口点
├── internal/            # 内部包（不导出）
│   ├── client/          # HTTP客户端功能
│   │   └── client.go    # HTTP请求处理
│   ├── constants/       # 常量定义
│   │   └── constants.go # 全局常量和变量
│   ├── core/            # 核心业务逻辑
│   │   └── core.go      # 主要处理流程
│   ├── models/          # 数据模型
│   │   └── models.go    # 数据结构定义
│   ├── parser/          # 解析功能
│   │   ├── parser.go    # HTML解析
│   │   └── js_engine.go # JavaScript加密实现
│   └── server/          # API服务器
│       └── server.go    # HTTP服务器实现
├── scripts/             # 构建脚本
│   ├── build.ps1        # Windows构建脚本
│   └── build.sh         # Linux/macOS构建脚本
├── dist/                # 构建输出目录
├── go.mod               # Go模块定义
├── go.sum               # 依赖校验
└── README.md            # 项目文档
```

## 构建标志说明

构建脚本使用以下构建标志：

- `-s -w`: 减小可执行文件大小
- `-X main.Version`: 设置版本号
- `-X main.buildDate`: 设置构建日期
- `-X ping0/internal/constants.UpdateDate`: 设置更新日期

## 免责声明

本工具仅用于学习和研究目的。使用者应自行承担使用风险，并遵守相关法律法规。 
