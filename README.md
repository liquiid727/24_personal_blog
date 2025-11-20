# Go-Blog 高性能博客系统

![CI](https://github.com/your-org/go_blog/actions/workflows/ci.yml/badge.svg)  
![Coverage](https://img.shields.io/codecov/c/github/your-org/go_blog)  
![License](https://img.shields.io/badge/license-MIT-blue.svg)

Go-Blog 是一个基于 Go + Gin + PostgreSQL/MySQL 构建的现代化博客系统，支持用户注册登录、文章发布、评论回复、文件上传、全文检索与管理员统计。项目采用分层架构（Handler → Service → Repository → Domain），内置 JWT 认证、Redis 缓存、CORS 与限流中间件，支持 Docker 一键部署，提供完整的 OpenAPI 文档与单元测试覆盖，适合个人博客或团队知识库场景。

## 技术栈

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.22+ | 后端语言 |
| Gin | v1.9 | Web 框架 |
| GORM | v1.25 | ORM 与迁移 |
| PostgreSQL | 15+ | 主数据库（推荐） |
| MySQL | 8.0+ | 兼容数据库 |
| Redis | 7+ | 缓存与视图计数 |
| JWT | golang-jwt/v5 | 无状态认证 |
| Zap | v1.27 | 结构化日志 |
| Viper | v1.19 | 配置管理 |
| Testify | v1.9 | 单元测试 |

## 系统要求

- **Go**: 1.22 及以上（启用 go modules）
- **数据库**: PostgreSQL 15+ 或 MySQL 8.0+
- **缓存**: Redis 7+（可选，用于视图计数）
- **容器**: Docker 24+ & docker-compose v2（可选，一键部署）

## 依赖项

核心依赖已在 `go.mod` 锁定，运行 `go mod download` 即可自动获取；主要模块：

```text
github.com/gin-gonic/gin v1.9.1
gorm.io/driver/postgres v1.5.4
gorm.io/driver/mysql v1.5.7
gorm.io/gorm v1.25.12
github.com/golang-jwt/jwt/v5 v5.2.1
github.com/redis/go-redis/v9 v9.5.1
go.uber.org/zap v1.27.0
github.com/spf13/viper v1.19.0
github.com/stretchr/testify v1.9.0
```

## 安装部署指南

### 1. 环境准备

克隆仓库并进入目录：

```bash
git clone https://github.com/your-org/go_blog.git
cd go_blog
```

准备数据库（以 PostgreSQL 为例）：

```bash
# 创建数据库与用户
psql -U postgres -c "CREATE DATABASE go_blog OWNER go_blog;"
# 复制环境模板
cp .env.example .env
# 编辑数据库连接串
vim .env
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 启动开发服务器

```bash
go run cmd/server/main.go
```

服务默认监听 `0.0.0.0:8080`，日志级别 `info`。访问：

- API 入口: http://localhost:8080/api/v1
- Swagger UI: http://localhost:8080/swagger/index.html

### 4. Docker 一键部署（可选）

```bash
docker-compose up -d
```

将自动启动：
- `app` 服务（端口 8080）
- `postgres` 数据库（端口 5432）
- `redis` 缓存（端口 6379）

## 使用说明

### 基础用法

1. 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"123456"}'
```

2. 登录获取 Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"123456"}'
# 返回 {"token":"eyJ...","user":{...}}
```

3. 创建文章
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer eyJ..." \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello Go","content":"# Hello","status":"published","tag_ids":[1],"category_ids":[1]}'
```

### 配置选项（.env）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| **APP_ENV** | development | 运行环境 |
| **APP_PORT** | 8080 | 服务端口 |
| **LOG_LEVEL** | info | 日志级别 |
| **DB_DRIVER** | postgres | 数据库驱动 |
| **DB_DSN** | - | 连接串 |
| **JWT_SECRET** | changeme | HS256 密钥（**生产环境务必修改**） |
| **JWT_TTL** | 24h | 令牌有效期 |
| **REDIS_ADDR** | localhost:6379 | Redis 地址 |
| **REDIS_PASS** | - | Redis 密码 |
| **REDIS_DB** | 0 | Redis DB 索引 |

### API 文档

完整接口说明见 [docs/openapi.json](docs/openapi.json)，也可通过 Swagger UI 在线浏览。主要端点：

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 |
| `/api/v1/auth/login` | POST | 用户登录 |
| `/api/v1/auth/me` | GET | 获取当前用户信息 |
| `/api/v1/posts` | GET/POST | 文章列表/创建 |
| `/api/v1/posts/{id}` | GET/PUT/DELETE | 文章详情/更新/删除 |
| `/api/v1/posts/{id}/publish` | POST | 发布文章 |
| `/api/v1/posts/{id}/views/incr` | POST | 递增浏览量 |
| `/api/v1/posts/{id}/comments` | GET/POST | 评论树/创建评论 |
| `/api/v1/comments/{id}` | PUT/DELETE | 更新/删除评论 |
| `/api/v1/files` | POST | 文件上传（图片） |
| `/api/v1/files/{id}` | GET/DELETE | 文件元数据/删除（管理员） |
| `/api/v1/admin/stats` | GET | 系统统计（管理员） |

## 贡献指南

### 代码规范

- 遵循 [docs/commenting-guidelines.md](docs/commenting-guidelines.md) 注释规范
- 使用 `gofmt` 格式化，`golangci-lint` 静态检查
- 单元测试覆盖率 ≥ 80%，新增代码需同步补充测试
- 提交信息格式：`type(scope): subject`，如 `feat(post): add fulltext search`

### 问题报告

请使用 GitHub Issues，模板如下：

```markdown
**问题描述**：简要说明现象与期望行为
**复现步骤**：1. 2. 3.
**环境信息**：Go 版本、数据库、日志片段
**附加信息**：截图或最小复现代码
```

### PR 流程

1. Fork 仓库并创建特性分支：`git checkout -b feat/your-feature`
2. 编写/更新测试，确保 `make test` 通过
3. 本地运行 `make lint` 修复所有警告
4. 提交 PR，填写模板并关联对应 Issue
5. 通过 CI 后由维护者合并

## 附加信息

### 许可证

MIT License，详见 [LICENSE](LICENSE)

### 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 Web 框架
- [GORM](https://gorm.io) - 强大 ORM
- [Zap](https://github.com/uber-go/zap) - 极速日志库
- [Viper](https://github.com/spf13/viper) - 配置管理利器

### 联系方式

- 维护者：Liquiid
- 邮箱：liquiid727@outlook.com
- 讨论区：[GitHub Discussions](https://github.com/your-org/go_blog/discussions)

---

⭐ 如果项目对你有帮助，欢迎 Star 支持！