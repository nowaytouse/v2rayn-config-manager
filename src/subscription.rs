use anyhow::{Context, Result};
use tokio::fs;

use crate::types::{Config, Subscription};

/// 订阅管理器（纯CLI模式）
pub struct SubscriptionManager;

impl SubscriptionManager {
    pub fn new() -> Self {
        Self
    }

    /// 下载单个订阅
    pub async fn download_subscription(&self, sub: &Subscription) -> Result<()> {
        println!("📥 开始下载订阅 [{}] 从 {}", sub.name, sub.url);

        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()?;

        let response = client
            .get(&sub.url)
            .send()
            .await
            .context("下载失败")?;

        if !response.status().is_success() {
            anyhow::bail!("下载失败，服务器返回状态码: {}", response.status());
        }

        let content = response.bytes().await.context("读取响应内容失败")?;

        if let Some(parent) = sub.save_path.parent() {
            fs::create_dir_all(parent)
                .await
                .context("创建目录失败")?;
        }

        fs::write(&sub.save_path, content)
            .await
            .context("保存文件失败")?;

        println!(
            "✅ 订阅 [{}] 下载成功，已保存至 {}",
            sub.name,
            sub.save_path.display()
        );

        Ok(())
    }

    /// 运行订阅下载器（下载所有订阅）
    pub async fn run_all(&self, config: &Config) -> Result<()> {
        println!("🔄 开始检查订阅更新...");

        for sub in &config.subscriptions {
            if let Err(e) = self.download_subscription(sub).await {
                eprintln!("❌ 处理订阅 [{}] 时发生错误: {}", sub.name, e);
            }
        }

        println!("✅ 订阅更新检查完成");
        Ok(())
    }
}

