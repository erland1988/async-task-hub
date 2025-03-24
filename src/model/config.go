package model

type Config struct {
	Id    int    `gorm:"type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`
	Key   string `gorm:"type:varchar(64);not null;uniqueIndex:key" json:"key"`
	Value string `gorm:"type:text;not null" json:"value"`
}

func (m *Config) TableName() string {
	return "configs"
}
