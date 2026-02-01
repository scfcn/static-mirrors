# 前端构建阶段
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json文件，利用Docker缓存
COPY package*.json ./

# 安装依赖，使用npm ci确保依赖版本一致性
RUN npm ci --no-audit --no-fund

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 前端运行阶段
FROM nginx:1.25-alpine AS runner

# 设置环境变量
ENV TZ=Asia/Shanghai

# 安装必要的依赖
RUN apk add --no-cache tzdata

# 复制构建结果到Nginx静态文件目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制Nginx配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 设置文件权限
RUN chown -R nginx:nginx /usr/share/nginx/html && chmod -R 755 /usr/share/nginx/html && \
    chown -R nginx:nginx /etc/nginx/conf.d && chmod -R 644 /etc/nginx/conf.d/*

# 配置Nginx
RUN mkdir -p /etc/nginx/ssl

# 暴露端口
EXPOSE 3000

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:3000 || exit 1

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]