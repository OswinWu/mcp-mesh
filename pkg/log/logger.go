package logger

import (
	"mcp-mesh/config"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_logger *zap.Logger
	once    sync.Once
)

// levelMap 将字符串日志级别映射到zapcore.Level
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLogLevel 获取日志级别
func getLogLevel(lvl string) zapcore.Level {
	if level, exists := levelMap[lvl]; exists {
		return level
	}
	// 默认Info级别
	return zapcore.InfoLevel
}

// Init 初始化日志系统
func Init(config config.LogConfig) {
	once.Do(func() {
		// 确保日志目录存在
		logDir := filepath.Dir(config.FilePath)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				panic("failed to create log directory: " + err.Error())
			}
		}

		// 设置日志级别
		level := getLogLevel(config.Level)

		// 文件输出配置
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,    // MB
			MaxBackups: config.MaxBackups, // 最大备份数
			MaxAge:     config.MaxAge,     // 天
			Compress:   config.Compress,   // 是否压缩
			LocalTime:  true,              // 使用本地时间
		})

		// 标准输出
		stdoutWriter := zapcore.AddSync(os.Stdout)

		// 编码器配置
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		// 创建Core
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			zap.NewAtomicLevelAt(level),
		)

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			stdoutWriter,
			zap.NewAtomicLevelAt(level),
		)

		// 合并多个Core
		core := zapcore.NewTee(fileCore, consoleCore)

		// 创建Logger
		_logger = zap.New(
			core,
			zap.AddCaller(),                       // 添加调用者信息
			zap.AddCallerSkip(1),                  // 跳过一层调用栈
			zap.AddStacktrace(zapcore.ErrorLevel), // Error及以上级别添加堆栈跟踪
		)

		// 替换全局logger
		zap.ReplaceGlobals(_logger)
	})
}

// getLogger 获取全局logger实例
func getLogger() *zap.Logger {
	if _logger == nil {
		panic("logger not initialized, please call Init() first")
	}
	return _logger
}

// Sync 同步缓冲数据到磁盘，应用程序退出前调用
func Sync() {
	if _logger != nil {
		_ = _logger.Sync()
	}
}

func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	getLogger().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getLogger().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}
