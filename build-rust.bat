@echo off
REM Sing-box Manager Rust 版本 Windows 编译脚本

set VERSION=2.0.0
set APP_NAME=singbox-manager
set BUILD_DIR=target\release-builds

echo ======================================
echo   Sing-box Manager Rust 版本编译
echo   版本: %VERSION%
echo ======================================
echo.

REM 检查 Rust 是否安装
where cargo >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未找到 Cargo，请先安装 Rust
    echo 访问 https://rustup.rs/ 安装 Rust
    exit /b 1
)

REM 清理旧的构建文件
echo 清理旧的构建文件...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM 编译当前平台
echo 编译 Windows 版本...
cargo build --release

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ======================================
    echo   编译成功！
    echo   二进制文件: target\release\%APP_NAME%.exe
    echo ======================================
    
    REM 复制到构建目录
    mkdir %BUILD_DIR%\windows
    copy target\release\%APP_NAME%.exe %BUILD_DIR%\windows\
    copy config.json.example %BUILD_DIR%\windows\
    
    echo.
    echo 文件已复制到: %BUILD_DIR%\windows\
) else (
    echo.
    echo 编译失败！
    exit /b 1
)

pause

