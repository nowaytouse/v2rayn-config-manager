# 项目架构说明

## 📁 项目结构

```
singbox-manager/
├── main.go                      # 主入口，处理命令行参数和模式切换
├── go.mod                       # Go 模块依赖定义
├── go.sum                       # 依赖校验和
├── config.json                  # 用户配置文件（运行时生成）
├── config.json.example          # 配置文件示例
├── pkg/                         # 包目录
│   ├── core/                    # 核心业务逻辑
│   │   ├── types.go            # 数据结构定义
│   │   ├── config.go           # 配置文件加载/保存
│   │   ├── subscription.go     # 订阅下载管理
│   │   └── updater.go          # 核心更新管理
│   └── gui/                     # GUI 界面
│       └── gui.go              # Fyne 图形界面实现
├── build.sh                     # Linux/macOS 编译脚本
├── build.bat                    # Windows 编译脚本
├── Makefile                     # Make 构建工具配置
├── .gitignore                   # Git 忽略文件配置
├── LICENSE                      # MIT 开源许可证
├── README.md                    # 项目说明文档
├── QUICK_START.md              # 快速开始指南
└── ARCHITECTURE.md             # 本文档

```

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────────────────────────────┐
│         用户界面层 (UI Layer)            │
│  ┌──────────┐         ┌──────────┐     │
│  │   GUI    │         │   CLI    │     │
│  │  模式    │         │  模式    │     │
│  └────┬─────┘         └────┬─────┘     │
└───────┼──────────────────────┼──────────┘
        │                      │
        └──────────┬───────────┘
                   │
┌──────────────────┼──────────────────────┐
│         业务逻辑层 (Business Logic)     │
│                  │                      │
│  ┌───────────────┴────────────────┐    │
│  │      核心业务模块 (pkg/core)    │    │
│  │                                 │    │
│  │  ┌──────────────────────────┐  │    │
│  │  │  配置管理 (config.go)    │  │    │
│  │  │  - 加载配置              │  │    │
│  │  │  - 保存配置              │  │    │
│  │  └──────────────────────────┘  │    │
│  │                                 │    │
│  │  ┌──────────────────────────┐  │    │
│  │  │  订阅管理                │  │    │
│  │  │  (subscription.go)       │  │    │
│  │  │  - 下载订阅配置          │  │    │
│  │  │  - 多源管理              │  │    │
│  │  └──────────────────────────┘  │    │
│  │                                 │    │
│  │  ┌──────────────────────────┐  │    │
│  │  │  核心更新 (updater.go)   │  │    │
│  │  │  - GitHub API 集成       │  │    │
│  │  │  - 版本检测              │  │    │
│  │  │  - 自动下载安装          │  │    │
│  │  └──────────────────────────┘  │    │
│  │                                 │    │
│  └─────────────────────────────────┘    │
└──────────────────────────────────────────┘
                   │
┌──────────────────┼──────────────────────┐
│         数据层 (Data Layer)             │
│                  │                      │
│  ┌───────────────┴────────────────┐    │
│  │  - config.json (配置持久化)    │    │
│  │  - 下载的配置文件              │    │
│  │  - 日志输出                    │    │
│  └────────────────────────────────┘    │
└──────────────────────────────────────────┘
```

## 🔄 工作流程

### 1. 启动流程

```
main.go (入口)
    ├─> 解析命令行参数
    ├─> 判断运行模式
    │   ├─> GUI 模式
    │   │   └─> 启动 pkg/gui.GUI
    │   │       ├─> 加载配置
    │   │       ├─> 创建界面
    │   │       └─> 事件循环
    │   │
    │   └─> CLI 模式
    │       └─> runCLI()
    │           ├─> 加载配置
    │           ├─> 执行更新任务
    │           └─> 启动定时器（可选）
```

### 2. 订阅更新流程

```
订阅更新请求
    ↓
SubscriptionManager.RunSubscriptionDownloader()
    ↓
遍历所有订阅源
    ↓
DownloadSubscription() (对每个订阅)
    ├─> HTTP GET 请求
    ├─> 接收响应
    ├─> 保存到文件
    └─> 日志记录
```

### 3. 核心更新流程

```
核心更新请求
    ↓
CoreUpdater.RunCoreUpdater()
    ↓
GetLatestRelease()
    ├─> 访问 GitHub API
    ├─> 获取版本列表
    └─> 筛选合适版本
    ↓
FindMatchingAsset()
    ├─> 检测系统平台
    ├─> 匹配对应资源
    └─> 返回下载链接
    ↓
