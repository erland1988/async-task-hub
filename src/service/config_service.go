package service

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"errors"
)

type configMap map[string]string

type ConfigService struct {
	keys      []string
	configMap []configMap
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		keys: []string{
			"notice",           // 公告
			"executor_timeout", // 执行器请求超时时间
			"clear_time",       // 清理间隔时间
		},
	}
}

func (s *ConfigService) GetConfigMap() (configMap, error) {
	var configs []model.Config
	if err := global.DB.Where("`key` IN (?)", s.keys).Find(&configs).Error; err != nil {
		return nil, err
	}
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}
	return configMap, nil
}

func (s *ConfigService) GetKeys() []string {
	return s.keys
}

func (s *ConfigService) GetConfig(key string) (string, error) {
	var config model.Config
	if err := global.DB.Where("`key` =?", key).First(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

func (s *ConfigService) UpdateConfig(key, value string) error {
	var config model.Config
	if key == "executor_timeout" {
		executorTimeout := common.Str2Int(value)
		if executorTimeout < 3 {
			return errors.New("执行器超时时间不能小于3秒")
		}
		if executorTimeout > 3600 {
			return errors.New("执行器超时时间不能大于3600秒")
		}
	}
	if key == "clear_time" {
		clearTime := common.Str2Int(value)
		if clearTime < 1 {
			return errors.New("清理间隔时间不能小于1小时")
		}
		if clearTime > 72 {
			return errors.New("清理间隔时间不能大于72小时")
		}
	}
	if err := global.DB.Where("`key` = ?", key).First(&config).Error; err != nil {
		return err
	}
	config.Value = value
	if err := global.DB.Where("`key` = ?", key).Save(&config).Error; err != nil {
		return err
	}
	return nil
}
