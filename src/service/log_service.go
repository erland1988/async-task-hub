package service

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"context"
	"errors"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) CreateLog(ctx context.Context, adminID int, operation string, data interface{}) error {
	if operation == "" {
		return errors.New("操作或详情不能为空")
	}
	details, _ := common.Struct2Json(data)

	Log := model.Log{
		AdminID:   adminID,
		Operation: operation,
		Details:   details,
	}

	if err := global.DB.Create(&Log).Error; err != nil {
		return err
	}
	return nil
}

func (s *LogService) GetLogById(ctx context.Context, id uint) (*model.Log, error) {
	var log model.Log
	if err := global.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (s *LogService) ListLogs(ctx context.Context, limit, offset int) ([]model.Log, error) {
	var logs []model.Log
	if err := global.DB.Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
