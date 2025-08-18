# 开发与维护指南 🛠️

## 概述 📋

本文档为博客系统的开发和维护提供详细指导，包括开发环境配置、编码规范、测试策略、代码提交流程、问题排查方法以及日常维护注意事项。旨在帮助开发团队保持代码质量，提高开发效率，确保系统稳定运行。

## 开发环境配置 💻

### 1. 必需工具安装

#### Go 开发工具

```bash
# 安装 Go 语言工具
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest

# 验证安装
goimports -version
golangci-lint version
swag -version
air -version
```

#### 编辑器配置 (VS Code)

推荐安装的 VS Code 扩展：

```json
// .vscode/extensions.json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode.vscode-typescript-next",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-eslint"
  ]
}
```

VS Code 工作区配置：

```json
// .vscode/settings.json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ],
  "go.formatTool": "goimports",
  "go.testFlags": [
    "-v",
    "-race"
  ],
  "go.testTimeout": "30s",
  "go.coverOnSave": true,
  "go.coverageDecorator": {
    "type": "gutter",
    "coveredHighlightColor": "rgba(64,128,128,0.5)",
    "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
  },
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/node_modules": true,
    "**/vendor": true,
    "**/*.exe": true
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

#### Git 配置

```bash
# 配置用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# 配置编辑器
git config --global core.editor "code --wait"

# 配置换行符处理
git config --global core.autocrlf input  # Linux/macOS
git config --global core.autocrlf true   # Windows

# 配置别名
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.lg "log --oneline --graph --decorate --all"
```

### 2. 项目初始化

```bash
# 克隆项目
git clone <repository-url>
cd blog-system-refactored

# 安装依赖
go mod tidy

# 复制配置文件
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
vim configs/config.yaml

# 初始化数据库
go run cmd/migrate/main.go

# 启动开发服务器
air
```

### 3. 开发工具配置

#### Air 热重载配置

创建 `.air.toml` 文件：

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs", "logs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
```

#### Makefile 配置

创建 `Makefile`：

```makefile
# 变量定义
APP_NAME=blog-system
VERSION=1.0.0
BUILD_DIR=build
CMD_DIR=cmd
MAIN_FILE=$(CMD_DIR)/main.go

# Go 相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建标志
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# 默认目标
.PHONY: all
all: clean deps lint test build

# 安装依赖
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 代码检查
.PHONY: lint
lint:
	golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w .

# 运行测试
.PHONY: test
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# 查看测试覆盖率
.PHONY: coverage
coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	open coverage.html

# 构建应用
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# 构建多平台版本
.PHONY: build-all
build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)

# 运行应用
.PHONY: run
run:
	$(GOCMD) run $(MAIN_FILE)

# 开发模式运行
.PHONY: dev
dev:
	air

# 生成 API 文档
.PHONY: docs
docs:
	swag init -g $(MAIN_FILE) -o ./docs/swagger

# 数据库迁移
.PHONY: migrate
migrate:
	$(GOCMD) run cmd/migrate/main.go

# 清理构建文件
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 安装开发工具
.PHONY: install-tools
install-tools:
	$(GOGET) golang.org/x/tools/cmd/goimports@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/cosmtrek/air@latest

# Docker 构建
.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

# Docker 运行
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# 帮助信息
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  deps         - 安装依赖"
	@echo "  lint         - 代码检查"
	@echo "  fmt          - 格式化代码"
	@echo "  test         - 运行测试"
	@echo "  coverage     - 查看测试覆盖率"
	@echo "  build        - 构建应用"
	@echo "  build-all    - 构建多平台版本"
	@echo "  run          - 运行应用"
	@echo "  dev          - 开发模式运行"
	@echo "  docs         - 生成 API 文档"
	@echo "  migrate      - 数据库迁移"
	@echo "  clean        - 清理构建文件"
	@echo "  install-tools- 安装开发工具"
	@echo "  docker-build - Docker 构建"
	@echo "  docker-run   - Docker 运行"
```

