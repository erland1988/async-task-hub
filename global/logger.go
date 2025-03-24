package global

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func InitializeLogger(config Config) {
	// 配置 encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 时间格式
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 日志级别
	var logLevel zapcore.Level
	switch config.APP_LOG_MODE {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel // 默认日志级别
	}

	var cores []zapcore.Core

	// 终端输出
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel)
	cores = append(cores, consoleCore)

	// 文件输出
	if config.APP_LOG_FILENAME != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.APP_LOG_FILENAME,
			MaxSize:    10,   // 每个日志文件最大 10 MB
			MaxBackups: 7,    // 保留 7 个备份
			MaxAge:     30,   // 文件最大保存 30 天
			Compress:   true, // 是否压缩备份日志
		})
		fileCore := zapcore.NewCore(encoder, fileWriter, logLevel)
		cores = append(cores, fileCore)
	}

	// 使用 zapcore.NewTee 将日志输出到多个目的地
	core := zapcore.NewTee(cores...)

	// 创建 Logger
	Logger = zap.New(core)
}
