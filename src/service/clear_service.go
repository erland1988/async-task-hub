package service

import (
	"async-task-hub/global"
	"async-task-hub/src/model"
	"context"
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
	ticker := time.NewTicker(s.clearTime)
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
	success, err := global.REDIS.SetNX(context.Background(), s.clearTimelockKey, "locked", s.clearTime).Result()
	if err != nil || !success {
		return
	}
	s.clearLogin()
	return
}

func (s *ClearService) clearLogin() {
	global.DB.Where("expires_at <?", time.Now()).Delete(&model.Login{})
	return
}
