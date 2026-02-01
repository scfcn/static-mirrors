package admin

import (
	"net/http"
	"time"

	"static-mirrors/pkg/config"

	"github.com/gin-gonic/gin"
)

// Admin 后台管理结构
type Admin struct {
	config config.Config
	// 这里可以添加其他依赖，如数据库连接等
}

// NewAdmin 创建新的后台管理实例
func NewAdmin(cfg config.Config) *Admin {
	return &Admin{
		config: cfg,
	}
}

// RegisterRoutes 注册后台管理路由
func (a *Admin) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		// 登录路由
		admin.POST("/login", a.login)

		// 需要认证的路由
		auth := admin.Group("/")
		auth.Use(a.authMiddleware())
		{
			// 仪表盘
			auth.GET("/dashboard", a.getDashboard)

			// URL封禁管理
			auth.POST("/block-url", a.blockUrl)
			auth.DELETE("/unblock-url", a.unblockUrl)
			auth.GET("/blocked-urls", a.getBlockedUrls)

			// 访问统计
			auth.GET("/stats", a.getStats)

			// 系统状态
			auth.GET("/system", a.getSystemStatus)
		}
	}
}

// login 登录处理
func (a *Admin) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 这里应该实现实际的认证逻辑
	// 为了演示，我们使用简单的硬编码验证
	if req.Username == "admin" && req.Password == "admin123" {
		c.JSON(http.StatusOK, gin.H{
			"token":   "dummy-token",
			"message": "登录成功",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
	}
}

// authMiddleware 认证中间件
func (a *Admin) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || token != "Bearer dummy-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getDashboard 获取仪表盘数据
func (a *Admin) getDashboard(c *gin.Context) {
	// 这里应该实现实际的仪表盘数据获取逻辑
	c.JSON(http.StatusOK, gin.H{
		"today_requests": 1234,
		"today_traffic":  "1.2 GB",
		"total_requests": 123456,
		"total_traffic":  "123.4 GB",
		"top_sources": []string{
			"cdn.jsdelivr.net",
			"unpkg.com",
			"cdnjs.cloudflare.com",
		},
		"system_status": "正常",
		"uptime":        "7天 12小时 34分钟",
	})
}

// blockUrl 封禁URL
func (a *Admin) blockUrl(c *gin.Context) {
	var req struct {
		Url    string `json:"url" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 这里应该实现实际的URL封禁逻辑
	// 为了演示，我们直接返回成功
	c.JSON(http.StatusOK, gin.H{
		"message":    "URL封禁成功",
		"url":        req.Url,
		"reason":     req.Reason,
		"blocked_at": time.Now().Format(time.RFC3339),
	})
}

// unblockUrl 解封URL
func (a *Admin) unblockUrl(c *gin.Context) {
	var req struct {
		Url string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 这里应该实现实际的URL解封逻辑
	// 为了演示，我们直接返回成功
	c.JSON(http.StatusOK, gin.H{
		"message":      "URL解封成功",
		"url":          req.Url,
		"unblocked_at": time.Now().Format(time.RFC3339),
	})
}

// getBlockedUrls 获取被封禁的URL列表
func (a *Admin) getBlockedUrls(c *gin.Context) {
	// 这里应该实现实际的获取被封禁URL列表的逻辑
	// 为了演示，我们返回一些模拟数据
	c.JSON(http.StatusOK, gin.H{
		"blocked_urls": []map[string]interface{}{
			{
				"url":        "https://example.com/malicious.js",
				"reason":     "恶意脚本",
				"blocked_at": "2026-01-01T00:00:00Z",
			},
			{
				"url":        "https://example.com/harmful.css",
				"reason":     "有害内容",
				"blocked_at": "2026-01-02T00:00:00Z",
			},
		},
		"total": 2,
	})
}

// getStats 获取详细的访问统计
func (a *Admin) getStats(c *gin.Context) {
	// 这里应该实现实际的获取详细访问统计的逻辑
	// 为了演示，我们返回一些模拟数据
	c.JSON(http.StatusOK, gin.H{
		"daily_stats": []map[string]interface{}{
			{
				"date":     "2026-01-29",
				"requests": 1234,
				"traffic":  "1.2 GB",
			},
			{
				"date":     "2026-01-30",
				"requests": 2345,
				"traffic":  "2.3 GB",
			},
			{
				"date":     "2026-01-31",
				"requests": 3456,
				"traffic":  "3.4 GB",
			},
			{
				"date":     "2026-02-01",
				"requests": 1234,
				"traffic":  "1.2 GB",
			},
		},
		"source_stats": []map[string]interface{}{
			{
				"source":     "cdn.jsdelivr.net",
				"requests":   4567,
				"percentage": "45.6%",
			},
			{
				"source":     "unpkg.com",
				"requests":   2345,
				"percentage": "23.4%",
			},
			{
				"source":     "cdnjs.cloudflare.com",
				"requests":   1234,
				"percentage": "12.3%",
			},
			{
				"source":     "ghcr.io",
				"requests":   1000,
				"percentage": "10.0%",
			},
			{
				"source":     "registry-1.docker.io",
				"requests":   854,
				"percentage": "8.7%",
			},
		},
	})
}

// getSystemStatus 获取系统状态
func (a *Admin) getSystemStatus(c *gin.Context) {
	// 这里应该实现实际的获取系统状态的逻辑
	// 为了演示，我们返回一些模拟数据
	c.JSON(http.StatusOK, gin.H{
		"system": map[string]interface{}{
			"version":      "1.0.0",
			"uptime":       "7天 12小时 34分钟",
			"cpu_usage":    "12.3%",
			"memory_usage": "45.6%",
			"disk_usage":   "67.8%",
			"network_in":   "123 MB/s",
			"network_out":  "456 MB/s",
		},
		"services": []map[string]interface{}{
			{
				"name":          "API服务",
				"status":        "正常",
				"response_time": "12ms",
			},
			{
				"name":          "反代服务",
				"status":        "正常",
				"response_time": "23ms",
			},
			{
				"name":          "缓存服务",
				"status":        "正常",
				"response_time": "5ms",
			},
			{
				"name":          "统计服务",
				"status":        "正常",
				"response_time": "8ms",
			},
		},
	})
}
