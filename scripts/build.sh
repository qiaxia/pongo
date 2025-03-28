#!/bin/bash
# 跨平台构建并打包Pong0工具的Shell脚本

echo -e "\e[32m开始构建Pong0 - Ping0.cc数据获取工具...\e[0m"

# 确保在项目根目录执行
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 切换到项目根目录
pushd "$PROJECT_ROOT" > /dev/null
echo -e "\e[90m工作目录: $(pwd)\e[0m"

# 确保dist目录存在
DIST_DIR="dist"
if [ ! -d "$DIST_DIR" ]; then
    mkdir -p "$DIST_DIR"
    echo -e "\e[90m创建dist目录\e[0m"
fi

# 版本信息
VERSION="1.0.0"
BUILD_DATE=$(date +"%Y-%m-%d")
UPDATE_DATE=$(date +"%Y-%m-%d")

# 构建信息
echo -e "\e[33m构建版本: $VERSION (构建日期: $BUILD_DATE, 更新日期: $UPDATE_DATE)\e[0m"

# 复制README.md到dist目录
if [ -f "README.md" ]; then
    cp README.md "$DIST_DIR/"
    echo -e "\e[90m已复制README.md到dist目录\e[0m"
else
    echo -e "\e[33m未找到README.md文件\e[0m"
fi

# 定义构建平台（新增 FreeBSD 支持）
declare -a PLATFORMS=(
    "windows:amd64:.exe:Windows 64位"
    "windows:386:.exe:Windows 32位"
    "linux:amd64::Linux 64位"
    "linux:386::Linux 32位"
    "linux:arm64::Linux ARM64"
    "darwin:amd64::macOS 64位"
    "darwin:arm64::macOS ARM64"
    # 新增 FreeBSD 的构建目标
    "freebsd:amd64::FreeBSD 64位"
    "freebsd:arm64::FreeBSD ARM64"
)

# 主程序路径
MAIN_PATH="cmd/pong0"

# 为每个平台构建
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r os arch suffix name <<< "$platform"
    
    echo -e "\e[36m正在构建: $name ($os/$arch)...\e[0m"
    
    # 检查主程序目录是否存在
    if [ ! -d "$MAIN_PATH" ]; then
        echo -e "  \e[31m- 错误: 主程序目录 $MAIN_PATH 不存在\e[0m"
        continue
    fi
    
    output_name="pong0_${VERSION}_${os}_${arch}${suffix}"
    output_path="$DIST_DIR/$output_name"
    
    # 构建二进制文件
    GOOS=$os GOARCH=$arch go build -o "$output_path" -ldflags "-s -w -X main.Version=$VERSION -X main.buildDate=$BUILD_DATE -X ping0/internal/constants.UpdateDate=$UPDATE_DATE" ./$MAIN_PATH
    
    if [ $? -eq 0 ]; then
        echo -e "  \e[32m- 构建成功: $output_name\e[0m"
        
        # 创建zip归档
        zip_name="pong0_${VERSION}_${os}_${arch}.zip"
        zip_path="$DIST_DIR/$zip_name"
        
        # 创建临时目录
        temp_dir="$DIST_DIR/temp_${os}_${arch}"
        mkdir -p "$temp_dir"
        
        # 复制文件到临时目录
        cp "$output_path" "$temp_dir/"
        if [ -f "$DIST_DIR/README.md" ]; then
            cp "$DIST_DIR/README.md" "$temp_dir/"
        fi
        
        # 创建zip文件
        (cd "$temp_dir" && zip -q -r "../$zip_name" .)
        
        # 删除临时目录和单独的二进制文件
        rm -rf "$temp_dir"
        rm -f "$output_path"
        
        echo -e "  \e[32m- 打包成功: $zip_name\e[0m"
    else
        echo -e "  \e[31m- 构建失败: $os/$arch\e[0m"
    fi
done

# 恢复原来的工作目录
popd > /dev/null

echo -e "\n\e[32m构建完成！所有文件已保存在$DIST_DIR目录中。\e[0m"
