package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"singbox-manager/pkg/core"
)

// GUI 主界面
type GUI struct {
	app             fyne.App
	window          fyne.Window
	config          *core.Config
	configPath      string
	logText         *widget.Entry
	statusLabel     *widget.Label
	autoUpdateCheck *widget.Check
	ticker          *time.Ticker
	isRunning       bool
}

// NewGUI 创建新的 GUI 实例
func NewGUI() *GUI {
	return &GUI{
		app:        app.New(),
		configPath: "config.json",
		isRunning:  false,
	}
}

// Run 启动 GUI
func (g *GUI) Run() {
	g.window = g.app.NewWindow("Sing-box 配置管理器")
	g.window.Resize(fyne.NewSize(800, 600))

	// 加载配置
	if err := g.loadConfig(); err != nil {
		dialog.ShowError(fmt.Errorf("加载配置失败: %w", err), g.window)
	}

	// 创建界面
	content := g.createMainContent()
	g.window.SetContent(content)

	g.window.ShowAndRun()
}

func (g *GUI) createMainContent() fyne.CanvasObject {
	// 顶部状态栏
	g.statusLabel = widget.NewLabel("就绪")

	// 日志显示区域
	g.logText = widget.NewMultiLineEntry()
	g.logText.SetPlaceHolder("日志信息将在这里显示...")
	g.logText.Disable()
	logScroll := container.NewScroll(g.logText)
	logScroll.SetMinSize(fyne.NewSize(750, 200))

	// 创建订阅列表
	subscriptionBox := g.createSubscriptionBox()

	// 创建核心更新配置
	coreUpdateBox := g.createCoreUpdateBox()

	// 创建定时更新配置
	intervalBox := g.createIntervalBox()

	// 操作按钮
	actionButtons := g.createActionButtons()

	// 组合所有内容
	mainContent := container.NewBorder(
		container.NewVBox(
			widget.NewCard("", "", g.statusLabel),
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			actionButtons,
		),
		nil,
		nil,
		container.NewVScroll(
			container.NewVBox(
				subscriptionBox,
				widget.NewSeparator(),
				coreUpdateBox,
				widget.NewSeparator(),
				intervalBox,
				widget.NewSeparator(),
				widget.NewCard("日志", "", logScroll),
			),
		),
	)

	return mainContent
}

func (g *GUI) createSubscriptionBox() fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			if g.config == nil {
				return 0
			}
			return len(g.config.Subscriptions)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("模板"),
				layout.NewSpacer(),
				widget.NewButton("编辑", func() {}),
				widget.NewButton("删除", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if g.config == nil || id >= len(g.config.Subscriptions) {
				return
			}
			sub := g.config.Subscriptions[id]
			box := item.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s - %s", sub.Name, sub.URL))

			// 更新编辑按钮
			box.Objects[2].(*widget.Button).OnTapped = func() {
				g.editSubscription(id)
			}

			// 更新删除按钮
			box.Objects[3].(*widget.Button).OnTapped = func() {
				g.deleteSubscription(id)
			}
		},
	)

	addButton := widget.NewButton("添加订阅", func() {
		g.addSubscription()
	})

	return widget.NewCard("订阅管理", "", container.NewBorder(
		nil,
		addButton,
		nil,
		nil,
		list,
	))
}

func (g *GUI) createCoreUpdateBox() fyne.CanvasObject {
	if g.config == nil {
		return container.NewVBox()
	}

	enabledCheck := widget.NewCheck("启用核心自动更新", func(checked bool) {
		g.config.SingboxCoreUpdate.Enabled = checked
		g.saveConfig()
	})
	enabledCheck.Checked = g.config.SingboxCoreUpdate.Enabled

	prereleaseCheck := widget.NewCheck("包含预发布版本", func(checked bool) {
		g.config.SingboxCoreUpdate.CheckPrerelease = checked
		g.saveConfig()
	})
	prereleaseCheck.Checked = g.config.SingboxCoreUpdate.CheckPrerelease

	installPathEntry := widget.NewEntry()
	installPathEntry.SetText(g.config.SingboxCoreUpdate.InstallPath)
	installPathEntry.OnChanged = func(text string) {
		g.config.SingboxCoreUpdate.InstallPath = text
		g.saveConfig()
	}

	browseButton := widget.NewButton("浏览...", func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				path := file.URI().Path()
				installPathEntry.SetText(path)
				g.config.SingboxCoreUpdate.InstallPath = path
				g.saveConfig()
				file.Close()
			}
		}, g.window)
	})

	form := container.NewVBox(
		enabledCheck,
		prereleaseCheck,
		widget.NewLabel("安装路径:"),
		container.NewBorder(nil, nil, nil, browseButton, installPathEntry),
	)

	return widget.NewCard("Sing-box 核心更新", "", form)
}

