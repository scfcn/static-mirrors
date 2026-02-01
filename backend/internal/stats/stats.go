package stats

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"static-mirrors/pkg/config"

	_ "github.com/mattn/go-sqlite3"
)

// Stats 统计接口
type Stats interface {
	RecordRequest(url string, source string, bytes int64, duration time.Duration)
	GetStats() (map[string]interface{}, error)
	GetTopSources() ([]string, error)
	GetTraffic() (int64, error)
	GetRequests() (int64, error)
}

// SQLiteStats SQLite统计实现
type SQLiteStats struct {
	db *sql.DB
}

// NewStats 创建新的统计实例
func NewStats(cfg config.Config) (Stats, error) {
	if !cfg.Stats.Enabled {
		return nil, nil
	}

	switch cfg.Stats.Type {
	case "sqlite":
		return NewSQLiteStats(cfg.Stats.SQLite.Path)
	default:
		return nil, fmt.Errorf("不支持的统计类型: %s", cfg.Stats.Type)
	}
}

// NewSQLiteStats 创建SQLite统计实例
func NewSQLiteStats(path string) (*SQLiteStats, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("打开SQLite数据库失败: %w", err)
	}

	// 创建表
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("创建表失败: %w", err)
	}

	log.Println("SQLite统计初始化成功")
	return &SQLiteStats{db: db}, nil
}

// createTables 创建SQLite表
func createTables(db *sql.DB) error {
	// 创建请求记录表
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		source TEXT,
		bytes INTEGER,
		duration INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	// 创建每日统计表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS daily_stats (
		date TEXT PRIMARY KEY,
		requests INTEGER DEFAULT 0,
		traffic INTEGER DEFAULT 0
	)
	`)
	if err != nil {
		return err
	}

	return nil
}

// RecordRequest 记录请求到SQLite
func (s *SQLiteStats) RecordRequest(url string, source string, bytes int64, duration time.Duration) {
	// 记录请求
	_, err := s.db.Exec(
		"INSERT INTO requests (url, source, bytes, duration) VALUES (?, ?, ?, ?)",
		url, source, bytes, duration.Milliseconds(),
	)
	if err != nil {
		log.Printf("记录请求失败: %v", err)
		return
	}

	// 更新每日统计
	date := time.Now().Format("2006-01-02")
	_, err = s.db.Exec(
		"INSERT INTO daily_stats (date, requests, traffic) VALUES (?, 1, ?) ON CONFLICT(date) DO UPDATE SET requests = requests + 1, traffic = traffic + ?",
		date, bytes, bytes,
	)
	if err != nil {
		log.Printf("更新每日统计失败: %v", err)
	}
}



// GetStats 获取SQLite统计信息
func (s *SQLiteStats) GetStats() (map[string]interface{}, error) {
	// 获取总请求数
	var requests int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM requests").Scan(&requests)
	if err != nil {
		return nil, err
	}

	// 获取总流量
	var traffic int64
	err = s.db.QueryRow("SELECT COALESCE(SUM(bytes), 0) FROM requests").Scan(&traffic)
	if err != nil {
		return nil, err
	}

	// 获取今日请求数
	var todayRequests int64
	date := time.Now().Format("2006-01-02")
	err = s.db.QueryRow("SELECT COALESCE(requests, 0) FROM daily_stats WHERE date = ?", date).Scan(&todayRequests)
	if err == sql.ErrNoRows {
		todayRequests = 0
	} else if err != nil {
		return nil, err
	}

	// 获取今日流量
	var todayTraffic int64
	err = s.db.QueryRow("SELECT COALESCE(traffic, 0) FROM daily_stats WHERE date = ?", date).Scan(&todayTraffic)
	if err == sql.ErrNoRows {
		todayTraffic = 0
	} else if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_requests": requests,
		"total_traffic":  traffic,
		"today_requests": todayRequests,
		"today_traffic":  todayTraffic,
	}, nil
}



// GetTopSources 获取SQLite热门源站
func (s *SQLiteStats) GetTopSources() ([]string, error) {
	rows, err := s.db.Query(
		"SELECT source, COUNT(*) as count FROM requests GROUP BY source ORDER BY count DESC LIMIT 5",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []string
	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}



// GetTraffic 获取SQLite流量
func (s *SQLiteStats) GetTraffic() (int64, error) {
	var traffic int64
	err := s.db.QueryRow("SELECT COALESCE(SUM(bytes), 0) FROM requests").Scan(&traffic)
	if err != nil {
		return 0, err
	}
	return traffic, nil
}



// GetRequests 获取SQLite请求数
func (s *SQLiteStats) GetRequests() (int64, error) {
	var requests int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM requests").Scan(&requests)
	if err != nil {
		return 0, err
	}
	return requests, nil
}


