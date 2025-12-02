@echo off
REM Sing-box Manager Windows 编译脚本

set VERSION=2.0.0
set APP_NAME=singbox-manager
set BUILD_DIR=build

echo ======================================
echo   Sing-box Manager 多平台编译
echo   版本: %VERSION%
echo ======================================
echo.

REM 清理旧的构建文件
echo 清理旧的构建文件...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM 获取依赖
echo 获取依赖...
go mod tidy
go mod download

REM 编译 Windows AMD64
echo 编译 Windows AMD64...
set GOOS=windows
set GOARCH=amd64
set OUTPUT_DIR=%BUILD_DIR%\%GOOS%-%GOARCH%
mkdir %OUTPUT_DIR%
set CGO_ENABLED=1
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\%APP_NAME%.exe" .
copy config.json "%OUTPUT_DIR%\config.json.example"
cd %BUILD_DIR%
tar -a -c -f "%APP_NAME%-%VERSION%-%GOOS%-%GOARCH%.zip" "%GOOS%-%GOARCH%"
cd ..

REM 编译 Windows 386
echo 编译 Windows 386...
set GOOS=windows
set GOARCH=386
set OUTPUT_DIR=%BUILD_DIR%\%GOOS%-%GOARCH%
mkdir %OUTPUT_DIR%
set CGO_ENABLED=1
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\%APP_NAME%.exe" .
copy config.json "%OUTPUT_DIR%\config.json.example"
cd %BUILD_DIR%
tar -a -c -f "%APP_NAME%-%VERSION%-%GOOS%-%GOARCH%.zip" "%GOOS%-%GOARCH%"
cd ..

REM 编译 Windows ARM64
echo 编译 Windows ARM64...
set GOOS=windows
set GOARCH=arm64
set OUTPUT_DIR=%BUILD_DIR%\%GOOS%-%GOARCH%
mkdir %OUTPUT_DIR%
set CGO_ENABLED=1
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\%APP_NAME%.exe" .
copy config.json "%OUTPUT_DIR%\config.json.example"
cd %BUILD_DIR%
tar -a -c -f "%APP_NAME%-%VERSION%-%GOOS%-%GOARCH%.zip" "%GOOS%-%GOARCH%"
cd ..

echo.
echo ======================================
echo   编译完成！
echo   输出目录: %BUILD_DIR%
echo ======================================

dir %BUILD_DIR%\*.zip

pause

