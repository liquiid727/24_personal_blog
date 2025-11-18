## 当前状态

* 已完成：Phase 1–5（骨架/配置/认证、文章模块、评论模块、文件上传、管理员统计），Phase 6（基础 Swagger 静态文档与集成测试）部分完成。

* 代码可构建与运行，集成测试通过。

## 待完成目标

* 完成 Phase 6（文档与测试）剩余部分；推进 Phase 7（性能与安全）。

* 提供 Docker 与 CI（可选但推荐）。

* 每一步里程碑完成后进行一次规范化 git commit（Conventional Commits）。

## Phase 6：Swagger 与测试完善

* 方案：集成 `swaggo/swag` 与 `gin-swagger`，在 Handler 上添加注释生成 `swagger.json`，仅在非生产环境暴露。

* 步骤：

  * 引入依赖与注释；生成 docs 包；注册 `/swagger/*` 路由。

  * 扩充文档：认证、文章、评论、文件、管理员统计的请求/响应模型与错误码。

  * 增加集成测试：文章更新/删除/发布、评论更新/删除、文件上传最小验证。

* 预期提交：

  * docs: phase 6 - integrate swaggo and annotate handlers

  * test: phase 6 - extend integration tests for posts/comments/files

## Phase 7：性能与安全增强

* Redis 缓存：

  * 新增 `internal/cache` 初始化 `go-redis/v9`；文章浏览量使用 Redis `INCR` 聚合，定时或请求阈值回写 DB；热门文章缓存（TTL 与失效策略）。

  * 提交：feat: phase 7 - integrate redis cache for views and hot posts

* 搜索优化（PostgreSQL 优先）：

  * 使用 `tsvector` 全文检索与 `GIN` 索引；新增迁移框架（`gormigrate`）执行 RAW SQL；MySQL 保持 LIKE 回退。

  * 提交：feat: phase 7 - add postgres fulltext search with migrations

* 安全加固：

  * CORS 白名单与严格方法；

  * 速率限制中间件（令牌桶），保护登录/评论/发布等；

  * 输入校验与统一错误码；

  * 提交：feat: phase 7 - add rate limiting and strict CORS; refactor: unify error responses

## DevOps：Docker 与 CI

* Docker：

  * 多阶段 Dockerfile（build → runtime）；`docker-compose.yml` 启动 `app+postgres+redis`；`.env.example` 更新。

  * 提交：chore: dockerize application and add docker-compose

* CI（GitHub Actions）：

  * 运行 `go fmt`, `go vet`, `go test`, 构建与（可选）镜像推送；

  * 提交：chore: add GitHub Actions CI for build and tests

## 管理员接口扩展（可选）

* 用户/文章/评论的管理端列表、启停、删除、筛选；权限校验沿用 `AdminOnly`。

* 提交：feat: admin - manage users/posts/comments endpoints

## 测试覆盖与质量

* 目标：核心模块覆盖率 ≥80%；

* 方法：服务层单元测试（仓储可用内存假实现）、路由集成测试（`httptest`）。

* 提交：test: increase coverage for services and routes

## 配置与中间件

* 扩充配置项：CORS、RateLimit、Redis、最大上传大小；

* 增加 RequestID 与结构化日志字段；

* 提交：feat: add config fields for cors/ratelimit/redis and request id middleware

## 交付清单

* 完整源码与注释、迁移脚本**（Postgres 全文检索）、Swagger 文档与生成方式说明。**

* **Docker 与 CI、测试报告。**

**确认后我将按上述里程碑顺序推进，每完成一个步骤即进行一次规范化 commit（不批量）。**
