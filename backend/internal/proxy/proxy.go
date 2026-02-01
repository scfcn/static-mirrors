package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"static-mirrors/pkg/config"

	"github.com/gin-gonic/gin"
)

// Proxy 反代服务结构
type Proxy struct {
	client          *http.Client
	config          config.Config
	purgeRecords    map[string]time.Time
	purgeMutex      sync.RWMutex
	purgeCount      int
	purgeCountMutex sync.Mutex
}

// NewProxy 创建新的反代服务实例
func NewProxy(cfg config.Config) *Proxy {
	return &Proxy{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		config:       cfg,
		purgeRecords: make(map[string]time.Time),
	}
}

// HandleMirror 处理镜像请求
func (p *Proxy) HandleMirror(c *gin.Context) {
	// 获取目标URL
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(400, gin.H{"error": "缺少URL参数"})
		return
	}

	// 解析URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的URL格式"})
		return
	}

	// 验证源站是否在白名单中
	if !p.isValidSource(parsedURL.Host) {
		c.JSON(403, gin.H{"error": "不支持的源站"})
		return
	}

	// 验证URL是否被封禁
	if p.isBlockedURL(targetURL) {
		c.JSON(403, gin.H{"error": "该URL已被封禁"})
		return
	}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "创建请求失败"})
		return
	}

	// 复制请求头
	for k, v := range c.Request.Header {
		if k != "Host" && k != "Content-Length" {
			req.Header[k] = v
		}
	}

	// 设置正确的Host头
	req.Host = parsedURL.Host

	// 发送请求到源站
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "连接源站失败", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for k, v := range resp.Header {
		if k != "Content-Length" {
			c.Header(k, strings.Join(v, ", "))
		}
	}

	// 设置缓存头
	p.setCacheHeaders(c, resp)

	// 设置响应状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		log.Printf("复制响应体失败: %v", err)
	}

	// 记录访问统计
	// 这里将来会实现统计功能
}

// HandlePathProxy 处理路径代理请求
func (p *Proxy) HandlePathProxy(c *gin.Context, source string, path string) {
	// 验证源站是否在白名单中
	if !p.isValidSource(source) {
		c.JSON(403, gin.H{"error": "不支持的源站"})
		return
	}

	// 构建目标URL
	targetURL := fmt.Sprintf("https://%s%s", source, path)

	// 验证URL是否被封禁
	if p.isBlockedURL(targetURL) {
		c.JSON(403, gin.H{"error": "该URL已被封禁"})
		return
	}

	// 验证路径是否被封禁
	if p.isBlockedPath(path) {
		c.JSON(403, gin.H{"error": "该路径已被封禁"})
		return
	}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "创建请求失败"})
		return
	}

	// 复制请求头
	for k, v := range c.Request.Header {
		if k != "Host" && k != "Content-Length" {
			req.Header[k] = v
		}
	}

	// 设置正确的Host头
	req.Host = source

	// 发送请求到源站
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "连接源站失败", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for k, v := range resp.Header {
		if k != "Content-Length" {
			c.Header(k, strings.Join(v, ", "))
		}
	}

	// 设置缓存头
	p.setCacheHeaders(c, resp)

	// 设置响应状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		log.Printf("复制响应体失败: %v", err)
	}
}

// isValidSource 验证源站是否在白名单中
func (p *Proxy) isValidSource(host string) bool {
	for _, source := range p.config.Sources {
		if source.Domain == host && source.Enabled {
			return true
		}
	}
	return false
}

// isBlockedURL 验证URL是否被封禁
func (p *Proxy) isBlockedURL(url string) bool {
	for _, blockedPattern := range p.config.Security.BlockedURLs {
		if strings.Contains(url, blockedPattern) {
			return true
		}
	}
	return false
}

// isBlockedPath 验证路径是否被封禁
func (p *Proxy) isBlockedPath(path string) bool {
	blockedPaths := []string{
		"/",
		"/login",
		"/signin",
		"/signup",
		"/register",
		"/account",
		"/user",
		"/profile",
		"/settings",
		"/admin",
		"/dashboard",
		"/auth",
		"/oauth",
	}

	for _, blockedPath := range blockedPaths {
		if strings.EqualFold(path, blockedPath) || strings.HasPrefix(path, blockedPath+"/") {
			return true
		}
	}
	return false
}

// setCacheHeaders 设置缓存头
func (p *Proxy) setCacheHeaders(c *gin.Context, resp *http.Response) {
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength
	var cacheControl string

	// 如果源站已经设置了缓存头，使用源站的缓存头
	if resp.Header.Get("Cache-Control") != "" {
		cacheControl = resp.Header.Get("Cache-Control")
	} else {
		// 根据文件大小和类型设置合理的缓存时间
		cacheControl = p.calculateCacheControl(contentType, contentLength)
	}

	c.Header("Cache-Control", cacheControl)

	// 设置Expires头
	if resp.Header.Get("Expires") == "" {
		expiresTime := p.calculateExpires(cacheControl)
		c.Header("Expires", expiresTime)
	}

	// 设置ETag
	if resp.Header.Get("ETag") != "" {
		c.Header("ETag", resp.Header.Get("ETag"))
	}

	// 设置Last-Modified
	if resp.Header.Get("Last-Modified") != "" {
		c.Header("Last-Modified", resp.Header.Get("Last-Modified"))
	}

	c.Header("X-Mirror-Cache", "MISS")
}