DownloadAndExtract()
    ├─> 下载压缩包
    ├─> 解压到临时目录
    └─> 提取可执行文件
    ↓
InstallFile()
    ├─> 复制到目标路径
    ├─> 设置执行权限
    └─> 完成安装
```

## 🎨 GUI 设计

### 界面组件结构

```
主窗口 (MainWindow)
├── 状态栏 (StatusBar)
│   └── 显示当前状态
│
├── 主内容区 (MainContent)
│   ├── 订阅管理卡片
│   │   ├── 订阅列表
│   │   └── 添加/编辑/删除按钮
│   │
│   ├── 核心更新配置卡片
│   │   ├── 启用开关
│   │   ├── 预发布版本开关
│   │   └── 安装路径设置
│   │
│   ├── 定时更新配置卡片
│   │   ├── 更新间隔设置
│   │   └── 自动更新开关
│   │
│   └── 日志显示区域
│       └── 多行文本框（只读）
│
└── 操作按钮栏
    ├── 立即更新订阅
    ├── 立即更新核心
    ├── 全部更新
    └── 清空日志
```

## 🔌 核心接口

### LogCallback
```go
type LogCallback func(message string)
```
用于将日志消息传递给 UI 层显示。

### SubscriptionManager
```go
type SubscriptionManager struct {
    logCallback LogCallback
}

func (sm *SubscriptionManager) DownloadSubscription(sub Subscription) error
func (sm *SubscriptionManager) RunSubscriptionDownloader(config *Config)
```

### CoreUpdater
```go
type CoreUpdater struct {
    logCallback LogCallback
}

func (cu *CoreUpdater) GetLatestRelease(checkPrerelease bool) (*GithubRelease, error)
func (cu *CoreUpdater) FindMatchingAsset(release *GithubRelease) *Asset
func (cu *CoreUpdater) DownloadAndExtract(asset *Asset) (string, error)
func (cu *CoreUpdater) InstallFile(sourcePath, destPath string) error
func (cu *CoreUpdater) RunCoreUpdater(config *Config) error
```

## 🔐 多平台支持

### 平台检测
使用 `runtime.GOOS` 和 `runtime.GOARCH` 检测当前平台：
- `windows`: Windows 系统
- `darwin`: macOS 系统
- `linux`: Linux 系统

### 平台特定处理

1. **路径分隔符**
   - Windows: `\` (但 Go 会自动处理 `/`)
   - Unix: `/`

2. **可执行文件扩展名**
   - Windows: `.exe`
   - Unix: 无扩展名

3. **文件权限**
   - Windows: 不需要设置执行权限
   - Unix: 需要 `chmod 0755`

4. **默认安装路径**
   - Windows: `C:\Program Files\sing-box\sing-box.exe`
   - macOS: `/usr/local/bin/sing-box`
   - Linux: `/usr/local/bin/sing-box`

## 🧩 依赖说明

### 核心依赖
- **fyne.io/fyne/v2**: 跨平台 GUI 框架
  - 支持 Windows/Linux/macOS
  - 原生外观
  - 易于使用的 API

### 标准库
- `encoding/json`: JSON 配置文件处理
- `net/http`: HTTP 请求（下载配置和核心）
- `archive/tar`, `archive/zip`: 压缩包解压
- `time`: 定时任务管理
- `flag`: 命令行参数解析
- `runtime`: 平台检测

## 🚀 性能优化

1. **并发安全**: GUI 操作在 goroutine 中执行，避免阻塞界面
2. **资源清理**: 使用 defer 确保资源正确释放
3. **临时文件**: 下载的压缩包使用临时目录，完成后自动清理
4. **增量更新**: 只在版本变化时下载核心

## 🔮 未来扩展方向

1. **多语言支持**: i18n 国际化
2. **主题支持**: 明暗主题切换
3. **版本回滚**: 保留历史版本，支持回滚
4. **配置备份**: 自动备份配置文件
5. **通知系统**: 更新完成后的系统通知
6. **代理支持**: 通过代理下载更新
7. **插件系统**: 支持自定义处理脚本

## 📖 代码规范

1. **包结构**: 按功能划分包，保持职责单一
2. **错误处理**: 所有错误都向上传递，由调用方决定如何处理
3. **日志记录**: 通过 LogCallback 统一日志输出
4. **配置管理**: 集中在 config.go，便于维护
5. **命名规范**: 遵循 Go 语言命名约定

---

如有疑问或建议，欢迎提交 Issue 讨论！

