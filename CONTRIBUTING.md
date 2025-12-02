# 贡献指南

感谢你考虑为 Sing-box Manager 做出贡献！

## 🤝 如何贡献

### 报告 Bug

如果你发现了 bug，请创建一个 Issue，并包含以下信息：

1. **Bug 描述**：清晰简洁的问题描述
2. **复现步骤**：详细的重现步骤
3. **预期行为**：你期望发生什么
4. **实际行为**：实际发生了什么
5. **环境信息**：
   - 操作系统（Windows/Linux/macOS）
   - 版本号
   - Go 版本（如果从源码运行）
6. **日志输出**：相关的错误日志
7. **截图**：如果适用

### 提出新功能

如果你有新功能建议：

1. 先查看 [Issues](../../issues) 确认是否已有类似建议
2. 创建新 Issue，使用 "Feature Request" 标签
3. 详细描述功能和使用场景
4. 解释为什么这个功能有用

### 提交代码

#### 准备工作

1. Fork 本仓库
2. 克隆你的 fork：
   ```bash
   git clone https://github.com/your-username/singbox-manager.git
   cd singbox-manager
   ```
3. 创建新分支：
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### 开发流程

1. **保持代码风格一致**
   ```bash
   # 格式化代码
   go fmt ./...
   
   # 检查代码质量
   go vet ./...
   ```

2. **编写测试**
   - 为新功能添加单元测试
   - 确保所有测试通过：
     ```bash
     go test -v ./...
     ```

3. **提交代码**
   - 使用清晰的提交信息
   - 提交信息格式：
     ```
     类型(范围): 简短描述
     
     详细描述（可选）
     
     关闭 #issue-number（如果适用）
     ```
   - 提交类型：
     - `feat`: 新功能
     - `fix`: Bug 修复
     - `docs`: 文档更新
     - `style`: 代码格式调整
     - `refactor`: 代码重构
     - `test`: 测试相关
     - `chore`: 构建/工具相关

4. **推送到你的 fork**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 描述你的更改
   - 链接相关的 Issue
   - 等待代码审查

#### Pull Request 检查清单

- [ ] 代码已格式化 (`go fmt`)
- [ ] 代码通过静态检查 (`go vet`)
- [ ] 添加了必要的测试
- [ ] 所有测试通过
- [ ] 更新了相关文档
- [ ] 提交信息清晰明确
- [ ] PR 描述完整

## 📝 代码规范

### Go 代码风格

遵循官方 [Effective Go](https://golang.org/doc/effective_go.html) 指南：

1. **命名**：
   - 包名：小写，简短，有意义
   - 导出函数/类型：PascalCase
   - 私有函数/类型：camelCase
   - 常量：PascalCase 或 UPPER_CASE

2. **注释**：
   - 所有导出的类型、函数都要有注释
   - 注释以名称开头
   ```go
   // LoadConfig 从指定路径加载配置文件
   func LoadConfig(path string) (*Config, error) {
   ```

3. **错误处理**：
   - 始终检查错误
   - 使用 `fmt.Errorf` 包装错误添加上下文
   ```go
   if err != nil {
       return fmt.Errorf("加载配置失败: %w", err)
   }
   ```

4. **函数设计**：
   - 保持函数简短
   - 单一职责
   - 最多 3-4 个参数

### 项目结构

```
singbox-manager/
├── main.go              # 入口文件，保持简洁
├── pkg/
│   ├── core/           # 核心业务逻辑
│   │   ├── types.go    # 数据结构定义
│   │   ├── *.go        # 功能模块
│   └── gui/            # GUI 相关代码
│       └── gui.go      # GUI 实现
└── docs/               # 额外文档（如果需要）
```

### 新增功能指南

1. **在 `pkg/core` 添加业务逻辑**
   - 创建新文件或扩展现有文件
   - 添加必要的测试

2. **在 `pkg/gui` 更新界面**
   - 更新 GUI 组件
   - 确保界面响应

3. **更新 `main.go`**
   - 如果需要新的命令行参数

4. **更新文档**
   - README.md
   - 相关的其他文档

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行特定包的测试
go test -v ./pkg/core

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 编写测试

```go
// subscription_test.go
package core

import "testing"

func TestDownloadSubscription(t *testing.T) {
    // 测试逻辑
}
```

## 📚 文档

如果你的更改影响到用户使用，请更新：

- `README.md` - 主要文档
- `QUICK_START.md` - 快速开始指南
- `ARCHITECTURE.md` - 架构说明（如果适用）
- `CHANGELOG.md` - 更新日志

## 🐛 调试

### 启用详细日志

```bash
# 设置日志级别
export LOG_LEVEL=debug
./singbox-manager -gui
```

### 使用调试器

```bash
# 使用 delve
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug
```

## 💡 开发技巧

1. **使用 Make 命令**
   ```bash
   make run-gui    # 快速测试 GUI
   make run        # 测试 CLI
   make build      # 构建
   ```

2. **监听文件变化自动重新编译**
   ```bash
   # 安装 air
   go install github.com/cosmtrek/air@latest
   
   # 运行
   air
   ```

3. **代码审查**
   - 使用 golangci-lint
   ```bash
   golangci-lint run
   ```

## 🎯 优先级

当前需要帮助的领域：

1. 🔴 高优先级：
   - Bug 修复
   - 安全问题
   - 性能问题

2. 🟡 中优先级：
   - 新功能
   - 文档改进
   - 代码重构

3. 🟢 低优先级：
   - 代码优化
   - UI/UX 改进

## 📞 联系方式

- GitHub Issues: [提交 Issue](../../issues)
- Pull Requests: [查看 PRs](../../pulls)

## ⚖️ 许可证

通过贡献代码，你同意你的贡献将在 MIT 许可证下授权。

---

再次感谢你的贡献！ 🎉

