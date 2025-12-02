mod cli;
mod config;
mod interactive;
mod subscription;
mod types;
mod updater;

use clap::Parser;
use std::path::PathBuf;

const VERSION: &str = env!("CARGO_PKG_VERSION");

/// Sing-box 配置和核心自动更新工具 - 纯 Rust CLI 版本
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// 交互式菜单模式
    #[arg(short, long)]
    interactive: bool,

    /// 配置文件路径
    #[arg(short, long, default_value = "config.json")]
    config: PathBuf,

    /// 仅执行一次更新任务后退出（非交互模式）
    #[arg(short, long)]
    once: bool,

    /// 显示版本信息
    #[arg(short, long)]
    version: bool,
}

#[tokio::main]
async fn main() {
    let args = Args::parse();

    // 显示版本信息
    if args.version {
        println!("Sing-box Manager v{}", VERSION);
        println!("纯 Rust CLI 版本");
        println!("平台: {}/{}", std::env::consts::OS, std::env::consts::ARCH);
        return;
    }

    // 初始化日志
    env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info"))
        .init();

    // 交互式模式
    if args.interactive {
        let cli = crate::interactive::InteractiveCLI::new(args.config);
        if let Err(e) = cli.run().await {
            eprintln!("程序运行出错: {}", e);
            std::process::exit(1);
        }
        return;
    }

    // 标准 CLI 模式
    if let Err(e) = cli::run_cli(args.config, args.once).await {
        eprintln!("程序运行出错: {}", e);
        std::process::exit(1);
    }
}