## 编码规范 📝

### 1. Go 代码规范

#### 命名规范

```go
// ✅ 正确的命名
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUserByID(ctx context.Context, id uint) (*User, error)
}

type userServiceImpl struct {
    repo UserRepository
    logger *logrus.Logger
}

// ❌ 错误的命名
type userservice interface {  // 应该使用 PascalCase
    createuser(ctx context.Context, req *CreateUserRequest) (*User, error)  // 应该使用 PascalCase
}

type UserServiceimpl struct {  // 应该使用 camelCase
    Repo UserRepository        // 私有字段应该使用 camelCase
}
```

#### 包结构规范

```go
// ✅ 正确的包结构
package handlers

import (
    "context"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    
    "blog-system/internal/models"
    "blog-system/internal/services"
)

// ❌ 错误的包结构
package handlers

import (
    "blog-system/internal/models"  // 应该按组分类
    "context"
    "github.com/gin-gonic/gin"     // 应该按组分类
    "net/http"
)
```

#### 错误处理规范

```go
// ✅ 正确的错误处理
func (s *userServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    if req == nil {
        return nil, errors.New("请求参数不能为空")
    }
    
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, fmt.Errorf("参数验证失败: %w", err)
    }
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        s.logger.WithError(err).Error("创建用户失败")
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    return user, nil
}

// ❌ 错误的错误处理
func (s *userServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, err  // 没有包装错误信息
    }
    return user, nil
}
```

#### 注释规范

```go
// ✅ 正确的注释

// UserService 定义用户服务接口
// 提供用户的增删改查等基本操作
type UserService interface {
    // CreateUser 创建新用户
    // 参数:
    //   ctx: 上下文对象
    //   req: 创建用户请求参数
    // 返回:
    //   *User: 创建的用户对象
    //   error: 错误信息
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
}

// validateEmail 验证邮箱格式是否正确
// 使用正则表达式进行验证
func validateEmail(email string) bool {
    // 邮箱正则表达式
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

// ❌ 错误的注释
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)  // 创建用户
}

func validateEmail(email string) bool {
    // 验证邮箱
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

### 2. 代码质量检查

#### golangci-lint 配置

创建 `.golangci.yml` 文件：

```yaml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  skip-dirs:
    - vendor
    - tmp
    - docs
  skip-files:
    - ".*\\.pb\\.go$"
    - ".*\\_gen\\.go$"

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  govet:
    check-shadowing: true
    enable-all: true
  
  gocyclo:
    min-complexity: 15
  
  dupl:
    threshold: 100
  
  goconst:
    min-len: 3
    min-occurrences: 3
  
  misspell:
    locale: US
  
  lll:
    line-length: 120
  
  goimports:
    local-prefixes: blog-system
  
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - dupImport
      - ifElseChain
      - octalLiteral
      - whyNoLint
      - wrapperFunc

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - exportloopref
    - exhaustive
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - lll
    - misspell
    - nakedret
    - noctx
    - nolintlint
    - rowserrcheck
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
  
  disable:
    - maligned
    - prealloc

issues:
  exclude-rules:
    - path: _test\\.go
      linters:
        - gomnd
        - funlen
        - gocyclo
    
    - path: cmd/
      linters:
        - gochecknoinits
  
  exclude:
    - "Error return value of .((os\\.)?std(out|err)\\..*|.*Close|.*Flush|os\\.Remove(All)?|.*printf?|os\\.(Un)?Setenv). is not checked"
  
  max-issues-per-linter: 0
  max-same-issues: 0
```

### 3. 测试规范

#### 单元测试示例

```go
// user_service_test.go
package services

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    
    "blog-system/internal/models"
    "blog-system/internal/repository/mocks"
)

// UserServiceTestSuite 用户服务测试套件
type UserServiceTestSuite struct {
    suite.Suite
    mockRepo    *mocks.UserRepository
    userService UserService
    ctx         context.Context
}

