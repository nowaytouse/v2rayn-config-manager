use serde::{Deserialize, Serialize};
use std::path::PathBuf;

/// 配置文件主结构
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    /// 订阅源列表
    pub subscriptions: Vec<Subscription>,
    /// 更新间隔（小时）
    pub update_interval_hours: u64,
    /// Sing-box 核心更新配置
    pub singbox_core_update: SingboxCoreUpdate,
}

/// 单个订阅源
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Subscription {
    /// 订阅名称
    pub name: String,
    /// 订阅 URL
    pub url: String,
    /// 保存路径
    pub save_path: PathBuf,
}

/// Sing-box 核心更新配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SingboxCoreUpdate {
    /// 是否启用自动更新
    pub enabled: bool,
    /// 是否检查预发布版本
    pub check_prerelease: bool,
    /// 安装路径
    pub install_path: PathBuf,
}

/// GitHub Release 信息
#[derive(Debug, Clone, Deserialize)]
pub struct GithubRelease {
    /// 标签名称
    pub tag_name: String,
    /// 是否为预发布版本
    pub prerelease: bool,
    /// 资源文件列表
    pub assets: Vec<Asset>,
}

/// GitHub Release 资源文件
#[derive(Debug, Clone, Deserialize)]
pub struct Asset {
    /// 文件名
    pub name: String,
    /// 下载链接
    pub browser_download_url: String,
}

impl Config {
    /// 创建默认配置
    pub fn default() -> Self {
        let install_path = if cfg!(windows) {
            PathBuf::from("C:\\Program Files\\sing-box\\sing-box.exe")
        } else {
            PathBuf::from("/usr/local/bin/sing-box")
        };

        Self {
            subscriptions: vec![Subscription {
                name: "默认订阅".to_string(),
                url: "https://example.com/your-singbox-config.json".to_string(),
                save_path: PathBuf::from("./singbox_config.json"),
            }],
            update_interval_hours: 24,
            singbox_core_update: SingboxCoreUpdate {
                enabled: true,
                check_prerelease: false,
                install_path,
            },
        }
    }
}
