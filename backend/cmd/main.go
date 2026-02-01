package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"static-mirrors/internal/admin"
	"static-mirrors/internal/cache"
	"static-mirrors/internal/proxy"
	"static-mirrors/internal/stats"
	"static-mirrors/pkg/config"
	"static-mirrors/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 全局变量
var (
	proxyService *proxy.Proxy
	cacheService cache.Cache
	statsService stats.Stats
	adminService *admin.Admin
)

func main() {
	// 加载配置文件
	configPath := "../../config/config.yaml"
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 获取配置
	cfg := config.GetConfig()

	// 初始化缓存服务
	var err error
	if cacheService, err = cache.NewCache(cfg); err != nil {
		log.Printf("初始化缓存服务失败: %v，将禁用缓存", err)
	}

	// 初始化统计服务
	if statsService, err = stats.NewStats(cfg); err != nil {
		log.Printf("初始化统计服务失败: %v，将禁用统计", err)
	}

	// 初始化反代服务
	proxyService = proxy.NewProxy(cfg)

	// 初始化后台管理服务
	adminService = admin.NewAdmin(cfg)

	// 设置Gin模式
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(utils.PerformanceMiddleware())
	r.Use(utils.SecurityMiddleware())
	r.Use(utils.ErrorHandlerMiddleware())

	// 添加CDN缓存配置中间件
	r.Use(utils.CDNCacheMiddleware())

	// 添加速率限制中间件（每分钟60个请求）
	r.Use(utils.RateLimitMiddleware(60, time.Minute))

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Printf("服务器启动在 %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
		os.Exit(1)
	}
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine) {
	// 注册后台管理路由
	if adminService != nil {
		adminService.RegisterRoutes(r.Group("/api"))
	}
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "前端文件公益镜像服务运行正常",
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "前端文件公益镜像服务",
			"version":     "1.0.0",
			"description": "为中国大陆开发者提供前端库镜像加速服务",
			"sources": []string{
				"cdn.jsdelivr.net",
				"cdnjs.cloudflare.com",
				"ghcr.io",
				"registry-1.docker.io",
				"unpkg.com",
			},
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// URL处理API
		api.POST("/process-url", func(c *gin.Context) {
			var req struct {
				URL string `json:"url" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "无效的URL"})
				return
			}

			// 生成加速后的URL
			acceleratedURL := fmt.Sprintf("/mirror?url=%s", url.QueryEscape(req.URL))

			c.JSON(200, gin.H{
				"original_url":    req.URL,
				"accelerated_url": acceleratedURL,
				"message":         "URL处理成功",
			})
		})

		// 延迟测试API
		api.POST("/test-latency", func(c *gin.Context) {
			var req struct {
				URL string `json:"url" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "无效的URL"})
				return
			}

			// 测试原地址延迟
			originalLatency, err := proxyService.TestLatency(req.URL)
			if err != nil {
				c.JSON(500, gin.H{"error": "测试原地址延迟失败", "details": err.Error()})
				return
			}

			// 生成加速后的URL
			acceleratedURL := fmt.Sprintf("http://localhost:1108/mirror?url=%s", url.QueryEscape(req.URL))

			// 测试加速后地址延迟
			acceleratedLatency, err := proxyService.TestLatency(acceleratedURL)
			if err != nil {
				c.JSON(500, gin.H{"error": "测试加速后地址延迟失败", "details": err.Error()})
				return
			}

			// 计算改进百分比
			improvement := float64(0)
			if originalLatency > 0 {
				improvement = (float64(originalLatency) - float64(acceleratedLatency)) / float64(originalLatency) * 100
			}

			c.JSON(200, gin.H{
				"original_url":        req.URL,
				"accelerated_url":     acceleratedURL,
				"original_latency":    originalLatency,
				"accelerated_latency": acceleratedLatency,
				"improvement":         fmt.Sprintf("%.1f%%", improvement),
			})
		})

		// 统计数据API
		api.GET("/stats", func(c *gin.Context) {
			if statsService == nil {
				c.JSON(200, gin.H{
					"requests":    0,
					"bandwidth":   "0 B",
					"top_sources": []string{},
				})
				return
			}

			// 添加错误处理
			defer func() {
				if r := recover(); r != nil {
					// 捕获到错误，返回默认的统计数据
					c.JSON(200, gin.H{
						"requests":       0,
						"bandwidth":      "0 B",
						"top_sources":    []string{},
						"today_requests": 0,
						"today_traffic":  "0 B",
					})
				}
			}()

			// 检查统计服务是否初始化
			if statsService == nil {
				c.JSON(200, gin.H{
					"requests":       0,
					"bandwidth":      "0 B",
					"top_sources":    []string{},
					"today_requests": 0,
					"today_traffic":  "0 B",
				})
				return
			}

			// 获取统计信息
			statsData, err := statsService.GetStats()
			if err != nil {
				c.JSON(500, gin.H{"error": "获取统计数据失败", "details": err.Error()})
				return
			}

			// 获取热门源站
			topSources, err := statsService.GetTopSources()
			if err != nil {
				topSources = []string{}
			}

			// 格式化流量数据
			traffic := int64(0)
			if t, ok := statsData["total_traffic"].(int64); ok {
				traffic = t
			}

			bandwidth := formatBytes(traffic)

			// 处理今日流量
			todayTraffic := int64(0)
			if t, ok := statsData["today_traffic"].(int64); ok {
				todayTraffic = t
			}

			c.JSON(200, gin.H{
				"requests":       statsData["total_requests"],
				"bandwidth":      bandwidth,
				"top_sources":    topSources,
				"today_requests": statsData["today_requests"],
				"today_traffic":  formatBytes(todayTraffic),
			})
		})
	}

	// 镜像服务路由
	r.Any("/mirror", func(c *gin.Context) {
		start := time.Now()

		// 获取目标URL
		targetURL := c.Query("url")
		if targetURL == "" {
			c.JSON(400, gin.H{"error": "缺少URL参数"})
			return
		}

		// 解析URL获取源站
		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(400, gin.H{"error": "无效的URL格式"})
			return
		}

		source := parsedURL.Host

		// 处理镜像请求
		proxyService.HandleMirror(c)

		// 记录访问统计
		if statsService != nil {
			duration := time.Since(start)
			// 这里简化处理，实际应该从响应中获取字节数
			bytes := int64(0)
			statsService.RecordRequest(targetURL, source, bytes, duration)
		}
	})
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
