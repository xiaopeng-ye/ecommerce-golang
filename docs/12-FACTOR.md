# 12-Factor App 最佳实践实现

本文档说明了 ecommerce-golang 项目如何实现 [12-Factor App](https://12factor.net/) 方法论的最佳实践。

## 目录
- [I. 基准代码 (Codebase)](#i-基准代码-codebase)
- [II. 依赖 (Dependencies)](#ii-依赖-dependencies)
- [III. 配置 (Config)](#iii-配置-config)
- [IV. 后端服务 (Backing Services)](#iv-后端服务-backing-services)
- [V. 构建、发布、运行 (Build, Release, Run)](#v-构建发布运行-build-release-run)
- [VI. 进程 (Processes)](#vi-进程-processes)
- [VII. 端口绑定 (Port Binding)](#vii-端口绑定-port-binding)
- [VIII. 并发 (Concurrency)](#viii-并发-concurrency)
- [IX. 易处理 (Disposability)](#ix-易处理-disposability)
- [X. 开发环境与线上环境等价 (Dev/Prod Parity)](#x-开发环境与线上环境等价-devprod-parity)
- [XI. 日志 (Logs)](#xi-日志-logs)
- [XII. 管理进程 (Admin Processes)](#xii-管理进程-admin-processes)

---

## I. 基准代码 (Codebase)

**原则**: 一份基准代码，多份部署

**实现**:
- 项目使用 Git 进行版本控制
- 单一代码仓库: `github.com/xiaopeng-ye/ecommerce-golang`
- 可以部署到多个环境 (开发、测试、预发布、生产)

**验证**:
```bash
git remote -v
```

---

## II. 依赖 (Dependencies)

**原则**: 显式声明依赖关系

**实现**:
- 使用 `go.mod` 和 `go.sum` 明确声明所有依赖
- 依赖版本锁定，确保可重复构建
- 不依赖系统级的隐式依赖

**相关文件**:
- `go.mod` - 依赖声明
- `go.sum` - 依赖校验和

**验证**:
```bash
go mod verify
go mod download
```

---

## III. 配置 (Config)

**原则**: 在环境中存储配置

**实现**:
- 所有配置通过环境变量注入
- 敏感信息（数据库密码、JWT密钥）必须通过环境变量提供
- 开发环境支持 `.env` 文件（不提交到版本控制）
- 生产环境通过容器编排平台注入环境变量

**相关文件**:
- `config/env.go` - 配置管理
- `.env.example` - 环境变量模板
- `k8s/configmap.yaml` - Kubernetes 配置
- `k8s/secret.yaml` - Kubernetes 敏感配置

**环境变量**:
```bash
# 应用配置
ENVIRONMENT=production
LOG_LEVEL=info
LOG_FORMAT=json

# 服务器配置
PORT=8080
PUBLIC_HOST=https://api.example.com

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=ecom

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRATION_IN_SECONDS=604800
```

---

## IV. 后端服务 (Backing Services)

**原则**: 把后端服务当作附加资源

**实现**:
- MySQL 数据库通过环境变量配置连接信息
- 可以轻松切换本地和远程数据库
- 不区分本地服务和第三方服务

**配置方式**:
```bash
# 本地数据库
DB_HOST=localhost
DB_PORT=3306

# 远程数据库
DB_HOST=prod-mysql.example.com
DB_PORT=3306
```

---

## V. 构建、发布、运行 (Build, Release, Run)

**原则**: 严格分离构建和运行

**实现**:

### 构建阶段
```bash
make build
# 或使用 Docker
make docker-build
```

### 发布阶段
- 构建时注入版本信息
- Docker 镜像打标签
```bash
VERSION=v1.2.3 make docker-build
```

### 运行阶段
- 使用不可变的镜像
- 环境变量在运行时注入
```bash
docker run -e DB_HOST=... -e DB_PASSWORD=... ecommerce-api:v1.2.3
```

**相关文件**:
- `Makefile` - 构建脚本
- `Dockerfile` - 多阶段构建
- `cmd/main.go:19-24` - 版本信息

---

## VI. 进程 (Processes)

**原则**: 以一个或多个无状态进程运行应用

**实现**:
- 应用进程是无状态的
- 不依赖本地文件系统存储持久化数据
- 所有持久化数据存储在 MySQL 数据库
- 可以水平扩展多个实例

**Kubernetes 水平扩展**:
```yaml
# k8s/deployment.yaml
replicas: 3  # 运行 3 个副本

# k8s/hpa.yaml
minReplicas: 2
maxReplicas: 10  # 自动扩展
```

---

## VII. 端口绑定 (Port Binding)

**原则**: 通过端口绑定提供服务

**实现**:
- 应用自包含 HTTP 服务器（不依赖 Apache/Nginx）
- 端口通过 `PORT` 环境变量配置
- 默认绑定到 8080 端口

**相关代码**:
- `cmd/main.go:56` - 端口绑定
- `cmd/api/api.go:55-61` - HTTP 服务器配置

**运行**:
```bash
PORT=3000 ./bin/ecommerce-golang
```

---

## VIII. 并发 (Concurrency)

**原则**: 通过进程模型进行扩展

**实现**:
- Go 原生支持并发（goroutines）
- HTTP 服务器自动处理并发请求
- 通过运行多个进程实例来扩展
- Kubernetes HPA 根据 CPU/内存自动扩展

**水平扩展**:
```bash
# Kubernetes
kubectl scale deployment ecommerce-api --replicas=5

# 或使用 HPA 自动扩展
# 配置见 k8s/hpa.yaml
```

---

## IX. 易处理 (Disposability)

**原则**: 快速启动和优雅终止

**实现**:

### 快速启动
- 编译的 Go 二进制文件启动快速
- Docker 镜像优化，使用 alpine 基础镜像
- 健康检查确保服务就绪

### 优雅关闭
- 监听 SIGTERM 和 SIGINT 信号
- 30 秒优雅关闭时间
- 等待进行中的请求完成
- 正确关闭数据库连接

**相关代码**:
- `cmd/main.go:68-98` - 优雅关闭实现
- `cmd/api/api.go:66-76` - Shutdown 方法

**Kubernetes 配置**:
```yaml
# k8s/deployment.yaml
terminationGracePeriodSeconds: 30
```

---

## X. 开发环境与线上环境等价 (Dev/Prod Parity)

**原则**: 尽可能保持开发、预发布、线上环境相同

**实现**:

### 相同的技术栈
- 开发和生产都使用 MySQL 8.0
- Docker Compose 用于本地开发
- Kubernetes 用于生产部署

### 相同的部署方式
```bash
# 开发环境
docker-compose up

# 生产环境（使用相同的镜像）
kubectl apply -f k8s/
```

### 时间差距小
- 使用 CI/CD 实现快速部署
- Git commit 可以快速部署到生产

### 人员差距小
- 开发人员可以访问日志和监控
- 使用相同的工具和流程

**相关文件**:
- `docker-compose.yml` - 本地开发环境
- `k8s/` - 生产环境配置

---

## XI. 日志 (Logs)

**原则**: 把日志当作事件流

**实现**:
- 结构化日志（JSON 格式）
- 日志输出到 stdout/stderr
- 不管理日志文件
- 由运行环境收集和聚合日志

### 日志格式
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Starting ecommerce API server",
  "service": "ecommerce-api",
  "fields": {
    "version": "v1.2.3",
    "environment": "production"
  },
  "file": "main.go",
  "line": 36
}
```

**相关文件**:
- `utils/logger.go` - 结构化日志实现
- `cmd/main.go:28` - 日志初始化

**配置**:
```bash
LOG_LEVEL=info    # debug, info, warn, error, fatal
LOG_FORMAT=json   # json 或 text
```

**查看日志**:
```bash
# Docker
docker-compose logs -f api

# Kubernetes
kubectl logs -f -n ecommerce -l app=ecommerce-api
```

---

## XII. 管理进程 (Admin Processes)

**原则**: 后台管理任务当作一次性进程运行

**实现**:
- 数据库迁移作为独立进程运行
- 使用相同的代码库和配置
- 在相同的环境中运行

### 数据库迁移
```bash
# 本地运行
make migrate-up
make migrate-down

# Docker 环境
docker-compose exec api go run cmd/migrate/main.go up

# Kubernetes（一次性 Job）
kubectl run -it --rm migrate \
  --image=ecommerce-api:latest \
  --restart=Never \
  --env-file=.env \
  -- go run cmd/migrate/main.go up
```

**相关文件**:
- `cmd/migrate/main.go` - 迁移工具
- `cmd/migrate/migrations/` - 迁移脚本

---

## 最佳实践总结

### ✅ 已实现的功能

1. **配置管理**: 完全通过环境变量配置
2. **结构化日志**: JSON 格式，输出到 stdout
3. **优雅关闭**: 处理系统信号，正确清理资源
4. **健康检查**: `/health` 和 `/ready` 端点
5. **容器化**: 多阶段 Dockerfile，非 root 用户运行
6. **Kubernetes**: 完整的部署配置，包括 HPA
7. **版本管理**: 构建时注入版本信息
8. **依赖管理**: Go modules 管理依赖

### 🚀 使用指南

#### 本地开发
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件设置密码等

# 运行开发环境
docker-compose up
```

#### 生产部署
```bash
# 构建镜像
VERSION=v1.0.0 make docker-build

# 推送镜像到仓库
docker push your-registry/ecommerce-api:v1.0.0

# 部署到 Kubernetes
kubectl apply -f k8s/
```

### 📚 相关资源

- [The Twelve-Factor App](https://12factor.net/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
