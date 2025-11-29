#!/usr/bin/env python3
"""
Config Manager CLI - v2rayN 内核和配置更新工具
用法:
  python cm.py core          # 更新所有内核(预览版)
  python cm.py core singbox  # 只更新singbox
  python cm.py core xray     # 只更新xray
  python cm.py geo           # 更新geofiles
  python cm.py conf          # 更新所有配置文件
  python cm.py all           # 更新内核+geo+配置
  python cm.py status        # 查看状态
"""

import os, sys, json, platform, tempfile, shutil, gzip, tarfile, urllib.request
from pathlib import Path

APP_DIR = Path(__file__).parent
CONFIG_FILE = APP_DIR / "cm_config.json"
UA = {"User-Agent": "ConfigManager/1.0"}

DEFAULT_CONFIG = {
    "v2rayn_bin_path": str(Path.home() / "Library/Application Support/v2rayN/bin"),
    "conf_save_path": str(Path.home() / "Library/Mobile Documents/com~apple~CloudDocs/Application/Conf/conf"),
    "cores": {
        "singbox": {"repo": "SagerNet/sing-box", "binary_name": "sing-box", "subdir": "sing_box"},
        "mihomo": {"repo": "MetaCubeX/mihomo", "binary_name": "mihomo", "subdir": "mihomo"},
        "xray": {"repo": "XTLS/Xray-core", "binary_name": "xray", "subdir": "xray"}
    },
    "geofiles": {
        "geoip.dat": "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat",
        "geosite.dat": "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
    },
    "configs": []
}

def get_arch():
    m = platform.machine().lower()
    return "amd64" if m in ("x86_64", "amd64") else "arm64"

def load_config():
    return json.loads(CONFIG_FILE.read_text()) if CONFIG_FILE.exists() else DEFAULT_CONFIG.copy()

def save_config(cfg):
    CONFIG_FILE.write_text(json.dumps(cfg, indent=2, ensure_ascii=False))

def fetch_prerelease(repo):
    """获取最新预览版release"""
    try:
        req = urllib.request.Request(f"https://api.github.com/repos/{repo}/releases", headers=UA)
        for rel in json.loads(urllib.request.urlopen(req, timeout=30).read()):
            if rel.get("prerelease"):
                return rel
        return None
    except Exception as e:
        print(f"   ❌ API请求失败: {e}")
        return None

def get_asset(release, core):
    """获取对应平台的asset"""
    arch = get_arch()
    for a in release.get("assets", []):
        n = a["name"].lower()
        if "darwin" in n and arch in n and "sha" not in n:
            if core == "singbox" and n.endswith(".tar.gz"):
                return a["browser_download_url"], a["name"]
            elif core == "mihomo" and n.endswith(".gz") and not n.endswith(".tar.gz"):
                return a["browser_download_url"], a["name"]
            elif core == "xray" and "macos" in n and n.endswith(".zip"):
                return a["browser_download_url"], a["name"]
    return None, None

def download(url, dest):
    """下载文件"""
    req = urllib.request.Request(url, headers=UA)
    Path(dest).write_bytes(urllib.request.urlopen(req, timeout=120).read())

def install_core(archive, dest, binary_name):
    """解压并安装内核，覆盖现有文件"""
    import zipfile, time
    dest = Path(dest)
    dest.parent.mkdir(parents=True, exist_ok=True)
    
    old_size = dest.stat().st_size if dest.exists() else 0
    
    with tempfile.TemporaryDirectory() as tmp:
        if archive.endswith(".tar.gz"):
            with tarfile.open(archive, "r:gz") as t:
                t.extractall(tmp)
        elif archive.endswith(".zip"):
            with zipfile.ZipFile(archive, 'r') as z:
                z.extractall(tmp)
        elif archive.endswith(".gz"):
            if dest.exists():
                dest.unlink()
            with gzip.open(archive, 'rb') as gz:
                dest.write_bytes(gz.read())
            os.chmod(dest, 0o755)
            os.utime(dest, (time.time(), time.time()))
            return old_size, dest.stat().st_size
        
        # 查找并复制二进制文件
        for f in Path(tmp).rglob(binary_name):
            if f.is_file():
                if dest.exists():
                    dest.unlink()
                shutil.copy2(f, dest)
                break
    
    os.chmod(dest, 0o755)
    os.utime(dest, (time.time(), time.time()))
    return old_size, dest.stat().st_size