// SetupTest 测试前置设置
func (suite *UserServiceTestSuite) SetupTest() {
    suite.mockRepo = &mocks.UserRepository{}
    suite.userService = NewUserService(suite.mockRepo, nil)
    suite.ctx = context.Background()
}

// TestCreateUser_Success 测试创建用户成功场景
func (suite *UserServiceTestSuite) TestCreateUser_Success() {
    // Arrange
    req := &CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    expectedUser := &models.User{
        ID:       1,
        Username: req.Username,
        Email:    req.Email,
        Status:   models.UserStatusActive,
        CreatedAt: time.Now(),
    }
    
    suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)
    
    // Act
    user, err := suite.userService.CreateUser(suite.ctx, req)
    
    // Assert
    assert.NoError(suite.T(), err)
    assert.NotNil(suite.T(), user)
    assert.Equal(suite.T(), expectedUser.Username, user.Username)
    assert.Equal(suite.T(), expectedUser.Email, user.Email)
    suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateUser_InvalidEmail 测试创建用户邮箱格式错误场景
func (suite *UserServiceTestSuite) TestCreateUser_InvalidEmail() {
    // Arrange
    req := &CreateUserRequest{
        Username: "testuser",
        Email:    "invalid-email",
        Password: "password123",
    }
    
    // Act
    user, err := suite.userService.CreateUser(suite.ctx, req)
    
    // Assert
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), user)
    assert.Contains(suite.T(), err.Error(), "邮箱格式不正确")
}

// TestUserServiceTestSuite 运行测试套件
func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}

// 基准测试示例
func BenchmarkUserService_CreateUser(b *testing.B) {
    mockRepo := &mocks.UserRepository{}
    userService := NewUserService(mockRepo, nil)
    ctx := context.Background()
    
    req := &CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    expectedUser := &models.User{
        ID:       1,
        Username: req.Username,
        Email:    req.Email,
    }
    
    mockRepo.On("Create", ctx, mock.Anything).Return(expectedUser, nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = userService.CreateUser(ctx, req)
    }
}
```

#### 集成测试示例

```go
// integration_test.go
// +build integration

package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    
    "blog-system/internal/config"
    "blog-system/internal/handlers"
    "blog-system/internal/routes"
)

// IntegrationTestSuite 集成测试套件
type IntegrationTestSuite struct {
    suite.Suite
    router *gin.Engine
    db     *gorm.DB
}

// SetupSuite 测试套件初始化
func (suite *IntegrationTestSuite) SetupSuite() {
    // 设置测试模式
    gin.SetMode(gin.TestMode)
    
    // 初始化测试数据库
    cfg := config.GetTestConfig()
    db, err := config.ConnectDatabase(cfg.Database)
    assert.NoError(suite.T(), err)
    
    suite.db = db
    
    // 初始化路由
    suite.router = routes.SetupRoutes(db, cfg)
}

// TearDownSuite 测试套件清理
func (suite *IntegrationTestSuite) TearDownSuite() {
    // 清理测试数据
    suite.db.Exec("DELETE FROM users")
    suite.db.Exec("DELETE FROM posts")
}

// TestCreateUser_Integration 测试创建用户接口
func (suite *IntegrationTestSuite) TestCreateUser_Integration() {
    // Arrange
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    
    jsonBody, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    // Act
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(suite.T(), http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), "success", response["status"])
    assert.NotNil(suite.T(), response["data"])
}

// TestIntegrationTestSuite 运行集成测试套件
func TestIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(IntegrationTestSuite))
}
```

## Git 工作流程 🔄

### 1. 分支管理策略

#### Git Flow 分支模型

```
master (生产分支)
├── develop (开发分支)
│   ├── feature/user-management (功能分支)
│   ├── feature/post-system (功能分支)
│   └── feature/comment-system (功能分支)
├── release/v1.0.0 (发布分支)
└── hotfix/critical-bug-fix (热修复分支)
```

#### 分支命名规范

```bash
# 功能分支
feature/user-authentication
feature/post-management
feature/comment-system

