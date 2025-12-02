package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// CoreUpdater sing-box 核心更新器
type CoreUpdater struct {
	logCallback LogCallback
}

// NewCoreUpdater 创建核心更新器
func NewCoreUpdater(logCallback LogCallback) *CoreUpdater {
	return &CoreUpdater{
		logCallback: logCallback,
	}
}

func (cu *CoreUpdater) log(message string) {
	if cu.logCallback != nil {
		cu.logCallback(message)
	}
}

// GetLatestRelease 获取最新版本
func (cu *CoreUpdater) GetLatestRelease(checkPrerelease bool) (*GithubRelease, error) {
	cu.log("正在从 GitHub API 获取最新版本信息...")
	resp, err := http.Get("https://api.github.com/repos/SagerNet/sing-box/releases")
	if err != nil {
		return nil, fmt.Errorf("访问 GitHub API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("访问 GitHub API 失败，状态码: %d", resp.StatusCode)
	}

	var releases []GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("解析 GitHub API 响应失败: %w", err)
	}

	for _, release := range releases {
		if !checkPrerelease && release.Prerelease {
			continue // 跳过测试版
		}
		cu.log(fmt.Sprintf("找到符合条件的最新版本: %s", release.TagName))
		return &release, nil // 返回第一个符合条件的版本
	}

	return nil, fmt.Errorf("未找到符合条件的版本")
}

// FindMatchingAsset 查找匹配的资源
func (cu *CoreUpdater) FindMatchingAsset(release *GithubRelease) *Asset {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	// 修正 arm64 到 aarch64 的名称差异
	if arch == "arm64" {
		arch = "aarch64"
	}

	// Windows 平台特殊处理
	var extensions []string
	if goos == "windows" {
		extensions = []string{".zip"}
	} else {
		extensions = []string{".tar.gz", ".zip"}
	}

	target := fmt.Sprintf("%s-%s", goos, arch)
	cu.log(fmt.Sprintf("正在为平台 %s 寻找匹配的资源...", target))

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, target) {
			for _, ext := range extensions {
				if strings.HasSuffix(asset.Name, ext) {
					cu.log(fmt.Sprintf("找到匹配资源: %s", asset.Name))
					return &asset
				}
			}
		}
	}
	return nil
}

// DownloadAndExtract 下载并解压
func (cu *CoreUpdater) DownloadAndExtract(asset *Asset) (string, error) {
	cu.log(fmt.Sprintf("正在下载资源: %s", asset.BrowserDownloadURL))
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载资源失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "singbox-")
	if err != nil {
		return "", err
	}

	cu.log(fmt.Sprintf("资源已下载，正在解压到临时目录: %s", tmpDir))

	if strings.HasSuffix(asset.Name, ".zip") {
		return cu.extractZip(body, tmpDir)
	} else if strings.HasSuffix(asset.Name, ".tar.gz") {
		return cu.extractTarGz(body, tmpDir)
	}

	return "", fmt.Errorf("不支持的压缩格式")
}

func (cu *CoreUpdater) extractZip(data []byte, destDir string) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	for _, f := range zipReader.File {
		// 查找 sing-box 可执行文件（支持 Windows 的 .exe）
		if strings.HasSuffix(f.Name, "sing-box") || strings.HasSuffix(f.Name, "sing-box.exe") {
			return cu.extractZipFile(f, destDir)
		}
	}

	return "", fmt.Errorf("在压缩包中未找到 sing-box 可执行文件")
}

func (cu *CoreUpdater) extractTarGz(data []byte, destDir string) (string, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(header.Name, "sing-box") {
			return cu.extractTarFile(tarReader, header, destDir)
		}
	}

	return "", fmt.Errorf("在压缩包中未找到 sing-box 可执行文件")
}

func (cu *CoreUpdater) extractTarFile(reader *tar.Reader, header *tar.Header, destDir string) (string, error) {
	fileName := "sing-box"
	if runtime.GOOS == "windows" {
		fileName = "sing-box.exe"
	}
	filePath := fmt.Sprintf("%s/%s", destDir, fileName)
	
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	
	if _, err := io.Copy(outFile, reader); err != nil {
		return "", err
	}
	
	// Unix 系统设置执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filePath, 0755); err != nil {
			return "", err
		}
	}
	
	cu.log(fmt.Sprintf("已解压 sing-box 到 %s", filePath))
	return filePath, nil
}

func (cu *CoreUpdater) extractZipFile(f *zip.File, destDir string) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	fileName := "sing-box"
	if runtime.GOOS == "windows" {
		fileName = "sing-box.exe"
	}
	filePath := fmt.Sprintf("%s/%s", destDir, fileName)
	
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return "", err
	}
	
	// Unix 系统设置执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filePath, 0755); err != nil {
			return "", err
		}
	}
	
	cu.log(fmt.Sprintf("已解压 sing-box 到 %s", filePath))
	return filePath, nil
}

// InstallFile 安装文件
func (cu *CoreUpdater) InstallFile(sourcePath, destPath string) error {
	cu.log(fmt.Sprintf("正在尝试将 %s 安装到 %s", sourcePath, destPath))

	// 确保目标目录存在
	if err := os.MkdirAll(getDir(destPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w (提示: 可能需要管理员权限)", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// Unix 系统设置执行权限
	if runtime.GOOS != "windows" {
		err = os.Chmod(destPath, 0755)
		if err != nil {
			return fmt.Errorf("设置执行权限失败: %w", err)
		}
	}

	cu.log(fmt.Sprintf("安装成功！ sing-box 已更新至 %s", destPath))
	return nil
}

// RunCoreUpdater 运行核心更新器
func (cu *CoreUpdater) RunCoreUpdater(config *Config) error {
	if !config.SingboxCoreUpdate.Enabled {
		cu.log("sing-box 核心更新已禁用")
		return nil
	}
	
	cu.log("--- 开始检查 sing-box 核心更新 ---")

	release, err := cu.GetLatestRelease(config.SingboxCoreUpdate.CheckPrerelease)
	if err != nil {
		return fmt.Errorf("获取版本信息失败: %w", err)
	}

	asset := cu.FindMatchingAsset(release)
	if asset == nil {
		return fmt.Errorf("未找到与当前系统匹配的资源文件")
	}

	tempBinaryPath, err := cu.DownloadAndExtract(asset)
	if err != nil {
		return fmt.Errorf("下载或解压失败: %w", err)
	}
	defer os.RemoveAll(getDir(tempBinaryPath)) // Clean up temp directory

	if err := cu.InstallFile(tempBinaryPath, config.SingboxCoreUpdate.InstallPath); err != nil {
		return fmt.Errorf("安装 sing-box 核心失败: %w", err)
	}

	cu.log("--- sing-box 核心更新完成 ---")
	return nil
}

