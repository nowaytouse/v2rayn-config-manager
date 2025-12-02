use anyhow::{Context, Result};
use flate2::read::GzDecoder;
use std::path::{Path, PathBuf};
use tar::Archive;
use tokio::fs;

use crate::types::{Asset, Config, GithubRelease};

/// 核心更新器（纯CLI模式）
pub struct CoreUpdater;

impl CoreUpdater {
    pub fn new() -> Self {
        Self
    }

    /// 获取最新的 GitHub Release
    pub async fn get_latest_release(&self, check_prerelease: bool) -> Result<GithubRelease> {
        println!("🔍 正在从 GitHub API 获取最新版本信息...");

        let client = reqwest::Client::builder()
            .user_agent("singbox-manager/2.0")
            .build()?;

        let response = client
            .get("https://api.github.com/repos/SagerNet/sing-box/releases")
            .send()
            .await
            .context("访问 GitHub API 失败")?;

        if !response.status().is_success() {
            anyhow::bail!("访问 GitHub API 失败，状态码: {}", response.status());
        }

        let releases: Vec<GithubRelease> = response
            .json()
            .await
            .context("解析 GitHub API 响应失败")?;

        for release in releases {
            if !check_prerelease && release.prerelease {
                continue;
            }
            println!("✅ 找到符合条件的最新版本: {}", release.tag_name);
            return Ok(release);
        }

        anyhow::bail!("未找到符合条件的版本")
    }

    /// 查找匹配当前平台的资源
    pub fn find_matching_asset(&self, release: &GithubRelease) -> Option<Asset> {
        let os = std::env::consts::OS;
        let arch = match std::env::consts::ARCH {
            "x86_64" => "amd64",
            "aarch64" => "arm64",
            "x86" => "386",
            other => other,
        };

        let target = format!("{}-{}", os, arch);
        println!("🔍 正在为平台 {} 寻找匹配的资源...", target);

        for asset in &release.assets {
            if asset.name.contains(&target) {
                if asset.name.ends_with(".tar.gz") || asset.name.ends_with(".zip") {
                    println!("✅ 找到匹配资源: {}", asset.name);
                    return Some(asset.clone());
                }
            }
        }

        None
    }

    /// 下载并解压资源
    pub async fn download_and_extract(&self, asset: &Asset) -> Result<PathBuf> {
        println!("📥 正在下载资源: {}", asset.browser_download_url);

        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(300))
            .build()?;

        let response = client
            .get(&asset.browser_download_url)
            .send()
            .await
            .context("下载资源失败")?;

        if !response.status().is_success() {
            anyhow::bail!("下载资源失败，状态码: {}", response.status());
        }

        let bytes = response.bytes().await.context("读取响应内容失败")?;

        // 创建临时目录
        let temp_dir = std::env::temp_dir().join(format!("singbox-{}", uuid::Uuid::new_v4()));
        std::fs::create_dir_all(&temp_dir).context("创建临时目录失败")?;

        println!("📦 资源已下载，正在解压到临时目录: {}", temp_dir.display());

        // 根据文件扩展名解压
        let extracted_path = if asset.name.ends_with(".tar.gz") {
            self.extract_tar_gz(&bytes, &temp_dir).await?
        } else if asset.name.ends_with(".zip") {
            self.extract_zip(&bytes, &temp_dir).await?
        } else {
            anyhow::bail!("不支持的压缩格式");
        };

