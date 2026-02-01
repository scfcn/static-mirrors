# 前端文件公益镜像服务

## 项目简介

前端文件公益镜像服务是一个为中国大陆开发者提供的公益项目，旨在解决前端资源访问慢、不稳定的问题。通过反代多个国际知名的前端资源CDN，为开发者提供高速、稳定的镜像加速服务。

## 功能特性

### 核心功能
- 为访问性较差的前端库提供镜像加速服务
- 支持反代以下目标源站：
  - cdn.jsdelivr.net
  - cdnjs.cloudflare.com
  - ghcr.io
  - registry-1.docker.io
  - unpkg.com

### 前端界面
- 现代、简洁、美观的UI设计
- 服务介绍、项目说明、服务优势等内容
- URL输入功能，允许用户提交需要加速的地址
- "检测"功能，对用户输入的地址进行处理并生成加速后的URL
- 集成延迟测试功能，可同时测试原地址和加速后地址的访问延迟并进行对比展示
- 提供请求数、流量数据等运营数据展示

### 后台管理系统
- 完整的后台管理功能
- URL手动封禁机制，可对违规URL进行管理
- 详细访问统计、流量监控等运营数据展示
- 系统状态监控

### 性能与安全
- 高可用性和稳定性
- 合理的缓存策略以提高访问速度并减轻源站压力
- 标准的缓存头，用于中国大陆CDN加速缓存配置
- 完善的安全机制，防止滥用和攻击

## 技术栈

### 后端
- Go语言
- Gin框架
- Redis（可选，用于缓存和统计）
- SQLite（可选，用于本地数据存储）

### 前端
- Vue.js 3
- Vite
- Element Plus
- ECharts
- Axios

### 部署
- Docker
- Docker Compose

## 系统要求

### 硬件要求
- CPU：至少2核
- 内存：至少2GB
- 磁盘：至少20GB
- 网络：稳定的互联网连接，建议带宽至少100Mbps

### 软件要求
- Go 1.20+
- Node.js 16+
- npm 8+
- Docker（可选，用于容器化部署）
- Redis（可选，用于缓存和统计）

## 部署步骤

### 方法一：本地部署

1. **克隆项目**
   ```bash
   git clone https://github.com/scfcn/static-mirrors.git
   cd static-mirrors
   ```

2. **安装后端依赖**
   ```bash
   cd backend
   go mod download
   ```

3. **安装前端依赖**
   ```bash
   cd ../frontend
   npm install
   ```

4. **构建前端项目**
   ```bash
   npm run build
   ```

5. **配置服务**
   编辑 `config/config.yaml` 文件，根据实际情况修改配置：
   ```yaml
   # 应用配置
   app:
     name: "前端文件公益镜像服务"
     version: "1.0.0"
     host: "0.0.0.0"
     port: 1108
     debug: false
   
   # 源站配置
   sources:
     - name: "jsdelivr"
       domain: "cdn.jsdelivr.net"
       enabled: true
     - name: "cdnjs"
       domain: "cdnjs.cloudflare.com"
       enabled: true
     - name: "ghcr"
       domain: "ghcr.io"
       enabled: true
     - name: "docker"
       domain: "registry-1.docker.io"
       enabled: true
     - name: "unpkg"
       domain: "unpkg.com"
       enabled: true
   
   # 缓存配置
   cache:
     enabled: true
     type: "memory"  # redis 或 memory
     redis:
       addr: "localhost:6379"
       password: ""
       db: 0
     memory:
       size: 1024  # MB
     ttl:
       default: 3600  # 秒
       max: 86400    # 秒
   
   # 统计配置
   stats:
     enabled: true
     type: "sqlite"  # sqlite 或 redis
     sqlite:
       path: "stats.db"
     redis:
       addr: "localhost:6379"
       password: ""
       db: 1
   
   # 安全配置
   security:
     rate_limit:
       enabled: true
       requests_per_minute: 60
     blocked_urls:
       - "*.malicious.com"
       - "/harmful/path"
   
   # 日志配置
   log:
     level: "info"
     format: "json"
   ```

6. **启动后端服务**
   ```bash
   cd ../backend/cmd
   go run main.go
   ```

7. **访问服务**
   打开浏览器，访问 `http://localhost:1108` 即可使用前端文件公益镜像服务。

### 方法二：Docker部署

1. **克隆项目**
   ```bash
   git clone https://github.com/scfcn/static-mirrors.git
   cd static-mirrors
   ```

