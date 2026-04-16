package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type logger struct {
	mu     sync.Mutex
	logDir string
	files  map[string]*os.File // "info", "error", "audit"
	done   chan struct{}
}

var l *logger

func Init(logDir string) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	l = &logger{
		logDir: logDir,
		files:  make(map[string]*os.File),
		done:   make(chan struct{}),
	}

	// 立即清理一次
	cleanOldLogs(logDir, 30)

	// 启动定时清理（每24小时）
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanOldLogs(logDir, 30)
			case <-l.done:
				return
			}
		}
	}()

	return nil
}

func Close() {
	if l == nil {
		return
	}
	close(l.done)
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, f := range l.files {
		f.Close()
	}
}

func Info(format string, args ...interface{}) {
	write("INFO", format, args...)
}

func Error(format string, args ...interface{}) {
	write("ERROR", format, args...)
}

func Audit(format string, args ...interface{}) {
	write("AUDIT", format, args...)
}

func write(level, format string, args ...interface{}) {
	if l == nil {
		// logger 未初始化，降级到 stdout
		ts := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("%s [%s] %s\n", ts, level, fmt.Sprintf(format, args...))
		return
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", ts, level, fmt.Sprintf(format, args...))

	// 输出到 stdout
	fmt.Print(msg)

	// 输出到文件
	l.mu.Lock()
	defer l.mu.Unlock()

	f := l.getFile(level)
	if f != nil {
		f.WriteString(msg)
	}
}

func (l *logger) getFile(level string) *os.File {
	key := strings.ToLower(level)
	today := time.Now().Format("2006-01-02")

	// 检查当前文件是否是今天的
	if f, ok := l.files[key]; ok {
		if f.Name() == filepath.Join(l.logDir, fmt.Sprintf("%s-%s.log", key, today)) {
			return f
		}
		// 日期变了，关闭旧文件
		f.Close()
	}

	// 打开今天的文件
	filename := filepath.Join(l.logDir, fmt.Sprintf("%s-%s.log", key, today))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file %s: %v\n", filename, err)
		return nil
	}

	l.files[key] = f
	return f
}

func cleanOldLogs(logDir string, days int) {
	cutoff := time.Now().AddDate(0, 0, -days)
	patterns := []string{"info-*.log", "error-*.log", "audit-*.log"}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(logDir, pattern))
		if err != nil {
			continue
		}
		for _, path := range matches {
			base := filepath.Base(path)
			// 从文件名提取日期，格式: level-2006-01-02.log
			parts := strings.SplitN(base, "-", 2)
			if len(parts) < 2 {
				continue
			}
			dateStr := strings.TrimSuffix(parts[1], ".log")
			t, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}
			if t.Before(cutoff) {
				os.Remove(path)
			}
		}
	}
}