        Ok(extracted_path)
    }

    /// 解压 tar.gz 文件（同步操作，使用 spawn_blocking）
    async fn extract_tar_gz(&self, data: &[u8], dest_dir: &Path) -> Result<PathBuf> {
        let data = data.to_vec();
        let dest_dir = dest_dir.to_path_buf();

        let result = tokio::task::spawn_blocking(move || {
            let decoder = GzDecoder::new(&data[..]);
            let mut archive = Archive::new(decoder);

            let exe_name = if cfg!(windows) {
                "sing-box.exe"
            } else {
                "sing-box"
            };

            for entry in archive.entries()? {
                let mut entry = entry?;
                let path = entry.path()?;

                if path.file_name().and_then(|n| n.to_str()) == Some(exe_name)
                    || path.to_str().unwrap_or("").ends_with(exe_name)
                {
                    let extract_path = dest_dir.join(exe_name);
                    entry.unpack(&extract_path)?;

                    // Unix 系统设置执行权限
                    #[cfg(unix)]
                    {
                        use std::os::unix::fs::PermissionsExt;
                        let mut perms = std::fs::metadata(&extract_path)?.permissions();
                        perms.set_mode(0o755);
                        std::fs::set_permissions(&extract_path, perms)?;
                    }

                    return Ok(extract_path);
                }
            }

            anyhow::bail!("在压缩包中未找到 sing-box 可执行文件")
        })
        .await??;

        println!("✅ 已解压 sing-box 到 {}", result.display());
        Ok(result)
    }

    /// 解压 zip 文件（同步操作，使用 spawn_blocking）
    async fn extract_zip(&self, data: &[u8], dest_dir: &Path) -> Result<PathBuf> {
        let data = data.to_vec();
        let dest_dir = dest_dir.to_path_buf();

        let result = tokio::task::spawn_blocking(move || {
            let reader = std::io::Cursor::new(&data);
            let mut archive = zip::ZipArchive::new(reader)?;

            let exe_name = if cfg!(windows) {
                "sing-box.exe"
            } else {
                "sing-box"
            };

            for i in 0..archive.len() {
                let mut file = archive.by_index(i)?;
                let file_name = file.name().to_string();

                if file_name.ends_with(exe_name) {
                    let extract_path = dest_dir.join(exe_name);
                    let mut outfile = std::fs::File::create(&extract_path)?;
                    std::io::copy(&mut file, &mut outfile)?;

                    // Unix 系统设置执行权限
                    #[cfg(unix)]
                    {
                        use std::os::unix::fs::PermissionsExt;
                        let mut perms = std::fs::metadata(&extract_path)?.permissions();
                        perms.set_mode(0o755);
                        std::fs::set_permissions(&extract_path, perms)?;
                    }

                    return Ok(extract_path);
                }
            }

            anyhow::bail!("在压缩包中未找到 sing-box 可执行文件")
        })
        .await??;

        println!("✅ 已解压 sing-box 到 {}", result.display());
        Ok(result)
    }

    /// 安装文件
    pub async fn install_file(&self, source: &Path, dest: &Path) -> Result<()> {
        println!("📦 正在尝试将 {} 安装到 {}", source.display(), dest.display());

        // 确保目标目录存在
        if let Some(parent) = dest.parent() {
            fs::create_dir_all(parent)
                .await
                .context("创建目标目录失败")?;
        }

        // 复制文件
        fs::copy(source, dest).await.context("复制文件失败")?;

        // Unix 系统设置执行权限
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mut perms = fs::metadata(dest).await?.permissions();
            perms.set_mode(0o755);
            fs::set_permissions(dest, perms).await?;
        }

        println!("✅ 安装成功！sing-box 已更新至 {}", dest.display());
        Ok(())
    }

    /// 运行核心更新器
    pub async fn run_all(&self, config: &Config) -> Result<()> {
        if !config.singbox_core_update.enabled {
            println!("⚠️  sing-box 核心更新已禁用");
            return Ok(());
        }

        println!("🔄 开始检查 sing-box 核心更新...");

        let release = self
            .get_latest_release(config.singbox_core_update.check_prerelease)
            .await?;

        let asset = self
            .find_matching_asset(&release)
            .context("未找到与当前系统匹配的资源文件")?;

        let temp_binary_path = self.download_and_extract(&asset).await?;

        self.install_file(&temp_binary_path, &config.singbox_core_update.install_path)
            .await?;

        // 清理临时文件
        if let Some(temp_dir) = temp_binary_path.parent() {
            let _ = fs::remove_dir_all(temp_dir).await;
        }

        println!("✅ sing-box 核心更新完成");
        Ok(())
    }
}
