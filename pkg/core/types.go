package core

// Config 结构对应 config.json 的结构
type Config struct {
	Subscriptions       []Subscription    `json:"subscriptions"`
	UpdateIntervalHours int               `json:"update_interval_hours"`
	SingboxCoreUpdate   SingboxCoreUpdate `json:"singbox_core_update"`
}

// Subscription 定义了单个订阅源
type Subscription struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	SavePath string `json:"save_path"`
}

// SingboxCoreUpdate 定义了 sing-box 核心更新的配置
type SingboxCoreUpdate struct {
	Enabled         bool   `json:"enabled"`
	CheckPrerelease bool   `json:"check_prerelease"`
	InstallPath     string `json:"install_path"`
}

// GithubRelease 定义了 GitHub Release API 响应的结构
type GithubRelease struct {
	URL        string  `json:"url"`
	TagName    string  `json:"tag_name"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
}

// Asset 定义了 Release 中包含的资源文件
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// LogCallback 日志回调函数类型
type LogCallback func(message string)

