package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"singbox-manager/pkg/core"
	"singbox-manager/pkg/gui"
)

const version = "2.0.0"

func main() {
	// 命令行参数
	showVersion := flag.Bool("version", false, "显示版本信息")
	showHelp := flag.Bool("help", false, "显示帮助信息")
	guiMode := flag.Bool("gui", false, "启动 GUI 模式")
	configPath := flag.String("config", "config.json", "配置文件路径")
	onceMode := flag.Bool("once", false, "仅执行一次更新任务（命令行模式）")

	flag.Parse()

	// 显示版本
	if *showVersion {
		fmt.Printf("Sing-box Manager v%s\n", version)
		fmt.Printf("Go version: %s\n", "1.21+")
		fmt.Printf("Platform: %s/%s\n", "跨平台", "支持 Windows/Linux/macOS")
		return
	}

	// 显示帮助
	if *showHelp {
		printHelp()
		return
	}

	// GUI 模式
	if *guiMode {
		g := gui.NewGUI()
		g.Run()
		return
	}

	// 命令行模式
	runCLI(*configPath, *onceMode)
}

func printHelp() {
	fmt.Println("Sing-box Manager - 配置和核心自动更新工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  singbox-manager [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -gui               启动图形界面模式")
	fmt.Println("  -config <path>     指定配置文件路径 (默认: config.json)")
	fmt.Println("  -once              仅执行一次更新任务后退出（命令行模式）")
	fmt.Println("  -version           显示版本信息")
	fmt.Println("  -help              显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  singbox-manager -gui                    # 启动 GUI 模式")
	fmt.Println("  singbox-manager                          # 命令行模式，按配置定时更新")
	fmt.Println("  singbox-manager -once                    # 仅执行一次更新")
	fmt.Println("  singbox-manager -config custom.json      # 使用自定义配置文件")
	fmt.Println()
	fmt.Println("配置文件示例请参考 config.json")
}

func runCLI(configPath string, onceMode bool) {
	log.SetFlags(log.LstdFlags)
	log.Printf("程序启动 (命令行模式)，正在加载配置文件: %s", configPath)

	// 加载配置
	config, err := core.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 创建日志回调
	logCallback := func(message string) {
		log.Println(message)
	}

	// 任务执行函数
	runTasks := func() {
		// 更新订阅
		sm := core.NewSubscriptionManager(logCallback)
		sm.RunSubscriptionDownloader(config)

		// 更新核心
		cu := core.NewCoreUpdater(logCallback)
		if err := cu.RunCoreUpdater(config); err != nil {
			log.Printf("核心更新出错: %v", err)
		}
	}

	// 仅执行一次
	if onceMode || config.UpdateIntervalHours <= 0 {
		log.Println("执行一次性更新任务...")
		runTasks()
		log.Println("所有任务完成，程序退出。")
		return
	}

	// 定时执行
	log.Println("立即执行第一次任务。")
	runTasks()

	interval := time.Duration(config.UpdateIntervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("已设置定时任务，每 %d 小时更新一次。按 Ctrl+C 停止程序。", config.UpdateIntervalHours)

	// 捕获系统信号以优雅退出
	for {
		select {
		case <-ticker.C:
			log.Println("定时任务触发，开始更新...")
			runTasks()
		}
	}
}
