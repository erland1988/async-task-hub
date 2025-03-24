package global

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	REDIS  *redis.Client
	ctx    = context.Background()
	prefix = ""
)

// InitializeRedis 用于初始化 Redis 客户端
func InitializeRedis(config Config, logger *zap.Logger) error {
	var err error
	maxRetries := 10                     // 最大重试次数
	retryDelay := 300 * time.Millisecond // 重试延迟时间

	for i := 0; i < maxRetries; i++ {
		REDIS = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
			Password: config.REDIS_PASSWORD,
			DB:       config.REDIS_DB,
		})

		_, err = REDIS.Ping(ctx).Result()
		if err == nil {
			// 连接成功
			return nil
		}

		// 连接失败，记录详细错误信息并等待重试
		logger.Warn("redis connect failed", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(retryDelay)
	}

	prefix = config.REDIS_PREFIX // 初始化时设置前缀
	return fmt.Errorf("redis connect failed after %d attempts: %v", maxRetries, err)
}

func CacheKey(key string) string {
	return prefix + ":" + key
}

// GetFromCache 从缓存中获取数据
func GetFromCache(key string, dest interface{}) error {
	val, err := REDIS.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return err
	}
	return nil
}

// SetToCache 将数据存入缓存
func SetToCache(key string, data interface{}, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := REDIS.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}
