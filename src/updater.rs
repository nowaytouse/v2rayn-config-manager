use anyhow::{Context, Result};
use flate2::read::GzDecoder;
use std::path::{Path, PathBuf};
use tar::Archive;
use tokio::fs;

use crate::types::{Asset, Config, GithubRelease};

/// Core Updater (Pure CLI Mode)
pub struct CoreUpdater;

impl CoreUpdater {
    pub fn new() -> Self {
        Self
    }

    /// Get the latest GitHub Release
    pub async fn get_latest_release(&self, check_prerelease: bool) -> Result<GithubRelease> {
        println!("🔍 Fetching latest version info from GitHub API...");

        let client = reqwest::Client::builder()
            .user_agent("singbox-manager/2.0")
            .build()?;

        let response = client
            .get("https://api.github.com/repos/SagerNet/sing-box/releases")
            .send()
            .await
            .context("Failed to access GitHub API")?;

        if !response.status().is_success() {
            anyhow::bail!("Failed to access GitHub API, status code: {}", response.status());
        }

        let releases: Vec<GithubRelease> = response
            .json()
            .await
            .context("Failed to parse GitHub API response")?;

        for release in releases {
            if !check_prerelease && release.prerelease {
                continue;
            }
            println!("✅ Found latest matching version: {}", release.tag_name);
            return Ok(release);
        }

        anyhow::bail!("No matching version found")
    }

    /// Find matching asset for current platform
    pub fn find_matching_asset(&self, release: &GithubRelease) -> Option<Asset> {
        let os = std::env::consts::OS;
        let arch = match std::env::consts::ARCH {
            "x86_64" => "amd64",
            "aarch64" => "arm64",
            "x86" => "386",
            other => other,
        };

        let target = format!("{}-{}", os, arch);
        println!("🔍 Looking for matching asset for platform {}...", target);

        for asset in &release.assets {
            if asset.name.contains(&target) {
                if asset.name.ends_with(".tar.gz") || asset.name.ends_with(".zip") {
                    println!("✅ Found matching asset: {}", asset.name);
                    return Some(asset.clone());
                }
            }
        }

        None
    }

    /// Download and extract asset
    pub async fn download_and_extract(&self, asset: &Asset) -> Result<PathBuf> {
        println!("📥 Downloading asset: {}", asset.browser_download_url);

        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(300))
            .build()?;

        let response = client
            .get(&asset.browser_download_url)
            .send()
            .await
            .context("Failed to download asset")?;

        if !response.status().is_success() {
            anyhow::bail!("Failed to download asset, status code: {}", response.status());
        }

        let bytes = response.bytes().await.context("Failed to read response content")?;

        // Create temporary directory
        let temp_dir = std::env::temp_dir().join(format!("singbox-{}", uuid::Uuid::new_v4()));
        std::fs::create_dir_all(&temp_dir).context("Failed to create temporary directory")?;

        println!("📦 Asset downloaded, extracting to temporary directory: {}", temp_dir.display());

        // Extract based on file extension
        let extracted_path = if asset.name.ends_with(".tar.gz") {
            self.extract_tar_gz(&bytes, &temp_dir).await?
        } else if asset.name.ends_with(".zip") {
            self.extract_zip(&bytes, &temp_dir).await?
        } else {
            anyhow::bail!("Unsupported compression format");
        };

        Ok(extracted_path)
    }

    /// Extract tar.gz file (synchronous operation, using spawn_blocking)
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

                    // Set execute permissions on Unix systems
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

            anyhow::bail!("sing-box executable not found in archive")
        })
        .await??;

        println!("✅ Extracted sing-box to {}", result.display());
        Ok(result)
    }

    /// Extract zip file (synchronous operation, using spawn_blocking)
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

                    // Set execute permissions on Unix systems
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

            anyhow::bail!("sing-box executable not found in archive")
        })
        .await??;

        println!("✅ Extracted sing-box to {}", result.display());
        Ok(result)
    }

    /// Install file
    pub async fn install_file(&self, source: &Path, dest: &Path) -> Result<()> {
        println!("📦 Installing {} to {}", source.display(), dest.display());

        // Ensure target directory exists
        if let Some(parent) = dest.parent() {
            fs::create_dir_all(parent)
                .await
                .context("Failed to create target directory")?;
        }

        // Copy file
        fs::copy(source, dest).await.context("Failed to copy file")?;

        // Set execute permissions on Unix systems
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mut perms = fs::metadata(dest).await?.permissions();
            perms.set_mode(0o755);
            fs::set_permissions(dest, perms).await?;
        }

        println!("✅ Installation successful! sing-box updated to {}", dest.display());
        Ok(())
    }

    /// Run core updater
    pub async fn run_all(&self, config: &Config) -> Result<()> {
        if !config.singbox_core_update.enabled {
            println!("⚠️  sing-box core update is disabled");
            return Ok(());
        }

        println!("🔄 Checking sing-box core updates...");

        let release = self
            .get_latest_release(config.singbox_core_update.check_prerelease)
            .await?;

        let asset = self
            .find_matching_asset(&release)
            .context("No matching asset found for current system")?;

        let temp_binary_path = self.download_and_extract(&asset).await?;

        self.install_file(&temp_binary_path, &config.singbox_core_update.install_path)
            .await?;

        // Clean up temporary files
        if let Some(temp_dir) = temp_binary_path.parent() {
            let _ = fs::remove_dir_all(temp_dir).await;
        }

        println!("✅ sing-box core update complete");
        Ok(())
    }
}
