# 前端构建阶段
FROM node:20-alpine as frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm install --legacy-peer-deps

# 复制前端源码
COPY frontend/ .

# 构建前端应用
RUN npm run build

# 后端构建阶段
FROM golang:1.24-alpine as backend-builder

WORKDIR /app/backend

# 复制后端依赖文件
COPY backend/go.mod backend/go.sum ./

# 下载后端依赖
RUN go mod download

# 复制后端源码
COPY backend/ .

# 构建后端应用
RUN go build -o static-mirrors ./cmd

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 复制前端构建结果
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# 复制后端构建结果
COPY --from=backend-builder /app/backend/static-mirrors /app/

# 复制配置文件
COPY config/ /app/config/

# 创建必要的目录
RUN mkdir -p /app/frontend/dist

# 安装运行时依赖
RUN apk add --no-cache ca-certificates wget

# 暴露端口
EXPOSE 1108

# 设置环境变量
ENV GIN_MODE=release

# 启动应用
CMD ["./static-mirrors"]
