package service

import (
	"async-task-hub/global"
	"async-task-hub/src/model"
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
	if err := global.DB.Where("`key` = ?", key).First(&config).Error; err != nil {
		return err
	}
	config.Value = value
	if err := global.DB.Where("`key` = ?", key).Save(&config).Error; err != nil {
		return err
	}
	return nil
}
