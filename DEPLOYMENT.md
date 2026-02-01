# 部署文档

## 1. 系统架构

### 1.1 技术栈
- **前端**: Vue.js 3 + Vite + Nginx
- **后端**: Go 1.20 + Gin + SQLite/Redis
- **缓存**: Redis
- **容器化**: Docker + Docker Compose

### 1.2 服务结构
- **backend**: 后端服务，负责反向代理和API处理
- **frontend**: 前端服务，负责用户界面
- **redis**: 缓存服务，用于提高性能

### 1.3 网络架构
- 自定义网络: `static-mirrors-network` (172.20.0.0/16)
- 服务间通过容器名称通信
- 端口映射:
  - backend: 1108 → 1108
  - frontend: 3000 → 3000
  - redis: 6379 → 6379 (可选，仅用于开发)

## 2. 部署准备

### 2.1 环境要求
- Docker 20.10.0+
- Docker Compose 1.29.0+
- Git (可选，用于获取提交信息)

### 2.2 目录结构
```
static-mirrors/
├── backend/           # 后端代码
├── frontend/          # 前端代码
├── config/            # 配置文件
├── docker/            # Docker配置
├── deploy.sh          # Linux部署脚本
├── deploy.bat         # Windows部署脚本
├── docker-compose.yml # 服务编排配置
└── .env               # 环境变量配置 (可选)
```

### 2.3 环境变量配置

创建 `.env` 文件，配置以下环境变量：

```env
# Redis配置
REDIS_PASSWORD=your_redis_password

# 缓存配置
CACHE_TYPE=redis       # redis 或 memory
CACHE_TTL=3600         # 缓存过期时间（秒）

# 构建信息
GIT_COMMIT=auto        # Git提交信息（自动获取）
BUILD_DATE=auto        # 构建日期（自动获取）
```

## 3. 部署步骤

### 3.1 使用部署脚本

#### Linux/macOS
```bash
# 赋予脚本执行权限
chmod +x deploy.sh

# 构建镜像
./deploy.sh build

# 启动服务
./deploy.sh up

# 查看服务状态
./deploy.sh status
```

#### Windows
```batch
# 构建镜像
deploy.bat build

# 启动服务
deploy.bat up

# 查看服务状态
deploy.bat status
```

### 3.2 手动部署

#### 构建镜像
```bash
# 构建前端镜像
docker build -t static-mirrors-frontend:latest ./frontend

# 构建后端镜像
docker build \
    --build-arg GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
    --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
    -t static-mirrors-backend:latest \
    ./backend
```

#### 启动服务
```bash
docker-compose up -d
```

#### 停止服务
```bash
docker-compose down
```

## 4. 服务管理

### 4.1 常用命令

| 命令 | 描述 | 示例 |
|------|------|------|
| build | 构建镜像 | `./deploy.sh build` |
| up | 启动服务 | `./deploy.sh up` |
| down | 停止服务 | `./deploy.sh down` |
| restart | 重启服务 | `./deploy.sh restart` |
| status | 查看服务状态 | `./deploy.sh status` |
| logs | 查看服务日志 | `./deploy.sh logs backend` |
| prune | 清理Docker资源 | `./deploy.sh prune` |
| healthcheck | 检查服务健康状态 | `./deploy.sh healthcheck` |

### 4.2 服务访问

- **前端**: http://localhost:3000
- **后端API**: http://localhost:1108
- **健康检查**: http://localhost:1108/health

### 4.3 配置管理

#### 后端配置
修改 `config/config.yaml` 文件，配置反向代理和其他服务设置。

#### 前端配置
修改 `frontend/nginx.conf` 文件，配置Nginx反向代理设置。

## 5. 性能优化

### 5.1 镜像优化
- **多阶段构建**: 减少最终镜像大小
- **缓存利用**: 合理使用Docker缓存层
- **Alpine基础镜像**: 使用轻量级基础镜像

### 5.2 运行时优化
- **资源限制**: 为每个服务设置合理的CPU和内存限制
- **健康检查**: 确保服务正常运行
- **重启策略**: 自动恢复失败的服务