func (g *GUI) createIntervalBox() fyne.CanvasObject {
	if g.config == nil {
		return container.NewVBox()
	}

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%d", g.config.UpdateIntervalHours))
	intervalEntry.OnChanged = func(text string) {
		var hours int
		fmt.Sscanf(text, "%d", &hours)
		if hours >= 0 {
			g.config.UpdateIntervalHours = hours
			g.saveConfig()
		}
	}

	g.autoUpdateCheck = widget.NewCheck("启用自动定时更新", func(checked bool) {
		if checked && g.config.UpdateIntervalHours > 0 {
			g.startAutoUpdate()
		} else {
			g.stopAutoUpdate()
		}
	})

	form := container.NewVBox(
		widget.NewLabel("更新间隔 (小时):"),
		intervalEntry,
		widget.NewLabel("(设置为 0 表示仅手动更新)"),
		g.autoUpdateCheck,
	)

	return widget.NewCard("定时更新设置", "", form)
}

func (g *GUI) createActionButtons() fyne.CanvasObject {
	updateSubBtn := widget.NewButton("立即更新订阅", func() {
		g.updateSubscriptions()
	})

	updateCoreBtn := widget.NewButton("立即更新核心", func() {
		g.updateCore()
	})

	updateAllBtn := widget.NewButton("全部更新", func() {
		g.updateAll()
	})

	clearLogBtn := widget.NewButton("清空日志", func() {
		g.logText.SetText("")
	})

	return container.NewHBox(
		updateSubBtn,
		updateCoreBtn,
		updateAllBtn,
		layout.NewSpacer(),
		clearLogBtn,
	)
}

// 订阅相关操作
func (g *GUI) addSubscription() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("订阅名称")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("订阅 URL")

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("保存路径")
	pathEntry.SetText("./singbox_config.json")

	form := dialog.NewForm("添加订阅", "添加", "取消", []*widget.FormItem{
		widget.NewFormItem("名称", nameEntry),
		widget.NewFormItem("URL", urlEntry),
		widget.NewFormItem("保存路径", pathEntry),
	}, func(submitted bool) {
		if submitted {
			g.config.Subscriptions = append(g.config.Subscriptions, core.Subscription{
				Name:     nameEntry.Text,
				URL:      urlEntry.Text,
				SavePath: pathEntry.Text,
			})
			g.saveConfig()
			g.window.SetContent(g.createMainContent())
		}
	}, g.window)

	form.Resize(fyne.NewSize(500, 300))
	form.Show()
}

func (g *GUI) editSubscription(id int) {
	if id >= len(g.config.Subscriptions) {
		return
	}

	sub := g.config.Subscriptions[id]

	nameEntry := widget.NewEntry()
	nameEntry.SetText(sub.Name)

	urlEntry := widget.NewEntry()
	urlEntry.SetText(sub.URL)

	pathEntry := widget.NewEntry()
	pathEntry.SetText(sub.SavePath)

	form := dialog.NewForm("编辑订阅", "保存", "取消", []*widget.FormItem{
		widget.NewFormItem("名称", nameEntry),
		widget.NewFormItem("URL", urlEntry),
		widget.NewFormItem("保存路径", pathEntry),
	}, func(submitted bool) {
		if submitted {
			g.config.Subscriptions[id] = core.Subscription{
				Name:     nameEntry.Text,
				URL:      urlEntry.Text,
				SavePath: pathEntry.Text,
			}
			g.saveConfig()
			g.window.SetContent(g.createMainContent())
		}
	}, g.window)

	form.Resize(fyne.NewSize(500, 300))
	form.Show()
}

