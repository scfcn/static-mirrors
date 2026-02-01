package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Sources  []SourceConfig `yaml:"sources"`
	Cache    CacheConfig    `yaml:"cache"`
	Stats    StatsConfig    `yaml:"stats"`
	Security SecurityConfig `yaml:"security"`
	Log      LogConfig      `yaml:"log"`
}

// AppConfig 应用基本配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Debug   bool   `yaml:"debug"`
}

// SourceConfig 源站配置
type SourceConfig struct {
	Name    string `yaml:"name"`
	Domain  string `yaml:"domain"`
	Enabled bool   `yaml:"enabled"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool           `yaml:"enabled"`
	Type    string         `yaml:"type"`
	Redis   RedisConfig    `yaml:"redis"`
	Memory  MemoryConfig   `yaml:"memory"`
	TTL     CacheTTLConfig `yaml:"ttl"`
}

// RedisConfig Redis缓存配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// MemoryConfig 内存缓存配置
type MemoryConfig struct {
	Size int `yaml:"size"`
}

// CacheTTLConfig 缓存过期时间配置
type CacheTTLConfig struct {
	Default int `yaml:"default"`
	Max     int `yaml:"max"`
}

// StatsConfig 统计配置
type StatsConfig struct {
	Enabled bool         `yaml:"enabled"`
	Type    string       `yaml:"type"`
	SQLite  SQLiteConfig `yaml:"sqlite"`
	Redis   RedisConfig  `yaml:"redis"`
}

// SQLiteConfig SQLite配置
type SQLiteConfig struct {
	Path string `yaml:"path"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	RateLimit   RateLimitConfig `yaml:"rate_limit"`
	BlockedURLs []string        `yaml:"blocked_urls"`
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// 全局配置实例
var GlobalConfig Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Println("配置文件加载成功")
	return nil
}

// GetConfig 获取全局配置
func GetConfig() Config {
	return GlobalConfig
}

// GetSourceConfig 根据域名获取源站配置
func GetSourceConfig(domain string) *SourceConfig {
	for _, source := range GlobalConfig.Sources {
		if source.Domain == domain && source.Enabled {
			return &source
		}
	}
	return nil
}
