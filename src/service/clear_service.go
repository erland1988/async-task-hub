package service

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"context"
	"go.uber.org/zap"
	"time"
)

type ClearService struct {
	clearTime        time.Duration
	clearTimelockKey string
}

func NewClearService() *ClearService {
	return &ClearService{
		clearTime:        24 * time.Hour,
		clearTimelockKey: global.CacheKey("clear_time_lock"),
	}
}

func (s *ClearService) StartClearMonitor(ctx context.Context) {
	s.Clear()
	//ticker := time.NewTicker(30 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done(): // 优雅退出
			global.Logger.Info("StartClearMonitor stopped.")
			return
		case <-ticker.C:
			s.Clear()
		}
	}
}

func (s *ClearService) Clear() {
	clearTimeVal, err := NewConfigService().GetConfig("clear_time")
	if err == nil {
		clearTimeInt := common.Str2Int(clearTimeVal)
		if clearTimeInt > 0 {
			clearTime := time.Duration(clearTimeInt) * time.Hour
			if clearTime != s.clearTime {
				s.clearTime = clearTime
				global.REDIS.Expire(context.Background(), s.clearTimelockKey, s.clearTime)
			}
		}
	}
	success, err := global.REDIS.SetNX(context.Background(), s.clearTimelockKey, "locked", s.clearTime).Result()
	if err != nil || !success {
		return
	}
	s.clearLogin()
	return
}

func (s *ClearService) clearLogin() {
	global.Logger.Info("clearLogin", zap.Int("clearTime", int(s.clearTime.Hours())))
	global.DB.Where("expires_at <?", time.Now()).Delete(&model.Login{})
	return
}