func (g *GUI) deleteSubscription(id int) {
	if id >= len(g.config.Subscriptions) {
		return
	}

	dialog.ShowConfirm("确认删除", fmt.Sprintf("确定要删除订阅 [%s] 吗？", g.config.Subscriptions[id].Name),
		func(confirm bool) {
			if confirm {
				g.config.Subscriptions = append(g.config.Subscriptions[:id], g.config.Subscriptions[id+1:]...)
				g.saveConfig()
				g.window.SetContent(g.createMainContent())
			}
		}, g.window)
}

// 更新操作
func (g *GUI) updateSubscriptions() {
	g.setStatus("正在更新订阅...")
	go func() {
		sm := core.NewSubscriptionManager(g.log)
		sm.RunSubscriptionDownloader(g.config)
		g.setStatus("订阅更新完成")
	}()
}

func (g *GUI) updateCore() {
	g.setStatus("正在更新 Sing-box 核心...")
	go func() {
		cu := core.NewCoreUpdater(g.log)
		if err := cu.RunCoreUpdater(g.config); err != nil {
			g.log(fmt.Sprintf("核心更新失败: %v", err))
			g.setStatus("核心更新失败")
		} else {
			g.setStatus("核心更新完成")
		}
	}()
}

func (g *GUI) updateAll() {
	g.setStatus("正在执行全部更新...")
	go func() {
		g.updateSubscriptions()
		time.Sleep(time.Second)
		g.updateCore()
		g.setStatus("全部更新完成")
	}()
}

// 自动更新
func (g *GUI) startAutoUpdate() {
	if g.isRunning {
		return
	}

	if g.config.UpdateIntervalHours <= 0 {
		g.log("警告: 更新间隔必须大于 0")
		g.autoUpdateCheck.Checked = false
		g.autoUpdateCheck.Refresh()
		return
	}

	g.isRunning = true
	interval := time.Duration(g.config.UpdateIntervalHours) * time.Hour
	g.ticker = time.NewTicker(interval)

	g.log(fmt.Sprintf("自动更新已启动，间隔: %d 小时", g.config.UpdateIntervalHours))

	go func() {
		for range g.ticker.C {
			g.log("定时任务触发...")
			g.updateAll()
		}
	}()
}

func (g *GUI) stopAutoUpdate() {
	if !g.isRunning {
		return
	}

	if g.ticker != nil {
		g.ticker.Stop()
		g.ticker = nil
	}
	g.isRunning = false
	g.log("自动更新已停止")
}

// 配置操作
func (g *GUI) loadConfig() error {
	config, err := core.LoadConfig(g.configPath)
	if err != nil {
		// 如果配置文件不存在，创建默认配置
		if os.IsNotExist(err) {
			g.config = g.createDefaultConfig()
			return g.saveConfig()
		}
		return err
	}
	g.config = config
	return nil
}

func (g *GUI) saveConfig() error {
	if g.config == nil {
		return fmt.Errorf("配置为空")
	}
	return core.SaveConfig(g.configPath, g.config)
}

func (g *GUI) createDefaultConfig() *core.Config {
	// 根据系统类型设置默认安装路径
	var defaultInstallPath string
	switch runtime.GOOS {
	case "windows":
		defaultInstallPath = filepath.Join(os.Getenv("PROGRAMFILES"), "sing-box", "sing-box.exe")
	case "darwin":
		defaultInstallPath = "/usr/local/bin/sing-box"
	default: // linux
		defaultInstallPath = "/usr/local/bin/sing-box"
	}

	return &core.Config{
		Subscriptions: []core.Subscription{
			{
				Name:     "默认订阅",
				URL:      "https://example.com/your-singbox-config.json",
				SavePath: "./singbox_config.json",
			},
		},
		UpdateIntervalHours: 24,
		SingboxCoreUpdate: core.SingboxCoreUpdate{
			Enabled:         true,
			CheckPrerelease: false,
			InstallPath:     defaultInstallPath,
		},
	}
}

// 日志和状态
func (g *GUI) log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)
	g.logText.SetText(g.logText.Text + logMessage)
	
	// 自动滚动到底部（通过将光标移到末尾）
	g.logText.CursorRow = len(g.logText.Text)
}

func (g *GUI) setStatus(status string) {
	g.statusLabel.SetText(status)
	g.log(status)
}

