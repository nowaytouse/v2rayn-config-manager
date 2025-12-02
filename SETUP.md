# Sing-box Manager 启用指南

## 🚀 快速启用

### 1. 编译 Rust 内核

```bash
cd config-manager
cargo build --release
```

编译完成后，二进制文件位于：`target/release/singbox-manager`

### 2. 配置文件设置

编辑 `cm_config.json`，确保以下路径正确：

```json
{
  "v2rayn_bin_path": "/Users/nyamiiko/Library/Application Support/v2rayN/bin",
  "conf_save_path": "/Users/nyamiiko/Library/Mobile Documents/com~apple~CloudDocs/Application/Conf/conf",
  "cores": {
    "singbox": {
      "repo": "SagerNet/sing-box",
      "binary_name": "sing-box",
      "subdir": "sing_box"
    },
    "mihomo": {
      "repo": "MetaCubeX/mihomo",
      "binary_name": "mihomo",
      "subdir": "mihomo"
    },
    "xray": {
      "repo": "XTLS/Xray-core",
      "binary_name": "xray",
      "subdir": "xray"
    }
  },
  "geofiles": {
    "geoip.dat": "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat",
    "geosite.dat": "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat",
    "Country.mmdb": "https://github.com/xream/geoip/releases/latest/download/ipinfo.country.mmdb",
    "geoip.metadb": "https://github.com/MetaCubeX/meta-rules-dat/releases/latest/download/geoip.metadb"
  },
  "configs": [
    {
      "name": "✨ 🏆Sing box配置合成1.13.0+🎊✨.conf",
      "url": "https://gist.githubusercontent.com/nowaytouse/dcf7620c59dffcb1528fa5af02ba5b08/raw/sing"
    }
  ]
}
```

### 3. 使用方法

#### 交互式模式（推荐）

```bash
./target/release/singbox-manager --interactive
```

#### 命令行模式

```bash
# 更新所有内核
./target/release/singbox-manager core

# 只更新 sing-box
./target/release/singbox-manager core singbox

# 更新 geofiles
./target/release/singbox-manager geo

# 更新配置文件
./target/release/singbox-manager conf

# 全部更新
./target/release/singbox-manager all

# 查看状态
./target/release/singbox-manager status
```

#### 一键更新脚本

创建 `update.sh`：

```bash
#!/bin/bash
cd "$(dirname "$0")"
./target/release/singbox-manager --once all
```

## 📋 功能清单

- ✅ 自动检测最新版本
- ✅ 下载并解压内核
- ✅ 更新 geofiles
- ✅ 下载配置文件
- ✅ 交互式菜单
- ✅ 日志记录
- ✅ 错误处理

## 🔧 配置说明

### 内核配置

每个内核需要配置：
- `repo`: GitHub 仓库（格式：owner/repo）
- `binary_name`: 二进制文件名
- `subdir`: 保存子目录

### Geofiles 配置

支持任意数量的 geofiles，只需在 `geofiles` 对象中添加：

```json
"文件名": "下载URL"
```

### 配置文件

支持多个配置文件，每个需要：
- `name`: 保存的文件名
- `url`: 下载 URL

## ⚠️ 注意事项

1. **备份重要文件** - 更新前请备份现有配置
2. **关闭 v2rayN** - 更新时请确保 v2rayN 未运行
3. **网络连接** - 需要稳定的网络连接
4. **权限** - 确保有写入权限到配置目录

## 🐛 故障排查

### 编译失败

```bash
# 清理并重新编译
cargo clean
cargo build --release
```

### 下载失败

- 检查网络连接
- 检查 GitHub 是否可访问
- 检查配置文件中的 URL 是否正确

### 权限错误

```bash
# 检查目录权限
ls -la ~/Library/Application\ Support/v2rayN/bin
chmod -R 755 ~/Library/Application\ Support/v2rayN/bin
```

## 📚 更多信息

- 查看 `README.md` 了解功能详情
- 查看 `cm_config.json` 了解配置选项
- 查看源代码了解实现细节

---

**版本**: 2.0.0  
**平台**: macOS  
**语言**: Rust  
**最后更新**: 2025-11-14
