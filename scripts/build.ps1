#!/usr/bin/env pwsh
# 跨平台构建并打包Pong0工具的PowerShell脚本

Write-Host "开始构建Pong0 - Ping0.cc数据获取工具..." -ForegroundColor Green

# 确保在项目根目录执行
$scriptDir = $PSScriptRoot
$projectRoot = Split-Path -Parent $scriptDir

# 切换到项目根目录
Push-Location $projectRoot
Write-Host "工作目录: $(Get-Location)" -ForegroundColor Gray

# 确保dist目录存在
$distDir = "dist"
If (-not (Test-Path -Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
    Write-Host "创建dist目录" -ForegroundColor Gray
}

# 版本信息
$version = "1.0.0"
$buildDate = Get-Date -Format "yyyy-MM-dd"
$updateDate = Get-Date -Format "yyyy-MM-dd"

# 构建信息
Write-Host "构建版本: $version (构建日期: $buildDate, 更新日期: $updateDate)" -ForegroundColor Yellow

# 复制README.md到dist目录
if (Test-Path -Path "README.md") {
    Copy-Item -Path "README.md" -Destination "$distDir/" -Force
    Write-Host "已复制README.md到dist目录" -ForegroundColor Gray
} else {
    Write-Host "未找到README.md文件" -ForegroundColor Yellow
}

# 构建配置
$platforms = @(
    @{GOOS = "windows"; GOARCH = "amd64"; Suffix = ".exe"; FriendlyName = "Windows 64位" },
    @{GOOS = "windows"; GOARCH = "386"; Suffix = ".exe"; FriendlyName = "Windows 32位" },
    @{GOOS = "linux"; GOARCH = "amd64"; Suffix = ""; FriendlyName = "Linux 64位" },
    @{GOOS = "linux"; GOARCH = "386"; Suffix = ""; FriendlyName = "Linux 32位" },
    @{GOOS = "linux"; GOARCH = "arm64"; Suffix = ""; FriendlyName = "Linux ARM64" },
    @{GOOS = "darwin"; GOARCH = "amd64"; Suffix = ""; FriendlyName = "macOS 64位" },
    @{GOOS = "darwin"; GOARCH = "arm64"; Suffix = ""; FriendlyName = "macOS ARM64" }
)

# 主程序路径
$mainPath = "cmd/pong0"

# 为每个平台构建
foreach ($platform in $platforms) {
    $env:GOOS = $platform.GOOS
    $env:GOARCH = $platform.GOARCH
    
    $outputName = "pong0_${version}_$($platform.GOOS)_$($platform.GOARCH)$($platform.Suffix)"
    $outputPath = "$distDir/$outputName"
    
    Write-Host "正在构建: $($platform.FriendlyName) ($($platform.GOOS)/$($platform.GOARCH))..." -ForegroundColor Cyan
    
    # 检查主程序目录是否存在
    if (-not (Test-Path -Path $mainPath)) {
        Write-Host "  - 错误: 主程序目录 $mainPath 不存在" -ForegroundColor Red
        continue
    }
    
    # 构建二进制文件
    & go build -o "$outputPath" -ldflags "-s -w -X main.Version=$version -X main.buildDate=$buildDate -X ping0/internal/constants.UpdateDate=$updateDate" ./$mainPath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  - 构建成功: $outputName" -ForegroundColor Green
        
        # 创建zip归档
        $zipName = "pong0_${version}_$($platform.GOOS)_$($platform.GOARCH).zip"
        $zipPath = "$distDir/$zipName"
        
        # 创建临时目录
        $tempDir = "$distDir/temp_$($platform.GOOS)_$($platform.GOARCH)"
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        
        # 复制文件到临时目录
        Copy-Item -Path $outputPath -Destination $tempDir/
        if (Test-Path -Path "$distDir/README.md") {
            Copy-Item -Path "$distDir/README.md" -Destination $tempDir/
        }
        
        # 创建zip文件
        Compress-Archive -Path "$tempDir/*" -DestinationPath $zipPath -Force
        
        # 删除临时目录和单独的二进制文件
        Remove-Item -Path $tempDir -Recurse -Force
        Remove-Item -Path $outputPath -Force
        
        Write-Host "  - 打包成功: $zipName" -ForegroundColor Green
    }
    else {
        Write-Host "  - 构建失败: $($platform.GOOS)/$($platform.GOARCH)" -ForegroundColor Red
    }
}

# 清理环境变量
$env:GOOS = ""
$env:GOARCH = ""

# 恢复原来的工作目录
Pop-Location

Write-Host "`n构建完成！所有文件已保存在$distDir目录中。" -ForegroundColor Green 