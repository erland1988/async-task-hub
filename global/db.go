package global

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitializeDB 初始化数据库
func InitializeDB(config Config, logger *zap.Logger) error {
	return initMySqlGorm(config, logger)
}

// 构建 DSN（数据源名称）
func buildDSN(config Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_DATABASE,
		config.DATABASE_CHARSET,
	)
}

// 使用全局 Logger 配置 GORM 的日志模式
func initMySqlGorm(config Config, zapLogger *zap.Logger) error {
	dsn := buildDSN(config)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	// 根据配置文件中的日志模式设置 LogLevel
	var logLevel logger.LogLevel
	switch config.DATABASE_LOG_MODE {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info // 默认使用 Info 级别
	}

	newLogger := logger.New(
		zapLoggerAdapter{zapLogger}, // 使用自定义的 zapLoggerAdapter
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel, // 使用配置中的日志级别
			IgnoreRecordNotFoundError: true,
			Colorful:                  false, // 使用日志文件，不使用彩色输出
		},
	)

	var db *gorm.DB
	var err error
	maxRetries := 10
	retryDelay := 300 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
			Logger:                                   newLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			break
		}
		zapLogger.Warn("mysql connect failed", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(retryDelay)
	}
	if err != nil {
		return fmt.Errorf("mysql connect failed after %d attempts: %v", maxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		zapLogger.Error("failed to get sqlDB from gorm DB", zap.Error(err))
		return err
	}
	sqlDB.SetMaxIdleConns(config.DATABASE_MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(config.DATABASE_MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

// 自定义的 zapLoggerAdapter 类型，封装 zap.Logger
type zapLoggerAdapter struct {
	*zap.Logger
}

func (z zapLoggerAdapter) Printf(format string, args ...interface{}) {
	z.Logger.Info(fmt.Sprintf(format, args...))
}

func (z zapLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	return z
}

func (z zapLoggerAdapter) Info(ctx context.Context, msg string, args ...interface{}) {
	z.Logger.Info(msg, zap.Any("args", args))
}

func (z zapLoggerAdapter) Warn(ctx context.Context, msg string, args ...interface{}) {
	z.Logger.Warn(msg, zap.Any("args", args))
}

func (z zapLoggerAdapter) Error(ctx context.Context, msg string, args ...interface{}) {
	z.Logger.Error(msg, zap.Any("args", args))
}

func (z zapLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, _ := fc()
		z.Logger.Error("sql error", zap.String("sql", sql), zap.Error(err))
	}
}