def update_core(name=None):
    """更新内核"""
    cfg = load_config()
    bin_path = Path(cfg["v2rayn_bin_path"])
    
    if not bin_path.exists():
        print(f"❌ v2rayN bin目录不存在: {bin_path}")
        return
    
    for core, info in cfg["cores"].items():
        if name and core != name:
            continue
        
        print(f"\n🔄 {core} (预览版)...")
        rel = fetch_prerelease(info["repo"])
        if not rel:
            print("   ❌ 获取版本失败")
            continue
        
        print(f"   📦 {rel['tag_name']}")
        url, fn = get_asset(rel, core)
        if not url:
            print(f"   ❌ 找不到darwin/{get_arch()}文件")
            continue
        
        with tempfile.NamedTemporaryFile(suffix=fn, delete=False) as tmp:
            try:
                print("   📥 下载中...")
                download(url, tmp.name)
                dest = bin_path / info.get("subdir", "") / info["binary_name"]
                old_size, new_size = install_core(tmp.name, str(dest), info["binary_name"])
                # 验证信息
                from datetime import datetime
                mtime = datetime.fromtimestamp(dest.stat().st_mtime).strftime("%Y-%m-%d %H:%M:%S")
                print(f"   ✅ 完成 → {dest}")
                print(f"   📊 大小: {old_size/1024/1024:.1f}MB → {new_size/1024/1024:.1f}MB | 时间: {mtime}")
            except Exception as e:
                print(f"   ❌ 失败: {e}")
            finally:
                os.unlink(tmp.name)

def update_geofiles():
    """更新geofiles"""
    cfg = load_config()
    bin_path = Path(cfg["v2rayn_bin_path"])
    geofiles = cfg.get("geofiles", {})
    
    if not geofiles:
        print("⚠️  没有配置geofiles")
        return
    
    print(f"\n🌍 更新geofiles → {bin_path}")
    for name, url in geofiles.items():
        try:
            print(f"   📥 {name}...")
            dest = bin_path / name
            old_size = dest.stat().st_size if dest.exists() else 0
            req = urllib.request.Request(url, headers=UA)
            dest.write_bytes(urllib.request.urlopen(req, timeout=120).read())
            new_size = dest.stat().st_size
            print(f"   ✅ {name} ({old_size/1024/1024:.1f}MB → {new_size/1024/1024:.1f}MB)")
        except Exception as e:
            print(f"   ❌ {name}: {e}")

def update_configs():
    """更新配置文件"""
    cfg = load_config()
    save_path = Path(cfg["conf_save_path"])
    save_path.mkdir(parents=True, exist_ok=True)
    
    configs = [c for c in cfg.get("configs", []) if c.get("url")]
    if not configs:
        print("⚠️  没有配置URL")
        return
    
    print(f"\n🔄 更新配置 → {save_path}")
    for c in configs:
        try:
            print(f"   📥 {c['name']}...")
            req = urllib.request.Request(c["url"], headers=UA)
            (save_path / c["name"]).write_bytes(urllib.request.urlopen(req, timeout=60).read())
            print(f"   ✅ {c['name']}")
        except Exception as e:
            print(f"   ❌ {c['name']}: {e}")

def status():
    """显示状态"""
    cfg = load_config()
    bin_path = Path(cfg["v2rayn_bin_path"])
    
    print("\n" + "="*50)
    print("📋 Config Manager")
    print("="*50)
    print(f"\n🗂️  bin: {bin_path}")
    
    if bin_path.exists():
        print("\n   内核:")
        for name, info in cfg["cores"].items():
            p = bin_path / info.get("subdir", "") / info["binary_name"]
            s = "✅" if p.exists() else "❌"
            print(f"   {s} {name}: {p.name}")
        
        print("\n   geofiles:")
        for name in cfg.get("geofiles", {}).keys():
            p = bin_path / name
            s = "✅" if p.exists() else "❌"
            print(f"   {s} {name}")
    else:
        print("   ❌ 目录不存在")
    
    print(f"\n📁 conf: {cfg['conf_save_path']}")
    print(f"\n📋 配置: {len([c for c in cfg.get('configs', []) if c.get('url')])} 个已配置URL")
    print(f"\n⚙️  {CONFIG_FILE}")
    print("="*50)

def main():
    if len(sys.argv) < 2:
        print(__doc__)
        return
    
    if not CONFIG_FILE.exists():
        save_config(DEFAULT_CONFIG)
        print(f"✅ 已创建: {CONFIG_FILE}")
        print("⚠️  请编辑配置文件添加订阅URL")
        if sys.argv[1] != "status":
            return
    
    cmd = sys.argv[1].lower()
    if cmd == "core":
        update_core(sys.argv[2].lower() if len(sys.argv) > 2 else None)
    elif cmd == "geo":
        update_geofiles()
    elif cmd == "conf":
        update_configs()
    elif cmd == "all":
        update_core()
        update_geofiles()
        update_configs()
    elif cmd == "status":
        status()
    else:
        print(__doc__)

if __name__ == "__main__":
    main()
