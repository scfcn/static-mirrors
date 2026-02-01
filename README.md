# Static Mirrors

![Static Mirrors](https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=GitHub%20repository%20banner%20for%20static-mirrors%20project%2C%20showing%20a%20network%20of%20mirror%20servers%20with%20data%20flow%2C%20modern%20tech%20design%2C%20blue%20and%20white%20color%20scheme&image_size=landscape_16_9)

## 项目简介

Static Mirrors 是一个前端文件公益镜像服务，旨在为开发者提供稳定、快速的前端资源访问。该项目通过镜像热门前端库和资源，帮助开发者在网络环境不佳时仍能快速获取所需资源。

## 主要功能

- 🚀 **多源镜像**：支持 jsdelivr、cdnjs、ghcr、docker、unpkg 等多个源站的镜像
- 📊 **统计分析**：提供请求量、流量、来源等统计数据
- ⚡ **缓存机制**：支持 Redis 和内存缓存，提高访问速度
- 🔧 **易于部署**：提供完整的 Docker 部署方案
- 📱 **管理界面**：内置简单的管理后台
- 🔒 **安全可靠**：支持 HTTPS，保障数据传输安全

## 技术栈

### 后端
- **语言**：Go 1.24
- **Web 框架**：Gin
- **缓存**：Redis / 内存缓存
- **数据库**：SQLite
- **配置管理**：Viper

### 前端
- **框架**：Vue 3
- **构建工具**：Vite
- **包管理**：npm
- **部署**：Nginx

### 部署
- **容器化**：Docker
- **编排**：Docker Compose
- **CI/CD**：GitHub Actions

## 快速开始

### 环境要求
- Docker
- Docker Compose

### 部署步骤

1. **克隆仓库**

```bash
git clone https://github.com/scfcn/static-mirrors.git
cd static-mirrors
```

2. **配置环境变量**

复制 `.env.example` 文件为 `.env` 并根据实际情况修改配置：

```bash
cp .env.example .env
# 编辑 .env 文件
```

3. **启动服务**

```bash
docker compose up -d
```

服务启动后，可通过 `http://localhost:1108` 访问。

## 配置说明

### 主要配置文件

- **`config/config.yaml`**：应用程序主配置文件
- **`.env`**：环境变量配置
- **`docker-compose.yml`**：Docker 部署配置

### 配置项说明

#### config.yaml

```yaml
# 服务器配置
server:
  port: 1108
  host: "0.0.0.0"
  timeout: 30s

# 缓存配置
cache:
  enabled: true
  type: "redis"  # redis 或 memory
  redis:
    addr: "redis:6379"
    password: ""
    db: 0
  memory:
    size: 100mb
  ttl:
    default: 24h
    min: 1h
    max: 72h

# 统计配置
stats:
  enabled: true
  type: "sqlite"  # sqlite 或 redis
  sqlite:
    path: "./data/stats.db"
  redis:
    addr: "redis:6379"
    password: ""
    db: 0

# 源站配置
sources:
  jsdelivr:
    enabled: true
    base_url: "https://cdn.jsdelivr.net"
  cdnjs:
    enabled: true
    base_url: "https://cdnjs.cloudflare.com"
  unpkg:
    enabled: true
    base_url: "https://unpkg.com"
  ghcr:
    enabled: true
    base_url: "https://ghcr.io"
  docker:
    enabled: true
    base_url: "https://registry-1.docker.io"

# 管理后台配置
admin:
  enabled: true
  username: "admin"
  password: "admin123"
  path: "/admin"
```

## 部署方式

### Docker Compose 部署

这是推荐的部署方式，包含了所有必要的服务：

```bash
docker compose up -d
```

### 手动部署

1. **构建后端**

```bash
cd backend
go build -o static-mirrors ./cmd/main.go
```

2. **构建前端**

```bash
cd frontend
npm install
npm run build
```

3. **启动服务**

```bash
# 启动后端
./backend/static-mirrors

# 启动前端（使用 Nginx 或其他 Web 服务器）
# 配置 Nginx 指向 frontend/dist 目录
```

## 开发指南

### 后端开发

```bash
cd backend
# 安装依赖
go mod tidy
# 运行开发服务器
go run ./cmd/main.go
```

### 前端开发

```bash
cd frontend
# 安装依赖
npm install
# 运行开发服务器
npm run dev
```

## CI/CD 流程

项目使用 GitHub Actions 进行 CI/CD，主要流程包括：

1. **版本生成**：使用 semantic-release 自动生成语义化版本
2. **Docker 构建**：构建并推送 Docker 镜像到 Docker Hub 和 GitHub Container Registry
3. **发布**：自动创建 GitHub Release

## 贡献指南

1. **Fork 仓库**
2. **创建分支**：`git checkout -b feature/your-feature`
3. **提交更改**：`git commit -m "feat: add your feature"`
4. **推送分支**：`git push origin feature/your-feature`
5. **创建 Pull Request**

## 许可证

本项目使用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 联系方式

- **GitHub**：[https://github.com/scfcn/static-mirrors](https://github.com/scfcn/static-mirrors)
- **Issues**：[https://github.com/scfcn/static-mirrors/issues](https://github.com/scfcn/static-mirrors/issues)

---

**Static Mirrors** - 让前端资源访问更加稳定可靠！