2. **创建Docker Compose配置文件**
   创建 `docker-compose.yml` 文件：
   ```yaml
   version: '3'
   services:
     backend:
       build:
         context: ./backend
         dockerfile: Dockerfile
       ports:
         - "1108:1108"
       volumes:
         - ./config:/app/config
         - ./backend/stats.db:/app/stats.db
       environment:
         - GIN_MODE=release
       restart: always
     
     frontend:
       build:
         context: ./frontend
         dockerfile: Dockerfile
       ports:
         - "3000:3000"
       depends_on:
         - backend
       restart: always
     
     redis:
       image: redis:alpine
       ports:
         - "6379:6379"
       restart: always
   ```

3. **创建Dockerfile**

   **后端Dockerfile** (`backend/Dockerfile`)：
   ```dockerfile
   FROM golang:1.20-alpine as builder
   
   WORKDIR /app
   
   COPY go.mod go.sum ./
   RUN go mod download
   
   COPY . .
   
   RUN go build -o static-mirrors ./cmd
   
   FROM alpine:latest
   
   WORKDIR /app
   
   COPY --from=builder /app/static-mirrors .
   COPY --from=builder /app/config /app/config
   
   EXPOSE 1108
   
   CMD ["./static-mirrors"]
   ```

   **前端Dockerfile** (`frontend/Dockerfile`)：
   ```dockerfile
   FROM node:16-alpine as builder
   
   WORKDIR /app
   
   COPY package.json package-lock.json ./
   RUN npm install
   
   COPY . .
   RUN npm run build
   
   FROM nginx:alpine
   
   COPY --from=builder /app/dist /usr/share/nginx/html
   COPY nginx.conf /etc/nginx/conf.d/default.conf
   
   EXPOSE 3000
   
   CMD ["nginx", "-g", "daemon off;"]
   ```

4. **创建Nginx配置文件**
   创建 `frontend/nginx.conf` 文件：
   ```nginx
   server {
     listen 3000;
     server_name localhost;
     
     location / {
       root /usr/share/nginx/html;
       index index.html index.htm;
       try_files $uri $uri/ /index.html;
     }
     
     location /api {
       proxy_pass http://backend:1108;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       proxy_set_header X-Forwarded-Proto $scheme;
     }
     
     location /mirror {
       proxy_pass http://backend:1108;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       proxy_set_header X-Forwarded-Proto $scheme;
     }
   }
   ```

5. **启动服务**
   ```bash
   docker-compose up -d
   ```

6. **访问服务**
   打开浏览器，访问 `http://localhost:3000` 即可使用前端文件公益镜像服务。

## 配置说明

### 配置文件结构

配置文件位于 `config/config.yaml`，包含以下主要部分：

1. **app**：应用基本配置，如名称、版本、监听地址等
2. **sources**：源站配置，可启用或禁用特定源站
3. **cache**：缓存配置，包括缓存类型、大小、过期时间等
4. **stats**：统计配置，包括统计类型、存储位置等
5. **security**：安全配置，包括速率限制、封禁URL等
6. **log**：日志配置，包括日志级别、格式等

### 关键配置项

- **app.port**：后端服务监听端口
- **sources**：需要反代的源站列表
- **cache.type**：缓存类型，可选值为 `redis` 或 `memory`
- **cache.redis.addr**：Redis服务器地址
- **stats.type**：统计存储类型，可选值为 `sqlite` 或 `redis`
- **security.rate_limit.requests_per_minute**：每分钟最大请求数
- **security.blocked_urls**：需要封禁的URL模式列表

## 使用方法

### 基本使用

1. **访问首页**
   打开浏览器，访问服务地址，即可看到首页，包含服务介绍、项目说明、服务优势等内容。

2. **加速URL**
   在 "加速工具" 部分，输入需要加速的URL，例如：`https://cdn.jsdelivr.net/npm/vue@3/dist/vue.global.js`，然后点击 "处理" 按钮，系统会生成加速后的URL。

3. **测试延迟**
   在 "加速工具" 部分，输入需要测试的URL，然后点击 "测试延迟" 按钮，系统会同时测试原地址和加速后地址的访问延迟并进行对比展示。

4. **查看运营数据**
   在 "数据" 部分，可以查看服务的运营数据，包括总请求数、总流量、今日请求数、今日流量、热门源站等。

### 后台管理

1. **登录**
   访问 `http://localhost:1108/api/admin/login`，使用默认账号密码登录：
   - 用户名：`admin`
   - 密码：`admin123`

2. **仪表盘**
   登录后，访问 `http://localhost:1108/api/admin/dashboard` 即可查看仪表盘，包含今日请求数、今日流量、总请求数、总流量、热门源站等信息。

