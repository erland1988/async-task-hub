package service

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"errors"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) CreateLog(adminID int, operation string, data interface{}) error {
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
