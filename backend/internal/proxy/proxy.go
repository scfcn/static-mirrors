package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"static-mirrors/pkg/config"

	"github.com/gin-gonic/gin"
)

// Proxy 反代服务结构
type Proxy struct {
	client *http.Client
	config config.Config
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
		config: cfg,
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

// setCacheHeaders 设置缓存头
func (p *Proxy) setCacheHeaders(c *gin.Context, resp *http.Response) {
	// 如果源站已经设置了缓存头，使用源站的缓存头
	if resp.Header.Get("Cache-Control") != "" {
		return
	}

	// 根据文件类型设置合理的缓存时间
	contentType := resp.Header.Get("Content-Type")
	var cacheControl string

	switch {
	case strings.Contains(contentType, "text/css"):
		cacheControl = "public, max-age=86400"
	case strings.Contains(contentType, "application/javascript"):
		cacheControl = "public, max-age=86400"
	case strings.Contains(contentType, "image/"):
		cacheControl = "public, max-age=604800"
	case strings.Contains(contentType, "font/"):
		cacheControl = "public, max-age=604800"
	default:
		cacheControl = fmt.Sprintf("public, max-age=%d", p.config.Cache.TTL.Default)
	}

	c.Header("Cache-Control", cacheControl)
	c.Header("X-Mirror-Cache", "MISS") // 将来会根据缓存状态修改
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
