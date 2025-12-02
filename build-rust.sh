#!/bin/bash

# Sing-box Manager Rust 版本多平台编译脚本

VERSION="2.0.0"
APP_NAME="singbox-manager"
BUILD_DIR="target/release-builds"

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Sing-box Manager Rust 版本编译${NC}"
echo -e "${BLUE}  版本: ${VERSION}${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 检查 Rust 是否安装
if ! command -v cargo &> /dev/null; then
    echo -e "${RED}错误: 未找到 Cargo，请先安装 Rust${NC}"
    echo "访问 https://rustup.rs/ 安装 Rust"
    exit 1
fi

# 检查 cross 是否安装（用于交叉编译）
if ! command -v cross &> /dev/null; then
    echo -e "${GREEN}安装 cross 用于交叉编译...${NC}"
    cargo install cross --git https://github.com/cross-rs/cross
fi

# 清理旧的构建文件
echo -e "${GREEN}清理旧的构建文件...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 构建函数
build() {
    local TARGET=$1
    local OUTPUT_NAME="${APP_NAME}"
    
    if [[ $TARGET == *"windows"* ]]; then
        OUTPUT_NAME="${APP_NAME}.exe"
    fi
    
    echo -e "${GREEN}编译目标: ${TARGET}...${NC}"
    
    # 使用 cross 进行交叉编译
    if cross build --release --target ${TARGET}; then
        local OUTPUT_DIR="${BUILD_DIR}/${TARGET}"
        mkdir -p ${OUTPUT_DIR}
        
        # 复制二进制文件
        cp "target/${TARGET}/release/${OUTPUT_NAME}" "${OUTPUT_DIR}/"
        
        # 复制配置文件示例
        cp config.json.example "${OUTPUT_DIR}/"
        
        # 创建压缩包
        cd ${BUILD_DIR}
        if [[ $TARGET == *"windows"* ]]; then
            zip -r "${APP_NAME}-${VERSION}-${TARGET}.zip" "${TARGET}"
        else
            tar -czf "${APP_NAME}-${VERSION}-${TARGET}.tar.gz" "${TARGET}"
        fi
        cd ../..
        
        echo -e "${GREEN}✓ ${TARGET} 编译成功${NC}"
    else
        echo -e "${RED}✗ ${TARGET} 编译失败${NC}"
    fi
}

# 当前平台编译（快速模式）
build_current() {
    echo -e "${GREEN}编译当前平台...${NC}"
    cargo build --release
    
    local OUTPUT_NAME="${APP_NAME}"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        OUTPUT_NAME="${APP_NAME}.exe"
    fi
    
    echo -e "${GREEN}✓ 当前平台编译成功${NC}"
    echo -e "二进制文件位置: target/release/${OUTPUT_NAME}"
}

# 检查参数
if [ "$1" == "--current" ]; then
    build_current
    exit 0
fi

# 编译所有目标平台
echo -e "${BLUE}开始多平台编译...${NC}"
echo ""

# Linux
build "x86_64-unknown-linux-gnu"
build "aarch64-unknown-linux-gnu"
build "i686-unknown-linux-gnu"

# macOS (需要在 macOS 上编译)
if [[ "$OSTYPE" == "darwin"* ]]; then
    build "x86_64-apple-darwin"
    build "aarch64-apple-darwin"
else
    echo -e "${RED}跳过 macOS 目标（需要在 macOS 上编译）${NC}"
fi

# Windows
build "x86_64-pc-windows-gnu"
build "i686-pc-windows-gnu"

echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  编译完成！${NC}"
echo -e "${BLUE}  输出目录: ${BUILD_DIR}/${NC}"
echo -e "${BLUE}======================================${NC}"

# 显示构建的文件
echo ""
echo "构建的文件:"
ls -lh ${BUILD_DIR}/*.{tar.gz,zip} 2>/dev/null || echo "未找到压缩包"

