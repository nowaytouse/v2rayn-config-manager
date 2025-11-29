# v2rayN Config Manager

简易的 v2rayN 内核和配置文件更新工具（macOS）

## 功能

- ✅ 自动更新内核（预览版）
  - sing-box
  - mihomo
  - xray
- ✅ 自动更新 geofiles
  - geoip.dat
  - geosite.dat
  - Country.mmdb
  - geoip.metadb
- ✅ 自动更新配置文件

## 使用方法

### 快速开始

双击 `更新.command` 即可一键更新所有内容（内核 + geofiles + 配置）

### 命令行

```bash
python3 cm.py core          # 更新所有内核
python3 cm.py core singbox  # 只更新 singbox
python3 cm.py core mihomo   # 只更新 mihomo
python3 cm.py core xray     # 只更新 xray
python3 cm.py geo           # 更新 geofiles
python3 cm.py conf          # 更新配置文件
python3 cm.py all           # 全部更新
python3 cm.py status        # 查看状态
```

## 配置

编辑 `cm_config.json`：

```json
{
  "v2rayn_bin_path": "~/Library/Application Support/v2rayN/bin",
  "conf_save_path": "~/Library/Mobile Documents/com~apple~CloudDocs/Application/Conf/conf",
  "cores": {
    "singbox": {...},
    "mihomo": {...},
    "xray": {...}
  },
  "geofiles": {
    "geoip.dat": "https://...",
    "geosite.dat": "https://...",
    "Country.mmdb": "https://...",
    "geoip.metadb": "https://..."
  },
  "configs": [
    {
      "name": "文件名.conf",
      "url": "https://..."
    }
  ]
}
```

### 添加配置文件

在 `configs` 数组中添加：

```json
{
  "name": "my-config.conf",
  "url": "https://raw.githubusercontent.com/..."
}
```

### 添加更多 geofiles

在 `geofiles` 对象中添加：

```json
"文件名": "下载URL"
```

## 文件说明

### geofiles 用途

- **geoip.dat / geosite.dat** - v2ray/xray 使用的路由规则
- **Country.mmdb** - sing-box/mihomo 使用的 GeoIP 数据库
- **geoip.metadb** - mihomo 专用的 GeoIP 数据库

这些文件会被 v2rayN 中的各个内核自动调用，用于：
- 分流规则（国内/国外）
- 广告拦截
- 域名路由

### 更新频率

建议：
- 内核：每周检查更新
- geofiles：每月更新
- 配置文件：按需更新

## 系统要求

- macOS
- Python 3.6+
- v2rayN 已安装

## 注意事项

1. 内核更新为**预览版**（prerelease），更稳定请手动选择正式版
2. 更新会覆盖现有文件，请确保 v2rayN 未运行
3. 首次运行会创建配置文件，需手动填入订阅 URL

## License

MIT