# 修复分支
bugfix/login-error
bugfix/database-connection

# 热修复分支
hotfix/security-vulnerability
hotfix/critical-performance-issue

# 发布分支
release/v1.0.0
release/v1.1.0
```

### 2. 提交信息规范

#### Conventional Commits 格式

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### 提交类型

- **feat**: 新功能
- **fix**: 修复 bug
- **docs**: 文档更新
- **style**: 代码格式调整（不影响功能）
- **refactor**: 代码重构
- **perf**: 性能优化
- **test**: 测试相关
- **chore**: 构建过程或辅助工具的变动
- **ci**: CI/CD 相关
- **build**: 构建系统或外部依赖的变动

#### 提交示例

```bash
# 新功能
git commit -m "feat(auth): 添加用户登录功能

- 实现用户名/邮箱登录
- 添加 JWT token 生成
- 集成密码加密验证

Closes #123"

# 修复 bug
git commit -m "fix(api): 修复用户创建时邮箱验证问题

修复了邮箱格式验证正则表达式错误导致的
合法邮箱被拒绝的问题。

Fixes #456"

# 文档更新
git commit -m "docs: 更新 API 文档和部署指南"

# 代码重构
git commit -m "refactor(service): 重构用户服务层代码结构"
```

### 3. 代码审查流程

#### Pull Request 模板

创建 `.github/pull_request_template.md`：

```markdown
## 变更描述

<!-- 请简要描述此 PR 的变更内容 -->

## 变更类型

- [ ] 新功能 (feature)
- [ ] 修复 bug (fix)
- [ ] 代码重构 (refactor)
- [ ] 性能优化 (perf)
- [ ] 文档更新 (docs)
- [ ] 测试相关 (test)
- [ ] 其他 (chore)

## 测试

- [ ] 单元测试已通过
- [ ] 集成测试已通过
- [ ] 手动测试已完成
- [ ] 代码覆盖率满足要求

## 检查清单

- [ ] 代码遵循项目编码规范
- [ ] 已添加必要的注释
- [ ] 已更新相关文档
- [ ] 已添加或更新测试用例
- [ ] 所有 CI 检查通过
- [ ] 已进行自我代码审查

## 相关 Issue

<!-- 关联的 Issue 编号，例如：Closes #123 -->

## 截图（如适用）

<!-- 如果有 UI 变更，请提供截图 -->

## 额外说明

<!-- 任何需要审查者注意的额外信息 -->
```

#### 代码审查检查点

1. **功能性**
   - 代码是否实现了预期功能
   - 是否有潜在的 bug
   - 边界条件是否处理正确

2. **可读性**
   - 代码逻辑是否清晰
   - 命名是否恰当
   - 注释是否充分

3. **性能**
   - 是否有性能问题
   - 数据库查询是否优化
   - 内存使用是否合理

4. **安全性**
   - 是否有安全漏洞
   - 输入验证是否充分
   - 权限控制是否正确

5. **测试**
   - 测试覆盖率是否足够
   - 测试用例是否合理
   - 是否有集成测试

### 4. 发布流程

#### 版本号规范 (Semantic Versioning)

```
MAJOR.MINOR.PATCH

例如：1.2.3
- MAJOR: 不兼容的 API 变更
- MINOR: 向后兼容的功能新增
- PATCH: 向后兼容的问题修正
```

#### 发布步骤

```bash
# 1. 创建发布分支
git checkout develop
git pull origin develop
git checkout -b release/v1.1.0

# 2. 更新版本号
vim version.go  # 或其他版本文件
git add version.go
git commit -m "chore: 更新版本号到 v1.1.0"

# 3. 运行完整测试
make test
make lint

# 4. 合并到 master
git checkout master
git merge --no-ff release/v1.1.0
git tag -a v1.1.0 -m "Release version 1.1.0"

# 5. 合并回 develop
git checkout develop
git merge --no-ff release/v1.1.0

