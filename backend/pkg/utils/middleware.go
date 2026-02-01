package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceMiddleware 性能优化中间件
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 添加响应时间头
		c.Header("X-Response-Time", fmt.Sprintf("%v", latency))

		// 记录请求处理时间
		// 这里可以添加日志记录，例如：
		// log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, latency)
	}
}

// SecurityMiddleware 安全防护中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全相关的HTTP头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 处理CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	// 简单的内存速率限制实现
	// 生产环境中建议使用Redis等分布式方案
	requests := make(map[string][]time.Time)

	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 清理过期的请求记录
		now := time.Now()
		if reqs, exists := requests[clientIP]; exists {
			var validReqs []time.Time
			for _, reqTime := range reqs {
				if now.Sub(reqTime) < duration {
					validReqs = append(validReqs, reqTime)
				}
			}
			requests[clientIP] = validReqs
		}

		// 检查请求数是否超过限制
		if len(requests[clientIP]) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 记录新的请求
		requests[clientIP] = append(requests[clientIP], now)

		c.Next()
	}
}

// CDNCacheMiddleware CDN缓存配置中间件
func CDNCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为静态资源设置CDN缓存头
		path := c.Request.URL.Path

		// 根据文件类型设置不同的缓存时间
		switch {
		case hasSuffix(path, ".js"), hasSuffix(path, ".css"):
			// JavaScript和CSS文件缓存1天
			c.Header("Cache-Control", "public, max-age=86400, stale-while-revalidate=3600")
		case hasSuffix(path, ".jpg"), hasSuffix(path, ".jpeg"), hasSuffix(path, ".png"), hasSuffix(path, ".gif"), hasSuffix(path, ".webp"):
			// 图片文件缓存7天
			c.Header("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")
		case hasSuffix(path, ".ico"), hasSuffix(path, ".woff"), hasSuffix(path, ".woff2"), hasSuffix(path, ".ttf"), hasSuffix(path, ".eot"):
			// 字体和图标文件缓存30天
			c.Header("Cache-Control", "public, max-age=2592000, stale-while-revalidate=604800")
		default:
			// 其他文件缓存1小时
			c.Header("Cache-Control", "public, max-age=3600, stale-while-revalidate=600")
		}

		// 添加CDN相关的头
		c.Header("X-Cache-Status", "MISS") // 将来可以根据实际缓存状态修改
		c.Header("X-CDN-Provider", "Static-Mirrors")

		c.Next()
	}
}

// hasSuffix 检查字符串是否以指定后缀结尾
func hasSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			// 记录错误
			for _, err := range c.Errors {
				fmt.Printf("Error: %v\n", err)
			}

			// 如果没有设置响应状态码，设置为500
			if c.Writer.Status() == 200 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "内部服务器错误",
				})
			}
		}
	}
}
