## 项目目标与技术选型

* 目标：实现高性能、可扩展的博客系统，覆盖用户认证、文章与评论管理、标签分类、文件上传、管理员后台、API 文档与测试。

* 语言与框架：`Go 1.22+`、`Gin`（HTTP 路由与中间件）、`GORM`（ORM）。

* 数据库：优先 `PostgreSQL`（支持全文检索与丰富类型），兼容 `MySQL`；使用 `AutoMigrate` 起步，后续可切换到 `goose/gormigrate` 管理迁移。

* 缓存：`Redis (go-redis/v9)` 用于热点数据与会话（可选）。

* 认证：`JWT (golang-jwt/jwt/v5)`；密码哈希：`bcrypt`。

* 配置：`Viper` + 环境变量（12-factor）。

* 日志：`Zap` 结构化日志。

* 文档：`Swagger (swaggo/gin-swagger)`。

* 测试：`testing` + `Testify`，`httptest` 做集成测试。

## 目录结构

* `cmd/server/main.go`：入口。

* `internal/config`：配置加载（Viper）。

* `internal/log`：日志初始化（Zap）。

* `internal/db`：数据库连接、迁移、事务。

* `internal/cache`：Redis 客户端与封装。

* `internal/auth`：JWT、密码哈希、RBAC。

* `internal/server`：Gin 引导、路由、全局中间件（CORS、Recovery、RequestID、RateLimit）。

* `internal/transport/http`：DTO、参数校验、Handlers。

* `internal/repository`：持久化接口与 GORM 实现。

* `internal/service`：业务服务层（聚合与领域逻辑）。

* `internal/domain/{user,post,comment,tag,category,file}`：领域模型与用例。

* `pkg/{utils,validator}`：通用工具与校验。

* `docs/swagger`：API 文档生成。

* `scripts/{migrate,seed}`：迁移与数据脚本。

* 根目录：`go.mod`、`Dockerfile`、`docker-compose.yml`、`Makefile`、`.env.example`。

## 数据库设计（ER 摘要）

* `users`：`id`、`username`(uniq)・`email`(uniq)・`password_hash`・`role`(`admin|user`)・`status`・`created_at`・`updated_at`。

* `posts`：`id`・`author_id(FK users)`・`title`・`slug`(uniq)・`content`・`status`(`draft|published|private`)・`views`・`published_at`・`created_at/updated_at`。

* `tags`：`id`・`name`(uniq)・`slug`(uniq)。

* `categories`：`id`・`name`(uniq)・`slug`(uniq)。

* `post_tags`：`post_id`・`tag_id`（联合唯一）。

* `post_categories`：`post_id`・`category_id`（联合唯一）。

* `comments`：`id`・`post_id`・`user_id`・`parent_id(NULL)`・`content`・`status`・`created_at/updated_at`（多层回复：邻接表）。

* `files`：`id`・`uploader_id`・`path`・`mime`・`size`・`created_at`。

* 索引：`posts.slug`、`tags.slug`、`categories.slug` 唯一；`comments(post_id,parent_id)`；`posts(status,views)`；需要的外键索引。

## 认证与授权

* 注册/登录/退出：注册写入哈希密码；登录颁发 `JWT`；退出前端删除令牌，后端可选加入黑名单。

* 中间件：验证 `Authorization: Bearer`，将用户上下文注入；RBAC：基于 `role` 控制路由与服务访问。

* 令牌：访问令牌（短时），可选刷新令牌（长时）。

## 博客文章与搜索

* CRUD：创建/编辑/发布/删除/查询；状态流转校验。

* 分类与标签：多对多（标签）与一对多/多对多（分类）绑定。

* 搜索：标题/内容/标签。起步用 `LIKE + trigram`（PostgreSQL），后续可扩展 `tsvector`。

* 浏览量：请求级原子递增；可选 Redis 缓存合并回写。

* Slug：按标题生成，碰撞解决。

## 评论系统

* CRUD：用户管理个人评论；管理员管理全部评论。

* 多层回复：`parent_id` 邻接表，查询时通过层次遍历返回树结构。

* 审核状态：`pending|approved|rejected`（可选）。

## 文件上传与管理

* 上传：限制 MIME 与大小；保存至 `uploads/`，记录元数据；提供访问 URL。

* 删除/查看：受权限控制。

* 安全：文件名清理、防路径穿越、仅允许图片类型；后续可扩展对象存储（S3/OSS）。

## 管理员功能与统计

* 管理：用户、文章、评论的列表、筛选、启停与删除。

* 统计：文章数、用户数、评论数、浏览量汇总；可缓存。

## 中间件与基础设施

* 日志：`Zap` 结构化日志（请求/响应摘要、错误栈）。

* Recovery：统一错误处理与响应格式。

* CORS：按环境配置来源与方法。

* RequestID：链路追踪；RateLimit（令牌桶）防滥用。

* 校验：`go-playground/validator`；统一请求 DTO 校验错误返回。

## API 约定与 Swagger

* 前缀：`/api/v1`；统一响应包：`{code,message,data}`。

* 分页：`page,page_size`；排序：`sort_by,order`；过滤统一参数。

* 鉴权：受保护路由需要 `Bearer`；管理员路由仅 `admin`。

* Swagger：注释生成，`/swagger/index.html` 暴露。

## 测试策略

* 单元测试：服务层与仓储层；使用内存或事务回滚。

* 集成测试：`httptest` 驱动路由；起本地 Postgres 测试库（或 docker-compose）。

* 覆盖率：核心模块≥80%；并在 CI 输出报告。

## 部署与配置

* Docker：多阶段构建（编译→精简运行时）；`docker-compose` 启动 `app+db+redis`。

* 配置：环境变量加载（端口、DB/Redis、JWT 秘钥、CORS）。

* 运行：`make dev`、`make test`、`make migrate`、`make swagger` 等目标。

* 安全：生产禁用 Swagger；只在非生产暴露。

## 性能与安全

* 性能：合理索引、避免 N+1（`Preload`）、分页查询、热点缓存。

* 安全：ORM 防 SQL 注入、XSS 过滤/输出转义、CORS 严格、速率限制、最小权限、错误信息最小化。

## 里程碑

* Phase 1：项目骨架（目录、go.mod）、配置与日志、数据库连接、用户注册/登录/JWT。

* Phase 2：文章模块（CRUD、状态、标签/分类、Slug、浏览量）。

* Phase 3：评论模块（CRUD、树形查询、权限）。

* Phase 4：文件上传/静态服务与文件管理。

* Phase 5：管理员功能与系统统计。

* Phase 6：Swagger 文档完善、单元与集成测试、Docker 化与部署脚本。

* Phase 7：Redis 缓存、搜索优化（tsvector）、安全加固与性能优化。

## 交付物

* 完整源码（遵循 gofmt/goimports）。

* 数据库 ER 图与初始化 SQL。

* Swagger API 文档与使用说明。

* 单元与集成测试代码，确保通过。

* 部署文档（Docker/本地运行）。

（当前仓库为空；确认计划后我将按上述里程碑自顶向下创建项目骨架、实现各模块与测试，并在关键步骤展示与验证。）
