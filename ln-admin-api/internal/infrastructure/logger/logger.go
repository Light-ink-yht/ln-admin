package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Light-ink-yht/ln-admin/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// Init 初始化日志
func Init() error {
	// 确定日志级别
	var level zapcore.Level
	switch config.Cfg.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"

	var encoder zapcore.Encoder
	if config.Cfg.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 写入器配置
	var writeSyncer zapcore.WriteSyncer

	// 创建日志目录
	if config.Cfg.Log.Output == "file" || config.Cfg.Log.Output == "both" {
		if err := os.MkdirAll(config.Cfg.Log.Path, os.ModePerm); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}

		logFile := filepath.Join(config.Cfg.Log.Path, config.Cfg.Log.Filename)
		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    config.Cfg.Log.MaxSize,
			MaxBackups: config.Cfg.Log.MaxBackups,
			MaxAge:     config.Cfg.Log.MaxAge,
			Compress:   config.Cfg.Log.Compress,
		}

		if config.Cfg.Log.Output == "both" {
			writeSyncer = zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(lumberjackLogger),
			)
		} else {
			writeSyncer = zapcore.AddSync(lumberjackLogger)
		}
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// Sync 同步日志
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Fields 返回字段（兼容性函数）
func Fields(fields ...zap.Field) []zap.Field {
	return fields
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
