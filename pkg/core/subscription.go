package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// SubscriptionManager 订阅管理器
type SubscriptionManager struct {
	logCallback LogCallback
}

// NewSubscriptionManager 创建订阅管理器
func NewSubscriptionManager(logCallback LogCallback) *SubscriptionManager {
	return &SubscriptionManager{
		logCallback: logCallback,
	}
}

func (sm *SubscriptionManager) log(message string) {
	if sm.logCallback != nil {
		sm.logCallback(message)
	}
}

// DownloadSubscription 下载单个订阅
func (sm *SubscriptionManager) DownloadSubscription(sub Subscription) error {
	sm.log(fmt.Sprintf("开始下载订阅 [%s] 从 %s", sub.Name, sub.URL))
	
	resp, err := http.Get(sub.URL)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，服务器返回状态码: %d", resp.StatusCode)
	}

	// 确保保存路径的目录存在
	if err := os.MkdirAll(getDir(sub.SavePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	out, err := os.Create(sub.SavePath)
	if err != nil {
		return fmt.Errorf("创建文件 %s 失败: %w", sub.SavePath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("保存文件 %s 失败: %w", sub.SavePath, err)
	}

	sm.log(fmt.Sprintf("订阅 [%s] 下载成功，已保存至 %s", sub.Name, sub.SavePath))
	return nil
}

// RunSubscriptionDownloader 运行订阅下载器
func (sm *SubscriptionManager) RunSubscriptionDownloader(config *Config) {
	sm.log("--- 开始检查订阅更新 ---")
	for _, sub := range config.Subscriptions {
		if err := sm.DownloadSubscription(sub); err != nil {
			sm.log(fmt.Sprintf("处理订阅 [%s] 时发生错误: %v", sub.Name, err))
		}
	}
	sm.log("--- 订阅更新检查完成 ---")
}

// getDir 获取文件路径的目录部分
func getDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}