3. **URL封禁管理**
   - 封禁URL：发送POST请求到 `http://localhost:1108/api/admin/block-url`，请求体为 `{"url": "https://example.com/malicious.js", "reason": "恶意脚本"}`
   - 解封URL：发送DELETE请求到 `http://localhost:1108/api/admin/unblock-url`，请求体为 `{"url": "https://example.com/malicious.js"}`
   - 查看封禁列表：发送GET请求到 `http://localhost:1108/api/admin/blocked-urls`

4. **访问统计**
   发送GET请求到 `http://localhost:1108/api/admin/stats` 即可查看详细的访问统计数据，包括每日统计和源站统计。

5. **系统状态**
   发送GET请求到 `http://localhost:1108/api/admin/system` 即可查看系统状态，包括系统信息、服务状态等。

## 维护指南

### 日志管理

- **后端日志**：默认输出到控制台，可以通过修改配置文件将日志输出到文件
- **前端日志**：前端日志主要通过浏览器控制台查看

### 监控与告警

- **系统监控**：建议使用Prometheus + Grafana对服务进行监控
- **告警**：可以配置基于CPU、内存、磁盘、网络等指标的告警

### 常见问题排查

1. **服务无法启动**
   - 检查端口是否被占用
   - 检查配置文件是否正确
   - 检查依赖是否安装完整

2. **加速效果不明显**
   - 检查网络连接是否稳定
   - 检查源站是否可访问
   - 检查缓存配置是否合理

3. **服务响应慢**
   - 检查服务器资源使用情况
   - 检查网络带宽是否足够
   - 检查缓存是否生效
   - 检查源站响应时间

4. **URL被封禁**
   - 检查URL是否包含违规内容
   - 联系管理员进行解封

### 定期维护

1. **备份数据**
   - 定期备份统计数据（如SQLite数据库文件）
   - 定期备份配置文件

2. **更新依赖**
   - 定期更新后端依赖
   - 定期更新前端依赖

3. **优化配置**
   - 根据实际使用情况，调整缓存配置
   - 根据流量情况，调整速率限制配置

## 常见问题

### Q: 服务支持哪些源站？
A: 服务支持以下源站：
- cdn.jsdelivr.net
- cdnjs.cloudflare.com
- ghcr.io
- registry-1.docker.io
- unpkg.com

### Q: 如何添加新的源站？
A: 在 `config/config.yaml` 文件中，修改 `sources` 部分，添加新的源站配置。

### Q: 服务的缓存策略是什么？
A: 服务采用多层缓存策略：
- 内存缓存：用于频繁访问的资源
- Redis缓存（可选）：用于分布式部署场景
- CDN缓存：通过设置标准的缓存头，支持中国大陆CDN加速缓存配置

### Q: 如何提高服务的性能？
A: 可以通过以下方式提高服务性能：
- 使用Redis作为缓存和统计存储
- 优化服务器硬件配置，增加CPU、内存、带宽
- 合理配置缓存策略，增加缓存大小和过期时间
- 使用CDN加速静态资源

### Q: 如何保证服务的安全性？
A: 服务通过以下方式保证安全性：
- 速率限制，防止滥用
- URL封禁机制，防止访问违规内容
- 安全HTTP头，防止XSS、CSRF等攻击
- 输入验证，防止恶意输入

## 贡献指南

### 提交代码

1. **Fork仓库**
   在GitHub上Fork项目仓库。

2. **克隆仓库**
   ```bash
   git clone https://github.com/scfcn/static-mirrors.git
   cd static-mirrors
   ```

3. **创建分支**
   ```bash
   git checkout -b feature/your-feature
   ```

4. **提交修改**
   ```bash
   git add .
   git commit -m "Add your feature"
   ```

5. **推送分支**
   ```bash
   git push origin feature/your-feature
   ```

6. **创建Pull Request**
   在GitHub上创建Pull Request，描述你的修改内容和目的。

### 代码规范

- **后端**：遵循Go语言代码规范，使用 `gofmt` 格式化代码
- **前端**：遵循Vue.js代码规范，使用ESLint检查代码
- **提交信息**：使用清晰、简洁的提交信息，描述修改内容和目的

## 许可证

前端文件公益镜像服务采用MIT许可证，详见 `LICENSE` 文件。

## 联系方式

- **项目地址**：https://github.com/scfcn/static-mirrors
- **问题反馈**：https://github.com/scfcn/static-mirrors/issues
- **贡献指南**：详见上文 "贡献指南" 部分

---

**前端文件公益镜像服务** - 为中国大陆开发者提供稳定、高效的前端资源加速服务。