# 6. 推送到远程
git push origin master
git push origin develop
git push origin v1.1.0

# 7. 删除发布分支
git branch -d release/v1.1.0
git push origin --delete release/v1.1.0
```

## 问题排查指南 🔍

### 1. 常见问题诊断

#### 应用启动问题

```bash
# 检查端口占用
netstat -tlnp | grep :8080
lsof -i :8080

# 检查配置文件
cat configs/config.yaml

# 检查环境变量
env | grep -E "(DB_|JWT_|SERVER_)"

# 检查日志
tail -f logs/app.log
journalctl -u blog-system -f
```

#### 数据库连接问题

```bash
# 测试数据库连接
mysql -h localhost -u blog_user -p blog_system

# 检查数据库服务状态
sudo systemctl status mysql

# 查看数据库错误日志
sudo tail -f /var/log/mysql/error.log

# 检查数据库配置
sudo cat /etc/mysql/mysql.conf.d/mysqld.cnf
```

#### 性能问题诊断

```bash
# 查看系统资源使用
top
htop
free -h
df -h

# 查看应用进程
ps aux | grep blog-system

# 查看网络连接
netstat -an | grep :8080
ss -tulpn | grep :8080

# 分析慢查询
mysql -u root -p -e "SHOW PROCESSLIST;"
mysql -u root -p -e "SELECT * FROM information_schema.processlist WHERE time > 5;"
```

### 2. 日志分析

#### 日志级别说明

```go
// 日志级别定义
const (
    DebugLevel = "debug"  // 调试信息
    InfoLevel  = "info"   // 一般信息
    WarnLevel  = "warn"   // 警告信息
    ErrorLevel = "error"  // 错误信息
    FatalLevel = "fatal"  // 致命错误
)
```

#### 日志查看命令

```bash
# 查看实时日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log
grep "FATAL" logs/app.log

# 查看特定时间段日志
grep "2024-01-15 10:" logs/app.log

# 查看特定用户操作日志
grep "user_id=123" logs/app.log

# 统计错误数量
grep -c "ERROR" logs/app.log

# 查看最近的错误
grep "ERROR" logs/app.log | tail -10
```

#### 日志分析脚本

创建 `scripts/analyze_logs.sh`：

```bash
#!/bin/bash

# 日志分析脚本
LOG_FILE="logs/app.log"
DATE=$(date +"%Y-%m-%d")

echo "=== 日志分析报告 ($DATE) ==="
echo

# 错误统计
echo "错误统计:"
echo "ERROR: $(grep -c "ERROR" $LOG_FILE)"
echo "WARN:  $(grep -c "WARN" $LOG_FILE)"
echo "FATAL: $(grep -c "FATAL" $LOG_FILE)"
echo

# 最近错误
echo "最近 10 个错误:"
grep "ERROR" $LOG_FILE | tail -10
echo

# API 请求统计
echo "API 请求统计:"
grep "method=" $LOG_FILE | awk '{print $NF}' | sort | uniq -c | sort -nr
echo

# 响应时间分析
echo "慢请求 (>1s):"
grep "duration=" $LOG_FILE | awk '$NF > 1000 {print}'
echo

# 用户活动统计
echo "活跃用户 TOP 10:"
grep "user_id=" $LOG_FILE | sed 's/.*user_id=\([0-9]*\).*/\1/' | sort | uniq -c | sort -nr | head -10
```

### 3. 性能监控

#### 应用性能指标

```go
// 性能监控中间件
func PerformanceMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\
",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}

// 健康检查端点
func HealthCheck(c *gin.Context) {
    start := time.Now()
    
    // 检查数据库连接
    db := database.GetDB()
    sqlDB, _ := db.DB()
    err := sqlDB.Ping()
    
    status := "ok"
    if err != nil {
        status = "error"
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status":    status,
        "timestamp": time.Now(),
        "version":   version.Version,
        "uptime":    time.Since(start),
        "database":  map[string]interface{}{
            "status": func() string {
                if err != nil {
                    return "disconnected"
                }
                return "connected"
            }(),
        },
    })
}
```

#### 监控脚本

创建 `scripts/monitor.sh`：

```bash
#!/bin/bash

