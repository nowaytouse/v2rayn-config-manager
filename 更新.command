#!/bin/bash
# 双击运行：更新内核 + 配置文件
cd "$(dirname "$0")"
python3 cm.py all
echo ""
echo "按任意键关闭..."
read -n 1
