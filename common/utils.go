package common

import (
	"encoding/json"
	"errors"
)

func Str2Int(str string) int {
	var num int
	for _, c := range str {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
	}
	return num
}

func Struct2Json(data interface{}) (string, error) {
	if data == nil {
		return "", errors.New("数据不能为空")
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func InArray[T comparable](needle T, haystack []T) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
