package common

import (
	"async-task-hub/src/types"
	"time"
)

// 将 int 时间戳转换为字符串
func FormatTimestamp(timestamp int) string {
	if timestamp <= 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}

// 格式化 *time.Time 类型
func FormatTime(t *types.Customtime) string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

// 格式化时间字符串
func FormatDatetime(datetime string) types.Customtime {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, time.Local)
	if err != nil {
		return types.Customtime{}
	}
	return types.Customtime(t)
}