# 系统监控脚本
APP_NAME="blog-system"
LOG_FILE="logs/monitor.log"
ALERT_EMAIL="admin@example.com"

# 检查应用是否运行
check_app_status() {
    if pgrep -f $APP_NAME > /dev/null; then
        echo "$(date): $APP_NAME is running" >> $LOG_FILE
        return 0
    else
        echo "$(date): $APP_NAME is not running" >> $LOG_FILE
        return 1
    fi
}

# 检查内存使用
check_memory() {
    MEMORY_USAGE=$(free | grep Mem | awk '{printf "%.2f", $3/$2 * 100.0}')
    echo "$(date): Memory usage: ${MEMORY_USAGE}%" >> $LOG_FILE
    
    if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
        echo "High memory usage: ${MEMORY_USAGE}%" | mail -s "Memory Alert" $ALERT_EMAIL
    fi
}

# 检查磁盘空间
check_disk() {
    DISK_USAGE=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    echo "$(date): Disk usage: ${DISK_USAGE}%" >> $LOG_FILE
    
    if [ $DISK_USAGE -gt 80 ]; then
        echo "High disk usage: ${DISK_USAGE}%" | mail -s "Disk Alert" $ALERT_EMAIL
    fi
}

# 检查 API 响应
check_api() {
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
    echo "$(date): API health check: $RESPONSE" >> $LOG_FILE
    
    if [ "$RESPONSE" != "200" ]; then
        echo "API health check failed: $RESPONSE" | mail -s "API Alert" $ALERT_EMAIL
    fi
}

# 主监控循环
main() {
    echo "$(date): Starting monitoring..." >> $LOG_FILE
    
    while true; do
        check_app_status
        check_memory
        check_disk
        check_api
        
        sleep 300  # 5 分钟检查一次
    done
}

main
```

## 日常维护 🔧

### 1. 定期维护任务

#### 每日维护

```bash
#!/bin/bash
# daily_maintenance.sh

echo "开始每日维护任务..."

# 检查应用状态
sudo systemctl status blog-system

# 检查日志错误
grep "ERROR\|FATAL" logs/app.log | tail -20

# 检查磁盘空间
df -h

# 检查内存使用
free -h

# 清理临时文件
find tmp/ -type f -mtime +7 -delete

# 检查数据库连接
mysql -u blog_user -p -e "SELECT 1" blog_system

echo "每日维护任务完成"
```

#### 每周维护

```bash
#!/bin/bash
# weekly_maintenance.sh

echo "开始每周维护任务..."

# 数据库优化
mysql -u root -p blog_system -e "OPTIMIZE TABLE users, posts, comments;"

# 日志轮转
logrotate -f /etc/logrotate.d/blog-system

# 清理旧备份
find /home/blog/backups -name "*.sql.gz" -mtime +30 -delete

# 更新系统包
sudo apt update && sudo apt upgrade -y

# 重启应用（如需要）
# sudo systemctl restart blog-system

echo "每周维护任务完成"
```

#### 每月维护

```bash
#!/bin/bash
# monthly_maintenance.sh

echo "开始每月维护任务..."

# 数据库完整备份
mysqldump -u root -p --all-databases > /backup/full_backup_$(date +%Y%m%d).sql

# 分析数据库性能
mysql -u root -p -e "SHOW ENGINE INNODB STATUS\G" > /tmp/innodb_status.txt

# 检查数据库表状态
mysql -u root -p blog_system -e "CHECK TABLE users, posts, comments;"

# 生成月度报告
scripts/generate_monthly_report.sh

# 安全更新
sudo apt update && sudo apt upgrade -y