// calculateCacheControl 根据文件类型和大小计算缓存控制
func (p *Proxy) calculateCacheControl(contentType string, contentLength int64) string {
	strategy := p.config.Cache.Strategy

	// 首先根据文件类型判断
	fileType := p.getFileType(contentType)
	if ttl, exists := strategy.FileTypes[fileType]; exists {
		return fmt.Sprintf("public, max-age=%d, s-maxage=%d, immutable", ttl, ttl)
	}

	// 根据文件大小判断
	if contentLength > strategy.LargeFileThreshold {
		return fmt.Sprintf("public, max-age=%d, s-maxage=%d, immutable", strategy.LargeFileTTL, strategy.LargeFileTTL)
	} else if contentLength > 0 && contentLength < 1048576 {
		return fmt.Sprintf("public, max-age=%d, s-maxage=%d", strategy.SmallFileTTL, strategy.SmallFileTTL)
	}

	// 默认使用普通文件缓存时间
	return fmt.Sprintf("public, max-age=%d, s-maxage=%d", strategy.NormalFileTTL, strategy.NormalFileTTL)
}

// getFileType 根据Content-Type获取文件类型
func (p *Proxy) getFileType(contentType string) string {
	switch {
	case strings.Contains(contentType, "text/css"):
		return "css"
	case strings.Contains(contentType, "application/javascript") || strings.Contains(contentType, "text/javascript"):
		return "js"
	case strings.Contains(contentType, "image/"):
		return "image"
	case strings.Contains(contentType, "font/") || strings.Contains(contentType, "application/font"):
		return "font"
	case strings.Contains(contentType, "video/"):
		return "video"
	case strings.Contains(contentType, "audio/"):
		return "audio"
	case strings.Contains(contentType, "application/zip") || strings.Contains(contentType, "application/x-tar") || strings.Contains(contentType, "application/x-rar"):
		return "archive"
	default:
		return ""
	}
}

// calculateExpires 根据Cache-Control计算Expires时间
func (p *Proxy) calculateExpires(cacheControl string) string {
	maxAge := p.extractMaxAge(cacheControl)
	if maxAge <= 0 {
		maxAge = int64(p.config.Cache.TTL.Default)
	}

	expiresTime := time.Now().Add(time.Duration(maxAge) * time.Second)
	return expiresTime.UTC().Format(time.RFC1123)
}

// extractMaxAge 从Cache-Control中提取max-age值
func (p *Proxy) extractMaxAge(cacheControl string) int64 {
	parts := strings.Split(cacheControl, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "max-age=") {
			ageStr := strings.TrimPrefix(part, "max-age=")
			var age int64
			fmt.Sscanf(ageStr, "%d", &age)
			return age
		}
	}
	return 0
}

// TestLatency 测试URL延迟
func (p *Proxy) TestLatency(url string) (int64, error) {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// 设置超时
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 读取响应体以确保完整请求
	_, _ = io.Copy(io.Discard, resp.Body)

	latency := time.Since(start).Milliseconds()
	return latency, nil
}

// HandlePurge 处理缓存刷新请求
func (p *Proxy) HandlePurge(c *gin.Context) {
	if !p.config.Cache.Purge.Enabled {
		c.JSON(403, gin.H{"error": "缓存刷新功能未启用"})
		return
	}

	// 获取要刷新的URL
	targetURL := c.Param("url")
	if targetURL == "" {
		c.JSON(400, gin.H{"error": "缺少URL参数"})
		return
	}

	// 解码URL
	targetURL, err := url.PathUnescape(targetURL)
	if err != nil {
		c.JSON(400, gin.H{"error": "URL解码失败"})
		return
	}

	// 验证URL格式
	_, err = url.Parse(targetURL)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的URL格式"})
		return
	}

	// 检查刷新频率限制
	if !p.checkPurgeRateLimit(targetURL) {
		c.JSON(429, gin.H{
			"error":   "刷新频率超限",
			"message": "每30分钟内最多允许执行一次缓存刷新操作",
		})
		return
	}

	// 检查总刷新次数限制
	if !p.checkPurgeCountLimit() {
		c.JSON(429, gin.H{
			"error":   "刷新次数超限",
			"message": "已达到最大刷新次数限制",
		})
		return
	}

	// 记录刷新操作
	p.recordPurge(targetURL)

	// 这里将来会实现实际的缓存刷新逻辑
	// 目前只是返回成功响应
	c.JSON(200, gin.H{
		"success": true,
		"message": "缓存刷新请求已处理",
		"url":     targetURL,
		"time":    time.Now().Format(time.RFC3339),
	})
}

// checkPurgeRateLimit 检查刷新频率限制
func (p *Proxy) checkPurgeRateLimit(url string) bool {
	p.purgeMutex.RLock()
	lastPurgeTime, exists := p.purgeRecords[url]
	p.purgeMutex.RUnlock()

	if !exists {
		return true
	}

	rateLimitDuration := time.Duration(p.config.Cache.Purge.RateLimitMinutes) * time.Minute
	return time.Since(lastPurgeTime) >= rateLimitDuration
}

// checkPurgeCountLimit 检查总刷新次数限制
func (p *Proxy) checkPurgeCountLimit() bool {
	p.purgeCountMutex.Lock()
	defer p.purgeCountMutex.Unlock()

	if p.purgeCount >= p.config.Cache.Purge.MaxPurgeCount {
		return false
	}

	p.purgeCount++
	return true
}

// recordPurge 记录刷新操作
func (p *Proxy) recordPurge(url string) {
	p.purgeMutex.Lock()
	defer p.purgeMutex.Unlock()

	p.purgeRecords[url] = time.Now()
}
