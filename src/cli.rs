use anyhow::Result;
use std::path::PathBuf;
use std::time::Duration;
use tokio::time::sleep;

use crate::config::load_config;
use crate::subscription::SubscriptionManager;
use crate::types::Config;
use crate::updater::CoreUpdater;

/// CLI 模式运行
pub async fn run_cli(config_path: PathBuf, once_mode: bool) -> Result<()> {
    println!("🚀 程序启动 (命令行模式)，正在加载配置文件: {}", config_path.display());

    // 加载配置
    let config = load_config(&config_path).await?;

    // 执行任务
    run_tasks(&config).await?;

    // 如果是一次性模式或更新间隔为0，退出
    if once_mode || config.update_interval_hours == 0 {
        println!("✅ 所有任务完成，程序退出。");
        return Ok(());
    }

    // 定时执行
    println!("⏰ 已设置定时任务，每 {} 小时更新一次。按 Ctrl+C 停止程序。", config.update_interval_hours);

    let interval = Duration::from_secs(config.update_interval_hours * 3600);

    loop {
        sleep(interval).await;
        println!("⏰ 定时任务触发，开始更新...");
        if let Err(e) = run_tasks(&config).await {
            eprintln!("❌ 任务执行出错: {}", e);
        }
    }
}

/// 执行更新任务
async fn run_tasks(config: &Config) -> Result<()> {
    // 更新订阅
    let sub_manager = SubscriptionManager::new();
    sub_manager.run_all(config).await?;

    // 更新核心
    let core_updater = CoreUpdater::new();
    core_updater.run_all(config).await?;

    Ok(())
}