echo "每月维护任务完成"
```

### 2. 监控和告警

#### 系统监控配置

```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  - job_name: 'blog-system'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']
```

#### 告警规则

```yaml
# monitoring/alert_rules.yml
groups:
  - name: blog-system-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is above 80%"

      - alert: DatabaseConnectionFailed
        expr: mysql_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection failed"
          description: "Cannot connect to MySQL database"
```

### 3. 安全维护

#### 安全检查清单

```bash
#!/bin/bash
# security_check.sh

echo "开始安全检查..."

# 检查系统更新
echo "检查系统更新:"
apt list --upgradable

# 检查开放端口
echo "检查开放端口:"
netstat -tlnp

# 检查登录日志
echo "检查登录日志:"
last -10

# 检查失败登录
echo "检查失败登录:"
lastb -10

# 检查文件权限
echo "检查应用文件权限:"
ls -la /home/blog/app/

# 检查配置文件权限
echo "检查配置文件权限:"
ls -la /home/blog/app/configs/

# 检查 SSL 证书有效期
echo "检查 SSL 证书:"
openssl x509 -in /etc/ssl/certs/your-domain.crt -text -noout | grep "Not After"

echo "安全检查完成"
```

#### 密码和密钥管理

```bash
# 生成强密码
openssl rand -base64 32

# 生成 JWT 密钥
openssl rand -hex 64

# 检查密钥文件权限
chmod 600 configs/jwt.key
chown blog:blog configs/jwt.key

# 定期轮换密钥
# 1. 生成新密钥
# 2. 更新配置
# 3. 重启应用
# 4. 验证功能正常
```

### 4. 备份和恢复

#### 自动备份脚本

```bash
#!/bin/bash
# auto_backup.sh

BACKUP_DIR="/backup/blog-system"
DATE=$(date +"%Y%m%d_%H%M%S")
RETENTION_DAYS=30

# 创建备份目录
mkdir -p $BACKUP_DIR

# 数据库备份
echo "备份数据库..."
mysqldump -u backup_user -p$BACKUP_PASSWORD blog_system | gzip > $BACKUP_DIR/db_$DATE.sql.gz

# 应用文件备份
echo "备份应用文件..."
tar -czf $BACKUP_DIR/app_$DATE.tar.gz -C /home/blog/app --exclude='logs' --exclude='tmp' .

# 配置文件备份
echo "备份配置文件..."
cp -r /home/blog/app/configs $BACKUP_DIR/configs_$DATE

# 上传到云存储（可选）
echo "上传到云存储..."
# aws s3 cp $BACKUP_DIR/db_$DATE.sql.gz s3://your-backup-bucket/
# aws s3 cp $BACKUP_DIR/app_$DATE.tar.gz s3://your-backup-bucket/

# 清理旧备份
echo "清理旧备份..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "configs_*" -mtime +$RETENTION_DAYS -exec rm -rf {} +

echo "备份完成: $DATE"
```

#### 恢复脚本

```bash
#!/bin/bash
# restore.sh

BACKUP_FILE="$1"
APP_BACKUP="$2"

if [ -z "$BACKUP_FILE" ] || [ -z "$APP_BACKUP" ]; then
    echo "用法: $0 <数据库备份文件> <应用备份文件>"
    exit 1
fi

echo "开始恢复..."

# 停止应用
echo "停止应用服务..."
sudo systemctl stop blog-system

# 恢复数据库
echo "恢复数据库..."
zcat $BACKUP_FILE | mysql -u root -p blog_system

# 恢复应用文件
echo "恢复应用文件..."
cd /home/blog
tar -xzf $APP_BACKUP

# 设置权限
echo "设置文件权限..."
chown -R blog:blog /home/blog/app
chmod +x /home/blog/app/blog-system

# 启动应用
echo "启动应用服务..."
sudo systemctl start blog-system

# 验证恢复
echo "验证应用状态..."
sleep 5
curl -s http://localhost:8080/health

echo "恢复完成"
```

---

**注意**：本文档提供了完整的开发和维护指南，请根据团队实际情况调整相关流程和规范。定期更新文档以保持与项目发展同步。