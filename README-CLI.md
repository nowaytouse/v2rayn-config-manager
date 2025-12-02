# Sing-box Manager - 纯 Rust CLI 版本

这是 Sing-box Manager 的纯 Rust 命令行版本，已完全移除 GUI 依赖，提供交互式菜单系统。

## 功能特性

- ✅ **纯 Rust 实现** - 无 GUI 依赖，轻量级
- ✅ **交互式 CLI 菜单** - 用户友好的命令行界面
- ✅ **订阅管理** - 添加、编辑、删除订阅源
- ✅ **自动更新** - 定时更新订阅和 Sing-box 核心
- ✅ **配置管理** - 灵活的配置选项
- ✅ **完整测试** - 6 个集成测试确保功能可靠

## 编译

```bash
cargo build --release
```

编译后的二进制文件位于 `target/release/singbox-manager`

## 使用方法

### 交互式模式（推荐）

```bash
./singbox-manager --interactive
# 或简写
./singbox-manager -i
```

这将启动交互式菜单，提供以下选项：

1. **查看配置** - 显示当前所有配置
2. **管理订阅** - 添加、编辑或删除订阅源
3. **更新订阅** - 立即下载所有订阅
4. **更新 Sing-box 核心** - 检查并更新 Sing-box 二进制文件
5. **执行所有更新** - 一次性执行订阅和核心更新
6. **配置设置** - 修改更新间隔和核心更新设置
7. **退出** - 退出程序

### 标准 CLI 模式

```bash
# 使用默认配置执行一次更新
./singbox-manager --once

# 使用自定义配置文件
./singbox-manager --config /path/to/config.json

# 定时执行更新（按配置文件中的间隔）
./singbox-manager

# 显示版本信息
./singbox-manager --version

# 显示帮助信息
./singbox-manager --help
```

## 配置文件

配置文件为 JSON 格式，默认位置为 `config.json`

### 配置文件示例

```json
{
  "subscriptions": [
    {
      "name": "我的订阅",
      "url": "https://example.com/config.json",
      "save_path": "./singbox_config.json"
    }
  ],
  "update_interval_hours": 24,
  "singbox_core_update": {
    "enabled": true,
    "check_prerelease": false,
    "install_path": "/usr/local/bin/sing-box"
  }
}
```

### 配置说明

- **subscriptions** - 订阅源列表
  - `name` - 订阅名称
  - `url` - 订阅 URL
  - `save_path` - 配置文件保存路径

- **update_interval_hours** - 自动更新间隔（小时），0 表示不自动更新

- **singbox_core_update** - Sing-box 核心更新配置
  - `enabled` - 是否启用自动更新
  - `check_prerelease` - 是否检查预发布版本
  - `install_path` - Sing-box 二进制文件安装路径

## 测试

运行所有测试：

```bash
cargo test
```

运行集成测试：

```bash
cargo test --test integration_tests
```

测试覆盖内容：

- ✅ 配置创建和默认值
- ✅ 订阅结构验证
- ✅ Sing-box 核心更新配置
- ✅ 配置序列化和反序列化
- ✅ 配置文件保存和加载

## 命令行参数

```
Sing-box 配置和核心自动更新工具 - 纯 Rust CLI 版本

USAGE:
    singbox-manager [OPTIONS]

OPTIONS:
  -i, --interactive          交互式菜单模式
  -c, --config <CONFIG>      配置文件路径 (默认: config.json)
  -o, --once                 仅执行一次更新任务后退出（非交互模式）
  -v, --version              显示版本信息
  -h, --help                 显示此帮助信息
```

## 项目结构

```
src/
├── main.rs           # 主程序入口
├── lib.rs            # 库导出
├── cli.rs            # 标准 CLI 模式
├── interactive.rs    # 交互式菜单系统
├── config.rs         # 配置文件处理
├── types.rs          # 数据结构定义
├── subscription.rs   # 订阅管理
└── updater.rs        # Sing-box 核心更新

tests/
└── integration_tests.rs  # 集成测试
```

## 依赖

- **tokio** - 异步运行时
- **reqwest** - HTTP 客户端
- **serde/serde_json** - JSON 序列化
- **anyhow** - 错误处理
- **tar/flate2/zip** - 压缩包处理
- **clap** - 命令行参数解析
- **uuid** - UUID 生成
- **chrono** - 时间处理

## 特性对比

| 功能 | Go 版本 | Rust CLI 版本 |
|------|--------|--------------|
| GUI | ✅ | ❌ |
| 交互式 CLI | ❌ | ✅ |
| 订阅管理 | ✅ | ✅ |
| 核心更新 | ✅ | ✅ |
| 定时更新 | ✅ | ✅ |
| 二进制大小 | ~50MB | ~10MB |
| 启动速度 | 中等 | 快速 |
| 内存占用 | 中等 | 低 |

## 故障排除

### 问题：权限不足无法更新 Sing-box

**解决方案**：使用 sudo 运行或确保有写入权限

```bash
sudo ./singbox-manager --interactive
```

### 问题：无法连接到 GitHub API

**解决方案**：检查网络连接或配置代理

### 问题：订阅下载失败

**解决方案**：
1. 检查订阅 URL 是否正确
2. 确保网络连接正常
3. 检查防火墙设置

## 许可证

MIT

## 更新日志

### v2.0.0 (Rust CLI 版本)
- 完全移除 GUI 依赖
- 新增交互式菜单系统
- 改进错误处理
- 添加完整的集成测试
- 优化性能和内存占用
