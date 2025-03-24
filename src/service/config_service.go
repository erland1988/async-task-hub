package service

import (
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"go.uber.org/zap"
)

type configMap map[string]string

type ConfigService struct {
	keys      []string
	configMap []configMap
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		keys: []string{
			"notice",
		},
	}
}

func (s *ConfigService) GetConfigMap() ([]configMap, error) {
	var configs []model.Config
	if err := global.DB.Where("`key` IN (?)", s.keys).Find(&configs).Error; err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		return nil, err
	}
	var configMap []configMap
	for _, config := range configs {
		configMap = append(configMap, map[string]string{
			config.Key: config.Value,
		})
	}
	return configMap, nil
}

func (s *ConfigService) GetKeys() []string {
	return s.keys
}

func (s *ConfigService) GetConfig(key string) (string, error) {
	var config model.Config
	if err := global.DB.Where("`key` =?", key).First(&config).Error; err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		return "", err
	}
	return config.Value, nil
}

func (s *ConfigService) UpdateConfig(key, value string) error {
	var config model.Config
	if err := global.DB.Where("`key` = ?", key).First(&config).Error; err != nil {
		global.Logger.Warn("获取配置失败", zap.Error(err))
		return err
	}
	config.Value = value
	if err := global.DB.Save(&config).Error; err != nil {
		global.Logger.Warn("更新配置失败", zap.Error(err))
		return err
	}
	return nil
}
