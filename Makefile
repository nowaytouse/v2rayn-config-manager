.PHONY: all build clean run run-gui install test deps

APP_NAME := singbox-manager
VERSION := 2.0.0
BUILD_DIR := build

# 默认目标
all: clean deps build

# 获取依赖
deps:
	@echo "获取依赖..."
	go mod tidy
	go mod download

# 开发环境运行（命令行模式）
run:
	@echo "运行命令行模式..."
	go run main.go

# 开发环境运行（GUI 模式）
run-gui:
	@echo "运行 GUI 模式..."
	go run main.go -gui

# 运行一次性更新
run-once:
	@echo "运行一次性更新..."
	go run main.go -once

# 当前平台编译
build:
	@echo "编译当前平台..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

# 多平台编译
build-all:
	@echo "多平台编译..."
	@chmod +x build.sh
	@./build.sh

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 安装到系统（需要 root 权限）
install: build
	@echo "安装到系统..."
	@sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "安装完成: /usr/local/bin/$(APP_NAME)"

# 卸载
uninstall:
	@echo "卸载..."
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "卸载完成"

# 测试
test:
	@echo "运行测试..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 帮助信息
help:
	@echo "可用的 make 命令："
	@echo "  make          - 清理、获取依赖并编译"
	@echo "  make deps     - 获取依赖"
	@echo "  make build    - 编译当前平台"
	@echo "  make build-all - 多平台编译"
	@echo "  make run      - 运行命令行模式"
	@echo "  make run-gui  - 运行 GUI 模式"
	@echo "  make run-once - 运行一次性更新"
	@echo "  make install  - 安装到系统"
	@echo "  make uninstall - 从系统卸载"
	@echo "  make clean    - 清理构建文件"
	@echo "  make test     - 运行测试"
	@echo "  make fmt      - 格式化代码"
	@echo "  make lint     - 代码检查"