### 5.3 缓存策略
- **Redis缓存**: 提高频繁访问资源的响应速度
- **缓存过期时间**: 根据资源类型设置合理的过期时间

## 6. 安全措施

### 6.1 容器安全
- 使用非root用户运行容器
- 限制容器权限
- 定期更新基础镜像

### 6.2 网络安全
- 使用自定义网络隔离服务
- 限制外部访问（生产环境）
- 配置防火墙规则

### 6.3 数据安全
- 使用数据卷持久化重要数据
- 配置Redis密码认证
- 定期备份数据

## 7. 监控与维护

### 7.1 日志管理
- **Docker日志**: 使用 `docker-compose logs` 查看服务日志
- **应用日志**: 后端服务日志输出到标准输出

### 7.2 健康检查
- **服务健康检查**: 定期检查服务状态
- **系统监控**: 监控CPU、内存和网络使用情况

### 7.3 常见问题排查

#### 服务启动失败
1. 检查Docker是否运行
2. 查看服务日志: `./deploy.sh logs backend`
3. 检查端口是否被占用

#### 反向代理失败
1. 检查源站是否可访问
2. 查看后端日志: `./deploy.sh logs backend`
3. 检查网络连接

#### 性能问题
1. 检查Redis连接是否正常
2. 监控系统资源使用情况
3. 调整缓存策略

## 8. 扩展与升级

### 8.1 水平扩展
- **后端服务**: 可通过修改 `docker-compose.yml` 增加实例数
- **Redis集群**: 生产环境建议使用Redis集群

### 8.2 版本升级
1. 拉取最新代码
2. 重新构建镜像: `./deploy.sh build`
3. 重启服务: `./deploy.sh restart`

### 8.3 环境迁移
1. 备份数据卷和配置文件
2. 在新环境中部署服务
3. 恢复数据和配置

## 9. 生产环境部署建议

### 9.1 配置调整
- **Redis**: 使用密码认证，配置持久化
- **资源限制**: 根据服务器配置调整CPU和内存限制
- **网络**: 配置域名和HTTPS

### 9.2 高可用性
- **多实例部署**: 部署多个后端实例
- **负载均衡**: 使用Nginx或其他负载均衡器
- **数据备份**: 定期备份重要数据

### 9.3 监控告警
- **Prometheus + Grafana**: 监控系统指标
- **ELK Stack**: 集中管理日志
- **告警系统**: 配置关键指标告警

## 10. 故障恢复

### 10.1 服务恢复
1. 停止异常服务: `./deploy.sh down`
2. 清理资源: `./deploy.sh prune`
3. 重新启动: `./deploy.sh up`

### 10.2 数据恢复
1. 从备份恢复数据
2. 重启服务
3. 验证数据完整性

## 11. 附录

### 11.1 Docker命令参考

```bash
# 查看容器状态
docker ps

# 进入容器
docker exec -it static-mirrors-backend sh

# 查看镜像
docker images

# 清理无用镜像
docker image prune -f
```

### 11.2 环境变量说明

| 环境变量 | 描述 | 默认值 |
|---------|------|-------|
| GIN_MODE | Gin运行模式 | release |
| TZ | 时区 | Asia/Shanghai |
| REDIS_ADDR | Redis地址 | redis:6379 |
| REDIS_PASSWORD | Redis密码 | (空) |
| REDIS_DB | Redis数据库 | 0 |
| STATS_DB_PATH | 统计数据库路径 | /app/data/stats.db |
| CACHE_TYPE | 缓存类型 | redis |
| CACHE_TTL | 缓存过期时间 | 3600 |

### 11.3 端口说明

| 服务 | 内部端口 | 外部端口 | 用途 |
|------|---------|---------|------|
| backend | 1108 | 1108 | 后端API和代理服务 |
| frontend | 3000 | 3000 | 前端Web界面 |
| redis | 6379 | 6379 | 缓存服务 |

## 12. 联系方式

如有问题或建议，请联系项目维护者。

---

**部署文档版本**: v1.0.0
**最后更新**: $(date -u +"%Y-%m-%d")