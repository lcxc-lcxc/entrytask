package utils

import (
	"strconv"
	"strings"
)

// ConvertRedisKeyToUintId 从redis key转换为uint
func ConvertRedisKeyToUintId(key string) (uint, error) {
	Uint64Id, err := strconv.ParseUint(key[strings.Index(key, ":")+1:len(key)], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(Uint64Id), nil
}

// ConvertUintIdToRedisKey 从uint转换为redis key
func ConvertUintIdToRedisKey(prefix string, id uint) string {
	return prefix + ":" + strconv.FormatUint(uint64(id), 10)
}

// ConvertRedisKeyToString 从redis key 转换为 string
func ConvertRedisKeyToString(key string) string {
	return key[strings.Index(key, ":")+1 : len(key)]

}

// ConvertStringToRedisKey 从string 转换为redis key
func ConvertStringToRedisKey(prefix string, str string) string {
	return prefix + ":" + str

}
