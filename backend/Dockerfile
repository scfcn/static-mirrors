# 后端构建阶段
FROM golang:1.20-alpine AS builder

# 设置构建参数
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# 安装必要的构建依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件，利用Docker缓存
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用，启用CGO，添加构建信息
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-w -s -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" \
    -o static-mirrors ./cmd

# 最终运行阶段
FROM alpine:3.18 AS runner

# 设置环境变量
ENV GIN_MODE=release
ENV TZ=Asia/Shanghai

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=builder /app/static-mirrors /app/

# 复制配置文件
COPY --from=builder /app/config /app/config

# 创建必要的目录和文件
RUN mkdir -p /app/data && touch /app/stats.db

# 设置文件权限
RUN chmod +x /app/static-mirrors && chmod 644 /app/stats.db

# 暴露端口
EXPOSE 1108

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:1108/health || exit 1

# 启动应用
CMD ["./static-mirrors"]