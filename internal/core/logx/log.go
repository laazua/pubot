package logx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"pubot/internal/core/config"
)

type Level = slog.Level

var LogLevelMap = map[string]Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

//
// ==== 文件切割日志处理器 ===
//

type RotateFileHandler struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

type RotateFileHandlerOptions struct {
	MaxSize         int64
	BackupCount     int
	ReplaceExisting bool
}

func NewRotateFileHandler(filePath string, opts RotateFileHandlerOptions) (*RotateFileHandler, error) {
	if opts.MaxSize <= 0 {
		opts.MaxSize = 10 * 1024 * 1024
	}
	if opts.BackupCount < 0 {
		opts.BackupCount = 5
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	var file *os.File
	var err error
	if opts.ReplaceExisting {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	} else {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("get file info failed: %w", err)
	}

	return &RotateFileHandler{
		file:        file,
		filePath:    filePath,
		maxSize:     opts.MaxSize,
		currentSize: info.Size(),
		backupCount: opts.BackupCount,
	}, nil
}

func (h *RotateFileHandler) Write(p []byte) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentSize+int64(len(p)) > h.maxSize {
		if err := h.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := h.file.Write(p)
	if err == nil {
		h.currentSize += int64(n)
		_ = h.file.Sync()
	}
	return n, err
}

func (h *RotateFileHandler) rotate() error {
	if err := h.file.Close(); err != nil {
		return fmt.Errorf("关闭当前日志文件失败: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.%s", h.filePath, timestamp)
	if err := os.Rename(h.filePath, backupPath); err != nil {
		return fmt.Errorf("rename file failed: %w", err)
	}

	file, err := os.OpenFile(h.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create new file failed: %w", err)
	}

	h.file = file
	h.currentSize = 0
	go h.cleanupOldBackups()
	return nil
}

func (h *RotateFileHandler) cleanupOldBackups() {
	if h.backupCount <= 0 {
		return
	}

	dir := filepath.Dir(h.filePath)
	base := filepath.Base(h.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasPrefix(name, base+".") {
				backups = append(backups, name)
			}
		}
	}

	if len(backups) > h.backupCount {
		sort.Strings(backups)
		for i := 0; i < len(backups)-h.backupCount; i++ {
			_ = os.Remove(filepath.Join(dir, backups[i]))
		}
	}
}

func (h *RotateFileHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.file.Close()
}

//
// ==== Logger 封装 ===
//

type Logger struct {
	*slog.Logger
	fileHandler *RotateFileHandler
	consoleOut  io.Writer
}

type LoggerOptions struct {
	FileOptions  RotateFileHandlerOptions
	LogConsoleOn bool
	LogFormat    string
	Level        slog.Level
	AddSource    bool
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: hs}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: hs}
}

func newLogger(filePath string, opts LoggerOptions) (*Logger, error) {
	fileHandler, err := NewRotateFileHandler(filePath, opts.FileOptions)
	if err != nil {
		return nil, err
	}

	fileLogHandler := slog.NewJSONHandler(fileHandler, &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
	})

	var handlers []slog.Handler
	handlers = append(handlers, fileLogHandler)

	if opts.LogConsoleOn {
		if opts.LogFormat == "json" {
			handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     opts.Level,
				AddSource: opts.AddSource,
			}))
		} else {
			handlers = append(handlers, slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     opts.Level,
				AddSource: opts.AddSource,
			}))
		}
	}

	var rootHandler slog.Handler
	if len(handlers) == 1 {
		rootHandler = handlers[0]
	} else {
		rootHandler = &multiHandler{handlers: handlers}
	}

	logger := slog.New(rootHandler)

	return &Logger{
		Logger:      logger,
		fileHandler: fileHandler,
		consoleOut:  os.Stdout,
	}, nil
}

//
// ==== 全局封装 ===
//

var (
	defaultLogger *Logger
	mu            sync.Mutex
)

func defaultOptions() LoggerOptions {
	return LoggerOptions{
		FileOptions: RotateFileHandlerOptions{
			MaxSize:     10 * 1024 * 1024,
			BackupCount: 7,
		},
		LogConsoleOn: true,
		LogFormat:    "text",
		Level:        LogLevelMap[config.Get().LogLevel],
		AddSource:    false, // 是否打印调用日志记录器的文件和具体行
	}
}

// Ensure 懒加载返回全局 Logger
func ensure() *Logger {
	if defaultLogger == nil {
		mu.Lock()
		defer mu.Unlock()
		if defaultLogger == nil {
			logger, err := newLogger("logs/app.log", defaultOptions())
			if err != nil {
				panic(err)
			}
			defaultLogger = logger
		}
	}
	return defaultLogger
}

// Init 自定义配置LoggerOptions,进行实例化
func Init(filePath string, opts LoggerOptions) error {
	mu.Lock()
	defer mu.Unlock()
	l, err := newLogger(filePath, opts)
	if err != nil {
		return err
	}
	defaultLogger = l
	return nil
}

func Close() error {
	if defaultLogger != nil {
		return ensure().fileHandler.Close()
	}
	return nil
}

//
// 简单调用:进一步封装
//

func Debug(msg string, args ...any)                         { ensure().Debug(msg, args...) }
func DebugCtx(ctx context.Context, msg string, args ...any) { ensure().DebugContext(ctx, msg, args...) }
func Info(msg string, args ...any)                          { ensure().Info(msg, args...) }
func InfoCtx(ctx context.Context, msg string, args ...any)  { ensure().InfoContext(ctx, msg, args...) }
func Warn(msg string, args ...any)                          { ensure().Warn(msg, args...) }
func WarnCtx(ctx context.Context, msg string, args ...any)  { ensure().WarnContext(ctx, msg, args...) }
func Error(msg string, args ...any)                         { ensure().Error(msg, args...) }
func ErrorCtx(ctx context.Context, msg string, args ...any) { ensure().ErrorContext(ctx, msg, args...) }

func String(key, value string) slog.Attr  { return slog.Attr{Key: key, Value: slog.StringValue(value)} }
func Int(key string, value int) slog.Attr { return slog.Attr{Key: key, Value: slog.IntValue(value)} }
func Int64(key string, value int64) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Int64Value(value)}
}
func Uint64(key string, value uint64) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Uint64Value(value)}
}
func Time(key string, v time.Time) slog.Attr { return slog.Attr{Key: key, Value: slog.TimeValue(v)} }
func Duration(key string, v time.Duration) slog.Attr {
	return slog.Attr{Key: key, Value: slog.DurationValue(v)}
}
func Any(key string, value any) slog.Attr { return slog.Attr{Key: key, Value: slog.AnyValue(value)} }
