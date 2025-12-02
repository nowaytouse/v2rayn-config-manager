#!/bin/bash

# Sing-box Manager 多平台编译脚本

VERSION="2.0.0"
APP_NAME="singbox-manager"
BUILD_DIR="build"

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Sing-box Manager 多平台编译${NC}"
echo -e "${BLUE}  版本: ${VERSION}${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 清理旧的构建文件
echo -e "${GREEN}清理旧的构建文件...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 获取依赖
echo -e "${GREEN}获取依赖...${NC}"
go mod tidy
go mod download

# 构建函数
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME="${APP_NAME}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${APP_NAME}.exe"
    fi
    
    local OUTPUT_DIR="${BUILD_DIR}/${GOOS}-${GOARCH}"
    mkdir -p ${OUTPUT_DIR}
    
    echo -e "${GREEN}编译 ${GOOS}/${GOARCH}...${NC}"
    
    CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w -X main.version=${VERSION}" -o "${OUTPUT_DIR}/${OUTPUT_NAME}" .
    
    if [ $? -eq 0 ]; then
        # 复制配置文件示例
        cp config.json "${OUTPUT_DIR}/config.json.example"
        
        # 创建压缩包
        cd ${BUILD_DIR}
        if [ "$GOOS" = "windows" ]; then
            zip -r "${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip" "${GOOS}-${GOARCH}"
        else
            tar -czf "${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "${GOOS}-${GOARCH}"
        fi
        cd ..
        
        echo -e "${GREEN}✓ ${GOOS}/${GOARCH} 编译成功${NC}"
    else
        echo -e "${RED}✗ ${GOOS}/${GOARCH} 编译失败${NC}"
    fi
}

# 编译各平台版本
# Linux
build "linux" "amd64"
build "linux" "arm64"
build "linux" "386"

# macOS
build "darwin" "amd64"
build "darwin" "arm64"

# Windows
build "windows" "amd64"
build "windows" "386"
build "windows" "arm64"

echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  编译完成！${NC}"
echo -e "${BLUE}  输出目录: ${BUILD_DIR}/${NC}"
echo -e "${BLUE}======================================${NC}"

# 显示构建的文件
echo ""
echo "构建的文件:"
ls -lh ${BUILD_DIR}/*.{tar.gz,zip} 2>/dev/null || echo "未找到压缩包"

